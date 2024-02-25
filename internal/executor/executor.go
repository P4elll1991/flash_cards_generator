package executor

import (
	"flash_cards/internal"
	flashcards "flash_cards/internal/flash_cards"
	openai "flash_cards/internal/openAI"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type executor struct {
	service service
	queue   chan internal.GenerateParams
	locker  *locker
}

type locker struct {
	table map[string]*sync.Mutex
	mu    sync.Mutex
}

func newLocker() *locker {
	return &locker{
		table: make(map[string]*sync.Mutex),
		mu:    sync.Mutex{},
	}
}

func (l *locker) Lock(hash string) {
	l.mu.Lock()
	if mu := l.table[hash]; mu == nil {
		l.table[hash] = &sync.Mutex{}
	}
	l.mu.Unlock()
	l.table[hash].Lock()
}

func (l *locker) UnLock(hash string) {
	l.mu.Lock()
	if mu := l.table[hash]; mu == nil {
		l.table[hash] = &sync.Mutex{}
	}
	l.mu.Unlock()
	l.table[hash].Unlock()
}

type service interface {
	GetNewParams() ([]internal.GenerateParams, error)
	SaveCards(taskId int64, cards map[string]internal.FlashCard) error
	GetExceptions(internal.GenerateParams) (map[string]internal.FlashCard, error)
	GetUnfinishedTasks() ([]internal.GenerateParams, error)
}

var (
	CONCURRENCY int64 = 1000
	LIMIT       int64 = 10
)

func init() {
	var err error
	CONCURRENCY, err = strconv.ParseInt(os.Getenv("FLASH_CARDS_GENERATOR_CONCURRENCY"), 0, 64)
	if err != nil {
		CONCURRENCY = 1000
	}

	LIMIT, err = strconv.ParseInt(os.Getenv("FLASH_CARDS_GENERATOR_LIMIT_WORDS_ONE_STEP"), 0, 64)
	if err != nil {
		LIMIT = 10
	}
}

func New(service service) *executor {
	flashcards.Init(openai.Request)
	return &executor{
		service: service,
		queue:   make(chan internal.GenerateParams, CONCURRENCY),
		locker:  newLocker(),
	}
}

func (exe *executor) Run() {
	go exe.check()
	for i := 0; i < int(CONCURRENCY); i++ {
		go exe.pop()
	}
}

func (exe *executor) check() {
	for i := 0; i < 10; i++ {
		params, err := exe.service.GetUnfinishedTasks()
		if err != nil {
			fmt.Println("check: ", err)
			time.Sleep(10 * time.Second)
			continue
		}
		if len(params) > 0 {
			exe.push(params)
		}
		break
	}

	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		fmt.Println("run check")
		params, err := exe.service.GetNewParams()
		if err != nil {
			fmt.Println("check: ", err)
			continue
		}
		exe.push(params)
	}
}

func (exe *executor) push(params []internal.GenerateParams) {
	tasks := []internal.GenerateParams{}
	for i := range params {
		tasks = append(tasks, params[i].Divide(LIMIT)...)
	}

	for i := range tasks {
		exe.queue <- tasks[i]
	}
}

func (exe *executor) pop() {
	for task := range exe.queue {
		exe.generate(task)
	}
}

func (exe *executor) generate(task internal.GenerateParams) {
	exe.locker.Lock(task.Hash())
	defer exe.locker.UnLock(task.Hash())
	fmt.Println("generate start", task.TaskId, task.Words)
	exceptions, err := exe.service.GetExceptions(task)
	if err != nil {
		fmt.Println("generate eror ", task.TaskId, err.Error())
		go func() {
			time.Sleep(10 * time.Second)
			exe.queue <- task
		}()
		return
	}

	cards, err := flashcards.Generate(task, exceptions)
	if err != nil {
		fmt.Println("generate eror ", task.TaskId, err.Error())
		go func() {
			time.Sleep(10 * time.Second)
			exe.queue <- task
		}()
		return
	}
	repeat := task.Words - int64(len(cards))

	if len(cards) > 0 {
		err := exe.service.SaveCards(task.TaskId, cards)
		if err != nil {
			fmt.Println("generate eror ", task.TaskId, err.Error())
			go func() {
				time.Sleep(10 * time.Second)
				exe.queue <- task
			}()
			return
		}
	}

	if repeat > 0 {
		fmt.Println("generate repeat", task.TaskId, repeat)
		go func() {
			repeatTask := task
			repeatTask.Words = repeat
			exe.queue <- repeatTask
		}()
	}
}
