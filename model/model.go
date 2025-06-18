package model

// структура правила заміни підрядка
type ReplaceRule struct {
	From string
	To   string
}

// основна структура правил
type Rule struct {
	Prefix    string
	Suffix    string
	Replace   ReplaceRule
	Extension string
	Lowercase bool
	Uppercase bool
}

// структура результату заміни
type RenameResult struct {
	OldName string
	NewName string
	Success bool
	Error   error
}
