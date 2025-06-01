package judge0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mods/config"
	"mods/interface/dto"
	"net/http"
	"strings"
)

func SubmitToJudge0(req dto.Judge0Request) (dto.Judge0Response, error) {
	cfg := config.LoadJudge0Config()

	url := "https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=true&wait=true&fields=*"
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return dto.Judge0Response{}, fmt.Errorf("failed to marshal request body: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return dto.Judge0Response{}, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-rapidapi-host", cfg.Host)
	httpReq.Header.Set("x-rapidapi-key", cfg.Key)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return dto.Judge0Response{}, fmt.Errorf("failed to send request to Judge0: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.Judge0Response{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var judgeResp dto.Judge0Response
	err = json.Unmarshal(respBody, &judgeResp)
	if err != nil {
		return dto.Judge0Response{}, fmt.Errorf("failed to parse Judge0 response: %v", err)
	}

	return judgeResp, nil
}

// func SubmitToJudge0Expected(req dto.Judge0SubmissionRequest) (dto.Judge0Response, error) {
// 	cfg := config.LoadJudge0Config()

// 	url := "https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=true&wait=true&fields=*"
// 	bodyBytes, err := json.Marshal(req)
// 	if err != nil {
// 		return dto.Judge0Response{}, fmt.Errorf("failed to marshal request body: %v", err)
// 	}

// 	fmt.Println("Judge0 JSON request body:", string(bodyBytes))

// 	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
// 	if err != nil {
// 		return dto.Judge0Response{}, fmt.Errorf("failed to create HTTP request: %v", err)
// 	}

// 	httpReq.Header.Set("Content-Type", "application/json")
// 	httpReq.Header.Set("x-rapidapi-host", cfg.Host)
// 	httpReq.Header.Set("x-rapidapi-key", cfg.Key)

// 	client := &http.Client{}
// 	resp, err := client.Do(httpReq)
// 	if err != nil {
// 		return dto.Judge0Response{}, fmt.Errorf("failed to send request to Judge0: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	respBody, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return dto.Judge0Response{}, fmt.Errorf("failed to read response body: %v", err)
// 	}

// 	fmt.Println("raw response body:", string(respBody)) // ADD THIS

// 	var judgeResp dto.Judge0Response
// 	err = json.Unmarshal(respBody, &judgeResp)
// 	if err != nil {
// 		return dto.Judge0Response{}, fmt.Errorf("failed to parse Judge0 response: %v", err)
// 	}

// 	return judgeResp, nil
// }

// func RunWithExpectedOutput(req dto.Judge0SubmissionRequest) (dto.Judge0SubmissionResult, error) {
// 	encodedSource := base64.StdEncoding.EncodeToString([]byte(req.SourceCode))
// 	encodedStdin := base64.StdEncoding.EncodeToString([]byte(req.Stdin + "\n"))
// 	encodedExpectedOutput := base64.StdEncoding.EncodeToString([]byte(req.ExpectedOutput + "\n"))


// 	request := dto.Judge0SubmissionRequest{
// 		LanguageID:     req.LanguageID,
// 		SourceCode:     encodedSource,
// 		Stdin:          encodedStdin,
// 		ExpectedOutput: encodedExpectedOutput,
// 	}

// 	fmt.Println("ini request:", request)

// 	resp, err := SubmitToJudge0Expected(request)
// 	if err != nil {
// 		return dto.Judge0SubmissionResult{}, err
// 	}

// 	fmt.Println("ini response:", resp)

// 	result := dto.Judge0SubmissionResult{
// 		StatusID:          resp.Status.ID,
// 		StatusDescription: resp.Status.Description,
// 	}

// 	fmt.Println("ini result:", result)

// 	if resp.Status.ID == 0 {
// 		return dto.Judge0SubmissionResult{}, fmt.Errorf("received empty status from Judge0: %+v", resp)
// 	}
	

// 	switch resp.Status.ID {
// 	case 3:
// 		result.Passed = true
// 	case 6, 7: // Compilation Error or similar
// 		result.CompilationError = true
// 	}

// 	return result, nil
// }

func SubmitToJudge0Batch(req dto.Judge0BatchSubmissionRequest) (dto.Judge0BatchSubmissionResponse, error) {
	cfg := config.LoadJudge0Config()
	url := "https://judge0-ce.p.rapidapi.com/submissions/batch?base64_encoded=true"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return dto.Judge0BatchSubmissionResponse{}, fmt.Errorf("marshal error: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return dto.Judge0BatchSubmissionResponse{}, fmt.Errorf("request error: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-rapidapi-host", cfg.Host)
	httpReq.Header.Set("x-rapidapi-key", cfg.Key)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return dto.Judge0BatchSubmissionResponse{}, fmt.Errorf("http error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.Judge0BatchSubmissionResponse{}, err
	}

	// Debug: print raw response
	fmt.Println("Judge0 raw response:", string(body))

	var batchResp dto.Judge0BatchSubmissionResponse
	if err := json.Unmarshal(body, &batchResp); err != nil {
		return dto.Judge0BatchSubmissionResponse{}, fmt.Errorf("unmarshal error: %v", err)
	}

	return batchResp, nil
}


func GetBatchResults(tokens []string) (dto.Judge0BatchResultResponse, error) {
	cfg := config.LoadJudge0Config()
	url := fmt.Sprintf("https://judge0-ce.p.rapidapi.com/submissions/batch?tokens=%s&base64_encoded=true", strings.Join(tokens, ","))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dto.Judge0BatchResultResponse{}, err
	}

	req.Header.Set("x-rapidapi-host", cfg.Host)
	req.Header.Set("x-rapidapi-key", cfg.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return dto.Judge0BatchResultResponse{}, err
	}
	defer resp.Body.Close()

	var res dto.Judge0BatchResultResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.Judge0BatchResultResponse{}, err
	}
	// Debug (optional)
	fmt.Println("Judge0 batch result raw response:", string(body))
	if err := json.Unmarshal(body, &res); err != nil {
		return dto.Judge0BatchResultResponse{}, err
	}
	return res, nil
}

