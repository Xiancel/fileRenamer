package renamer

import (
	mod "filerenamer/model"
	"os"
	"path/filepath"
	"strings"
)

// пошук файлів в деректорії
func FindFile(directory string) ([]string, error) {
	// змінна для зберегання файлів
	var files []string

	// пошук і додавання файлів
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// змінна назви файлу
func RenameFile(pattern, action, value string) ([]mod.RenameResult, error) {
	// пошук файлів
	files, err := FindFile(pattern)
	if err != nil {
		return nil, err
	}

	// змінна для хранніня правил на основі вибору користувачем
	rule := actions(action, value)

	// змінна результату зміни
	result := make([]mod.RenameResult, 0, len(files))
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
		newName = strings.ToLower(newName)
	}

	if rule.Uppercase {
		newName = strings.ToUpper(newName)
	}

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
