package renamer

import (
	mod "filerenamer/model"
	"os"
	"path/filepath"
	"strings"
)

func FindFile(directory string) ([]string, error) {
	var files []string
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

func RenameFile(pattern, action, value string) ([]mod.RenameResult, error) {
	files, err := FindFile(pattern)
	if err != nil {
		return nil, err
	}

	rule := actions(action, value)

	result := make([]mod.RenameResult, 0, len(files))
	for _, file := range files {
		res, err := rename(file, rule)
		if err != nil {
			res = mod.RenameResult{
				OldName: file,
				NewName: "",
				Success: false,
				Error:   err,
			}
		}
		result = append(result, res)
	}
	return result, nil
}

func actions(action, value string) mod.Rule {
	var rule mod.Rule

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

	return rule
}

func rename(path string, rule mod.Rule) (mod.RenameResult, error) {
	dir := filepath.Dir(path)
	oldName := filepath.Base(path)
	newName := oldName

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

	newPath := filepath.Join(dir, newName)

	if path == newPath {
		return mod.RenameResult{
			OldName: path,
			NewName: newPath,
			Success: true,
			Error:   nil,
		}, nil
	}

	err := os.Rename(path, newPath)
	result := mod.RenameResult{
		OldName: path,
		NewName: newPath,
		Success: err == nil,
		Error:   err,
	}

	return result, err
}
