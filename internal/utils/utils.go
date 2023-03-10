package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

func MD5[T string | []byte](v T) string {
	b := md5.Sum([]byte(v))
	return hex.EncodeToString(b[:])
}

func ConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	configDir = filepath.Join(configDir, "toybox")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0755); err != nil {
			panic(err)
		}
	}
	return configDir
}
