package dto

type UserCodeDraftRequest struct {
	UserID    string `json:"user_id"`
	ExamID    string `json:"exam_id"`
	ProblemID string `json:"problem_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
}

type UserCodeDraftResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ExamID    string `json:"exam_id"`
	ProblemID string `json:"problem_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
}
