package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_SendURL(t *testing.T) {
	t.Skip("need network")
	var cfg ClientInfo
	cfg.Certificate.CA = "../../example/native/var/db/oasis/localhub-cert-only-for-test/ca.pem"
	cli, err := NewClient(cfg)
	assert.NoError(t, err)
	url := "https://www.baidu.com/"
	res, err := cli.SendURL("GET", url, nil, nil)
	assert.NoError(t, err)
	defer res.Close()
}
