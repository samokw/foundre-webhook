package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func computrHMACSHA256Hex(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyGitHubSignature(secret string, body []byte, sigHeader string) bool {
	if secret == "" {
		return false
	}
	if !strings.HasPrefix(sigHeader, "sha256=") {
		return false
	}
	gotHex := strings.TrimPrefix(sigHeader, "sha256=")

	wantHex := computrHMACSHA256Hex(secret, body)

	got, err := hex.DecodeString(gotHex)

	if err != nil {
		return false
	}

	want, _ := hex.DecodeString(wantHex)
	return hmac.Equal(got, want)
}
