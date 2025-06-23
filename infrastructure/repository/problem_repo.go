package repository

import (
	"context"
	"errors"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type problemRepository struct {
	db *gorm.DB
}

func NewProblemRepository(db *gorm.DB) domain.ProblemRepository {
	return &problemRepository{
		db: db,
	}
}

func (r *problemRepository) GetByTitle(ctx context.Context, tx *gorm.DB, title string) (entity.Problem, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	var problem entity.Problem
	if err := db.WithContext(ctx).Where("title = ?", title).First(&problem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Problem{}, err
		}
		return entity.Problem{}, err
	}

	return problem, nil
}

func (pr *problemRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Problem, error) {
	if tx == nil {
		tx = pr.db
	}

	var problem entity.Problem
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&problem).Error; err != nil {
		return entity.Problem{}, err
	}

	return problem, nil
}

func (pr *problemRepository) GetByExamID(ctx context.Context, tx *gorm.DB, examID string) ([]entity.ExamProblem, error) {
	if tx == nil {
		tx = pr.db
	}

	var examProblems []entity.ExamProblem
	err := tx.WithContext(ctx).
		Preload("Problem").
		Where("exam_id = ?", examID).
		Find(&examProblems).Error

	if err != nil {
		return nil, err
	}

	return examProblems, nil
}

func (pr *problemRepository) GetByExamIDStudent(ctx context.Context, tx *gorm.DB, examID string) ([]dto.ProblemWithStatusResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var results []dto.ProblemWithStatusResponse

	rawSQL := `
		SELECT 
			ep.exam_id,
			ep.problem_id,
			ep.created_at,
			p.id, p.title, p.description, p.constraints, p.sample_input, p.sample_output, 
			p.cpu_time_limit, p.memory_limit,
			CASE 
				WHEN EXISTS (
					SELECT 1 FROM submissions s 
					WHERE s.problem_id = ep.problem_id AND s.exam_id = ep.exam_id AND s.status_id = 2
				) THEN 'accepted'
				WHEN EXISTS (
					SELECT 1 FROM submissions s 
					WHERE s.problem_id = ep.problem_id AND s.exam_id = ep.exam_id AND s.status_id NOT IN (1,2)
				) THEN 'wrong answer'
				ELSE ''
			END as status
		FROM exam_problems ep
		JOIN problems p ON p.id = ep.problem_id
		WHERE ep.exam_id = ?
		ORDER BY ep.created_at
	`

	err := tx.Raw(rawSQL, examID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (pr *problemRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Problem, error) {
	if tx == nil {
		tx = pr.db
	}

	var problems []entity.Problem
	if err := tx.WithContext(ctx).Find(&problems).Error; err != nil {
		return nil, err
	}

	return problems, nil
}

func (pr *problemRepository) Create(ctx context.Context, tx *gorm.DB, problem entity.Problem) (entity.Problem, error) {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).Create(&problem).Error; err != nil {
		return entity.Problem{}, err
	}

	return problem, nil
}

func (pr *problemRepository) Update(ctx context.Context, tx *gorm.DB, problem entity.Problem) (entity.Problem, error) {
	if tx == nil {
		tx = pr.db
	}

	

	if err := tx.WithContext(ctx).Model(&entity.Problem{}).Where("id = ?", problem.ID).Updates(map[string]interface{}{
		"title":          problem.Title,
		"description":    problem.Description,
		"constraints":    problem.Constraints,
		"sample_input":   problem.SampleInput,
		"sample_output":  problem.SampleOutput,
		"cpu_time_limit": problem.CpuTimeLimit,
		"memory_limit":   problem.MemoryLimit,
		}).Error; err != nil {
		return entity.Problem{}, err
	}

	return problem, nil
}

func (pr *problemRepository) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	problem := entity.Problem{ID: uuid.MustParse(id)}
	if err := tx.WithContext(ctx).Select("TestCase").Delete(problem).Error; err != nil {
		return err
	}

	return nil
}

func (pr *problemRepository) IsUserInProblemClass(ctx context.Context, tx *gorm.DB, userID, problemID string) (bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var count int64
	err := tx.WithContext(ctx).
		Table("user_classes").
		Joins("JOIN exams ON user_classes.class_id = exams.class_id").
		Joins("JOIN problems ON exams.id = problems.exam_id").
		Where("user_classes.user_id = ? AND problems.id = ?", userID, problemID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}


func (pr *problemRepository) IsUserInExamClass(ctx context.Context, tx *gorm.DB, userId, examId string) (bool, error) {
	if tx == nil {
		tx = pr.db
	}

	// First, get the class_id from the exam
	var exam entity.Exam
	if err := tx.WithContext(ctx).Select("class_id").Where("id = ?", examId).First(&exam).Error; err != nil {
		return false, err
	}

	// Now check if user exists in that class
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.UserClass{}).
		Where("user_id = ? AND class_id = ?", userId, exam.ClassID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

