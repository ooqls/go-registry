package registry

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type TokenConfiguration struct {
	Audience                []string  `yaml:"audience"`
	Issuer                  string  `yaml:"issuer"`
	IdGenType               string  `yaml:"id_gen_type"`
	ValidityDurationSeconds float64 `yaml:"validity_duration_seconds"`
}

func (tc *TokenConfiguration) GenerateId() string {
	defaultIdGen := uuid.NewString
	var id string
	switch tc.IdGenType {
	case "uuid":
		id = uuid.NewString()
	default:
		id = defaultIdGen()
	}

	return id
}

type TLSConfig struct {
	Enabled               bool   `yaml:"enabled"`
	CertPath              string `yaml:"cert_path"`
	KeyPath               string `yaml:"key_path"`
	CaPath                string `yaml:"ca_path"`
	InsecureSkipTLSVerify bool   `yaml:"insecure_skip_tls_verify"`
}

func (cfg *TLSConfig) TLSConfig() (*tls.Config, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	defaultConfig := &tls.Config{}
	if transport, ok := http.DefaultTransport.(*http.Transport); ok {
		if transport.TLSClientConfig != nil {
			defaultConfig = transport.TLSClientConfig
		}
	}

	cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
	if err != nil {
		return nil, err
	}
	defaultConfig.Certificates = append(defaultConfig.Certificates, cert)

	if cfg.CaPath != "" {
		caCert, err := os.ReadFile(cfg.CaPath)
		if err != nil {
			return nil, err
		}

		if defaultConfig.RootCAs == nil {
			defaultConfig.RootCAs = x509.NewCertPool()
		}
		defaultConfig.RootCAs.AppendCertsFromPEM(caCert)
	}

	defaultTLSConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: cfg.InsecureSkipTLSVerify,
	}
	return defaultTLSConfig, nil
}

type Auth struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Topics struct {
	Messages string `yaml:"messages"`
}

type Server struct {
	Name  string                 `yaml:"name"`
	Host  string                 `yaml:"host"`
	Port  int                    `yaml:"port"`
	TLS   *TLSConfig             `yaml:"tls,omitempty"`
	Auth  Auth                   `yaml:"auth"`
	Extra map[string]interface{} `yaml:"extra"`
}

type Database struct {
	Server   `yaml:",inline"`
	Database string `yaml:"database"`
}

type MessageBroker struct {
	Server `yaml:",inline"`
	Topics []string `yaml:"topics"`
}

func (s *Server) GetConnectionString() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type Registry struct {
	TokenConfiguration TokenConfiguration `yaml:"token_configuration"`
	Kafka              *MessageBroker     `yaml:"kafka,omitempty"`
	Nats               *MessageBroker     `yaml:"nats,omitempty"`

	Redis    *Database `yaml:"redis,omitempty"`
	Postgres *Database `yaml:"postgres,omitempty"`
	Mongo    *Database `yaml:"mongo,omitempty"`
}
