package ssh

import (
	"fmt"
	"net"
	"os"
	"time"

	"port_forwarder/internal/config"

	"golang.org/x/crypto/ssh"
)

func NewClient(server config.SSHServer) (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	if server.KeyFile != "" {
		key, err := os.ReadFile(server.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if server.Password != "" {
		authMethods = append(authMethods, ssh.Password(server.Password))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method provided for server %s", server.Name)
	}

	sshConfig := &ssh.ClientConfig{
		User:            server.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Implement proper host key verification
		Timeout:         5 * time.Second,
	}

	addr := net.JoinHostPort(server.Host, server.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return client, nil
}
