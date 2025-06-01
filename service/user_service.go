package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"mods/domain/entity"
	domainrepo "mods/domain/repository"
	"mods/interface/dto"
	"os"
	"strings"
	"time"

	"mods/utils"

	"gopkg.in/yaml.v3"
)

type userService struct {
	userDomain domainrepo.UserRepository
	jwtService  JWTService
	
}

type UserService interface {
	// functional
	Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]dto.UserResponse, error)
	GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error)
	Verify(ctx context.Context, loginDTO dto.UserLoginRequest) (dto.UserLoginResponse, error)
	Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error)
	UpdateMe(ctx context.Context, req dto.UserUpdateEmailRequest, userId string) (dto.UserUpdateResponse, error)
	Delete(ctx context.Context, userId string) error
	GetUserById(ctx context.Context, userId string) (dto.UserResponse, error)
	SendForgotPasswordEmail(ctx context.Context, req dto.SendResetPasswordRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (dto.ResetPasswordResponse, error)

	RegisterUsersFromYAML(ctx context.Context, fileHeader *multipart.FileHeader) (map[string]interface{}, error)
	RegisterUsersFromCSV(ctx context.Context, fileHeader *multipart.FileHeader) (map[string]interface{}, error)
}

func NewUserService(ur domainrepo.UserRepository, jwtService JWTService) UserService {
	return &userService{
		userDomain: ur,
		jwtService:  jwtService,
	}
}

const (
	LOCAL_URL          = "http://localhost:3000"
	RESET_PASSWORD_ROUTE = "reset_password"
)

func (us *userService) Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error) {
	_, flag, _ := us.userDomain.CheckNoid(ctx, nil, req.Noid)
	if flag {
		return dto.UserResponse{}, dto.ErrNoidAlreadyExists
	}


	user := entity.User{
		Username:   req.Username,
		Name:       req.Name,
		Noid:       req.Noid,
		RoleID:     uint(req.RoleId),
		Email:      req.Email,
		Password:   req.Password,
	}

	userReg, err := us.userDomain.RegisterUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	return dto.UserResponse{
		ID:         userReg.ID.String(),
		Username:   userReg.Username,
		Name:       userReg.Name,
		Noid: 		userReg.Noid,
		RoleID:       userReg.RoleID,
		Email:      userReg.Email,
		CreatedAt:  userReg.CreatedAt.String(),
	}, nil
}

func (us *userService) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	Users, err := us.userDomain.GetAllUsers(ctx, nil)
	if err != nil {
		return nil, err
	}

	var responses []dto.UserResponse
	for _, User := range Users {
		responses = append(responses, dto.UserResponse{
			ID:            	User.ID.String(),
			Username:       User.Username,
			Name: 			User.Name,
			Email: 			User.Email,
			RoleID: 		User.RoleID,
			Noid: 			User.Noid,
			CreatedAt:      User.CreatedAt.String(),
		})
	}

	return responses, nil
}

func (us *userService) GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := us.userDomain.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	var datas []dto.UserResponse
	for _, user := range dataWithPaginate.Users {
		data := dto.UserResponse{
			ID:         user.ID.String(),
			Username:   user.Username,
			Name:       user.Name,
			Email:      user.Email,
			RoleID:       user.RoleID,
			Noid:       user.Noid,
		}

		datas = append(datas, data)
	}

	return dto.UserPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}


func (us *userService) Verify(ctx context.Context, loginDTO dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	check, flag, err := us.userDomain.CheckUsername(ctx,nil,  loginDTO.Username)
	if err != nil || !flag {
		return dto.UserLoginResponse{}, dto.ErrUsernameNotFound
	}

	checkPassword, err := utils.CheckPassword(check.Password, []byte(loginDTO.Password))
	if err != nil || !checkPassword {
		return dto.UserLoginResponse{}, dto.ErrPasswordNotMatch
	}

	token := us.jwtService.GenerateToken(check.ID.String(), check.Role.Name)

	return dto.UserLoginResponse{
		Token: token,
		RoleID:  check.RoleID,
	}, nil
}

func (us *userService) GetUserById(ctx context.Context, userId string) (dto.UserResponse, error) {
	user, err := us.userDomain.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserById
	}

	return dto.UserResponse{
		ID:         user.ID.String(),
		Username:   user.Username,
		Name:       user.Name,
		Noid:       user.Noid,
		RoleID:       user.RoleID,
		Email:      user.Email,

	}, nil
}

func (us *userService) Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error) {
	user, err := us.userDomain.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUserNotFound
	}

	// if(user.RoleID != 0){
	// 	return dto.UserUpdateResponse{}, dto.ErrDeniedAccess
	// }

	data := entity.User{
		ID:         user.ID,
		Username:   req.Username,
		Name:       req.Name,
		RoleID:       user.RoleID,
		Email:      req.Email,
		Noid:       req.Noid,
	}

	userUpdate, err := us.userDomain.UpdateUser(ctx, nil, data)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUpdateUser
	}

	return dto.UserUpdateResponse{
		ID:         userUpdate.ID.String(),
		Username:   userUpdate.Username,
		Name:       userUpdate.Name,
		Noid:       userUpdate.Noid,
		RoleID:       userUpdate.RoleID,
		Email:      userUpdate.Email,
	}, nil
}

func (us *userService) UpdateMe(ctx context.Context, req dto.UserUpdateEmailRequest, userId string) (dto.UserUpdateResponse, error) {
	user, err := us.userDomain.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUserNotFound
	}

	data := entity.User{
		ID:         user.ID,
		Email:      req.Email,
	}

	userUpdate, err := us.userDomain.UpdateUser(ctx, nil, data)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUpdateUser
	}

	return dto.UserUpdateResponse{
		ID:         userUpdate.ID.String(),
		Username:   userUpdate.Username,
		Name:       userUpdate.Name,
		Noid:       userUpdate.Noid,
		RoleID:       userUpdate.RoleID,
		Email:      userUpdate.Email,
	}, nil
}

func (us *userService) Delete(ctx context.Context, userId string) error {
	user, err := us.userDomain.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// if(user.RoleID != 0){
	// 	return dto.ErrDeniedAccess
	// }

	err = us.userDomain.DeleteUser(ctx, nil, user.ID.String())
	if err != nil {
		return dto.ErrDeleteUser
	}

	return nil
}

func (us *userService) RegisterUsersFromYAML(ctx context.Context, fileHeader *multipart.FileHeader) (map[string]interface{}, error) {
	file, err := fileHeader.Open()
    if err != nil {
        return nil, fmt.Errorf("failed to open uploaded YAML file: %w", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read YAML file: %w", err)
    }

    if len(data) == 0 {
        return nil, fmt.Errorf("uploaded YAML file is empty")
    }

    var yamlData dto.UserFileList
    if err := yaml.Unmarshal(data, &yamlData); err != nil {
        return nil, fmt.Errorf("invalid YAML format: %w", err)
    }

    if len(yamlData.Users) == 0 {
        return nil, fmt.Errorf("no users found in YAML")
    }

	var createdUsers []dto.UserResponse = []dto.UserResponse{}
	var failedUsers []dto.FailedUserResponse = []dto.FailedUserResponse{}


	if len(yamlData.Users) == 0 {
		return nil, fmt.Errorf("YAML file contains no users")
	}

	for _, user := range yamlData.Users {
		if user.Username == "" || user.Noid == "" || user.Name == "" || user.Email == "" || user.Password == "" {
            failedUsers = append(failedUsers, dto.FailedUserResponse{
                Noid:   user.Noid,
                Email:  user.Email,
                Reason: "Missing required fields",
            })
            continue
        }
		_, flag, _ := us.userDomain.CheckNoid(ctx, nil, user.Noid)
		if flag {
			failedUsers = append(failedUsers, dto.FailedUserResponse{
				Noid:   user.Noid,
				Email:  user.Email,
				Reason: dto.ErrNoidAlreadyExists.Error(),
			})
			continue
		}

		req := entity.User{
			Username: user.Username,
			Name:     user.Name,
			Noid:     user.Noid,
			RoleID:   2,
			Email:    user.Email,
			Password: user.Password,
		}

		userReg, err := us.userDomain.RegisterUser(ctx, nil, req)
		resp := dto.UserResponse{
			ID:    userReg.ID.String(),
			Username: userReg.Username,
			Name:  userReg.Name,
			Noid:  userReg.Noid,
			RoleID:  userReg.RoleID,
			Email: userReg.Email,
		}

		if err != nil {
			failedUsers = append(failedUsers, dto.FailedUserResponse{
                Noid:   user.Noid,
                Email:  user.Email,
                Reason: fmt.Sprintf("Database error: %v", err.Error()),
            })
            continue
		} else {
			createdUsers = append(createdUsers, resp)
		}
	}

	return map[string]interface{}{
		"created_users": createdUsers,
		"failed_users":  failedUsers,
	}, nil
}

func (us *userService) RegisterUsersFromCSV(ctx context.Context, fileHeader *multipart.FileHeader) (map[string]interface{}, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) < 1 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	headers := records[0]
	fieldMap := map[string]int{}
	for i, h := range headers {
		fieldMap[strings.ToLower(h)] = i
	}

	var createdUsers []dto.UserResponse
	var failedUsers []dto.FailedUserResponse

	for _, record := range records[1:] {
		get := func(key string) string {
			if idx, ok := fieldMap[strings.ToLower(key)]; ok && idx < len(record) {
				return strings.TrimSpace(record[idx])
			}
			return ""
		}

		user := dto.UserFile{
			Username: get("username"),
			Name:     get("name"),
			Email:    get("email"),
			Noid:     get("noid"),
			Password: get("password"),
		}

		if user.Username == "" || user.Noid == "" || user.Name == "" || user.Email == "" || user.Password == "" {
			failedUsers = append(failedUsers, dto.FailedUserResponse{
				Noid:   user.Noid,
				Email:  user.Email,
				Reason: "Missing required fields",
			})
			continue
		}

		_, exists, _ := us.userDomain.CheckNoid(ctx, nil, user.Noid)
		if exists {
			failedUsers = append(failedUsers, dto.FailedUserResponse{
				Noid:   user.Noid,
				Email:  user.Email,
				Reason: dto.ErrNoidAlreadyExists.Error(),
			})
			continue
		}

		newUser := entity.User{
			Username: user.Username,
			Name:     user.Name,
			Noid:     user.Noid,
			RoleID:   2,
			Email:    user.Email,
			Password: user.Password,
		}

		userReg, err := us.userDomain.RegisterUser(ctx, nil, newUser)
		if err != nil {
			failedUsers = append(failedUsers, dto.FailedUserResponse{
				Noid:   user.Noid,
				Email:  user.Email,
				Reason: fmt.Sprintf("Database error: %v", err.Error()),
			})
			continue
		}

		createdUsers = append(createdUsers, dto.UserResponse{
			ID:       userReg.ID.String(),
			Username: userReg.Username,
			Name:     userReg.Name,
			Noid:     userReg.Noid,
			RoleID:   userReg.RoleID,
			Email:    userReg.Email,
		})
	}

	return map[string]interface{}{
		"created_users": createdUsers,
		"failed_users":  failedUsers,
	}, nil
}



func makeForgotPasswordEmail(receiverEmail string) (map[string]string, error) {
	expired := time.Now().Add(time.Hour * 1).Format("2006-01-02 15:04:05") // token valid for 1 hour
	plainText := receiverEmail + "_" + expired
	token, err := utils.AESEncrypt(plainText)
	if err != nil {
		return nil, err
	}

	resetLink := LOCAL_URL + "/" + RESET_PASSWORD_ROUTE + "?token=" + token

	readHtml, err := os.ReadFile("utils/email-template/forgot_password_mail.html")
	if err != nil {
		return nil, err
	}

	data := struct {
		Email string
		Reset string
	}{
		Email: receiverEmail,
		Reset: resetLink,
	}

	tmpl, err := template.New("forgot-password").Parse(string(readHtml))
	if err != nil {
		return nil, err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return nil, err
	}

	draftEmail := map[string]string{
		"subject": "Reset Your Password",
		"body":    strMail.String(),
	}

	return draftEmail, nil
}

func (s *userService) SendForgotPasswordEmail(ctx context.Context, req dto.SendResetPasswordRequest) error {
	user, err := s.userDomain.GetUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.ErrEmailNotFound
	}

	draftEmail, err := makeForgotPasswordEmail(user.Email)
	if err != nil {
		return err
	}

	err = utils.SendMail(user.Email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (dto.ResetPasswordResponse, error) {
	decryptedToken, err := utils.AESDecrypt(req.Token)
	if err != nil {
		return dto.ResetPasswordResponse{}, dto.ErrTokenInvalid
	}

	if !strings.Contains(decryptedToken, "_") {
		return dto.ResetPasswordResponse{}, dto.ErrTokenInvalid
	}

	decryptedTokenSplit := strings.Split(decryptedToken, "_")
	email := decryptedTokenSplit[0]
	expired := decryptedTokenSplit[1]

	now := time.Now()
	layout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	expiredTime, err := time.ParseInLocation(layout, expired, loc)
	if err != nil {
		return dto.ResetPasswordResponse{}, dto.ErrTokenInvalid
	}

	fmt.Println("Is expired?", expiredTime.Before(now))
	if expiredTime.Before(now) {
		return dto.ResetPasswordResponse{}, dto.ErrTokenExpired
	}

	user, err := s.userDomain.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.ResetPasswordResponse{}, dto.ErrUserNotFound
	}

	// Update password (you should hash the password first!)
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return dto.ResetPasswordResponse{}, dto.ErrHashPassword
	}
	_, err = s.userDomain.UpdateUser(ctx, nil, entity.User{
		ID:       user.ID,
		Password: hashedPassword,
	})
	if err != nil {
		return dto.ResetPasswordResponse{}, dto.ErrUpdateUser
	}

	return dto.ResetPasswordResponse{
		Email: user.Email,
	}, nil
}
