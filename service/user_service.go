package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"mods/constants"
	"mods/dto"
	"mods/entity"
	"mods/repository"

	"mods/utils"

	"gopkg.in/yaml.v3"
)

type userService struct {
	userRepository repository.UserRepository
	jwtService  JWTService
	
}

type UserService interface {
	// functional
	Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error)
	GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error)
	Verify(ctx context.Context, loginDTO dto.UserLoginRequest) (dto.UserLoginResponse, error)
	Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error)
	UpdateMe(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error)
	Delete(ctx context.Context, userId string) error
	GetUserById(ctx context.Context, userId string) (dto.UserResponse, error)

	RegisterUsersFromYAML(ctx context.Context, fileHeader *multipart.FileHeader) (map[string]interface{}, error)
}

func NewUserService(ur repository.UserRepository, jwtService JWTService) UserService {
	return &userService{
		userRepository: ur,
		jwtService:  jwtService,
	}
}

func (us *userService) Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error) {
	_, flag, _ := us.userRepository.CheckNoId(ctx, nil, req.Noid)
	if flag {
		return dto.UserResponse{}, dto.ErrNoidAlreadyExists
	}


	user := entity.User{
		Name:       req.Name,
		Noid:       req.Noid,
		Role:       constants.ENUM_ROLE_USER,
		Email:      req.Email,
		Password:   req.Password,
	}

	userReg, err := us.userRepository.RegisterUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	return dto.UserResponse{
		ID:         userReg.ID.String(),
		Name:       userReg.Name,
		Noid: 		userReg.Noid,
		Role:       userReg.Role,
		Email:      userReg.Email,
	}, nil
}


func (us *userService) GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := us.userRepository.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	var datas []dto.UserResponse
	for _, user := range dataWithPaginate.Users {
		data := dto.UserResponse{
			ID:         user.ID.String(),
			Name:       user.Name,
			Email:      user.Email,
			Role:       user.Role,
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
	check, flag, err := us.userRepository.CheckNoId(ctx,nil,  loginDTO.Noid)
	if err != nil || !flag {
		return dto.UserLoginResponse{}, dto.ErrNRPNotFound
	}

	checkPassword, err := utils.CheckPassword(check.Password, []byte(loginDTO.Password))
	if err != nil || !checkPassword {
		return dto.UserLoginResponse{}, dto.ErrPasswordNotMatch
	}

	token := us.jwtService.GenerateToken(check.ID.String(), check.Role)

	return dto.UserLoginResponse{
		Token: token,
		Role:  check.Role,
	}, nil
}

func (us *userService) GetUserById(ctx context.Context, userId string) (dto.UserResponse, error) {
	user, err := us.userRepository.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserById
	}

	return dto.UserResponse{
		ID:         user.ID.String(),
		Name:       user.Name,
		Noid:       user.Noid,
		Role:       user.Role,
		Email:      user.Email,

	}, nil
}

func (us *userService) Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error) {
	user, err := us.userRepository.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUserNotFound
	}

	if(user.Role != "user"){
		return dto.UserUpdateResponse{}, dto.ErrDeniedAccess
	}

	data := entity.User{
		ID:         user.ID,
		Name:       req.Name,
		Role:       user.Role,
		Email:      req.Email,
		Noid:       req.Noid,
	}

	userUpdate, err := us.userRepository.UpdateUser(ctx, nil, data)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUpdateUser
	}

	return dto.UserUpdateResponse{
		ID:         userUpdate.ID.String(),
		Name:       userUpdate.Name,
		Noid:       userUpdate.Noid,
		Role:       userUpdate.Role,
		Email:      userUpdate.Email,
	}, nil
}

func (us *userService) UpdateMe(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error) {
	user, err := us.userRepository.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUserNotFound
	}

	data := entity.User{
		ID:         user.ID,
		Name:       req.Name,
		Role:       user.Role,
		Email:      req.Email,
		Noid:       req.Noid,
	}

	userUpdate, err := us.userRepository.UpdateUser(ctx, nil, data)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUpdateUser
	}

	return dto.UserUpdateResponse{
		ID:         userUpdate.ID.String(),
		Name:       userUpdate.Name,
		Noid:       userUpdate.Noid,
		Role:       userUpdate.Role,
		Email:      userUpdate.Email,
	}, nil
}

func (us *userService) Delete(ctx context.Context, userId string) error {
	user, err := us.userRepository.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.ErrUserNotFound
	}

	if(user.Role != "user"){
		return dto.ErrDeniedAccess
	}

	err = us.userRepository.DeleteUser(ctx, nil, user.ID.String())
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

    var yamlData dto.UserYAMLList
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
		if user.Noid == "" || user.Name == "" || user.Email == "" || user.Password == "" {
            failedUsers = append(failedUsers, dto.FailedUserResponse{
                Noid:   user.Noid,
                Email:  user.Email,
                Reason: "Missing required fields",
            })
            continue
        }
		_, flag, _ := us.userRepository.CheckNoId(ctx, nil, user.Noid)
		if flag {
			failedUsers = append(failedUsers, dto.FailedUserResponse{
				Noid:   user.Noid,
				Email:  user.Email,
				Reason: dto.ErrNoidAlreadyExists.Error(),
			})
			continue
		}

		req := entity.User{
			Name:     user.Name,
			Noid:     user.Noid,
			Role:     constants.ENUM_ROLE_USER,
			Email:    user.Email,
			Password: user.Password,
		}

		userReg, err := us.userRepository.RegisterUser(ctx, nil, req)
		resp := dto.UserResponse{
			ID:    userReg.ID.String(),
			Name:  userReg.Name,
			Noid:  userReg.Noid,
			Role:  userReg.Role,
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
