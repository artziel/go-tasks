package GoTasks

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrTaskUUIDNotFound error = errors.New("task UUID not found")
var ErrTaskUUIDDuplicated error = errors.New("duplicated task UUID")
var ErrTaskInvalidStatus error = errors.New("invalid task status")
var ErrTaskAlreadyFinish error = errors.New("task already finish")

const (
	TaskRuning  = "RUNING"
	TaskSuccess = "SUCCESS"
	TaskFail    = "FAIL"
)

type Task struct {
	UUID     string
	Status   string
	State    string
	Progress float32
	Data     interface{}
	Errors   []string
	StartAt  time.Time
	EndAt    time.Time
}

type Group struct {
	tasksLock *sync.Mutex
	Tasks     map[string]*Task
}

func (g *Group) NewTask(data interface{}) (string, error) {

	uid := uuid.New().String()

	g.tasksLock.Lock()
	defer g.tasksLock.Unlock()

	if _, found := g.Tasks[uid]; found {
		return "", ErrTaskUUIDDuplicated
	}

	g.Tasks[uid] = &Task{
		UUID:     uid,
		Data:     data,
		Status:   TaskRuning,
		Progress: 0,
		Errors:   []string{},
		StartAt:  time.Now(),
	}
	return uid, nil
}

func (g *Group) SetTaskProgress(uid string, progress float32) error {
	g.tasksLock.Lock()
	defer g.tasksLock.Unlock()

	if t, found := g.Tasks[uid]; !found {
		return ErrTaskUUIDNotFound
	} else {
		t.Progress = progress
	}

	return nil
}

func (g *Group) AddErrorsToTask(uid string, errs []string) error {
	g.tasksLock.Lock()
	defer g.tasksLock.Unlock()

	if t, found := g.Tasks[uid]; !found {
		return ErrTaskUUIDNotFound
	} else {
		t.Errors = append(t.Errors, errs...)
	}

	return nil
}

func (g *Group) AddErrorToTask(uid string, errs string) error {
	g.tasksLock.Lock()
	defer g.tasksLock.Unlock()

	if t, found := g.Tasks[uid]; !found {
		return ErrTaskUUIDNotFound
	} else {
		t.Errors = append(t.Errors, errs)
	}

	return nil
}

func (g *Group) SetTaskStatus(uid string, status string) error {
	g.tasksLock.Lock()
	defer g.tasksLock.Unlock()

	if t, found := g.Tasks[uid]; !found {
		return ErrTaskUUIDNotFound
	} else {
		t.Status = status
	}

	return nil
}

func (g *Group) SetTaskState(uid string, state string) error {
	g.tasksLock.Lock()
	defer g.tasksLock.Unlock()

	if t, found := g.Tasks[uid]; !found {
		return ErrTaskUUIDNotFound
	} else {
		t.State = state
	}

	return nil
}

func (g *Group) FinishTask(uid string, status string) error {

	if status != TaskFail && status != TaskSuccess {
		return ErrTaskInvalidStatus
	}

	g.tasksLock.Lock()
	defer g.tasksLock.Unlock()

	if t, found := g.Tasks[uid]; !found {
		return ErrTaskUUIDNotFound
	} else {
		if t.Status != TaskRuning {
			return ErrTaskAlreadyFinish
		}
		t.Status = status
		t.EndAt = time.Now()
		g.Tasks[uid] = t
	}

	return nil
}

func (g *Group) GetTask(uid string) (Task, error) {
	g.tasksLock.Lock()
	defer g.tasksLock.Unlock()

	t, found := g.Tasks[uid]
	if !found {
		return Task{}, ErrTaskUUIDNotFound
	}

	return Task{
		Status:   t.Status,
		Progress: t.Progress,
		State:    t.State,
		Errors:   t.Errors,
		StartAt:  t.StartAt,
		EndAt:    t.EndAt,
		Data:     t.Data,
	}, nil
}

func (g *Group) GetTasks() []Task {
	g.tasksLock.Lock()
	defer g.tasksLock.Unlock()

	tasks := []Task{}

	for _, t := range g.Tasks {
		tasks = append(tasks, Task{
			Status:   t.Status,
			UUID:     t.UUID,
			Progress: t.Progress,
			State:    t.State,
			Errors:   t.Errors,
			StartAt:  t.StartAt,
			EndAt:    t.EndAt,
			Data:     t.Data,
		})
	}

	return tasks
}

func (g *Group) RemoveTask(uid string) error {
	g.tasksLock.Lock()
	defer g.tasksLock.Unlock()

	if _, found := g.Tasks[uid]; !found {
		return ErrTaskUUIDNotFound
	} else {
		delete(g.Tasks, uid)
	}

	return nil
}

func NewGroup() Group {
	return Group{
		Tasks:     map[string]*Task{},
		tasksLock: &sync.Mutex{},
	}
}
