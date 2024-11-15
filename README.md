# Dorisdump

Dump and replay queries from Doris database.

[![demo](https://asciinema.org/a/6MIhuruC668RvElND8RiMFnH9.svg)](https://asciinema.org/a/6MIhuruC668RvElND8RiMFnH9)

## Install

```sh
curl -sSL https://raw.githubusercontent.com/Thearas/dorisdump/master/install.sh | bash
```

## Usage

By default, only `SELECT` statments will be dumped. Use `--only-select=false` to dump all.

```sh
# Dump
dorisdump dump --help

# Dump schemas of database db1 and db2
dorisdump dump --dump-schema --host <host> --port <port> --user root --password '******' --dbs db1,db2

# Also dump queries from db1, queries will be extracted from audit logs
# Hint: Use '*' like '/path/to/fe.audit.log*' to match multiple logs
dorisdump dump --dump-schema --dump-query --dbs db1 --audit-logs '/path/to/fe.audit.log,/path/to/fe.audit.log.20240802-1'

# Auto download audit log from remote FE (require SSH password or private key)
dorisdump dump --dump-query --host <fe host> --port <fe port> --ssh-password '******'

# Dump queries from audit log table instead of files
# Need to enable audit plugin on FE, see <https://doris.apache.org/docs/admin-manual/audit-plugin>
dorisdump dump --dump-query --host <fe host> --port <fe port> --audit-log-table=audit_db.audit_table


# Replay
dorisdump replay --help

# Replay queries from dump sql file
dorisdump replay --host <host> --port <port> --user root --password '******' -f /path/to/dump.sql

# Replay with args
dorisdump replay -f /path/to/dump.sql \
    --from '2024-09-20 08:00:00' --to '2024-09-20 09:00:00' \ # from time to time
    --users 'readonly,root' --dbs 'db1,db2' \                 # filter sql by users and databases
    --count 100 \                                             # max replay sql count
    --speed 0.5 \                                             # replay speed
    --result-dir replay1


# Diff replay result
dorisdump diff --help

# Print diff of two replay result directories
dorisdump diff replay1 replay2
```

### Config

You may want to pass parameters by config file or environment, see `dorisdump --help` and [example](./example/example.dorisdump.yaml).

## Build

1. Install **optional** dependences:

    - On macOS: [vectorscan](https://github.com/VectorCamp/vectorscan) with Chimera support
    - On Linux: [hyperscan](https://intel.github.io/hyperscan) with Chimera support

2. Run `make` (or `make build-hyper` if the dependences in step 1 are installed)

## Update Doris Parser

```sh
make gen
```
