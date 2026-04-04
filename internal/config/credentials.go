package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/pbkdf2"
)

const (
	keyringService = "adguard-cli"
	encryptedFile  = "credentials.enc"
	pbkdf2Iter     = 100000
	keyLen         = 32
)

// CredentialStore abstracts credential storage.
type CredentialStore interface {
	Get(instance string) (string, error)
	Set(instance, password string) error
	Delete(instance string) error
}

// NewCredentialStore returns the best available credential store.
// Tries system keyring first, falls back to encrypted file.
func NewCredentialStore() CredentialStore {
	if keyringAvailable() {
		return &keyringStore{}
	}
	return &encryptedFileStore{}
}

func keyringAvailable() bool {
	testKey := "adguard-cli-keyring-test"
	err := keyring.Set(keyringService, testKey, "test")
	if err != nil {
		return false
	}
	_ = keyring.Delete(keyringService, testKey)
	return true
}

// --- Keyring Store ---

type keyringStore struct{}

func (k *keyringStore) Get(instance string) (string, error) {
	pass, err := keyring.Get(keyringService, instance)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}
		return "", fmt.Errorf("keyring get: %w", err)
	}
	return pass, nil
}

func (k *keyringStore) Set(instance, password string) error {
	return keyring.Set(keyringService, instance, password)
}

func (k *keyringStore) Delete(instance string) error {
	return keyring.Delete(keyringService, instance)
}

// --- Encrypted File Store (fallback for headless servers) ---

type encryptedFileStore struct{}

func encFilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, encryptedFile), nil
}

func deriveKey(instance string) []byte {
	hostname, _ := os.Hostname()
	salt := []byte(fmt.Sprintf("adguard-cli:%s:%s", hostname, instance))
	return pbkdf2.Key(salt, salt, pbkdf2Iter, keyLen, sha256.New)
}

func (e *encryptedFileStore) Get(instance string) (string, error) {
	path, err := encFilePath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path + "." + instance)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return "", fmt.Errorf("decoding credentials: %w", err)
	}

	key := deriveKey(instance)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypting credentials: %w", err)
	}

	return string(plaintext), nil
}

func (e *encryptedFileStore) Set(instance, password string) error {
	path, err := encFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	key := deriveKey(instance)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return os.WriteFile(path+"."+instance, []byte(encoded), 0600)
}

func (e *encryptedFileStore) Delete(instance string) error {
	path, err := encFilePath()
	if err != nil {
		return err
	}
	return os.Remove(path + "." + instance)
}
