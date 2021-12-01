package tunnel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddr(t *testing.T) {
	assert := assert.New(t)

	dummyIP := Addr{
		Proto: "TCP",
		IP:    "127.0.0.1",
	}

	_ = dummyIP.Network()
	_ = dummyIP.String()

	assert.Equal(dummyIP.Proto, dummyIP.Network())
	assert.Equal(dummyIP.IP, dummyIP.String())
}
