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
	fmt.Println("Вкажіть шаблон файлів (наприклад, *.jpg):")
	templ := strings.TrimSpace(getInput("> "))

	// приймає дію від користувача
	fmt.Println("Оберіть дію (prefix, suffix, replace, extension, lowercase, uppercase):")
	action := getInput("> ")
	action = strings.TrimSpace(action) // видалення зайвих пробілів

	// перевірка та визов функції ValueToAdd
	var value string
	if action != "lowercase" && action != "uppercase" {
		value = ValueToAdd(action)
	}
	if value != "" || action == "lowercase" || action == "uppercase" {
		pattern := filepath.Join(dir, templ)       // обьєднує діректорію і розширення
		renamer.RenameFile(pattern, action, value) // виклик функції rename
	}
}

// повертає значення для додавання
func ValueToAdd(action string) string {
	var input string
	switch action {
	case "prefix":
		fmt.Println("Введіть значення для додавання як префікс:")
		input = getInput("> ")
	case "suffix":
		fmt.Println("Введіть значення для додавання як суфікс:")
		input = getInput("> ")
	case "replace":
		fmt.Println("Введіть значення для замінни(через пробіл, наприклад text new):")
		input = getInput("> ")
	case "extension":
		fmt.Println("Введіть значення для змінненя розширення(без крапки):")
		input = getInput("> ")
	default:
		fmt.Println("Такой дії не існує")
		return ""
	}
	return input
}
