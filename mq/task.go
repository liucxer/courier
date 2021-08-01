package mq

import (
	"github.com/google/uuid"
)

func NewTask(subject string, argv []byte, id string) *Task {
	if id == "" {
		id = uuid.New().String()
	}
	return &Task{
		Subject: subject,
		Argv:    argv,
		Id:      id,
	}
}
