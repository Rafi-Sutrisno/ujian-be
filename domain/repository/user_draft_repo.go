package repository

import (
	"context"
	"mods/domain/entity"

	"gorm.io/gorm"
)

type UserDraftRepository interface {
	Create(ctx context.Context, tx *gorm.DB, userDraft entity.UserCodeDraft) (entity.UserCodeDraft, error)
	Update(ctx context.Context, tx *gorm.DB, draft entity.UserCodeDraft) (entity.UserCodeDraft, error)
	GetByIdentifiers(ctx context.Context, userID, examID, problemID, language string) (entity.UserCodeDraft, error)
}
