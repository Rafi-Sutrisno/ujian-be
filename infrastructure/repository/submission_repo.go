package repository

import (
	"context"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"

	"gorm.io/gorm"
)



type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) domain.SubmissionRepository {
	return &submissionRepository{
		db: db,
	}
}

// submission_repository.go
func (r *submissionRepository) GetCorrectSubmissionStatsByExam(ctx context.Context, examID string) ([]dto.ExamUserCorrectDTO, error) {
	var results []dto.ExamUserCorrectDTO

	query := `
		SELECT 
			u.id AS user_id,
			u.name AS user_name,
			u.noid AS user_no_id,
			COUNT(DISTINCT CASE WHEN s.status_id = 2 THEN s.problem_id END) AS total_correct,
			(
				SELECT COUNT(*) 
				FROM exam_problems ep 
				WHERE ep.exam_id = ?
			) AS total_problem,
			es.status AS status,
			es.finished_at AS finished_at
		FROM exam_sesssions es
		JOIN users u ON es.user_id = u.id
		LEFT JOIN submissions s ON s.user_id = u.id AND s.exam_id = es.exam_id
		WHERE es.exam_id = ?
		GROUP BY u.id, u.name, u.noid, es.status, es.finished_at
		ORDER BY u.name;
	`

	if err := r.db.Raw(query, examID, examID).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}


func (r *submissionRepository) GetCorrectSubmissionStatsByExamandStudent(ctx context.Context, examID, userID string) (dto.ExamUserCorrectDTO, error) {
	var result dto.ExamUserCorrectDTO

	query := `
		SELECT 
			u.id AS user_id,
			u.name AS user_name,
			u.noid AS user_no_id,
			COUNT(DISTINCT CASE WHEN s.status_id = 2 THEN s.problem_id END) AS total_correct,
			(
				SELECT COUNT(*) 
				FROM exam_problems ep 
				WHERE ep.exam_id = ?
			) AS total_problem
		FROM exam_sesssions es
		JOIN users u ON es.user_id = u.id
		LEFT JOIN submissions s ON s.user_id = u.id AND s.exam_id = es.exam_id
		WHERE es.exam_id = ? AND u.id = ?
		GROUP BY u.id, u.name, u.noid
		ORDER BY u.name;
	`

	if err := r.db.Raw(query, examID, examID, userID).Scan(&result).Error; err != nil {
		return dto.ExamUserCorrectDTO{}, err
	}

	return result, nil
}

func (r *submissionRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submission entity.Submission
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&submission).Error; err != nil {
		return entity.Submission{}, err
	}

	return submission, nil
}

func (r *submissionRepository) Update(ctx context.Context, tx *gorm.DB, sub entity.Submission) (entity.Submission, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).Model(&sub).Updates(map[string]interface{}{
		"status_id":  sub.StatusId,
		"time":    sub.Time,
		"memory":  sub.Memory,
	}).Error

	return sub, err
}

func (r *submissionRepository) GetPendingSubmissions(ctx context.Context) ([]entity.Submission, error) {
	var subs []entity.Submission
	err := r.db.WithContext(ctx).Where("status_id = ?", 1).Find(&subs).Error
	return subs, err
}


func (r *submissionRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submissions []entity.Submission
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (r *submissionRepository) GetByExamIDandUserID(ctx context.Context, tx *gorm.DB, examID string, userID string) ([]entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	// fmt.Println("ini exam dan user id:", examID, userID)

	var submissions []entity.Submission
	if err := tx.WithContext(ctx).Where("exam_id = ? AND user_id = ?", examID, userID).Preload("Problem").Preload("Language").Preload("Status").Find(&submissions).Error; err != nil {
		return nil, err
	}
	// fmt.Println("ini hasil repo:", submissions)

	return submissions, nil
}

func (r *submissionRepository) GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submissions []entity.Submission
	if err := tx.WithContext(ctx).Where("exam_id = ?", examID).Preload("Problem").Preload("Language").Preload("User").Preload("Status").Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (r *submissionRepository) GetByProblemID(ctx context.Context, tx *gorm.DB, problemID string) ([]entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submissions []entity.Submission
	if err := tx.WithContext(ctx).Where("problem_id = ?", problemID).Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (r *submissionRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	var submissions []entity.Submission
	if err := tx.WithContext(ctx).Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (r *submissionRepository) Create(ctx context.Context, tx *gorm.DB, submission entity.Submission) (entity.Submission, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&submission).Error; err != nil {
		return entity.Submission{}, err
	}

	return submission, nil
}