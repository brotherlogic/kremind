package db

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	rstore_client "github.com/brotherlogic/rstore/client"

	pb "github.com/brotherlogic/kremind/proto"
	rspb "github.com/brotherlogic/rstore/proto"
)

const (
	REMINDER_KEY = "reminders/reminder/"
)

type DB struct {
	rclient rstore_client.RStoreClient
}

func GetTestDB() *DB {
	return &DB{
		rclient: rstore_client.GetTestClient(),
	}
}

func (d *DB) SaveReminder(ctx context.Context, r *pb.Reminder) error {

	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	_, err = d.rclient.Write(ctx, &rspb.WriteRequest{
		Key:   fmt.Sprintf("%v%v", REMINDER_KEY, r.GetId()),
		Value: &anypb.Any{Value: data},
	})
	return err
}

func (d *DB) LoadReminders(ctx context.Context) ([]*pb.Reminder, error) {
	keys, err := d.rclient.GetKeys(ctx, &rspb.GetKeysRequest{
		Prefix: REMINDER_KEY,
	})
	if err != nil {
		return nil, err
	}

	var reminders []*pb.Reminder
	for _, key := range keys.GetKeys() {
		val, err := d.rclient.Read(ctx, &rspb.ReadRequest{
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
