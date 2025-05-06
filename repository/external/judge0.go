package judge0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mods/config"
	"mods/dto"
	"net/http"
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
