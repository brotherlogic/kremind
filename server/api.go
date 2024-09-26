package server

import (
	"context"
	"time"

	db "github.com/brotherlogic/kremind/db"

	pb "github.com/brotherlogic/kremind/proto"
)

type Server struct {
	db db.DB
}

func (s *Server) AddReminder(ctx context.Context, req *pb.AddReminderRequest) (*pb.AddReminderResponse, error) {
	id := time.Now().UnixNano()
	r := &pb.Reminder{
		Id:              id,
		StartTime:       req.GetStartTime(),
		RepeatInSeconds: req.GetRepeatInSeconds(),
		Reminder:        req.GetReminder(),
		Source:          req.GetSource(),
	}
	return &pb.AddReminderResponse{Id: id}, s.db.SaveReminder(ctx, r)
}

func (s *Server) ListReminders(ctx context.Context, req *pb.ListRemindersRequest) (*pb.ListRemindersResponse, error) {
	reminders, err := s.db.LoadReminders(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ListRemindersResponse{Reminders: reminders}, nil
}
