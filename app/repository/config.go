package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/aicacia/pickup/app/env"
	"github.com/lib/pq"
)

type ConfigRowST struct {
	Id        int32     `db:"id"`
	Key       string    `db:"key"`
	Value     string    `db:"value"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

func GetConfigs() ([]ConfigRowST, error) {
	return All[ConfigRowST]("SELECT * FROM configs;")
}

func CreateListener(channel string) (*pq.Listener, error) {
	listener := pq.NewListener(env.GetDatabaseUrl(), 10*time.Second, time.Minute, listenerEventCallback)
	if listener == nil {
		return nil, fmt.Errorf("listener is nil")
	}
	err := listener.Listen(channel)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func listenerEventCallback(event pq.ListenerEventType, err error) {
	if err != nil {
		log.Printf("ListenerEvent: %s\n", err)
		return
	}
	switch event {
	case pq.ListenerEventConnected:
		log.Println("Connected")
	case pq.ListenerEventDisconnected:
		log.Println("Disconnected")
	case pq.ListenerEventReconnected:
		log.Println("Reconnected")
	case pq.ListenerEventConnectionAttemptFailed:
		log.Println("Connection attempt failed")
	}
}
