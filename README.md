# Dorisdump (Support Doris >= 2.0.14)

Main features:

1. **Dump** schema and query
2. **Generate and import** data for table
3. **Replay** dump query
4. **Anonymize** database, table and column name in SQL

> [!IMPORTANT]
> **See [Introduction & FAQ](./introduction-zh.md) for more details.**

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

# dump schemas of database db1 and db2
dorisdump dump --dump-schema --dbs db1,db2 --host <host> --port <port> --user root --password '***' 

# also dump queries from audit logs of db1 and db2
dorisdump dump --dump-schema --dump-query --dbs db1,db2 --audit-logs 'fe.audit.log,fe.audit.log.20240802-1'

# dump queries from audit log table instead of files, need enable <https://doris.apache.org/docs/admin-manual/audit-plugin>
dorisdump dump --dump-query --audit-log-table <db.table> --from '2024-11-14 18:45:25' --to '2024-11-14 18:45:26'


# Create dump schemas in another DB server
dorisdump create --help

# create all tables and views of db1 and db2, it auto finds dump schemas under 'output/' dir
dorisdump create --dbs db1,db2 --host <host> --port <port> --user root --password '***'

# run any create table/view SQL in db1
dorisdump create --ddl 'dir/*.sql' --db db1


# Generate data (Totally offline!)
dorisdump gendata --help

# gen data for db1 and db2, it auto finds dump schemas under 'output/' dir
dorisdump gendata --dbs db1,db2 --host <host> --port <port> --user root --password '***'

# gen data for t1 and t2 in db1 with config
dorisdump gendata --dbs db1 --table t1,t2 --genconf example/gendata.yaml

# gen data from any create-table SQL
dorisdump gendata --ddl create.sql


# Import data (Require curl command)
dorisdump import --help

# import data for db1, it auto finds generated data under 'output/' dir
dorisdump import --dbs db1,db2 --host <host> --port <port> --user root --password '***'

# import data for t1 and t2 in db1
dorisdump import --dbs db1 --table t1,t2

# import data from any CSV file
dorisdump import --tables db1.t1 --data data.csv


# Replay
dorisdump replay --help

# replay queries from dump sql file
dorisdump replay --host <host> --port <port> --user root --password '***' -f output/sql/q0.sql

# replay with args
dorisdump replay -f output/sql/q0.sql \
    --from '2024-09-20 08:00:00' --to '2024-09-20 09:00:00' \
    --users 'readonly,root' --dbs 'db1,db2' \   # filter sql by users and databases
    --speed 0.5 \                               # increase(< 1.0) or decrease(> 1.0) the time between two serial sqls proportionally, default 1
    --result-dir output/replay \
    --clean                                     # clean 'output/replay' dir before replay


# Diff replay result
dorisdump diff --help

# diff replay result which is slower more than 200ms than original
dorisdump diff --min-duration-diff 200ms --original-sqls 'output/sql/*.sql' output/replay

# diff of two replay result directories
dorisdump diff replay1/ replay2/
```

### Generate Data

Generate CSV data from create-table SQLs. Totally offline!

> [!Tip]
> Not only for Doris, other create-table statements with similar syntax are also supported (like Hive SQL). See [introduction](./introduction-zh.md#生成和导入数据) for more.

Example:

```sh
echo 'create table t1 (
    a varchar(2),
    b struct<foo:tinyint>,
    c date
)' > t1.sql

dorisdump gendata --ddl t1.sql --rows 5

cat output/gendata/t1/*
sO☆{"foo":-66}☆2020-07-23
lg☆{"foo":-121}☆2021-06-15
4☆{"foo":-117}☆2015-06-17
8h☆{"foo":-83}☆2024-09-06
KW☆{"foo":7}☆2019-02-02
```

### Anonymize

**This feature is experimental, case-insensitive, which means `table1` and `TABLE1` will have the same result.** Two ways:

- Use `dorisdump anonymize`:

    ```bash
    echo "select * from table1" | dorisdump anonymize -f -
    ```

- Use `--anonymize` flag while dumping:

    ```bash
    dorisdump dump ... --anonymize
    ```

> [!NOTE]
> Keep `./dorisdump_hashdict.yaml` if you want the result to be consistent (put it at current directory, or specify by `--anonymize-minihash-dict`).

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
