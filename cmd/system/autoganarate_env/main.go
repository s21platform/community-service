package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/s21platform/community-service/internal/config"
)

func main() {
	cfg := &config.Config{} // Инициализируем вашу структуру конфигурации
	envLines := []string{}

	processStruct(cfg, &envLines)

	// Записываем данные в .env файл
	err := os.WriteFile(".env", []byte(strings.Join(envLines, "\n")), 0644)
	if err != nil {
		fmt.Printf("Failed to write .env file: %v\n", err)
		return
	}

	fmt.Println(".env file successfully generated")
}

func processStruct(s interface{}, envLines *[]string) {
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Если это структура, рекурсивно обрабатываем ее
		if fieldValue.Kind() == reflect.Struct {
			processStruct(fieldValue.Addr().Interface(), envLines)
		} else {
			// Извлекаем теги
			tag := field.Tag.Get("env")
			if tag != "" {
				*envLines = append(*envLines, fmt.Sprintf("%s=", tag))
			}
		}
	}
}
