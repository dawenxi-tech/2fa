package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/denisbrodbeck/machineid"
)

const bundleId = "tech.dawenxi.2fa"

var uniqueId = func() string {
	id, err := machineid.ProtectedID(bundleId)
	if err != nil {
		return bundleId
	}
	return id
}()

var encryptKey = func() []byte {
	sum := md5.Sum([]byte(uniqueId))
	return sum[:]
}()

func aesEncrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	in := Pad([]byte(plaintext), aes.BlockSize)
	cipherText := make([]byte, len(in))
	sum := md5.Sum(key)
	iv := sum[:]
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText, in)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func aesDecrypt(key []byte, ct string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	sum := md5.Sum(key)
	iv := sum[:]
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(ciphertext, ciphertext)
	ciphertext, err = Unpad(ciphertext)
	if err != nil {
		return "", err
	}
	return string(ciphertext), nil
}

var errPKCS7Padding = errors.New("pkcs7pad: bad padding")

func Pad(buf []byte, size int) []byte {
	if size < 1 || size > 255 {
		panic(fmt.Sprintf("pkcs7pad: inappropriate block size %d", size))
	}
	i := size - (len(buf) % size)
	return append(buf, bytes.Repeat([]byte{byte(i)}, i)...)
}

func Unpad(buf []byte) ([]byte, error) {
	if len(buf) == 0 {
		return nil, errPKCS7Padding
	}
	padLen := buf[len(buf)-1]
	toCheck := 255
	good := 1
	if toCheck > len(buf) {
		toCheck = len(buf)
	}
	for i := 0; i < toCheck; i++ {
		b := buf[len(buf)-1-i]

		outOfRange := subtle.ConstantTimeLessOrEq(int(padLen), i)
		equal := subtle.ConstantTimeByteEq(padLen, b)
		good &= subtle.ConstantTimeSelect(outOfRange, 1, equal)
	}

	good &= subtle.ConstantTimeLessOrEq(1, int(padLen))
	good &= subtle.ConstantTimeLessOrEq(int(padLen), len(buf))

	if good != 1 {
		return nil, errPKCS7Padding
	}

	return buf[:len(buf)-int(padLen)], nil
}

type EncryptData string

func NewEncryptData(val string) *EncryptData {
	ed := EncryptData(val)
	return &ed
}

func (e *EncryptData) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	data, err := aesDecrypt(encryptKey, string(text))
	if err != nil {
		return err
	}
	*e = EncryptData(data)
	return nil
}

func (e *EncryptData) MarshalText() (text []byte, err error) {
	if e == nil {
		return nil, nil
	}
	data, err := aesEncrypt(encryptKey, string(*e))
	if err != nil {
		return
	}
	return []byte(data), err
}

func (e *EncryptData) Val() string {
	if e == nil {
		return ""
	}
	return string(*e)
}
