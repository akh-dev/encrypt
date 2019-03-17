package engine

// based on https://github.com/gtank/cryptopasta/blob/master/encrypt.go

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
)

type AESEngine struct{}

func NewAESEngine() (*AESEngine, error) {
	return &AESEngine{}, nil
}

func (*AESEngine) GenerateNewKey() (*[32]byte, error) {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a new key")
	}

	return &key, nil
}

func (*AESEngine) Encrypt(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new cipher.Block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create NewGCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new random nonce")
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (*AESEngine) Decrypt(ciphertext []byte, key *[32]byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new cipher.Block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create NewGCM")
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}
