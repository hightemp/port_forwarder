package forwarder

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"port_forwarder/internal/config"

	"golang.org/x/crypto/ssh"
)

type Forwarder struct {
	Config  *config.Config
	Clients map[string]*ssh.Client
	wg      sync.WaitGroup
}

func NewForwarder(cfg *config.Config, clients map[string]*ssh.Client) *Forwarder {
	return &Forwarder{
		Config:  cfg,
		Clients: clients,
	}
}

func (f *Forwarder) Start() {
	for _, tunnel := range f.Config.Tunnels {
		client, ok := f.Clients[tunnel.ServerName]
		if !ok {
			log.Printf("Server %s not found for tunnel %s -> %s", tunnel.ServerName, tunnel.LocalAddr, tunnel.RemoteAddr)
			continue
		}

		f.wg.Add(1)
		go func(t config.Tunnel, c *ssh.Client) {
			defer f.wg.Done()
			if err := f.startTunnel(t, c); err != nil {
				log.Printf("Tunnel error (%s): %v", t.Mode, err)
			}
		}(tunnel, client)
	}
	f.wg.Wait()
}

func (f *Forwarder) startTunnel(t config.Tunnel, client *ssh.Client) error {
	switch t.Mode {
	case "local", "L":
		return f.startLocalForwarding(t, client)
	case "remote", "R":
		return f.startRemoteForwarding(t, client)
	default:
		return fmt.Errorf("unknown mode: %s", t.Mode)
	}
}

func (f *Forwarder) startLocalForwarding(t config.Tunnel, client *ssh.Client) error {
	listener, err := net.Listen("tcp", t.LocalAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on local address %s: %w", t.LocalAddr, err)
	}
	defer listener.Close()

	log.Printf("Listening on local %s, forwarding to remote %s via %s", t.LocalAddr, t.RemoteAddr, t.ServerName)

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept local connection: %v", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			remoteConn, err := client.Dial("tcp", t.RemoteAddr)
			if err != nil {
				log.Printf("Failed to dial remote address %s: %v", t.RemoteAddr, err)
				return
			}
			defer remoteConn.Close()

			copyConn(conn, remoteConn)
		}(localConn)
	}
}

func (f *Forwarder) startRemoteForwarding(t config.Tunnel, client *ssh.Client) error {
	listener, err := client.Listen("tcp", t.RemoteAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on remote address %s: %w", t.RemoteAddr, err)
	}
	defer listener.Close()

	log.Printf("Listening on remote %s, forwarding to local %s via %s", t.RemoteAddr, t.LocalAddr, t.ServerName)

	for {
		remoteConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept remote connection: %v", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			localConn, err := net.Dial("tcp", t.LocalAddr)
			if err != nil {
				log.Printf("Failed to dial local address %s: %v", t.LocalAddr, err)
				return
			}
			defer localConn.Close()

			copyConn(conn, localConn)
		}(remoteConn)
	}
}

func copyConn(local, remote net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(local, remote)
		local.Close() // Close local to signal we are done writing to it
	}()

	go func() {
		defer wg.Done()
		io.Copy(remote, local)
		remote.Close() // Close remote to signal we are done writing to it
	}()

	wg.Wait()
}
