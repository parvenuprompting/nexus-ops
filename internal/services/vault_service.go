package services

import (
	"encoding/json"
	"os"
	"sync"
)

type Secret struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VaultService struct {
	mu       sync.Mutex
	FilePath string
}

func NewVaultService() *VaultService {
	return &VaultService{
		FilePath: "nexus-vault.json",
	}
}

func (s *VaultService) LoadSecrets() ([]Secret, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.FilePath)
	if os.IsNotExist(err) {
		return []Secret{}, nil
	}
	if err != nil {
		return nil, err
	}

	var secrets []Secret
	if err := json.Unmarshal(data, &secrets); err != nil {
		return nil, err
	}
	return secrets, nil
}

func (s *VaultService) AddSecret(key, value string) error {
	secrets, err := s.LoadSecrets()
	if err != nil {
		return err
	}

	// Remove if exists (upsert)
	newSecrets := []Secret{}
	for _, secret := range secrets {
		if secret.Key != key {
			newSecrets = append(newSecrets, secret)
		}
	}
	newSecrets = append(newSecrets, Secret{Key: key, Value: value})

	return s.save(newSecrets)
}

func (s *VaultService) RemoveSecret(key string) error {
	secrets, err := s.LoadSecrets()
	if err != nil {
		return err
	}

	newSecrets := []Secret{}
	for _, secret := range secrets {
		if secret.Key != key {
			newSecrets = append(newSecrets, secret)
		}
	}

	return s.save(newSecrets)
}

func (s *VaultService) save(secrets []Secret) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.FilePath, data, 0600)
}
