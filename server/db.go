package server

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	pb "github.com/brotherlogic/kremind/proto"
	rspb "github.com/brotherlogic/rstore/proto"
)

const (
	REMINDER_KEY = "reminders/reminder/"
)

func (s *Server) saveReminder(ctx context.Context, r *pb.Reminder) error {

	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	_, err = s.rclient.Write(ctx, &rspb.WriteRequest{
		Key:   fmt.Sprintf("%v%v", REMINDER_KEY, r.GetId()),
		Value: &anypb.Any{Value: data},
	})
	return err
}

func (s *Server) loadReminders(ctx context.Context) ([]*pb.Reminder, error) {
	keys, err := s.rclient.GetKeys()
	if err != nil {
		return nil, err
	}

	var reminders []*pb.Reminder
	for _, key := range keys {
		val, err := s.rclient.Read(ctx, &rspb.ReadRequest{
			Key: key,
		})
		if err != nil {
			return nil, err
		}
		r := &pb.Reminder{}
		err = proto.Unmarshal(val.GetValue().GetValue(), r)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, r)
	}
	return reminders, nil
}
