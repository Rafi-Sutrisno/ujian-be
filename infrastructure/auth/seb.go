package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidateSEBRequest handles all SEB header logic and validation
func ValidateSEBRequest(ctx *gin.Context, sebBrowserKey, sebConfigKey string, isRestricted bool) error {
	
	if !isRestricted {
		return nil // No SEB restrictions
	}

	userAgent := ctx.Request.UserAgent()
	requestHash := ctx.GetHeader("X-SafeExamBrowser-RequestHash")
	configKeyHash := ctx.GetHeader("X-Safeexambrowser-Configkeyhash")

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	fullURL := fmt.Sprintf("%s://%s%s", scheme, ctx.Request.Host, ctx.Request.RequestURI)

	// If not in header, fall back to request body (already parsed in controller)
	var body struct {
		BrowserExamKey string `json:"browser_exam_key"`
		ConfigKey      string `json:"config_key"`
	}
	if err := ctx.ShouldBind(&body); err == nil {
		if requestHash == "" {
			requestHash = body.BrowserExamKey
		}
		if configKeyHash == "" {
			configKeyHash = body.ConfigKey
		}
	}

	// Validate based on provided keys
	if sebBrowserKey != "" {
		if !validateSEBRequest(fullURL, sebBrowserKey, requestHash) {
			return errors.New("unauthorized SEB request: browser exam key hash mismatch")
		}
	}

	if sebConfigKey != "" {
		if !validateSEBRequest(fullURL, sebConfigKey, configKeyHash) {
			return errors.New("unauthorized SEB request: config key hash mismatch")
		}
	}

	// Fallback to user agent check if no keys provided
	if sebBrowserKey == "" && sebConfigKey == "" {
		if !strings.Contains(userAgent, "SEB") {
			return errors.New("unauthorized SEB request: user agent mismatch")
		}
	}

	return nil
}

func validateSEBRequest(url string, key string, recvHash string) bool {

    hasher := sha256.New()

	hasher.Write([]byte(url))
	hasher.Write([]byte(key))
    
    finalHash := hasher.Sum(nil)
    hashHex := hex.EncodeToString(finalHash)

    fmt.Println("BEK/ConfigKey: Expected Hash:", hashHex)
    fmt.Println("BEK/ConfigKey: Received Hash:", recvHash)

    return hashHex == recvHash
}
