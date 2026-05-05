package keychain

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

const servicePrefix = "envchain"

// ErrNotFound is returned when a secret is not found in the keychain.
var ErrNotFound = errors.New("secret not found in keychain")

// Client provides access to the OS keychain for storing and retrieving secrets.
type Client struct {
	project string
}

// New creates a new keychain Client scoped to the given project.
func New(project string) *Client {
	return &Client{project: project}
}

func (c *Client) serviceName() string {
	return fmt.Sprintf("%s/%s", servicePrefix, c.project)
}

// Set stores a secret in the OS keychain under the given key.
func (c *Client) Set(key, value string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	if err := keyring.Set(c.serviceName(), key, value); err != nil {
		return fmt.Errorf("keychain set %q: %w", key, err)
	}
	return nil
}

// Get retrieves a secret from the OS keychain by key.
// Returns ErrNotFound if the key does not exist.
func (c *Client) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("key must not be empty")
	}
	val, err := keyring.Get(c.serviceName(), key)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("keychain get %q: %w", key, err)
	}
	return val, nil
}

// Delete removes a secret from the OS keychain by key.
func (c *Client) Delete(key string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	if err := keyring.Delete(c.serviceName(), key); err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("keychain delete %q: %w", key, err)
	}
	return nil
}
