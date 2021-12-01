package tunnel

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSSHTunnel(t *testing.T) {
	assert := assert.New(t)
	privKey := PrivateKeyFile(os.Getenv("CWD") + "/ssh_tunnel/test_files/id_rsa")
	pubKey := PublicKey(privKey)
	tun := NewSSHTunnel(
		3309,
		"gavinyap@192.168.86.205",
		"localhost:3309",
		privKey,
	)
	assert.Equal("localhost", tun.Local.Host)
	assert.Equal(3309, tun.Local.Port)
	assert.Equal("", tun.Local.User)

	assert.Equal("192.168.86.205", tun.Server.Host)
	assert.Equal(22, tun.Server.Port)
	assert.Equal("gavinyap", tun.Server.User)

	assert.Equal("localhost", tun.Remote.Host)
	assert.Equal(3309, tun.Remote.Port)
	assert.Equal("", tun.Remote.User)

	tun1 := NewSSHTunnel(
		0,
		"gavinyap@192.168.86.205:2202",
		"localhost:3310",
		privKey,
	)
	assert.Equal("localhost", tun1.Local.Host)
	assert.Equal(0, tun1.Local.Port)
	assert.Equal("", tun1.Local.User)

	assert.Equal("192.168.86.205", tun1.Server.Host)
	assert.Equal(2202, tun1.Server.Port)
	assert.Equal("gavinyap", tun1.Server.User)

	assert.Equal("localhost", tun1.Remote.Host)
	assert.Equal(3310, tun1.Remote.Port)
	assert.Equal("", tun1.Remote.User)

	// To pass callback test
	dummyIP := Addr{
		Proto: "TCP",
		IP:    "127.0.0.1",
	}
	assert.Nil(tun.Config.HostKeyCallback("192.168.86.205", dummyIP, pubKey))

	_ = assert
}

func TestNewSSHTunnelLogger(t *testing.T) {
	assert := assert.New(t)

	privKey := PrivateKeyFile(os.Getenv("CWD") + "/ssh_tunnel/test_files/id_rsa")
	tun := NewSSHTunnel(
		3309,
		"gavinyap@192.168.86.205",
		"localhost:3309",
		privKey,
	)

	output := captureStdoutOutput(func() {
		tun.logf("Hello World")
	})

	assert.Equal("Hello World", output)

	var buf bytes.Buffer

	tun.Log = log.New(&buf, "", log.Ldate|log.Lmicroseconds)

	tun.logf("Hello World")

	assert.Regexp("Hello World\n$", buf.String())
}

func captureStdoutOutput(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return fmt.Sprintf("%s", out)
}
