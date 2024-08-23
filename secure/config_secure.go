package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"regexp"

	"golang.org/x/crypto/scrypt"
)

// 判断一个字符串是 `S(cryp.salt)` 格式的，并提取括号中的内容
// 其中 aaaa 是加密后的密文，bbbb 为密钥 salt
func IsEncryptedConfItem(str string) (bool, string, string) {
	// 定义正则表达式，匹配 S(base64.base64) 格式
	re, err := regexp.Compile(`^S\(([A-Za-z0-9+/=]+)\.([A-Za-z0-9+/=]+)\)$`)
	if err != nil {
		// 正则表达式编译失败，返回默认值
		return false, "", ""
	}

	// 使用正则表达式匹配字符串
	matches := re.FindStringSubmatch(str)
	if len(matches) > 0 {
		// 提取括号中的内容
		cryp := matches[1]
		salt := matches[2]
		// 返回值: 是否加密, 加密内容, salt
		return true, cryp, salt
	}

	return false, "", ""
}

// 生成指定长度的随机salt
func RandomSalt(length int) (string, error) {
	// 创建一个字节切片来存储随机字节
	salt := make([]byte, length)

	// 使用 crypto/rand 包生成随机字节
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return "", err
	}

	// 将字节切片编码为Base64字符串，方便存储和传输
	return base64.StdEncoding.EncodeToString(salt), nil
}

// DeriveKey 函数使用scrypt算法从密码和盐值中派生密钥
//
// 参数：
// password string - 密码
// salt string - 盐值
//
// 返回值：
// []byte - 派生出的密钥
// error - 如果派生密钥过程中发生错误，则返回非nil的error
func DeriveKey(password string, salt string) ([]byte, error) {
	// derived key for e.g. AES-256 (which needs a 32-byte key)
	return scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
}

// Encrypt 函数用于将明文进行AES加密，并返回加密后的密文字符串
// plainText: 待加密的明文数据，以字节切片形式传入
// key: 加密密钥，以字节切片形式传入
// 返回值：
// string: 加密后的密文字符串
// error: 如果加密过程中发生错误，则返回错误信息
func Encrypt(plainText, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(cipherText string, key []byte) ([]byte, error) {
	cipherData, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := cipherData[:aes.BlockSize]
	cipherData = cipherData[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherData, cipherData)

	return cipherData, nil
}

func EncryptConfItem(str string, password string, salt string) (string, error) {
	key, err := DeriveKey(password, salt)
	if err != nil {
		return "", err
	}

	// 加密明文
	cipherText, err := Encrypt([]byte(str), key)
	if err != nil {
		return "", err
	}

	saltText := base64.StdEncoding.EncodeToString([]byte(salt))
	// 返回加密后的密文字符串
	return "S(" + cipherText + "." + saltText + ")", nil
}

func DecryptIfEncryptedConfItem(str string, password string) string {
	isEncrypted, cryp, salt := IsEncryptedConfItem(str)
	if !isEncrypted {
		return str
	}

	saltData, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return str
	}

	key, err := DeriveKey(password, string(saltData))
	if err != nil {
		return str
	}

	plainText, err := Decrypt(cryp, key)
	if err != nil {
		return str
	}

	return string(plainText)
}
