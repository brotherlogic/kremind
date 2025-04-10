package runner

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	ghbclient "github.com/brotherlogic/githubridge/client"
	db "github.com/brotherlogic/kremind/db"

	ghbpb "github.com/brotherlogic/githubridge/proto"
	pb "github.com/brotherlogic/kremind/proto"
)

const (
	BACKOFF = time.Minute
)

type ATimer struct {
	abandon bool
	timer   *time.Timer
}

func NewTimer(t time.Duration) *ATimer {
	log.Printf("Waiting %v", t)
	tr := &ATimer{
		abandon: false,
		timer:   time.NewTimer(t),
	}

	return tr
}

func (a *ATimer) Wait() {
	<-a.timer.C
}

type Runner struct {
	lastActivity int64
	tMap         map[int64]*ATimer
	mapLock      *sync.RWMutex
	ghclient     ghbclient.GithubridgeClient
	db           *db.DB
}

func GetTestRunner() *Runner {
	return &Runner{
		ghclient: ghbclient.GetTestClient(),
		db:       db.GetTestDB(),
		mapLock:  &sync.RWMutex{},
		tMap:     make(map[int64]*ATimer),
	}
}

func getNextRunTime(re *pb.Reminder, now time.Time) time.Time {
	if re.GetLastRunTime() == 0 {
		if re.GetStartTime() > 0 {
			return time.Unix(re.GetStartTime(), 0)
		}
		return now
	}

	// If we previously failed, wait the backoff and retry
	if re.GetLastFailure() != "" {
		return time.Unix(re.GetLastRunTime(), 0).Add(BACKOFF)
	}

	return time.Unix(re.GetLastRunTime(), 0).Add(time.Duration(re.GetRepeatInSeconds()) * time.Second)
}

func (r *Runner) runReminder(ctx context.Context, re *pb.Reminder) error {
	log.Printf("Running %v", re)
	_, err := r.ghclient.CreateIssue(ctx, &ghbpb.CreateIssueRequest{
		Title: re.GetReminder(),
		Body:  "From your reminders",
		Repo:  re.GetSource(),
		User:  "brotherlogic",
	})
	return err
}

func (r *Runner) Stop() {
	r.mapLock.Lock()
	for _, val := range r.tMap {
		val.abandon = true
	}
}

func (r *Runner) DeleteReminder(ctx context.Context, re int64) error {
	r.mapLock.Lock()
	defer r.mapLock.Unlock()

	if val, ok := r.tMap[re]; ok {
		val.abandon = true
		r.db.DeleteReminder(ctx, re)
		return nil
	}

	return nil
}

func (r *Runner) AddReminder(ctx context.Context, now time.Time, re *pb.Reminder) error {
	if timer, ok := r.tMap[re.GetId()]; ok {
		// Cancel the existing timer
		timer.abandon = true
	}

	// Create and store a new timer
	r.mapLock.Lock()
	log.Printf("%v -> %v and %v", re, now, getNextRunTime(re, now))
	r.tMap[re.GetId()] = NewTimer((getNextRunTime(re, now).Sub(now)))
	go func() {
		// Wait for the reminder to trigger and then delete it
		r.tMap[re.GetId()].Wait()
		val := r.tMap[re.GetId()]
		delete(r.tMap, re.GetId())

		log.Printf("Reminder: %+v", val)

		// Only run this if we've not abandoned it
		if !val.abandon {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			err := r.runReminder(ctx, re)
			log.Printf("Ran Reminder %v -> %v", re, err)

			// Adjust the reminder and reup accordingly
			re.LastRunTime = now.Unix()
			if err != nil {
				re.LastFailure = fmt.Sprintf("%v", err)
			} else {
				re.LastFailure = ""
			}

			log.Printf("Saving %v", re)
			// There's not much we can do if the save fails, so ignore return here
			r.db.SaveReminder(ctx, re)

			// Re up the reminder
			r.AddReminder(ctx, time.Now(), re)
		}
	}()
	r.mapLock.Unlock()

	return nil
}
