package engine

type Interface interface {
	GenerateNewKey() (*[32]byte, error)
	Encrypt(plaintext []byte, key *[32]byte) ([]byte, error)
}
