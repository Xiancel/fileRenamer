package model

type ReplaceRule struct {
	From string
	To   string
}

type Rule struct {
	Prefix    string
	Suffix    string
	Replace   ReplaceRule
	Extension string
	Lowercase bool
	Uppercase bool
}

type RenameResult struct {
	OldName string
	NewName string
	Success bool
	Error   error
}
