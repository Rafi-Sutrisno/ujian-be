package service

import (
	"archive/zip"
	"context"
	"io"
	"mods/domain/entity"
	domain "mods/domain/repository"
	"mods/interface/dto"
	"path/filepath"
	"strings"
)

type (
	testCaseService struct {
		repo domain.TestCaseRepository
	}

	TestCaseService interface {
		GetByID(ctx context.Context, id string) (dto.TestCaseResponse, error)
		GetByProblemID(ctx context.Context, problemID string, userId string) ([]dto.TestCaseResponse, error)
		GetAll(ctx context.Context) ([]dto.TestCaseResponse, error)
		Create(ctx context.Context, req dto.TestCaseCreateRequest, userId string) (dto.TestCaseResponse, error)
		CreateFromZip(ctx context.Context, zipPath, problemID, userId string) (int, error) 
		Update(ctx context.Context, req dto.TestCaseUpdateRequest, id string, userId string) (dto.TestCaseUpdateResponse, error)
		Delete(ctx context.Context, id string, userId string) error
	}
)

func NewTestCaseService(repo domain.TestCaseRepository) TestCaseService {
	return &testCaseService{
		repo: repo,
	}
}

func (ts *testCaseService) GetByID(ctx context.Context, id string) (dto.TestCaseResponse, error) {
	testCase, err := ts.repo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.TestCaseResponse{}, err
	}

	return dto.TestCaseResponse{
		ID:            testCase.ID.String(),
		ProblemID:     testCase.ProblemID,
		InputData:     testCase.InputData,
		ExpectedOutput: testCase.ExpectedOutput,
	}, nil
}

func (ts *testCaseService) GetByProblemID(ctx context.Context, problemID, userId string) ([]dto.TestCaseResponse, error) {
	testCases, err := ts.repo.GetByProblemID(ctx, nil, problemID)
	if err != nil {
		return nil, err
	}

	var responses []dto.TestCaseResponse
	for _, testCase := range testCases {
		responses = append(responses, dto.TestCaseResponse{
			ID:            testCase.ID.String(),
			ProblemID:     testCase.ProblemID,
			InputData:     testCase.InputData,
			ExpectedOutput: testCase.ExpectedOutput,
			CreatedAt: testCase.CreatedAt.String(),
		})
	}

	return responses, nil
}

func (ts *testCaseService) GetAll(ctx context.Context) ([]dto.TestCaseResponse, error) {
	testCases, err := ts.repo.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	var responses []dto.TestCaseResponse
	for _, testCase := range testCases {
		responses = append(responses, dto.TestCaseResponse{
			ID:            testCase.ID.String(),
			ProblemID:     testCase.ProblemID,
			InputData:     testCase.InputData,
			ExpectedOutput: testCase.ExpectedOutput,
		})
	}

	return responses, nil
}

func (ts *testCaseService) Create(ctx context.Context, req dto.TestCaseCreateRequest, userId string) (dto.TestCaseResponse, error) {
	testCase := entity.TestCase{
		ProblemID:      req.ProblemID,
		InputData:      req.InputData,
		ExpectedOutput: req.ExpectedOutput,
	}

	createdTestCase, err := ts.repo.Create(ctx, nil, testCase)
	if err != nil {
		return dto.TestCaseResponse{}, err
	}

	return dto.TestCaseResponse{
		ID:            createdTestCase.ID.String(),
		ProblemID:     createdTestCase.ProblemID,
		InputData:     createdTestCase.InputData,
		ExpectedOutput: createdTestCase.ExpectedOutput,
	}, nil
}

func (ts *testCaseService) CreateFromZip(ctx context.Context, zipPath, problemID, userId string) (int, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	inputs := map[string]string{}
	outputs := map[string]string{}

	// Extract files and separate .in/.out
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return 0, err
		}
		content, _ := io.ReadAll(rc)
		rc.Close()

		name := strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
		if strings.HasSuffix(f.Name, ".in") {
			inputs[name] = string(content)
		} else if strings.HasSuffix(f.Name, ".out") {
			outputs[name] = string(content)
		}
	}

	// Match and insert test cases
	count := 0
	for name, inData := range inputs {
		outData, exists := outputs[name]
		if !exists {
			continue // skip unmatched
		}
		tc := entity.TestCase{
			ProblemID:      problemID,
			InputData:      inData,
			ExpectedOutput: outData,
		}
		_, err := ts.repo.Create(ctx, nil, tc)
		if err != nil {
			return count, err // or continue if partial success is okay
		}
		count++
	}

	return count, nil
}


func (ts *testCaseService) Update(ctx context.Context, req dto.TestCaseUpdateRequest, id string, userId string) (dto.TestCaseUpdateResponse, error) {
	testCase, err := ts.repo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.TestCaseUpdateResponse{}, dto.ErrTestCaseNotFound
	}

	updatedTestCase := entity.TestCase{
		ID:             testCase.ID,
		ProblemID:      testCase.ProblemID,
		InputData:      req.InputData,
		ExpectedOutput: req.ExpectedOutput,
	}

	updated, err := ts.repo.Update(ctx, nil, updatedTestCase)
	if err != nil {
		return dto.TestCaseUpdateResponse{}, err
	}

	return dto.TestCaseUpdateResponse{
		ID:            updated.ID.String(),
		ProblemID:     updated.ProblemID,
		InputData:     updated.InputData,
		ExpectedOutput: updated.ExpectedOutput,
	}, nil
}

func (ts *testCaseService) Delete(ctx context.Context, id string, userId string) error {
	if err := ts.repo.Delete(ctx, nil, id); err != nil {
		return dto.ErrDeleteTestCase
	}
	return nil
}
