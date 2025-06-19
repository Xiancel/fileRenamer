package renamer

import (
	mod "filerenamer/model"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// пошук файлів в деректорії
func FindFile(directory string) ([]string, error) {
	// змінна для зберегання файлів
	files, err := filepath.Glob(directory)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// змінна назви файлу
func RenameFile(pattern, action, value string) ([]mod.RenameResult, error) {
	// пошук файлів
	files, err := FindFile(pattern)
	if err != nil {
		return nil, err
	}

	// перевірка діректорії на наявність файлів
	if len(files) == 0 {
		fmt.Println("\n❌ Не знайдено файлів у деректорії")
	} else {
		fmt.Printf("\n🔍 Знайдено файлів: %d\n\n", len(files))
	}
	// змінна для хранніня правил на основі вибору користувачем
	rule := actions(action, value)
	// змінна результату зміни
	result := make([]mod.RenameResult, 0, len(files))
	count := 0

	// зміна файлів
	for _, file := range files {
		// визов функції змінни
		res, err := rename(file, rule)
		// перевірка на помилку
		if err != nil {
			res = mod.RenameResult{
				OldName: file,
				NewName: "",
				Success: false,
				Error:   err,
			}
		}
		// додавання файлу до змінни result
		result = append(result, res)
		oldName := filepath.Base(res.OldName)
		newName := filepath.Base(res.NewName)

		fmt.Printf("%s - %s\n", oldName, newName)
		if res.Success {
			count++
		}
	}
	if count > 0 {
		fmt.Printf("\n✅ Успішно перейменовано: %d файли(ів)\n", count)
	}

	return result, nil
}

// функція вибору дії
func actions(action, value string) mod.Rule {
	// змінна для храніння правил дії
	var rule mod.Rule

	// перевірка вибору дії і передавання її до змінни
	switch action {
	case "prefix":
		rule.Prefix = value
	case "suffix":
		rule.Suffix = value
	case "replace":
		parts := strings.SplitN(value, " ", 2)
		if len(parts) == 2 {
			rule.Replace = mod.ReplaceRule{
				From: parts[0],
				To:   parts[1],
			}
		}
	case "extension":
		rule.Extension = value
	case "lowercase":
		rule.Lowercase = true
	case "uppercase":
		rule.Uppercase = true
	}

	// повертає вибір діє
	return rule
}

// функція перейменування файлів
func rename(path string, rule mod.Rule) (mod.RenameResult, error) {
	// змінни для хранніння
	dir := filepath.Dir(path)      // директорії
	oldName := filepath.Base(path) // старого імені файлу
	newName := oldName             // нового імені файлу

	// перевірка і змінна імені файлів по правилу дії вибране користувачем

	if rule.Prefix != "" {
		newName = rule.Prefix + newName
	}

	if rule.Suffix != "" {
		ext := filepath.Ext(newName)
		name := strings.TrimSuffix(newName, ext)
		newName = name + rule.Suffix + ext
	}

	if rule.Replace.From != "" {
		newName = strings.ReplaceAll(newName, rule.Replace.From, rule.Replace.To)
	}

	if rule.Extension != "" {
		ext := filepath.Ext(newName)
		name := strings.TrimSuffix(newName, ext)
		newName = name + "." + rule.Extension
	}

	if rule.Lowercase {
		ext := filepath.Ext(newName)
		name := strings.TrimSuffix(newName, ext)
		name = strings.ToLower(name)
		newName = name + ext
	}

	if rule.Uppercase {
		ext := filepath.Ext(newName)
		name := strings.TrimSuffix(newName, ext)
		name = strings.ToUpper(name)
		newName = name + ext
	}

	newName = strings.TrimSpace(newName)
	newName = filepath.Base(newName)
	// змінна для храніння нового путі файлу
	newPath := filepath.Join(dir, newName)
	// перевірка на совпадіння старого путі і нового
	if path == newPath {
		return mod.RenameResult{
			OldName: path,
			NewName: newPath,
			Success: true,
			Error:   nil,
		}, nil
	}

	// перейменівання файлу
	err := os.Rename(path, newPath)
	// змінна для ініціальзації результату
	result := mod.RenameResult{
		OldName: path,
		NewName: newPath,
		Success: err == nil,
		Error:   err,
	}

	// повернення зміненог файлу
	return result, err
}
