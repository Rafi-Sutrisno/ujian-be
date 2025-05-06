package dto

type Judge0Request struct {
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	Stdin      string `json:"stdin,omitempty"`
}

type Judge0Status struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Judge0Response struct {
	Token         string       `json:"token"`
	Stdout        string       `json:"stdout"`
	Stderr        string       `json:"stderr"`
	CompileOutput string       `json:"compile_output"`
	Status        Judge0Status `json:"status"`
	Time          string       `json:"time"`
	Memory        int          `json:"memory"`
}
