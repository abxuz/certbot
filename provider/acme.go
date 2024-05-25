package provider

import (
	"crypto/sha256"
	"encoding/base64"
)

func ChallengeHost(host string) string {
	if host == "" {
		return "_acme-challenge"
	}
	return "_acme-challenge." + host
}

func ChallengeValue(key string) string {
	keyAuthShaBytes := sha256.Sum256([]byte(key))
	return base64.RawURLEncoding.EncodeToString(keyAuthShaBytes[:sha256.Size])
}
