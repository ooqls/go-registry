package registry

import (
	"fmt"

	"github.com/google/uuid"
)

type TokenConfiguration struct {
	Audience                string  `yaml:"audience"`
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
	Enabled  bool   `yaml:"enabled"`
	CertPath string `yaml:"cert_path"`
	KeyPath  string `yaml:"key_path"`
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
	TLS   TLSConfig              `yaml:"tls"`
	Auth  Auth                   `yaml:"auth"`
	Extra map[string]interface{} `yaml:"extra"`
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

	Redis    *Server `yaml:"redis,omitempty"`
	Postgres *Server `yaml:"postgres,omitempty"`
	Mongo    *Server `yaml:"mongo,omitempty"`
}
