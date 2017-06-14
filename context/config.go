package context

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

const (
	configIndent = "\t"
)

// Config holds service configuration inforamtion.
type Config struct {
	// ServiceName represents the service's identifier
	ServiceName string

	// DiscoveryEndpoint specifies the Discovery service endpoint to use.
	DiscoveryEndpoint string

	// WebServerAddr specifies the address to bind to for the webserver.
	WebServerAddr string

	// StreamServerAddr specifies the address to bind to for the streaming server.
	StreamServerAddr string

	// PrivateKeyFile specifies a path to the global private key.
	PrivateKeyFile string
}

// LoadConfig loads a configuration file.
func LoadConfig(filename string) (*Config, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	var c Config
	err = json.NewDecoder(fd).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, err
}

// Save saves the configuration to a file.
func (c *Config) Save(filename string) error {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	enc := json.NewEncoder(fd)
	enc.SetIndent("", configIndent)
	return enc.Encode(c)
}

// DefaultConfig returns a Config template.
func DefaultConfig() *Config {
	var c Config
	c.ServiceName = uuid.New().String()
	c.DiscoveryEndpoint = "ws://example.com/discovery"
	c.WebServerAddr = ":13000"
	c.StreamServerAddr = ":14000"
	c.PrivateKeyFile = "privateKey.pem"
	return &c
}
