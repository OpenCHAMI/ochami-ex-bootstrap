# ochami-ex-bootstrap

A small CLI tool to discover bootable node NICs via BMC Redfish and produce a YAML inventory (`bmcs[]` and `nodes[]`).

This repository contains a lightweight Go implementation originally based on small scripts for creating an initial BMC inventory and then discovering node NICs through Redfish.

## Features

- Generate an initial `inventory.yaml` with a `bmcs` list (xname, MAC, IP) using `--init-bmcs`.
- Discover bootable NICs via Redfish on each BMC and allocate IPs from a given subnet.
- Output file format: a single YAML file with two top-level keys:
  - `bmcs`: list of management controllers (xname, mac, ip)
  - `nodes`: list of discovered node network records (xname, mac, ip)

## Layout

- `main.go` — minimal entrypoint that invokes the Cobra CLI.
- `cmd/` — Cobra commands:
  - `init-bmcs` — generate initial inventory with BMC entries
  - `discover` — discover bootable NICs via Redfish and update nodes[]
- `internal/` — code split by concern:
  - `inventory/` — YAML types (`Entry`, `FileFormat`)
  - `redfish/` — minimal Redfish client and bootable NIC heuristics
  - `netalloc/` — IP allocation using `github.com/metal-stack/go-ipam`
  - `xname/` — xname helpers and conversions
  - `initbmcs/` — helpers used by the `init-bmcs` command
  - `discover/` — discovery orchestration (Redfish + IP allocation)
- `examples/` — sample files (e.g., `inventory.yaml`).

## Build

This project uses Go modules. From the repo root:

```bash
# fetch modules and build
go mod tidy
go build -o ochami_bootstrap .
```

## Usage 

Show help:

```bash
./ochami_bootstrap --help
```

### 1) Generate an initial BMC inventory

```bash
./ochami_bootstrap init-bmcs --file examples/inventory.yaml \
  --chassis "x9000c1=02:23:28:01,x9000c3=02:23:28:03" \
  --bmc-subnet 192.168.100 \
  --nodes-per-chassis 32 \
  --nodes-per-bmc 2 \
  --start-nid 1
```

Writes `examples/inventory.yaml` with a `bmcs:` list and `nodes: []`.

### 2) Discover bootable NICs and allocate IPs

The discovery flow reads the YAML `--file` (must contain non-empty `bmcs[]`) and writes back the same file with updated `nodes[]`.

Required env vars:
- `REDFISH_USER` — Redfish username
- `REDFISH_PASSWORD` — Redfish password

Example:

```bash
export REDFISH_USER=admin
export REDFISH_PASSWORD=secret
./ochami_bootstrap discover \
  --file examples/inventory.yaml \
  --subnet 10.42.0.0/24 \
  --timeout 12s \
  --insecure \
  --ssh-pubkey ~/.ssh/id_rsa.pub   # optional: set AuthorizedKeys on each BMC
```

Notes:
- The program makes simple heuristic decisions about which NIC is bootable (UEFI path hints, DHCP addresses, or a MAC on an enabled interface).
- IP allocation is done with `github.com/metal-stack/go-ipam`. The code reserves `.1` (first host) as a gateway and avoids network/broadcast implicitly.
- If `--ssh-pubkey` is provided, the tool attempts a Redfish PATCH to `/redfish/v1/Managers/BMC/NetworkProtocol` with an OEM payload setting `SSHAdmin.AuthorizedKeys` to the contents of the file.

## Dependencies

- Go (module aware). The project will download dependencies with `go mod tidy`.
- `github.com/metal-stack/go-ipam` — used for IP allocation.
- `gopkg.in/yaml.v3` — YAML parsing and writing.

## Contributing / Next steps

- Add unit tests for the xname / MAC generation helpers and the Redfish parsing heuristics.
- Add input validation for chassis/macro formats if you require stricter MAC formatting.
- Consider adding a `--dry-run` mode for discovery to avoid writing changes while testing.

## License

Pick an appropriate license for your project. This repo currently has none specified.

If you'd like, I can also add a small `README` section that includes an example `inventory.yaml` and a quick test script to validate the discovery flow without real BMCs (e.g., a mock HTTP server).