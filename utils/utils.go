package utils

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// Magic string to verify correct decryption
	encryptionMarker = "AL_ENCRYPTED_NOTE"
)

// LevenshteinDistance calculates the Levenshtein distance between two strings
func LevenshteinDistance(a, b string) int {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}

	return matrix[len(a)][len(b)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// FindSimilarStrings finds strings similar to the target
func FindSimilarStrings(target string, candidates []string, maxDistance int) []string {
	var similar []string
	for _, candidate := range candidates {
		distance := LevenshteinDistance(target, candidate)
		if distance <= maxDistance {
			similar = append(similar, candidate)
		}
	}
	return similar
}

// ReadPassword reads a password from stdin without echoing
func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// Encrypt encrypts plaintext with the given password
func Encrypt(plaintext, password string) (string, error) {
	// Add marker to verify decryption
	data := []byte(encryptionMarker + "\n" + plaintext)

	// Derive key from password
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Combine salt and ciphertext
	result := append(salt, ciphertext...)

	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts ciphertext with the given password
func Decrypt(ciphertext, password string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(data) < 32 {
		return "", fmt.Errorf("invalid ciphertext")
	}

	// Extract salt
	salt := data[:32]
	cipherData := data[32:]

	// Derive key from password
	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(cipherData) < gcm.NonceSize() {
		return "", fmt.Errorf("invalid ciphertext")
	}

	nonce := cipherData[:gcm.NonceSize()]
	cipherData = cipherData[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", fmt.Errorf("incorrect password or corrupted data")
	}

	// Verify marker
	text := string(plaintext)
	if !strings.HasPrefix(text, encryptionMarker+"\n") {
		return "", fmt.Errorf("incorrect password")
	}

	return strings.TrimPrefix(text, encryptionMarker+"\n"), nil
}

// OpenEditor opens the default editor (vim) with the given file
func OpenEditor(filepath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CopyToClipboard copies text to the system clipboard
func CopyToClipboard(text string) error {
	return clipboard.WriteAll(text)
}

// AskConfirmation asks the user for yes/no confirmation
func AskConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", prompt)
	
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// TruncateString truncates a string to the given length
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}
