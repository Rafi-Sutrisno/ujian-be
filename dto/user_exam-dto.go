package dto

type (
	UserExamCreateRequest struct {
		UserID        string        `json:"user_id" binding:"required"`
		ExamID   	  string        `json:"exam_id" binding:"required"`
		Role 		  string        `json:"role" binding:"required"`
	}

	UserExamResponse struct {
		ID            string        `json:"id"`
		UserID        string        `json:"user_id" binding:"required"`
		ExamID   	  string        `json:"exam_id" binding:"required"`
		Role 		  string        `json:"role" binding:"required"`
	}
)