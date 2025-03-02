# Dorisdump

Main features:

1. **Dump** schemas, column stats and queries from Doris database
2. **Replay** dump queries with speed control
3. **Anonymize** database, table and column names in SQLs

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

# Also dump queries of db1 from audit logs
# Hint: Use '*' like '/path/to/fe.audit.log*' to match multiple files
dorisdump dump --dump-schema --dump-query --dbs db1 --audit-logs '/path/to/fe.audit.log,/path/to/fe.audit.log.20240802-1'

# Dump queries from audit log table instead of files, need enable <https://doris.apache.org/docs/admin-manual/audit-plugin>
dorisdump dump --dump-query --audit-log-table <db.table> --from '2024-11-14 18:45:25' --to '2024-11-14 18:45:26'


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

### Anonymize

> Note: This feature is experimental, **only works properly for case-insensitive names, which means `table1` and `TABLE1` will have the same result.**

Two ways:

1. Use `dorisdump anonymize`:

    ```bash
    echo "select * from table1" | dorisdump anonymize -f -
    ```

2. Use `--anonymize` flag while dumping:

    ```bash
    dorisdump dump <some flags...> --anonymize
    ```

Remember to keep `./dorisdump_hashdict.yaml` if you want the result to be consistent (default to find it at current directory, or specify by `--anonymize-minihash-dict`).

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
