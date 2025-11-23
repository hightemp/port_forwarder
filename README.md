# Port Forwarder

A simple Go tool for SSH port forwarding. Supports both Local (`-L`) and Remote (`-R`) port forwarding.

## Features

*   Configuration via YAML file.
*   Support for multiple SSH servers.
*   Authentication via key or password.
*   Local Forwarding (Local port -> Remote service).
*   Remote Forwarding (Remote port -> Local service).

## Installation and Usage

### Requirements

*   Go 1.21+
*   Make (optional)

### Build

```bash
make build
```

### Run

1.  Create a configuration file `config.yaml` (see example below).
2.  Run the application:

```bash
./port_forwarder -config config.yaml
```

Or use `make`:

```bash
make run
```

## Configuration

Example `config.yaml`:

```yaml
ssh_servers:
  - name: "production"
    host: "example.com"
    port: "22"
    user: "admin"
    key_file: "/home/user/.ssh/id_rsa"

tunnels:
  # Local Forwarding: Listen on localhost:8080, forward to localhost:80 on remote server
  - server_name: "production"
    local_addr: "localhost:8080"
    remote_addr: "localhost:80"
    mode: "local"

  # Remote Forwarding: Listen on 0.0.0.0:9090 on remote server, forward to localhost:3000 locally
  - server_name: "production"
    local_addr: "localhost:3000"
    remote_addr: "0.0.0.0:9090"
    mode: "remote"
```

## Makefile Commands

*   `make build`: Build the binary.
*   `make run`: Run with default config.
*   `make clean`: Clean up binaries.
*   `make help`: Show help.