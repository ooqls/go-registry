package registry

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cfg string = `
postgres:
  name: pg
	host: pg
	port: 10
	tls:
    enabled: false
  auth:
    enabled: false
kafka:
  "group.id": aaa
token_config:
	audience: test
	issuer: test
	validity_duration_seconds: 100
`

func TestInit(t *testing.T) {
	f, err := ioutil.TempFile("/tmp/", "registry-test")
	if err != nil {
		log.Print(err.Error())
		t.FailNow()
	}

	_, err = f.Write([]byte(strings.ReplaceAll(cfg, "\t", "  ")))
	if err != nil {
		log.Print(err.Error())
		t.FailNow()
	}

	err = Init(f.Name())
	if err != nil {
		t.Error(err)
	}

	assert.Len(t, Get().Kafka, 1)

}
