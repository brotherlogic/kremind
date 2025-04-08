package server

import (
	"context"
	"testing"
	"time"

	"github.com/brotherlogic/kremind/db"

	pb "github.com/brotherlogic/kremind/proto"
)

func TestDeleteReminder(t *testing.T) {
	ctx := context.Background()

	s := NewServer(db.GetTestDB())

	r, err := s.AddReminder(ctx, &pb.AddReminderRequest{
		RepeatInSeconds: 200,
		StartTime:       time.Now().Unix(),
		Reminder:        "Test",
		Source:          "test",
	})
	if err != nil {
		t.Fatalf("Unable to add reminder: %v", err)
	}

	rs, err := s.ListReminders(ctx, &pb.ListRemindersRequest{})
	if err != nil {
		t.Fatalf("Unable to list reminders: %v", err)
	}

	if len(rs.GetReminders()) != 1 || rs.GetReminders()[0].GetReminder() != "Test" {
		t.Fatalf("Reminder was not added: %v", rs)
	}

	_, err = s.DeleteReminder(ctx, &pb.DeleteReinderRequest{
		Id: r.GetId(),
	})

	rs, err = s.ListReminders(ctx, &pb.ListRemindersRequest{})
	if err != nil {
		t.Fatalf("Unable to get reminders: %v", err)
	}

	if len(rs.GetReminders()) != 0 {
		t.Errorf("Reminder was not deleted (%v) -> %v", r, rs)
	}
}
