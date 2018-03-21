package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Encode(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))
	return hex.EncodeToString(hash.Sum(nil))
}
