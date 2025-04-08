package db

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	pstore_client "github.com/brotherlogic/pstore/client"

	pb "github.com/brotherlogic/kremind/proto"
	pspb "github.com/brotherlogic/pstore/proto"
)

const (
	REMINDER_KEY = "reminders/reminder/"
)

type DB struct {
	pclient pstore_client.PStoreClient
	lock    sync.Mutex
}

func GetDB() *DB {
	client, err := pstore_client.GetClient()
	if err != nil {
		return nil
	}
	return &DB{pclient: client}
}

func GetTestDB() *DB {
	return &DB{
		pclient: pstore_client.GetTestClient(),
	}
}

func (d *DB) SaveReminder(ctx context.Context, r *pb.Reminder) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	_, err = d.pclient.Write(ctx, &pspb.WriteRequest{
		Key:   fmt.Sprintf("%v%v", REMINDER_KEY, r.GetId()),
		Value: &anypb.Any{Value: data},
	})
	return err
}

func (d *DB) DeleteReminder(ctx context.Context, id int64) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	_, err := d.pclient.Delete(ctx, &pspb.DeleteRequest{
		Key: fmt.Sprintf("%v%v", REMINDER_KEY, id),
	})
	return err
}

func (d *DB) LoadReminders(ctx context.Context) ([]*pb.Reminder, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	keys, err := d.pclient.GetKeys(ctx, &pspb.GetKeysRequest{
		Prefix: REMINDER_KEY,
	})
	if err != nil {
		return nil, err
	}

	var reminders []*pb.Reminder
	for _, key := range keys.GetKeys() {
		val, err := d.pclient.Read(ctx, &pspb.ReadRequest{
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
