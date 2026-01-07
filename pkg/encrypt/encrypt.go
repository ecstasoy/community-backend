package encrypt

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/zeromicro/go-zero/core/codec"
)

const (
	passwordEncryptSalt = "s0m37H1ng1N7h3w4y"
	mobileEncryptSalt   = "od9f8a7s6d5f4g3h2j1k0l9m8n7b6v5c"
)

func EncryptPassword(password string) string {
	return Md5Sum([]byte(strings.TrimSpace(password + passwordEncryptSalt)))
}

func EncryptMobile(mobile string) (string, error) {
	data, err := codec.EcbEncrypt([]byte(mobileEncryptSalt), []byte(mobile))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func DecryptMobile(mobile string) (string, error) {
	originalData, err := base64.StdEncoding.DecodeString(mobile)
	if err != nil {
		return "", err
	}
	data, err := codec.EcbDecrypt([]byte(mobileEncryptSalt), originalData)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Md5Sum calculates the MD5 checksum of the given data and returns it as a hexadecimal string
// md5.Sum returns a [16]byte array, which we convert to a byte slice for encoding
func Md5Sum(data []byte) string {
	return hex.EncodeToString(byte16ToBytes(md5.Sum(data)))
}

func byte16ToBytes(in [16]byte) []byte {
	res := make([]byte, 16)
	for i, v := range in {
		res[i] = v
	}
	return res
}
