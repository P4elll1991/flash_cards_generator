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
	// flashcards.Init(openai.Request)

	// inputFile := flag.String("input", "", "Path to input file")
	// outputFile := flag.String("output", "", "Path to output file")
	// flag.Parse()

	// if *inputFile == "" || *outputFile == "" {
	// 	fmt.Println("Usage: go run main.go -input <input_file_path> -output <output_file_path>")
	// 	return
	// }

	// // Чтение данных из файла
	// inputBytes, err := ioutil.ReadFile(*inputFile)
	// if err != nil {
	// 	fmt.Println("Error reading input file:", err)
	// 	return
	// }

	// var params []flashcards.GenerateParams
	// err = json.Unmarshal(inputBytes, &params)
	// if err != nil {
	// 	fmt.Println("Error parsing input JSON:", err)
	// 	return
	// }

	// totalCards := []flashcards.FlashCard{}
	// mu := sync.Mutex{}
	// wg := sync.WaitGroup{}
	// for index := range params {
	// 	wg.Add(1)
	// 	go func(param flashcards.GenerateParams) {
	// 		defer wg.Done()
	// 		cards, err := flashcards.Generate(param)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		mu.Lock()
	// 		totalCards = append(totalCards, cards...)
	// 		mu.Unlock()
	// 	}(params[index])
	// }

	// wg.Wait()
	// // Преобразование структуры Flashcard обратно в JSON
	// outputBytes, err := json.MarshalIndent(totalCards, "", "\t")
	// if err != nil {
	// 	fmt.Println("Error marshaling JSON:", err)
	// 	return
	// }

	// // Запись данных в выходной файл
	// err = ioutil.WriteFile(*outputFile, outputBytes, 0644)
	// if err != nil {
	// 	fmt.Println("Error writing output file:", err)
	// 	return
	// }

	// fmt.Println("Data written to", *outputFile)

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
