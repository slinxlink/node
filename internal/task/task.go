package task

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID   string
	ch   chan string
	done bool
	mu   sync.Mutex
}

var (
	tasks = map[string]*Task{}
	mu    sync.Mutex
)

func New(id string) *Task {
	t := &Task{
		ID: id,
		ch: make(chan string, 100),
	}
	mu.Lock()
	tasks[id] = t
	mu.Unlock()
	return t
}

func Get(id string) *Task {
	mu.Lock()
	defer mu.Unlock()
	return tasks[id]
}

func (t *Task) Log(level, msg string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		return
	}
	line := time.Now().Format("2006-01-02 15:04:05") + " " + level + " " + msg
	t.ch <- line
}

func (t *Task) Done() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.done = true
	close(t.ch)
	go func() {
		time.Sleep(30 * time.Second)
		mu.Lock()
		delete(tasks, t.ID)
		mu.Unlock()
	}()
}

func (t *Task) Chan() <-chan string {
	return t.ch
}

// ── Logger ────────────────────────────────────────────────────────────────────

type Logger struct {
	t *Task
}

func NewLogger(t *Task) *Logger {
	return &Logger{t: t}
}

func (l *Logger) Fatal(args ...interface{}) {
	l.t.Log("ERROR", fmt.Sprint(args...))
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.t.Log("ERROR", fmt.Sprintln(args...))
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.t.Log("ERROR", fmt.Sprintf(format, args...))
}

func (l *Logger) Print(args ...interface{}) {
	l.t.Log("INFO", fmt.Sprint(args...))
}

func (l *Logger) Println(args ...interface{}) {
	l.t.Log("INFO", fmt.Sprintln(args...))
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.t.Log("INFO", fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.t.Log("INFO", fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.t.Log("WARN", fmt.Sprintf(format, args...))
}
