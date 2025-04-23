# Dorisdump

Main features:

1. **Dump** schemas, column stats and queries from Doris database
2. **Replay** dump queries with speed control
3. **Anonymize** database, table and column names in SQLs

See [Introduction & FAQ](./introduction-zh.md) for more details.

[![demo](https://asciinema.org/a/706093.svg)](https://asciinema.org/a/706093)

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
# Hint: Use '*' like 'fe.audit.log*' to match multiple files
dorisdump dump --dump-schema --dump-query --dbs db1 --audit-logs 'fe.audit.log,fe.audit.log.20240802-1'

# Dump queries from audit log table instead of files, need enable <https://doris.apache.org/docs/admin-manual/audit-plugin>
dorisdump dump --dump-query --audit-log-table <db.table> --from '2024-11-14 18:45:25' --to '2024-11-14 18:45:26'


# Replay
dorisdump replay --help

# Replay queries from dump sql file
dorisdump replay --host <host> --port <port> --user root --password '******' -f output/sql/q0.sql

# Replay with args
dorisdump replay -f output/sql/q0.sql \
    --from '2024-09-20 08:00:00' --to '2024-09-20 09:00:00' \
    --users 'readonly,root' --dbs 'db1,db2' \   # filter sql by users and databases
    --speed 0.5 \                               # increase(< 1.0) or decrease(> 1.0) the time between two serial sqls proportionally, default 1
    --result-dir output/replay


# Diff replay result
dorisdump diff --help

# Print diff of replay result which is slower more than 200ms than original
dorisdump diff --min-duration-diff 200ms --original-sqls output/sql/q0.sql output/replay

# Print diff of two replay result directories
dorisdump diff replay1/ replay2/
```

### Config

You may want to pass parameters by config file or environment, see `dorisdump --help` and [example](./example/example.dorisdump.yaml).

### Anonymize

> Note: This feature is experimental, **case-insensitive, which means `table1` and `TABLE1` will have the same result.**

Two ways:

1. Use `dorisdump anonymize`:

    ```bash
    echo "select * from table1" | dorisdump anonymize -f -
    ```

2. Use `--anonymize` flag while dumping:

    ```bash
    dorisdump dump <some flags...> --anonymize
    ```

Keep `./dorisdump_hashdict.yaml` if you want the result to be consistent (put it at current directory, or specify by `--anonymize-minihash-dict`).

## Build

1. Install **optional** dependences:

    - On macOS: [vectorscan](https://github.com/VectorCamp/vectorscan) with Chimera support
    - On Linux: [hyperscan](https://intel.github.io/hyperscan) with Chimera support

2. Run `make` (or `make build-hyper` if the dependences in step 1 are installed)

## Update Doris Parser

```sh
make gen
```
