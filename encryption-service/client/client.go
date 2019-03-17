package client

import (
	"github.com/akh-dev/encrypt/encryption-service/config"
	"github.com/akh-dev/encrypt/encryption-service/engine"
	"github.com/akh-dev/encrypt/encryption-service/service"

	"github.com/pkg/errors"
)

type EncryptionClient struct {
	service *service.Service
}

func New() (*EncryptionClient, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	encryptionEngine, err := engine.NewAESEngine()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise encryption service")
	}

	s, err := service.New(cfg, encryptionEngine)

	client := &EncryptionClient{service: s}

	return client, nil
}

func (c *EncryptionClient) Store(id, payload []byte) (aesKey []byte, err error) {
	return c.service.ProcessStore(id, payload)
}

func (c *EncryptionClient) Retrieve(id, aesKey []byte) (payload []byte, err error) {
	return c.service.ProcessRetrieve(id, payload)
}
