package dto

type SubmissionRequest struct {
	ProblemID  string `json:"problem_id" binding:"required"`
	ExamID     string `json:"exam_id" binding:"required"`
	LanguageID int    `json:"language_id" binding:"required"`
	SourceCode string `json:"source_code" binding:"required"`
}

type Judge0SubmissionRequest struct {
	LanguageID     int    `json:"language_id"`
	SourceCode     string `json:"source_code"`
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output"`
}

// type Judge0SubmissionResult struct {
// 	Passed            bool
// 	CompilationError  bool
// 	StatusID          int
// 	StatusDescription string
// }

// dto/judge0.go

type Judge0BatchSubmissionRequest struct {
	Submissions []Judge0SubmissionRequest `json:"submissions"`
}

type Judge0BatchSubmissionResponse []struct {
	Token string `json:"token"`
}

type Judge0BatchResultResponse struct {
	Submissions []Judge0SubmissionResult `json:"submissions"`
}

type Judge0SubmissionResult struct {
	Token  string       `json:"token"`
	Status Judge0Status `json:"status"`
	Time   string       `json:"time"`
	Memory int          `json:"memory"`
}

type Judge0Status struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Judge0Request struct {
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	Stdin      string `json:"stdin,omitempty"`
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
