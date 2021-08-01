package mq

type TaskMgr interface {
	Push(channel string, task *Task) error
	Shift(channel string) (*Task, error)
	Remove(channel string, id string) error
	Destroy(channel string) error
}
