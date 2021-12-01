package tunnel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEndpoint(t *testing.T) {
	assert := assert.New(t)

	// Test endpoint with user specify but without port
	// defaults port to 22
	ep1 := NewEndpoint("gavinyap@192.168.86.205")
	assert.Equal("gavinyap", ep1.User)
	assert.Equal("192.168.86.205", ep1.Host)
	assert.Equal(22, ep1.Port)

	// Test endpoint without user but with port
	ep2 := NewEndpoint("localhost:3309")
	assert.Equal("", ep2.User)
	assert.Equal("localhost", ep2.Host)
	assert.Equal(3309, ep2.Port)

	// Test endpoint without user and port
	// defaults port to 22
	ep3 := NewEndpoint("localhost")
	assert.Equal("", ep3.User)
	assert.Equal("localhost", ep3.Host)
	assert.Equal(22, ep3.Port)
}

func TestNewEndpointStringFunc(t *testing.T) {
	assert := assert.New(t)

	ep := NewEndpoint("gavinyap@192.168.86.205")
	assert.Equal("192.168.86.205:22", ep.String())
}
