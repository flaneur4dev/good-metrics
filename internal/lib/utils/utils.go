package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	GaugeName   = "gauge"
	CounterName = "counter"
	GaugeTmpl   = "%s:gauge:%f"
	CounterTmpl = "%s:counter:%d"
)

func StringEnv(name, defaultV string) (res string) {
	if val, ok := os.LookupEnv(name); ok {
		res = val
	} else {
		res = defaultV
	}
	return
}

func BoolEnv(name string, defaultV bool) (res bool) {
	if val, ok := os.LookupEnv(name); ok {
		b, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatal("Incorrect parameter: ", val)
		}
		res = b
	} else {
		res = defaultV
	}
	return
}

func Sign256(msg, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}

func IsEqualSign256(msg, hash, key string) bool {
	data, err := hex.DecodeString(hash)
	if err != nil {
		fmt.Println(err)
		return false
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(msg))
	return hmac.Equal(data, h.Sum(nil))
}
