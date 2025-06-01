package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"
	"strings"

	"gopkg.in/yaml.v3"
)

type (
	userClassService struct {
		repo domain.UserClassRepository
		UserDomain domain.UserRepository
	}

	UserClassService interface {
		AssignUsersFromYAML(ctx context.Context, class_id string, fileHeader *multipart.FileHeader) (map[string]interface{}, error)
		AssignUsersFromCSV(ctx context.Context, class_id string, fileHeader *multipart.FileHeader) (map[string]interface{}, error)
		GetByUserID(ctx context.Context, userID string) ([]dto.UserClassResponse, error)
		GetByClassID(ctx context.Context, classID string, userId string) ([]dto.UserClassResponse, error)
		GetUnassignedUsersByClassID(ctx context.Context, classID string, userId string) ([]dto.UserResponse, error)
		Create(ctx context.Context, req dto.UserClassCreateRequest, userId string) (dto.UserClassResponse, error)
		CreateMany(ctx context.Context, reqs []dto.UserClassCreateRequest, userId string) error
		Delete(ctx context.Context, id string, userId string) error
	}
)

func NewUserClassService(repo domain.UserClassRepository, UserDomain domain.UserRepository) UserClassService {
	return &userClassService{
		repo: repo,
		UserDomain:  UserDomain,
	}
}

func (ucs *userClassService) AssignUsersFromYAML(ctx context.Context, class_id string, fileHeader *multipart.FileHeader) (map[string]interface{}, error) {
	fmt.Println("masuk service yaml assign file")
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

	fmt.Println("check 1")

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
		userCheck, flag, _ := ucs.UserDomain.CheckNoid(ctx, nil, user.Noid)
		
		if flag {
			fmt.Println("sudah ada user: ", userCheck)
			userClass := entity.UserClass{
				UserID:    userCheck.ID.String(),
				ClassID:   class_id,
			}

			flag, _ := ucs.repo.CheckExist(ctx, nil, userClass.UserID, userClass.ClassID)
			if flag {
				failedUsers = append(failedUsers, dto.FailedUserResponse{
					Noid:   user.Noid,
					Email:  user.Email,
					Reason: "This User is Already Assigned in this Class",
				})
				continue
			}
		
			_, err := ucs.repo.Create(ctx, nil, userClass)
			resp := dto.UserResponse{
				ID:    userCheck.ID.String(),
				Username: userCheck.Username,
				Name:  userCheck.Name,
				Noid:  userCheck.Noid,
				RoleID:  userCheck.RoleID,
				Email: userCheck.Email,
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
			continue
		}

		fmt.Println("belum ada user: ", user)
		req := entity.User{
			Username: user.Username,
			Name:     user.Name,
			Noid:     user.Noid,
			RoleID:   2,
			Email:    user.Email,
			Password: user.Password,
		}

		userReg, err := ucs.UserDomain.RegisterUser(ctx, nil, req)
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
			userClass := entity.UserClass{
				UserID:    resp.ID,
				ClassID:   class_id,
				
			}
			_, err := ucs.repo.Create(ctx, nil, userClass)
	
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
	}

	return map[string]interface{}{
		"created_users": createdUsers,
		"failed_users":  failedUsers,
	}, nil
}

func (ucs *userClassService) AssignUsersFromCSV(ctx context.Context, class_id string, fileHeader *multipart.FileHeader) (map[string]interface{}, error) {
	fmt.Println("masuk service csv assign file")
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

		userCheck, exists, _ := ucs.UserDomain.CheckNoid(ctx, nil, user.Noid)
		if exists {
			fmt.Println("sudah ada user: ", userCheck)
			userClass := entity.UserClass{
				UserID:    userCheck.ID.String(),
				ClassID:   class_id,
			}

			flag, _ := ucs.repo.CheckExist(ctx, nil, userClass.UserID, userClass.ClassID)
			if flag {
				failedUsers = append(failedUsers, dto.FailedUserResponse{
					Noid:   user.Noid,
					Email:  user.Email,
					Reason: "This User is Already Assigned in this Class",
				})
				continue
			}
		
			_, err := ucs.repo.Create(ctx, nil, userClass)
			resp := dto.UserResponse{
				ID:    userCheck.ID.String(),
				Username: userCheck.Username,
				Name:  userCheck.Name,
				Noid:  userCheck.Noid,
				RoleID:  userCheck.RoleID,
				Email: userCheck.Email,
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
			continue
		}

		
		fmt.Println("belum ada user: ", user)
		req := entity.User{
			Username: user.Username,
			Name:     user.Name,
			Noid:     user.Noid,
			RoleID:   2,
			Email:    user.Email,
			Password: user.Password,
		}

		userReg, err := ucs.UserDomain.RegisterUser(ctx, nil, req)
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
			userClass := entity.UserClass{
				UserID:    resp.ID,
				ClassID:   class_id,
				
			}
			_, err := ucs.repo.Create(ctx, nil, userClass)
	
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
	}

	return map[string]interface{}{
		"created_users": createdUsers,
		"failed_users":  failedUsers,
	}, nil
}


func (ucs *userClassService) GetByUserID(ctx context.Context, userID string) ([]dto.UserClassResponse, error) {
	userClasses, err := ucs.repo.GetByUserID(ctx, nil, userID)
	if err != nil {
		return nil, dto.ErrGetAllUserClassByUserId
	}

	var responses []dto.UserClassResponse
	for _, uc := range userClasses {
		responses = append(responses, dto.UserClassResponse{
			ID:        uc.ID.String(),
			UserID:    uc.UserID,
			ClassID:   uc.ClassID,
			
		})
	}

	return responses, nil
}

func (ucs *userClassService) GetByClassID(ctx context.Context, classID string, userId string) ([]dto.UserClassResponse, error) {
	exists, err := ucs.repo.IsUserInClass(ctx, nil, userId, classID)
	if err != nil {
		return nil, dto.ErrAuthorize
	}

	if !exists {
		return nil, dto.ErrAuthorize 
	}

	userClasses, err := ucs.repo.GetByClassID(ctx, nil, classID)
	if err != nil {
		return nil, dto.ErrGetAllUserClassByClassId
	}

	var responses []dto.UserClassResponse
	for _, uc := range userClasses {
		user := &dto.UserResponse{
			Name:       uc.User.Name,
			Noid: 		uc.User.Noid,
			RoleID:     uc.User.RoleID,
			Email:      uc.User.Email,
		}

		responses = append(responses, dto.UserClassResponse{
			ID:        uc.ID.String(),
			UserID:    uc.UserID,
			ClassID:   uc.ClassID,
			User:      user,
		})
	}

	return responses, nil
}

func (ucs *userClassService) GetUnassignedUsersByClassID(ctx context.Context, classID string, userId string) ([]dto.UserResponse, error) {
	exists, err := ucs.repo.IsUserInClass(ctx, nil, userId, classID)
	if err != nil {
		return nil, dto.ErrAuthorize
	}

	if !exists {
		return nil, dto.ErrAuthorize 
	}

	allStudents, err := ucs.UserDomain.GetAllStudents(ctx, nil)
	if err != nil {
		return nil, err
	}

	assignedUserClasses, err := ucs.repo.GetByClassID(ctx, nil, classID)
	if err != nil {
		return nil, err
	}

	assignedMap := make(map[string]bool)
	for _, uc := range assignedUserClasses {
		assignedMap[uc.UserID] = true
	}

	var unassignedUsers []dto.UserResponse
	for _, student := range allStudents {
		if !assignedMap[student.ID.String()] {
			unassignedUsers = append(unassignedUsers, dto.UserResponse{
				ID:    student.ID.String(),
				Name:  student.Name,
				Noid: 		student.Noid,
				RoleID:       student.RoleID,
				Email: student.Email,
				
			})
		}
	}

	return unassignedUsers, nil
}


func (ucs *userClassService) Create(ctx context.Context, req dto.UserClassCreateRequest, userId string) (dto.UserClassResponse, error) {
	exists, err := ucs.repo.IsUserInClass(ctx, nil, userId, req.ClassID)
	if err != nil {
		return dto.UserClassResponse{}, dto.ErrAuthorize
	}

	if !exists {
		return dto.UserClassResponse{}, dto.ErrAuthorize // or any custom error you want
	}

	userClass := entity.UserClass{
		UserID:    req.UserID,
		ClassID:   req.ClassID,
		
	}

	createdUserClass, err := ucs.repo.Create(ctx, nil, userClass)
	if err != nil {
		return dto.UserClassResponse{}, dto.ErrCreateUserClass
	}

	return dto.UserClassResponse{
		ID:        createdUserClass.ID.String(),
		UserID:    createdUserClass.UserID,
		ClassID:   createdUserClass.ClassID,
		
	}, nil
}

func (ucs *userClassService) CreateMany(ctx context.Context, reqs []dto.UserClassCreateRequest, userId string) error {
	exists, err := ucs.repo.IsUserInClass(ctx, nil, userId, reqs[0].ClassID)
	if err != nil {
		return dto.ErrAuthorize
	}

	if !exists {
		return dto.ErrAuthorize // or any custom error you want
	}
	
	var userClasses []entity.UserClass
	for _, req := range reqs {
		userClasses = append(userClasses, entity.UserClass{
			UserID:    req.UserID,
			ClassID:   req.ClassID,
			
		})
	}

	if err := ucs.repo.CreateMany(ctx, nil, userClasses); err != nil {
		return dto.ErrCreateUserClass
	}

	return nil
}

func (ucs *userClassService) Delete(ctx context.Context, id string, userId string) error {
	userClasses, err := ucs.repo.GetById(ctx, nil, id)
	if err != nil {
		return dto.ErrGetAllUserClass
	}

	exists, err := ucs.repo.IsUserInClass(ctx, nil, userId, userClasses.ClassID)
	if err != nil {
		return  dto.ErrAuthorize
	}

	if !exists {
		return  dto.ErrAuthorize // or any custom error you want
	}
	userClass, err := ucs.repo.GetById(ctx, nil, id)
	if err != nil {
		return dto.ErrUserClassNotFound
	}

	if err := ucs.repo.Delete(ctx, nil, userClass.ID.String()); err != nil {
		return dto.ErrDeleteUserClass
	}
	return nil
}
