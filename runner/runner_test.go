package runner

import (
	"context"
	"testing"
	"time"

	"github.com/brotherlogic/kremind/proto"
)

func TestRunReminder(t *testing.T) {
	runner := GetTestRunner()

	// This should run immediately
	ti := time.Now()
	r := &proto.Reminder{
		Reminder:        "testing",
		Id:              123,
		RepeatInSeconds: 5,
	}
	runner.db.SaveReminder(context.Background(), r)
	runner.AddReminder(context.Background(), ti, r)

	time.Sleep(time.Second * 2)

	reminders, err := runner.db.LoadReminders(context.Background())
	if err != nil {
		t.Fatalf("Bad load: %v", err)
	}

	if len(reminders) == 0 {
		t.Fatalf("Reminder was not saved: %v", reminders)
	}

	if reminders[0].GetLastRunTime() < ti.Unix() {
		t.Errorf("Reminder was not run: %v", reminders[0])
	}

	// See if it runs again
	time.Sleep(time.Second * 5)

	reminders2, err := runner.db.LoadReminders(context.Background())
	if err != nil {
		t.Fatalf("Bad load: %v", err)
	}

	if len(reminders2) == 0 {
		t.Fatalf("Reminder was not saved: %v", reminders2)
	}

	if reminders2[0].GetLastRunTime() == reminders[0].GetLastRunTime() {
		t.Errorf("Reminder was not rerun: %v -> %v", reminders[0], reminders2[0])
	}
}
