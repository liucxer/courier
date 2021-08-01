package memtaskmgr

import (
	"container/list"
	"sync"

	"github.com/liucxer/courier/mq"
)

func NewMemTaskMgr() *MemTaskMgr {
	return &MemTaskMgr{
		list: list.New(),
		m:    map[string]*list.Element{},
	}
}

type MemTaskMgr struct {
	m    map[string]*list.Element
	list *list.List
	rw   sync.RWMutex
}

var _ mq.TaskMgr = (*MemTaskMgr)(nil)

func (mgr *MemTaskMgr) Push(channel string, task *mq.Task) error {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()

	id := toKey(channel, task.Id)

	mgr.m[id] = mgr.list.PushBack(task)
	return nil
}

func (mgr *MemTaskMgr) Shift(channel string) (*mq.Task, error) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()

	e := mgr.list.Front()
	if e == nil {
		return nil, nil
	}
	mgr.list.Remove(e)

	task, ok := e.Value.(*mq.Task)
	if !ok {
		return nil, nil
	}

	id := toKey(channel, task.Id)

	if _, ok := mgr.m[id]; !ok {
		return nil, nil
	}

	return task, mgr.remove(id)
}

func (mgr *MemTaskMgr) Remove(channel, id string) error {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	return mgr.remove(toKey(channel, id))
}

func (mgr *MemTaskMgr) Destroy(channel string) error {
	*mgr = *NewMemTaskMgr()
	return nil
}

func (mgr *MemTaskMgr) remove(id string) error {
	e := mgr.m[id]
	if e != nil {
		mgr.list.Remove(e)
		delete(mgr.m, id)
	}
	return nil
}

func toKey(channel, id string) string {
	return channel + "::" + id
}
