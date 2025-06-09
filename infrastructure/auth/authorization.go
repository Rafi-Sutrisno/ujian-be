package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"mods/domain/entity"
	domainrepo "mods/domain/repository"
	"strings"
	"time"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domainrepo.AuthRepo {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) IsUserInExamClass(ctx context.Context, userId, examId string) (bool, error) {
	var exam entity.Exam
	if err := r.db.WithContext(ctx).Select("class_id").Where("id = ?", examId).First(&exam).Error; err != nil {
		return false, err
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.UserClass{}).
		Where("user_id = ? AND class_id = ?", userId, exam.ClassID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *authRepository) GetUserById(ctx context.Context, userId string) (entity.User, error) {

	var user entity.User
	if err := r.db.WithContext(ctx).Where("id = ?", userId).Take(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *authRepository) GetExamSessionBySessionId(ctx context.Context, sessionId string) (entity.ExamSesssion, error) {

	var examSession entity.ExamSesssion
	if err := r.db.WithContext(ctx).Where("session_id = ?", sessionId).Take(&examSession).Error; err != nil {
		return entity.ExamSesssion{}, err
	}

	return examSession, nil
}

func (r *authRepository) GetExamSessionByExamIdandUserId(ctx context.Context, examId, userId string) (entity.ExamSesssion, error) {

	var examSession entity.ExamSesssion
	if err := r.db.WithContext(ctx).Where("exam_id = ? AND user_id = ?", examId, userId).Take(&examSession).Error; err != nil {
		return entity.ExamSesssion{}, err
	}

	return examSession, nil
}

func (r *authRepository) IsExamActive(ctx context.Context, examId string) (bool, int64, error) {
	var exam entity.Exam
	if err := r.db.WithContext(ctx).Select("start_time", "end_time").Where("id = ?", examId).First(&exam).Error; err != nil {
		return false, 0, err
	}
	

	now := time.Now()
	isActive := now.After(exam.StartTime) && now.Before(exam.EndTime)

	var timeLeftSeconds int64
	if isActive {
		timeLeft := exam.EndTime.Sub(now)
		timeLeftSeconds = int64(timeLeft.Seconds())
	}

	return isActive, timeLeftSeconds, nil
}

func (r *authRepository) HasExamSession(ctx context.Context, examId, sessionId string) (bool, error) {
	var session entity.ExamSesssion
	err := r.db.WithContext(ctx).
		Where("session_id = ? AND exam_id = ?", sessionId, examId).
		First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (r *authRepository) CanStartExam(ctx context.Context, userAgent,requestHash, configKeyHash, fullURL, userId, examId string) (int64, error) {
	inClass, err := r.IsUserInExamClass(ctx, userId, examId)
	if err != nil {
		return 0, err
	}
	if !inClass {
		return 0, errors.New("user is not in the exam class")
	}

	active, timeleft, err := r.IsExamActive(ctx, examId)
	if err != nil {
		return 0, err
	}
	if !active {
		return 0, errors.New("exam is not active (not in time interval)")
	}

	if err := r.ValidateSEBRequest( ctx, userAgent,requestHash, configKeyHash, fullURL, examId); err != nil {
		return 0, fmt.Errorf("SEB validation failed: %w", err)
	}

	return timeleft, nil
}

func (r *authRepository) CanAccessExam(ctx context.Context, userAgent,requestHash, configKeyHash, fullURL, sessionID, userId, examId string) error {
	inClass, err := r.IsUserInExamClass(ctx, userId, examId)
	if err != nil {
		return err
	}
	if !inClass {
		return errors.New("user is not in the exam class")
	}

	user, err := r.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user.RoleID == 1 {
		return nil
	}
	
	active, _, err := r.IsExamActive(ctx, examId)
	if err != nil {
		return err
	}
	if !active {
		return errors.New("exam is not active")
	}

	// sessionID, err := ginCtx.Cookie("session_id")
	// if err != nil {
	// 	return errors.New("you don't have session id to access the exam")
	// }
	// fmt.Println("ini session id dari cookie:", sessionID)

	hasSession, err := r.HasExamSession(ctx, examId, sessionID)
	if err != nil {
		fmt.Println("ini error:", err)
		return err
	}
	if !hasSession {
		fmt.Println("masuk sini:")
		return errors.New("you don't have access to exam (Session Invalid)")
	}

	if err := r.ValidateSEBRequest(ctx, userAgent,requestHash, configKeyHash, fullURL, examId); err != nil {
		return fmt.Errorf("SEB validation failed: %w", err)
	}

	return nil
}

func (r *authRepository) CanSeeExamResult(ctx context.Context, userId, examId string) error {
	inClass, err := r.IsUserInExamClass(ctx, userId, examId)
	if err != nil {
		return err
	}
	if !inClass {
		return errors.New("user is not in the exam class")
	}

	user, err := r.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user.RoleID == 1 {
		return nil
	}
	
	session, err := r.GetExamSessionByExamIdandUserId(ctx, examId, userId)
	if err != nil {
		return err
	}

	if session.Status != 1{
		return errors.New("you have not finished this exam")
	}

	return nil
}

func (r *authRepository) CanAccessProblem(ctx context.Context,userAgent,requestHash, configKeyHash, fullURL, sessionID, userId, problemId string) error {
	user, err := r.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user.RoleID == 1 {
		return nil
	}

	// sessionID, err := ginCtx.Cookie("session_id")
	// if err != nil {
	// 	return errors.New("you don't have session id to access the problem")
	// }
	// fmt.Println("ini session id dari cookie:", sessionID)

    examSession, err := r.GetExamSessionBySessionId(ctx, sessionID)
	if err != nil {
		return err
	}

	inClass, err := r.IsUserInExamClass(ctx, userId, examSession.ExamID)
	if err != nil {
		return err
	}
	if !inClass {
		return errors.New("user is not in the exam class")
	}

	
	
	active, _, err := r.IsExamActive(ctx, examSession.ExamID)
	if err != nil {
		return err
	}
	if !active {
		return errors.New("exam is not active")
	}

	

	hasSession, err := r.HasExamSession(ctx, examSession.ExamID, sessionID)
	if err != nil {
		return err
	}
	if !hasSession {
		return errors.New("you don't have access to exam")
	}

	if err := r.ValidateSEBRequest( ctx,userAgent,requestHash, configKeyHash, fullURL, examSession.ExamID); err != nil {
		return fmt.Errorf("SEB validation failed: %w", err)
	}

	return nil
}

func (r *authRepository) ValidateSEBRequest(ctx context.Context, userAgent,requestHash, configKeyHash, fullURL, examId string) error {
	var exam entity.Exam
	if err := r.db.WithContext(ctx).Where("id = ?", examId).Take(&exam).Error; err != nil {
		return err
	}
	
	if !exam.IsSEBRestricted {
		return nil // No SEB restrictions
	}

	// userAgent := GinCtx.Request.UserAgent()
	// requestHash := GinCtx.GetHeader("X-SafeExamBrowser-RequestHash")
	// configKeyHash := GinCtx.GetHeader("X-Safeexambrowser-Configkeyhash")

	// scheme := "http"
	// if GinCtx.Request.TLS != nil {
	// 	scheme = "https"
	// }
	// fullURL := fmt.Sprintf("%s://%s%s", scheme, GinCtx.Request.Host, GinCtx.Request.RequestURI)

	// If not in header, fall back to request body (already parsed in controller)
	// var body struct {
	// 	BrowserExamKey string `json:"browser_exam_key"`
	// 	ConfigKey      string `json:"config_key"`
	// }
	// if err := GinCtx.ShouldBind(&body); err == nil {
	// 	if requestHash == "" {
	// 		requestHash = body.BrowserExamKey
	// 	}
	// 	if configKeyHash == "" {
	// 		configKeyHash = body.ConfigKey
	// 	}
	// }

	// Validate based on provided keys
	if exam.SEBBrowserKey != "" {
		if !r.validateSEBHash(fullURL, exam.SEBBrowserKey, requestHash) {
			return errors.New("unauthorized SEB request: browser exam key hash mismatch")
		}
	}

	if exam.SEBConfigKey != "" {
		if !r.validateSEBHash(fullURL, exam.SEBConfigKey, configKeyHash) {
			return errors.New("unauthorized SEB request: config key hash mismatch")
		}
	}

	// Fallback to user agent check if no keys provided
	if exam.SEBBrowserKey == "" && exam.SEBConfigKey == "" {
		if !strings.Contains(userAgent, "SEB") {
			return errors.New("unauthorized SEB request: user agent mismatch")
		}
	}

	return nil
}

func (r *authRepository) validateSEBHash(url string, key string, recvHash string) bool {

    hasher := sha256.New()

	hasher.Write([]byte(url))
	hasher.Write([]byte(key))
    
    finalHash := hasher.Sum(nil)
    hashHex := hex.EncodeToString(finalHash)

    fmt.Println("BEK/ConfigKey: Expected Hash:", hashHex)
    fmt.Println("BEK/ConfigKey: Received Hash:", recvHash)

    return hashHex == recvHash
}
