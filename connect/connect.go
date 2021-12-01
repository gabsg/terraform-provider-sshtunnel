package connect

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/maclarensg/terraform-provider-sshtunnel/instance_connect"
	"github.com/maclarensg/terraform-provider-sshtunnel/sshkeygen"
	"github.com/maclarensg/terraform-provider-sshtunnel/tunnel"
)

type ConnectionInfo struct {
	ListenPort   int
	Jumphost     string
	JumphostPort int
	AwsProfile   string
	AwsRegion    string
	DBHost       string
	DBPort       int
	User         string
	Tunnel       *tunnel.SSHTunnel
}

func New(lport int, jh string, jhp int, awsp string, awsr string, dbh string, dbp int, user string) *ConnectionInfo {
	con := &ConnectionInfo{
		ListenPort:   lport,
		Jumphost:     jh,
		JumphostPort: jhp,
		AwsProfile:   awsp,
		AwsRegion:    awsr,
		DBHost:       dbh,
		DBPort:       dbp,
		User:         user,
	}
	cwd, _ := os.Getwd()
	path := cwd + "/.sshkey"

	// Check if KeyPair exist or Generate one
	if err := sshkeygen.HasKeyPair(path); err != nil {
		log.Println(err)
		log.Println("Generating Key ")
		if err2 := sshkeygen.GenerateKeyPair(path); err2 != nil {
			log.Fatalf("%s", err.Error())
		}
	}

	// Read public sshkey from file
	sshkey, err := sshkeygen.ReadKeyFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Send instance connect to the host
	// Creates a instance_connect client
	instance_connect, err := instance_connect.New(con.AwsProfile, con.AwsRegion)
	if err != nil {
		log.Fatal(err)
	}
	// Send public sshkey via instance connect
	info, err := instance_connect.Send(con.Jumphost, sshkey, con.User)
	if err != nil {
		log.Fatal(err)
	}

	if info == nil {
		log.Fatal("info is empty!")
	}

	time.Sleep(time.Second * 2)

	// Lauch Tunnel
	target := fmt.Sprintf("%s@%s:%d", info.User, info.PublicDNS, con.JumphostPort)
	destination := fmt.Sprintf("%s:%d", con.DBHost, con.DBPort)

	con.Tunnel = tunnel.NewSSHTunnel(
		con.ListenPort,
		target,
		destination,
		tunnel.PrivateKeyFile(path+"/id_rsa"),
	)

	go con.Tunnel.Start()

	return con
}
