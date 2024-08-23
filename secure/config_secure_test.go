package secure

import (
	"testing"
)

func TestIsEncrypted(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		cryp     string
		salt     string
	}{
		{
			name:     "Valid encrypted string",
			input:    "S(content.key)",
			expected: true,
			cryp:     "content",
			salt:     "key",
		},
		{
			name:     "Valid encrypted string ",
			input:    "S(aczACZ089+/=.KJLBHJH+J/ttg765+=/)",
			expected: true,
			cryp:     "aczACZ089+/=",
			salt:     "KJLBHJH+J/ttg765+=/",
		},
		{
			name:     "Invalid encrypted string",
			input:    "S(content.key.extra)",
			expected: false,
			cryp:     "",
			salt:     "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
			cryp:     "",
			salt:     "",
		},
		{
			name:     "No match string",
			input:    "NoMatchString",
			expected: false,
			cryp:     "",
			salt:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, cryp, salt := IsEncryptedConfItem(tt.input)
			if got != tt.expected {
				t.Errorf("IsEncrypted() got = %v, expected %v", got, tt.expected)
			}
			if cryp != tt.cryp {
				t.Errorf("IsEncrypted() cryp = %v, expected %v", cryp, tt.cryp)
			}
			if salt != tt.salt {
				t.Errorf("IsEncrypted() salt = %v, expected %v", salt, tt.salt)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {

	key, _ := DeriveKey("password", "salt")

	cipherText, _ := Encrypt([]byte("content"), key)

	plainText, _ := Decrypt(cipherText, key)

	if string(plainText) != "content" {
		t.Errorf("Decrypted text is not the same as original content")
	}
}

func TestDecrypt(t *testing.T) {

	key, _ := DeriveKey("password", "salt")
	plainText, _ := Decrypt("kbNhTOgrAzku65UnmYtGGOSDA2Zl/Iw=", key)

	if string(plainText) != "content" {
		t.Errorf("Decrypted text is not the same as original content")
	}
}

func TestConfItem(t *testing.T) {

	salt, _ := RandomSalt(8)
	cipherText, _ := EncryptConfItem("abc#$%^&*d", "pass", salt)
	print(cipherText)

	plainText := DecryptIfEncryptedConfItem(cipherText, "pass")
	print(plainText)
}
