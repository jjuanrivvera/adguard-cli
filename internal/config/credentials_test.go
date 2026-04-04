package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptedFileStore_SetAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	store := &encryptedFileStore{}

	// Store a password
	err := store.Set("test-instance", "my-secret-password")
	require.NoError(t, err)

	// Retrieve it
	password, err := store.Get("test-instance")
	require.NoError(t, err)
	assert.Equal(t, "my-secret-password", password)
}

func TestEncryptedFileStore_GetNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	store := &encryptedFileStore{}

	password, err := store.Get("nonexistent")
	assert.NoError(t, err)
	assert.Empty(t, password)
}

func TestEncryptedFileStore_OverwritePassword(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	store := &encryptedFileStore{}

	err := store.Set("instance1", "password-v1")
	require.NoError(t, err)

	err = store.Set("instance1", "password-v2")
	require.NoError(t, err)

	password, err := store.Get("instance1")
	require.NoError(t, err)
	assert.Equal(t, "password-v2", password)
}

func TestEncryptedFileStore_MultipleInstances(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	store := &encryptedFileStore{}

	err := store.Set("primary", "pass1")
	require.NoError(t, err)
	err = store.Set("secondary", "pass2")
	require.NoError(t, err)

	p1, err := store.Get("primary")
	require.NoError(t, err)
	assert.Equal(t, "pass1", p1)

	p2, err := store.Get("secondary")
	require.NoError(t, err)
	assert.Equal(t, "pass2", p2)
}

func TestEncryptedFileStore_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	store := &encryptedFileStore{}

	err := store.Set("deleteme", "password")
	require.NoError(t, err)

	err = store.Delete("deleteme")
	require.NoError(t, err)

	password, err := store.Get("deleteme")
	assert.NoError(t, err)
	assert.Empty(t, password)
}

func TestEncryptedFileStore_SpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	store := &encryptedFileStore{}

	// Passwords with special chars, unicode, long strings
	passwords := []string{
		`p@$$w0rd!#%^&*()`,
		`contraseña con espacios y ñ`,
		`{"json": "injection"}`,
		`a very long password that exceeds typical lengths to ensure our encryption handles arbitrary input sizes correctly without truncation or corruption`,
	}

	for i, pw := range passwords {
		instance := "inst" + string(rune('0'+i))
		err := store.Set(instance, pw)
		require.NoError(t, err)

		got, err := store.Get(instance)
		require.NoError(t, err)
		assert.Equal(t, pw, got)
	}
}

func TestNewCredentialStore_Fallback(t *testing.T) {
	// On a headless server without a keyring, this should return the encrypted file store
	store := NewCredentialStore()
	assert.NotNil(t, store)
	// We can't assert the type externally, but it should not panic
}
