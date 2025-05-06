package dto

import (
	"errors"
	"mime/multipart"

	"mods/entity"
)

const (
	// Failed
	MESSAGE_FAILED_GET_DATA_FROM_BODY      = "failed get data from body"
	MESSAGE_FAILED_REGISTER_USER           = "failed create user"
	MESSAGE_FAILED_GET_LIST_USER           = "failed get list user"
	MESSAGE_FAILED_GET_USER_TOKEN          = "failed get user token"
	MESSAGE_FAILED_TOKEN_NOT_VALID         = "token not valid"
	MESSAGE_FAILED_TOKEN_NOT_FOUND         = "token not found"
	MESSAGE_FAILED_GET_USER                = "failed get user"
	MESSAGE_FAILED_LOGIN                   = "failed login"
	MESSAGE_FAILED_WRONG_EMAIL_OR_PASSWORD = "wrong email or password"
	MESSAGE_FAILED_UPDATE_USER             = "failed update user"
	MESSAGE_FAILED_DELETE_USER             = "failed delete user"
	MESSAGE_FAILED_PROSES_REQUEST          = "failed proses request"
	MESSAGE_FAILED_DENIED_ACCESS           = "denied access"
	MESSAGE_FAILED_VERIFY_EMAIL            = "failed verify email"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER           = "success create user"
	MESSAGE_SUCCESS_GET_LIST_USER           = "success get list user"
	MESSAGE_SUCCESS_GET_USER                = "success get user"
	MESSAGE_SUCCESS_LOGIN                   = "success login"
	MESSAGE_SUCCESS_UPDATE_USER             = "success update user"
	MESSAGE_SUCCESS_DELETE_USER             = "success delete user"
	MESSAGE_SEND_FORGOT_PASSWORD_SUCCESS    = "success send forgot password email"
	MESSAGE_SUCCESS_RESET_PASSWORD            = "success reset password"
)

var (
	ErrCreateUser             = errors.New("failed to create user")
	ErrGetAllUser             = errors.New("failed to get all user")
	ErrGetUserById            = errors.New("failed to get user by id")
	ErrGetUserByEmail         = errors.New("failed to get user by email")
	ErrEmailAlreadyExists     = errors.New("email already exist")
	ErrNoidAlreadyExists      = errors.New("NRP already exist")
	ErrUpdateUser             = errors.New("failed to update user")
	ErrUserNotAdmin           = errors.New("user not admin")
	ErrUserNotFound           = errors.New("user not found")
	ErrEmailNotFound          = errors.New("email not found")
	ErrUsernameNotFound       = errors.New("username not found")
	ErrDeleteUser             = errors.New("failed to delete user")
	ErrDeniedAccess           = errors.New("denied access")
	ErrPasswordNotMatch       = errors.New("password not match")
	ErrEmailOrPassword        = errors.New("wrong email or password")
	ErrAccountNotVerified     = errors.New("account not verified")
	ErrTokenInvalid           = errors.New("token invalid")
	ErrTokenExpired           = errors.New("token expired")
	ErrAccountAlreadyVerified = errors.New("account already verified")
	ErrHashPassword           = errors.New("err hash password")
)

type (
	UserCreateRequest struct {
		Username   string `json:"username" binding:"required"`
		Name       string `json:"name" binding:"required"`
		Password   string `json:"password" binding:"required"`
		Email      string `json:"email" binding:"required"`
		Noid       string `json:"noid" binding:"required"`
		RoleId     int `json:"role_id" binding:"required"`
	}

	UserResponse struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		RoleID    uint   `json:"role_id"`
		Noid      string `json:"noid"`
		CreatedAt string `json:"created_at"`
	}

	UserFileUploadRequest struct {
		File *multipart.FileHeader `form:"file" binding:"required"`
	}

	 UserFile struct {
		Username       string `yaml:"username"`
		Name       string `yaml:"name"`
		Email      string `yaml:"email"`
		Noid 	   string `yaml:"noid"`
		Password   string `yaml:"password"`  
	}
	
	 UserFileList struct {
		Users []UserFile `yaml:"users"`
	}

	FailedUserResponse struct {
		Noid   string `json:"noid"`
		Email  string `json:"email"`
		Reason string `json:"reason"`
	}
	

	UserPaginationResponse struct {
		Data []UserResponse `json:"data"`
		PaginationResponse
	}

	GetAllUserRepositoryResponse struct {
		Users []entity.User `json:"users"`
		PaginationResponse
	}

	UserUpdateRequest struct {
		Username       string `json:"username" form:"username"`
		Name       string `json:"name" form:"name"`
		Email      string `json:"email" form:"email"`
		Noid       string `json:"noid" form:"noid"`
	}
	UserUpdateEmailRequest struct {
		Email      string `json:"email" form:"email"`
	}

	UserUpdateResponse struct {
		ID         string `json:"id"`
		Username  string `json:"username"`
		Name       string `json:"name"`
		RoleID       uint `json:"role_id"`
		Email      string `json:"email"`
		Noid       string `json:"noid" `
	}

	SendResetPasswordRequest struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	ResetPasswordRequest struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	

	ResetPasswordResponse struct {
		Email string `json:"email"`
	}
	
	

	UserLoginRequest struct {
		Username     string `json:"username" form:"username" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserLoginResponse struct {
		Token string `json:"token"`
		RoleID  uint `json:"role_id"`
	}

	UpdateStatusIsVerifiedRequest struct {
		UserId     string `json:"user_id" form:"user_id" binding:"required"`
		IsVerified bool   `json:"is_verified" form:"is_verified"`
	}
)