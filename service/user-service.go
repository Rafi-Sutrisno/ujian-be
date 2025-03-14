package service

import (
	"context"
	"mods/dto"
	"mods/entity"
	"mods/repository"
	"mods/utils"
)

type userService struct {
	userRepository repository.UserRepository
	jwtService  JWTService
}

type UserService interface {
	// functional
	CreateUser(ctx context.Context, userDTO dto.UserCreateRequest) (entity.User, error)
	GetAllUser(ctx context.Context) ([]entity.User, error)
	Verify(ctx context.Context, loginDTO dto.UserLoginRequest) (dto.UserLoginResponse, error)
}

func NewUserService(ur repository.UserRepository, jwtService JWTService) UserService {
	return &userService{
		userRepository: ur,
		jwtService:  jwtService,
	}
}

func (us *userService) CreateUser(ctx context.Context, userDTO dto.UserCreateRequest) (entity.User, error) {
	newUser := entity.User{
		Name:  userDTO.Name,
		Password: userDTO.Password,
		Email: userDTO.Email,
		NoId: userDTO.NoId,
		Role: "user",
	}

	return us.userRepository.CreateUser(ctx, newUser)
}

func (us *userService) GetAllUser(ctx context.Context) ([]entity.User, error) {
	return us.userRepository.GetAllUser(ctx)
}

func (us *userService) Verify(ctx context.Context, loginDTO dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	check, flag, err := us.userRepository.CheckNoId(ctx, loginDTO.NoId)
	if err != nil || !flag {
		return dto.UserLoginResponse{}, dto.ErrEmailNotFound
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
