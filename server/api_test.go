package server

import (
	"context"
	"testing"
	"time"

	"github.com/brotherlogic/kremind/db"

	pb "github.com/brotherlogic/kremind/proto"
)

func TestListReminderWithTIme(t *testing.T) {
	ctx := context.Background()

	s := NewServer(db.GetTestDB())

	_, err := s.AddReminder(ctx, &pb.AddReminderRequest{
		RepeatInSeconds: 200,
		StartTime:       time.Now().Add(time.Hour).Unix(),
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

	if len(rs.GetReminders()) != 1 {
		t.Fatalf("Reminder was not added")
	}

	rs, err = s.ListReminders(ctx, &pb.ListRemindersRequest{TimestampSeconds: time.Now().Unix()})
	if err != nil {
		t.Fatalf("Unable to list reminders")
	}
	if len(rs.GetReminders()) != 0 {
		t.Errorf("Should have listed zero reminders: %v", rs)
	}

	rs, err = s.ListReminders(ctx, &pb.ListRemindersRequest{TimestampSeconds: time.Now().Add(time.Hour * 2).Unix()})
	if err != nil {
		t.Fatalf("Unable to list reminders")
	}
	if len(rs.GetReminders()) != 1 {
		t.Errorf("Should have listed one reminders: %v", rs)
	}

}

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
