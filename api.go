package main

import (
	"context"
	"time"
)

func (s *Server) AddReminder(ctx context.Context, req *pb.AddReminderRequest) (*pb.AddReminderResponse, error) {
	reminders, err := s.loadReminders(ctx)
	if err != nil {
		return nil, err
	}

	for _, reminder := range reminders {
		if reminder.GetReminder() == req.GetReminder() {
			return nil, status.Errorf(codes.AlreadyExists, "already exists in the db (%v)", reminder)
		}
	}

	id := time.Now().UnixNano()
	reminders = append(reminders, &pb.Reminder{
		Id:              id,
		StartTime:       req.GetStartTime(),
		RepeatInSeconds: req.GetRepeatInSeconds(),
		Reminder:        req.GetReminder(),
		Source:          req.GetSource(),
	})
	return &pb.AddReminderResponse{Id: id}, s.saveReminders(ctx, reminders)
}

func (s *Server) ListReminders(ctx context.Context, req *pb.ListRemindersRequest) (*pb.ListRemindersResponse, error) {
	reminders, err := s.loadReminders(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ListRemindersResponse{Reminders: reminders}, nil
}
