package service

import (
	"flash_cards/internal"
	"flash_cards/internal/excel"
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"github.com/tealeg/xlsx"
)

type service struct {
	taskRepo       taskRepo
	taskParamsRepo taskParamsRepo
	flashCardRepo  flashCardRepo
}

type taskRepo interface {
	Create() (internal.Task, error)
	Get(id int64) (internal.Task, error)
	SetStatus(id int64, status string) error
}

type taskParamsRepo interface {
	Set(id int64, params []internal.GenerateParams) error
	Get(id int64) ([]internal.GenerateParams, error)
	GetNewParams() ([]internal.GenerateParams, error)
	GetUnfinishedTasks() ([]internal.GenerateParams, error)
}

type flashCardRepo interface {
	CountWords(taskId int64) (int64, error)
	Search(params map[string]interface{}) ([]internal.FlashCard, error)
	Update(cards []internal.FlashCard) error
	SaveCards(taskId int64, cards map[string]internal.FlashCard) error
	GetExceptions(internal.GenerateParams) (map[string]internal.FlashCard, error)
}

type Config struct {
	TaskRepo       taskRepo
	TaskParamsRepo taskParamsRepo
	FlashCardRepo  flashCardRepo
}

func New(config Config) *service {
	return &service{
		taskRepo:       config.TaskRepo,
		taskParamsRepo: config.TaskParamsRepo,
		flashCardRepo:  config.FlashCardRepo,
	}
}

func (servcie *service) SetTask(task internal.Task) (internal.Task, error) {
	newTask, err := servcie.taskRepo.Create()
	if err != nil {
		return task, err
	}

	if err := servcie.taskParamsRepo.Set(newTask.Id, task.Params); err != nil {
		return task, err
	}

	return servcie.GetTask(newTask.Id)
}

func (service *service) GetTask(id int64) (internal.Task, error) {
	task, err := service.taskRepo.Get(id)
	if err != nil {
		return task, err
	}
	params, err := service.taskParamsRepo.Get(id)
	if err != nil {
		return task, err
	}

	task.Params = params

	var totalWords int64
	for i := range task.Params {
		totalWords += task.Params[i].Words
	}

	var progress int64 = 100
	if totalWords > 0 {
		readyCards, err := service.flashCardRepo.CountWords(id)
		if err == nil {
			progress = int64((float64(readyCards) / float64(totalWords)) * 100)
		}
	}

	task.Progress = progress
	if task.Progress >= 100 && task.Status != "Completed" {
		task.Status = "Completed"
		if err := service.taskRepo.SetStatus(id, "Completed"); err != nil {
			fmt.Println("SetStatus", err)
		}
	}
	task.RouteResult = fmt.Sprintf("/generator/cards?task_id=%d", id)
	fmt.Println(task.RouteResult)

	return task, nil
}

func (service *service) SearchCards(params map[string]interface{}) ([]internal.FlashCard, error) {
	return service.flashCardRepo.Search(params)
}

func (service *service) GetNewParams() ([]internal.GenerateParams, error) {
	return service.taskParamsRepo.GetNewParams()
}

func (service *service) SaveCards(taskId int64, cards map[string]internal.FlashCard) error {
	return service.flashCardRepo.SaveCards(taskId, cards)
}

func (service *service) GetExceptions(params internal.GenerateParams) (map[string]internal.FlashCard, error) {
	return service.flashCardRepo.GetExceptions(params)
}

func (service *service) GetUnfinishedTasks() ([]internal.GenerateParams, error) {
	params, err := service.taskParamsRepo.GetUnfinishedTasks()
	if err != nil {
		return nil, err
	}
	unfinishedTasks := []internal.GenerateParams{}
	for i := range params {
		searchParam := map[string]interface{}{
			"task_id":       params[i].TaskId,
			"native_lang":   params[i].NativeLang,
			"learning_lang": params[i].LearningLang,
			"topic":         params[i].Topic,
			"level":         params[i].Level,
		}
		cards, err := service.flashCardRepo.Search(searchParam)
		if err != nil {
			return nil, err
		}
		if len(cards) < int(params[i].Words) {
			task := params[i]
			task.Words -= int64(len(cards))
			unfinishedTasks = append(unfinishedTasks, task)
		}
	}

	return unfinishedTasks, nil
}

func (service *service) UpdateCards(cards []internal.FlashCard) ([]internal.FlashCard, error) {
	if err := service.flashCardRepo.Update(cards); err != nil {
		return nil, err
	}

	return cards, nil
}

func (service *service) ExelExport(params map[string]interface{}) (*xlsx.File, error) {
	cards, err := service.SearchCards(params)
	if err != nil {
		return nil, err
	}

	file, err := excel.Marshal(cards)
	if err != nil {
		return nil, err
	}

	return file, err
}

func (service *service) ExelImport(file multipart.File) ([]internal.FlashCard, error) {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(fileBytes))
	xlFile, err := xlsx.OpenBinary(fileBytes)
	if err != nil {
		return nil, err
	}

	cards, err := excel.Unmarshal(xlFile)
	if err != nil {
		return nil, err
	}
	if len(cards) > 0 {
		if err := service.flashCardRepo.Update(cards); err != nil {
			return nil, err
		}
	}

	return cards, nil
}
