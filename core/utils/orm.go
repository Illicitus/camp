package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"unicode"
)

func First(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return GormErr.NotFound
	}
	return err
}

func splitByUppercase(value string) []string {
	var result []string

	l := 0
	for s := value; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(s)
		}
		result = append(result, s[:l])
	}
	return result
}

func excludeModelWord(values []string) []string {
	for i, s := range values {
		if s == "Model" || s == "model" {
			values[i] = ""
		}
	}
	return values
}

func joinStringsArrayWithUnderscoreSeparator(values []string) string {
	var result string
	for i, v := range values {
		if i == 0 && v != "" {
			result += v
		} else {
			if v != "" {
				result += "_" + v
			}
		}
	}
	return result
}

func NormalizeModelName(subAppName string, modelName string) string {
	splittedWords := splitByUppercase(modelName)
	splittedWords = excludeModelWord(splittedWords)

	finalName := joinStringsArrayWithUnderscoreSeparator(splittedWords)
	finalName = strings.ToLower(finalName)
	finalName = strings.TrimSpace(finalName)
	finalName = fmt.Sprintf("%ss", finalName)

	return fmt.Sprintf("%s_%s", subAppName, finalName)

}
