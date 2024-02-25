package excel

import (
	"flash_cards/internal"
	"fmt"

	"github.com/tealeg/xlsx"
)

func Unmarshal(xlFile *xlsx.File) ([]internal.FlashCard, error) {
	fmt.Println(xlFile)
	if len(xlFile.Sheets) > 0 {
		sheet := xlFile.Sheets[0]
		if len(sheet.Rows) > 1 {
			cards := make([]internal.FlashCard, 0, len(sheet.Rows)-1)
			for _, row := range sheet.Rows[1:] {
				cards = append(cards, internal.FlashCard{
					TaskId:             getInt64(0, row),
					NativeLang:         getString(1, row),
					LearningLang:       getString(2, row),
					Topic:              getString(3, row),
					Level:              getString(4, row),
					Word:               getString(5, row),
					Pronunciation:      getString(6, row),
					PhoneticRespelling: getString(7, row),
					Definition:         getString(8, row),
					Translation:        getString(9, row),
					Example:            getString(10, row),
					ExampleTranslation: getString(11, row),
				})
			}
			return cards, nil
		}
	}
	return nil, nil
}

func getString(index int, row *xlsx.Row) string {
	result := ""
	if len(row.Cells) > index {
		result = row.Cells[index].String()
	}
	return result
}

func getInt64(index int, row *xlsx.Row) int64 {
	var result int64
	if len(row.Cells) > index {
		result, _ = row.Cells[index].Int64()
	}
	return result
}
