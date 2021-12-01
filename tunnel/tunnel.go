package tunnel

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	Local    *Endpoint
	Server   *Endpoint
	Remote   *Endpoint
	Config   *ssh.ClientConfig
	Log      *log.Logger
	Listener string
}

func NewSSHTunnel(localPort int, tunnel string, destination string, auth ssh.AuthMethod /*ssh key*/) *SSHTunnel {
	localEndpoint := &Endpoint{
		Host: "localhost",
		Port: localPort, //if localPort is 0 a random Port will be generated
	}

	server := NewEndpoint(tunnel)

	sshTunnel := &SSHTunnel{
		Config: &ssh.ClientConfig{
			User: server.User,
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				// Always accept key.
				return nil
			},
		},
		Local:  localEndpoint,
		Server: server,
		Remote: NewEndpoint(destination),
	}
	return sshTunnel
}

func (tunnel *SSHTunnel) logf(format string, args ...interface{}) {
	if tunnel.Log != nil { // Prints with Logger
		tunnel.Log.Printf(format, args...)
	} else { // Prints to stdout
		fmt.Printf(format, args...)
	}
}

func (tunnel *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	tunnel.Listener = listener.Addr().String()
	tunnel.logf("Listening: %s ", tunnel.Listener)
	if err != nil {
		return err
	}
	defer listener.Close()
	tunnel.Local.Port = listener.Addr().(*net.TCPAddr).Port
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		tunnel.logf("accepted connection")
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		tunnel.logf("server dial error: %s", err)
		return
	}
	tunnel.logf("connected to %s (1 of 2)\n", tunnel.Server.String())
	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		tunnel.logf("remote dial error: %s", err)
		return
	}
	tunnel.logf("connected to %s (2 of 2)\n", tunnel.Remote.String())
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			tunnel.logf("io.Copy error: %s", err)
		}
	}
	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func (tunnel *SSHTunnel) GetListenerPort() int {
	parts := strings.Split(tunnel.Listener, ":")

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}
	return port
}
