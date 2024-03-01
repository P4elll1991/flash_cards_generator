package main

import (
	"flash_cards/internal/executor"
	"flash_cards/internal/handler"
	"flash_cards/internal/repository"
	"flash_cards/internal/service"
	"fmt"
	"log"
)

func main() {
	fmt.Println("START")

	repo, err := repository.New()
	if err != nil {
		log.Fatal(err)
	}

	service := service.New(service.Config{
		TaskRepo:       repo.Tasks,
		TaskParamsRepo: repo.TaskParams,
		FlashCardRepo:  repo.FlashCards,
	})

	executor.New(service).Run()
	handler.New(service).Run()
}
