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
	abanson bool
	timer   *time.Timer
}

func NewTimer(t time.Duration) *ATimer {
	return &ATimer{
		abanson: false,
		timer:   time.NewTimer(t),
	}
}

func (a *ATimer) Wait() {
	<-a.timer.C
}

type Runner struct {
	lastActivity int64
	tMap         map[int64]*ATimer
	mapLock      *sync.RWMutex
	ghclient     ghbclient.GithubridgeClient
	db           db.DB
}

func getNextRunTime(re *pb.Reminder, now time.Time) time.Time {
	if re.GetLastRunTime() == 0 {
		return time.Unix(re.GetStartTime(), 0)
	}

	// If we previously failed, wait the backoff and retry
	if re.GetLastFailure() != "" {
		return now.Add(BACKOFF)
	}

	return time.Unix(re.GetLastRunTime(), 0).Add(time.Duration(re.GetRepeatInSeconds()) * time.Second)
}

func (r *Runner) runReminder(ctx context.Context, re *pb.Reminder) error {
	_, err := r.ghclient.CreateIssue(ctx, &ghbpb.CreateIssueRequest{
		Title: re.GetReminder(),
		Body:  "From your reminders",
		Repo:  re.GetSource(),
		User:  "brotherlogic",
	})
	return err
}

func (r *Runner) AddReminder(ctx context.Context, now time.Time, re *pb.Reminder) error {
	if timer, ok := r.tMap[re.GetId()]; ok {
		// Cancel the existing timer
		timer.abanson = true
	}

	// Create and store a new timer
	r.mapLock.Lock()
	r.tMap[re.GetId()] = NewTimer(now.Sub(getNextRunTime(re, now)))
	go func() {
		r.tMap[re.GetId()].Wait()

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

		// There's not much we can do if the save fails, so ignore return here
		r.db.SaveReminder(ctx, re)
	}()
	r.mapLock.Unlock()

	return nil
}
