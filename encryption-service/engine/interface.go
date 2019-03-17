package engine

type Interface interface {
	GenerateNewKey() (*[32]byte, error)
	Encrypt(plaintext []byte, key *[32]byte) ([]byte, error)
	Decrypt(ciphertext []byte, key *[32]byte) (plaintext []byte, err error)
}
