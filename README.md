# Dorisdump

Dump schemas and queries from Doris database with anonymization.

[![demo](https://asciinema.org/a/6MIhuruC668RvElND8RiMFnH9.svg)](https://asciinema.org/a/6MIhuruC668RvElND8RiMFnH9)

## Install

```sh
curl -sSL https://raw.githubusercontent.com/Thearas/dorisdump/master/install.sh | bash
```

## Usage

```sh
dorisdump --help

# Dump schemas of database db1 and db2
dorisdump dump --host <host> --port <port> --user root --password '******' --dbs db1,db2 --dump-schema

# Also dump queries from db1, queries will be extracted from audit logs
# Hint: Use '*' like '/path/to/fe.audit.log*' to match multiple logs
dorisdump dump --dbs db1 --dump-schema --dump-query --audit-logs '/path/to/fe.audit.log,/path/to/fe.audit.log.20240802-1'

# Dump with anonymization
dorisdump dump --dbs db1 --dump-schema --dump-query --audit-logs '/path/to/fe.audit.log' --anonymize

# Auto download audit log from remote (require SSH password or private key)
dorisdump dump --dbs db1 --dump-schema --dump-query --anonymize --ssh-password '******'
```

## Build

1. Install **optional** dependences:

    - On macOS: [vectorscan](https://github.com/VectorCamp/vectorscan) with Chimera support
    - On Linux: [hyperscan](https://intel.github.io/hyperscan) with Chimera support

2. Run `make` (or `make build-hyper` if the dependences in step 1 are installed)

## Update Doris Parser

```sh
make gen
```
