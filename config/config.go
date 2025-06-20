package config

import (
	"bufio"
	"filerenamer/renamer"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// читачь строки
var reader *bufio.Reader = bufio.NewReader(os.Stdin)

// приймає текстове введеня
func getInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// обробка інформації від користувача
func Inputs() {
	// приймає діректорію від користувача
	fmt.Println("Вкажіть директорію для перейменування (наприклад, ./test_files/):")
	dir := strings.TrimSpace(getInput("> "))

	// приймає шаблон від користувача
	fmt.Println("\nВкажіть шаблон файлів (наприклад, *.jpg):")
	templ := strings.TrimSpace(getInput("> "))

	// приймає дію від користувача
	fmt.Println("\nОберіть дію (prefix, suffix, replace, extension, lowercase, uppercase):")
	action := getInput("> ")
	action = strings.TrimSpace(action) // видалення зайвих пробілів

	// перевірка та визов функції ValueToAdd
	var value string
	if action != "lowercase" && action != "uppercase" {
		value = ValueToAdd(action)
	}
	if value != "" || action == "lowercase" || action == "uppercase" {
		pattern := filepath.Join(dir, templ)                 // обьєднує діректорію і розширення
		_, err := renamer.RenameFile(pattern, action, value) // виклик функції rename
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// повертає значення для додавання
func ValueToAdd(action string) string {
	var prompt string
	switch action {
	case "prefix":
		prompt = "\nВведіть значення для додавання як префікс:"
	case "suffix":
		prompt = "\nВведіть значення для додавання як суфікс:"
	case "replace":
		prompt = "\nВведіть значення для замінни(через пробіл, наприклад text new):"
	case "extension":
		prompt = "\nВведіть значення для змінненя розширення(без крапки):"
	default:
		fmt.Println("\nТакой дії не існує")
		return ""
	}
	fmt.Println(prompt)
	input := getInput("> ")
	return input
}
