package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// func Md5(str string) string {
// 	h := md5.New()
// 	h.Write([]byte(str))
// 	return hex.EncodeToString(h.Sum(nil))
// }

func Sha1Sign(s string) string {
	// The pattern for generating a hash is `sha1.New()`,
	// `sha1.Write(bytes)`, then `sha1.Sum([]byte{})`.
	// Here we start with a new hash.
	h := sha1.New()

	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	_, _ = h.Write([]byte(s))

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	bs := h.Sum(nil)

	// SHA1 values are often printed in hex, for example
	// in git commits. Use the `%x` format verb to convert
	// a hash results to a hex string.
	return fmt.Sprintf("%x", bs)
}

// AesDecrypt AES-CBC解密,PKCS#7,传入密文和密钥，[]byte
func AesDecrypt(src, key []byte) (dst []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	dst = make([]byte, len(src))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(dst, src)

	return PKCS7UnPad(dst), nil
}

// PKCS7UnPad PKSC#7解包
func PKCS7UnPad(msg []byte) []byte {
	length := len(msg)
	padlen := int(msg[length-1])
	return msg[:length-padlen]
}

// AesEncrypt AES-CBC加密+PKCS#7打包，传入明文和密钥
func AesEncrypt(src []byte, key []byte) ([]byte, error) {
	k := len(key)
	if len(src)%k != 0 {
		src = PKCS7Pad(src, k)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	dst := make([]byte, len(src))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(dst, src)

	return dst, nil
}

// PKCS7Pad PKCS#7打包
func PKCS7Pad(msg []byte, blockSize int) []byte {
	if blockSize < 1<<1 || blockSize >= 1<<8 {
		panic("unsupported block size")
	}
	padlen := blockSize - len(msg)%blockSize
	padding := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(msg, padding...)
}

// SortSha1 排序并sha1，主要用于计算signature
func SortSha1(s ...string) string {
	sort.Strings(s)
	h := sha1.New()
	h.Write([]byte(strings.Join(s, "")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SortMd5 排序并md5，主要用于计算sign
func SortMd5(s ...string) string {
	sort.Strings(s)
	h := md5.New()
	h.Write([]byte(strings.Join(s, "")))
	return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
}
