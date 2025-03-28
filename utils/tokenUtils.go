package utils

import (
    "crypto/rand"
    "encoding/base64"
)

func GenerateRandomToken() string {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if (err != nil) {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}