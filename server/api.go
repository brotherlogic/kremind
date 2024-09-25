package server

import (
	"context"
	"time"

	ghbclient "github.com/brotherlogic/githubridge/client"
	rstore_client "github.com/brotherlogic/rstore/client"

	pb "github.com/brotherlogic/kremind/proto"
)

type Server struct {
	rclient  rstore_client.RStoreClient
	ghclient ghbclient.GithubridgeClient
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
	return &pb.AddReminderResponse{Id: id}, s.saveReminder(ctx, r)
}

func (s *Server) ListReminders(ctx context.Context, req *pb.ListRemindersRequest) (*pb.ListRemindersResponse, error) {
	reminders, err := s.loadReminders(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ListRemindersResponse{Reminders: reminders}, nil
}
