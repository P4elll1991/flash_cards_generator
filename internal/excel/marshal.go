package excel

import (
	"flash_cards/internal"

	"github.com/tealeg/xlsx"
)

func Marshal(cards []internal.FlashCard) (*xlsx.File, error) {

	// Создаем новый файл Excel
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}

	// Создаем стиль для заголовков столбцов
	headerStyle := xlsx.NewStyle()
	headerStyle.Font.Bold = true

	// Записываем заголовки столбцов с применением стиля
	header := sheet.AddRow()
	header.AddCell().SetString("Task Id")
	header.AddCell().SetString("Native Lang")
	header.AddCell().SetString("Learning Lang")
	header.AddCell().SetString("Topic")
	header.AddCell().SetString("Level")
	header.AddCell().SetString("Word")
	header.AddCell().SetString("Pronunciation")
	header.AddCell().SetString("Phonetic Respelling")
	header.AddCell().SetString("Definition")
	header.AddCell().SetString("Translation")
	header.AddCell().SetString("Example")
	header.AddCell().SetString("Example Translation")
	for _, cell := range header.Cells {
		cell.SetStyle(headerStyle)
	}

	// Создаем стиль для данных
	dataStyle := xlsx.NewStyle()
	dataStyle.Alignment.Horizontal = "center"

	// Записываем данные из структур в Excel файл с применением стиля
	for _, card := range cards {
		row := sheet.AddRow()
		row.AddCell().SetInt64(card.TaskId)
		row.AddCell().SetString(card.NativeLang)
		row.AddCell().SetString(card.LearningLang)
		row.AddCell().SetString(card.Topic)
		row.AddCell().SetString(card.Level)
		row.AddCell().SetString(card.Word)
		row.AddCell().SetString(card.Pronunciation)
		row.AddCell().SetString(card.PhoneticRespelling)
		row.AddCell().SetString(card.Definition)
		row.AddCell().SetString(card.Translation)
		row.AddCell().SetString(card.Example)
		row.AddCell().SetString(card.ExampleTranslation)

		for _, cell := range row.Cells {
			cell.SetStyle(dataStyle)
		}
	}

	return file, nil
}
