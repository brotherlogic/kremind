package server

import (
	"context"
	"time"

	db "github.com/brotherlogic/kremind/db"

	pb "github.com/brotherlogic/kremind/proto"
)

type Server struct {
	db *db.DB
}

func NewServer(db *db.DB) *Server {
	if db == nil {
		return nil
	}
	return &Server{db: db}
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

	var rr []*pb.Reminder
	for _, r := range reminders {
		if req.GetTimestampSeconds() == 0 {
			rr = append(rr, r)
		} else {
			if r.GetLastRunTime() > 0 && r.GetLastRunTime()+r.GetRepeatInSeconds() < req.GetTimestampSeconds() {
				rr = append(rr, r)
			} else if r.GetLastRunTime() == 0 && r.GetStartTime() < req.GetTimestampSeconds() {
				rr = append(rr, r)
			}
		}
	}

	return &pb.ListRemindersResponse{Reminders: rr}, nil
}

func (s *Server) DeleteReminder(ctx context.Context, req *pb.DeleteReinderRequest) (*pb.DeleteReminderResponse, error) {
	return &pb.DeleteReminderResponse{}, s.db.DeleteReminder(ctx, req.GetId())
}
