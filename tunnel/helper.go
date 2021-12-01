package tunnel

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func PrivateKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func PublicKey(auth ssh.AuthMethod) ssh.PublicKey {
	key, _ := ssh.NewPublicKey(auth)
	return key
}
