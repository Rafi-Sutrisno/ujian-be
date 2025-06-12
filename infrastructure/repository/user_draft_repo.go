package repository

import (
	"context"
	"fmt"
	"mods/domain/entity"
	domain "mods/domain/repository"

	"gorm.io/gorm"
)

type userDraftRepository struct {
	db *gorm.DB
}

func NewUserDraftRepository(db *gorm.DB) domain.UserDraftRepository {
	return &userDraftRepository{
		db: db,
	}
}

func (udr *userDraftRepository) Create(ctx context.Context, tx *gorm.DB, userDraft entity.UserCodeDraft) (entity.UserCodeDraft, error) {
	if tx == nil {
		tx = udr.db
	}

	if err := tx.WithContext(ctx).Create(&userDraft).Error; err != nil {
		return entity.UserCodeDraft{}, err
	}

	return userDraft, nil
}

func (udr *userDraftRepository) Update(ctx context.Context, tx *gorm.DB, draft entity.UserCodeDraft) (entity.UserCodeDraft, error) {
	if tx == nil {
		tx = udr.db
	}

	if err := tx.WithContext(ctx).Save(&draft).Error; err != nil {
		return entity.UserCodeDraft{}, err
	}

	return draft, nil
}


func (udr *userDraftRepository) GetByIdentifiers(ctx context.Context, userID, examID, problemID, language string) (entity.UserCodeDraft, error) {
	var draft entity.UserCodeDraft

	err := udr.db.WithContext(ctx).
		Where("user_id = ? AND exam_id = ? AND problem_id = ? AND language = ?", userID, examID, problemID, language).
		First(&draft).Error // <-- this was missing

	if err != nil {
		return entity.UserCodeDraft{}, err
	}
	fmt.Println("ini draft dari db:", draft)

	return draft, nil
}

