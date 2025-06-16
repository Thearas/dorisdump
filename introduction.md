# Introduction

- [Workflow](#workflow)
- [Dump](#dump)
  - [Dump Tables and Views](#dump-tables-and-views)
  - [Dump Queries](#dump-queries)
  - [Other Dump Parameters](#other-dump-parameters)
- [Create Tables and Views](#create-tables-and-views)
- [Generate and Import Data](#generate-and-import-data)
  - [Default Generation Rules](#default-generation-rules)
  - [Custom Generation Rules](#custom-generation-rules)
- [Replay](#replay)
  - [Replay Speed and Concurrency](#replay-speed-and-concurrency)
  - [Other Replay Parameters](#other-replay-parameters)
- [Compare Replay Results](#compare-replay-results)
- [Best Practices](#best-practices)
  - [Command-line Prompts and Autocompletion](#command-line-prompts-and-autocompletion)
  - [Environment Variables and Configuration Files](#environment-variables-and-configuration-files)
  - [Monitoring Dump/Replay Process](#monitoring-dumpreplay-process)
  - [Multi-FE Replay](#multi-fe-replay)
  - [Large-scale Batch Replay](#large-scale-batch-replay)
  - [Find SQLs with Replay Duration Exceeding 1s](#find-sqls-with-replay-duration-exceeding-1s)
  - [Automation](#automation)
- [Anonymization](#anonymization)
- [FAQ](#faq)
  - [How to provide the tool to customers, and is there any impact on the production environment?](#how-to-provide-the-tool-to-customers-and-is-there-any-impact-on-the-production-environment)
  - [The number of dump SQLs is less than in the audit log](#the-number-of-dump-sqls-is-less-than-in-the-audit-log)
  - [Dump SQL has syntax errors](#dump-sql-has-syntax-errors)
  - [Dump statistics do not match the actual ones](#dump-statistics-do-not-match-the-actual-ones)
  - [Replay error: too many connections](#replay-error-too-many-connections)

## Workflow

There are two types of workflows, with each step representing a `dorisdump` command:

- No data generation needed: `Dump -> Replay -> Compare Replay Results`
- Data generation needed: `Dump -> Create Tables and Views (Optional) -> Generate and Import Data -> Replay -> Compare Replay Results`

## Dump

`dorisdump dump --help`

This is divided into two parts: "Dump Tables and Views" and "Dump Queries". Both can be combined into a single `dorisdump` command.

### Dump Tables and Views

`dorisdump dump --dump-schema`

Dumps `CREATE` statements for tables and views from a Doris database. By default, it also dumps table statistics. If the statistics differ significantly from the actual data, specifying `--analyze` is recommended. See [Dump statistics do not match the actual ones](#dump-statistics-do-not-match-the-actual-ones).

```sh
# Dump all tables and views from db1 and db2
dorisdump dump --dump-schema --host xxx --port xxx --user xxx --password xxx --dbs db1,db2

# Default dump to output/ddl:
output
└── ddl
    ├── db1.t1.table.sql
    ├── db1.stats.yaml
    ├── db2.t2.table.sql
    └── db2.stats.yaml
```

### Dump Queries

`dorisdump dump --dump-query`

Queries can be dump from an audit log table or file. By default, only `SELECT` statements are dump. You can add `--only-select=false` to dump other statements as well.

```sh
# Dump from an audit log table, table name is usually __internal_schema.audit_log
dorisdump dump --dump-query --audit-log-table <db.table> --from '2024-11-14 17:00:00' --to '2024-11-14 18:00:00' --host xxx --port xxx --user xxx --password xxx

# Dump from audit log files, '*' matches multiple files (note the single quotes)
dorisdump dump --dump-query --audit-logs 'fe.audit.log,fe.audit.log.20240802*'

# Default dump to output/sql:
output
└── sql
    ├── q0.sql
    └── q1.sql
```

> [!NOTE]
>
> - When dumping from log files, `q0.sql` corresponds to the first log file, `q1.sql` to the second, and so on. However, when dumping from a log table, all queries are written to `q0.sql`.
> - Dump results are appended to previously dump SQL files unless `--clean` is specified, which deletes the previous `output/ddl` and `output/sql` directories.

### Other Dump Parameters

- `--analyze`: Automatically runs `ANALYZE TABLE <table> WITH SYNC` before dumping a table to make statistics more accurate. Default is off.
- `--parallel`: Controls the dump concurrency. Increasing it speeds up the dump; decreasing it uses fewer resources. Default is `min(machine_cores, 10)`.
- `--dump-stats`: Also dumps table statistics when dumping tables. Statistics are dump to `output/ddl/db.stats.yaml`. Default is on.
- `--only-select`: Whether to dump only `SELECT` statements. Default is on.
- `--from` and `--to`: Dump SQL within a specified time range.
- `--query-min-duration`: Minimum execution duration for dump SQL.
- `--query-states`: States of the SQL to be dump, can be `ok`, `eof`, and `err`.
- `-s, --strict`: Validates SQL syntax correctness when dumping from audit logs.
- `--audit-log-encoding`: Audit log file encoding. Default is auto-detect.
- `--anonymize`: Anonymizes data during dump, e.g., `select * from table1` becomes `select * from a`.
- `--anonymize-xxx`: Other anonymization parameters, see [Anonymization](#anonymization).

## Create Tables and Views

`dorisdump create --help`

You need to first [Dump Tables and Views](#dump-tables-and-views) locally, then create them in another Doris instance:

```sh
# Create all dump tables and views for db1 and db2
dorisdump create --dbs db1,db2

# Create dump table1 and table2
dorisdump create --dbs db1 --tables table1,table2

# Run any create table/view SQL in db1
dorisdump create --ddl 'dir/*.sql' --db db1
```

## Generate and Import Data

`dorisdump gendata --help`/`dorisdump import --help`

You need to first [Dump Tables and Views](#dump-tables-and-views) locally, then generate and import data:

```sh
# Generate data for all dump tables in db1 and db2
dorisdump gendata --dbs db1,db2

# Generate data for dump table1
dorisdump gendata --tables db1.table1 # or --dbs db1 --tables table1

# Data can also be generated for any create table SQL without prior dump
# P.S. It might not necessarily be Doris; other databases like Hive might also work, but it hasn't been tested ;)
dorisdump gendata --ddl my_create_table.sql


# Import data for all tables with generated data in db1 and db2
dorisdump import --dbs db1,db2

# Import data for table1 with generated data
dorisdump import --tables db1.table1 # or --dbs db1 --tables table1

# Import any CSV data file into table1
dorisdump import --tables db1.table1 --data my_data.csv
```

In implementation, the tool performs these actions in two stages based on the `--dbs` and `--tables` parameters:

1. In the generation stage:

    1. Scans the dump directory `output/ddl/` for matching `<db>.<table>.table.sql` files. The dump directory can be specified with `--ddl`.
    2. Combines the corresponding statistics file `<db>.stats.yaml` with the custom generation rules file (specified by `--genconf`) to determine the final generation rules.
    3. Generates CSV files into the data generation directory `output/gendata/<db>.<table>/` according to the generation rules.
2. In the import stage:

    1. Scans the data generation directory `output/gendata/` for matching `<db>.<table>/*` data files. The data generation directory can be specified with `--data`.
    2. Uses the `curl` command to run StreamLoad for data import.

> [!TIP]
>
> - A maximum of one million rows of data can be generated per table.
> - Specifying `-Ldebug` during import shows the specific `curl` command, which is helpful for reproducing and troubleshooting issues.

### Default Generation Rules

By default, `NULL` values are not generated. This can be changed by specifying `null_frequency` in [Custom Generation Rules](#custom-generation-rules).

Default generation rules for various types:

| Type | Length | Min - Max | Structure |
| --- | --- | --- | --- |
| ARRAY | 1 - 3 |  |  |
| MAP | 1 - 3 |  |  |
| JSON/JSONB |  |  | STRUCT<col1:SMALLINT, col2:SMALLINT> |
| VARIANT |  |  | STRUCT<col1:SMALLINT, col2:SMALLINT> |
| BITMAP | 5 | element: 0 - MaxInt32 |  |
| TEXT/STRING/VARCHAR | 1 - 10 |  |  |
| TINYINT |  | MinInt8 - MaxInt8 |  |
| SMALLINT |  | MinInt16 - MaxInt16 |  |
| INT |  | MinInt32 - MaxInt32 |  |
| BIGINT |  | MinInt32 - MaxInt32 |  |
| LARGEINT |  | MinInt32 - MaxInt32 |  |
| FLOAT |  | MinInt16 - MaxInt16 | |
| DOUBLE |  | MinInt32 - MaxInt32 |  |
| DECIMAL |  | MinInt32 - MaxInt32 |  |
| DATE |  | 10 years ago - now |  |
| DATETIME |  | 10 years ago - now |  |

### Custom Generation Rules

Specify with `--genconf gendata.yaml` during data generation. See [example/gendata.yaml](./example/gendata.yaml) for an example.

## Replay

`dorisdump replay --help`

You need to first [Dump Queries](#dump-queries), then replay based on the dump SQL files.

```sh
# Dump
dorisdump dump --dump-query --audit-logs fe.audit.log

# Replay, results are placed in the `output/replay` directory by default. Each file represents a client, and each line in the file represents the result of a SQL query.
dorisdump replay -f output/q0.sql
```

> [!NOTE]
> Executing the replay command multiple times will append results to the previous results, unless `--clean` is specified, which deletes the previous `output/replay` directory.

---

### Replay Speed and Concurrency

The principle of replay is that SQL from different clients runs concurrently, while SQL from the same client runs serially with an interval, strictly calculated according to the audit log:

```sh
# sql1 and sql2 are two consecutive SQLs executed by the same client
Interval duration = sql2 start time - sql1 start time - sql1 execution duration
```

#### Custom Speed and Concurrency

Controlled by the following parameters:

- `--speed`: Controls the replay speed, affecting the "interval duration" mentioned above. For example, `--speed 0.5` means slowing down by half, while `--speed 2` means speeding up twice. The principle is to proportionally increase or decrease the interval duration between consecutive SQLs from the same client. Note that if the SQL execution time itself is too long, `--speed` may not be effective.
- `--client-count int`: Resets the number of clients, and all SQLs will be evenly distributed among the clients to run in parallel. **Setting this value makes the replay effect unpredictable!** By default, it's the same as the number of clients in the log to achieve the same effect as online.

> [!TIP]
> If you only want to replay with 50 concurrency without intervals, and each SQL is independent, you can set `--speed 999999 --client-count 50`.

---

### Other Replay Parameters

- `-c, --cluster`: The cluster for replay, only useful in Cloud mode.
- `--result-dir`: Replay result directory, default `output/replay`.
- `--users`: Only replay SQL initiated by these users, default is to replay for all users.
- `--from` and `--to`: Replay SQL within a specified time range.
- `--max-hash-rows`: Maximum number of hash result rows to record during replay, used to compare if two replay results are consistent. Default is no hashing.
- `--max-conn-idle-time`: Maximum idle time for a client connection. If the interval duration between consecutive SQLs from the same client exceeds this value, the connection will be recycled. Default is `5s`.

## Compare Replay Results

`dorisdump diff --help`

There are two ways:

1. Compare two replay results:

    ```sh
    dorisdump diff output/replay1 output/replay2
    ```

2. Compare dump SQL with its replay result:

    ```sh
    dorisdump diff --min-duration-diff 2s --original-sqls 'output/sql/*.sql' output/replay
    ```

> `--min-duration-diff` means print SQLs whose execution duration difference exceeds this value. Default is `100ms`.

## Best Practices

### Command-line Prompts and Autocompletion

`dorisdump completion --help`

When the installation is complete or when you execute the command above, it will provide instructions on how to enable autocompletion.

---

### Environment Variables and Configuration Files

`dorisdump --help`

Besides command-line arguments, there are two other ways:

1. Pass parameters through uppercase environment variables prefixed with `DORIS_xxx`, e.g., `DORIS_HOST=xxx` is equivalent to `--host xxx`.
2. Pass parameters through a configuration file, e.g., `dorisdump --config-file xxx.yaml`. See [example](./example/example.dorisdump.yaml).

Parameter priority from high to low:

1. Command-line arguments
2. Environment variables
3. Configuration file specified by `--config`
4. Default configuration file `~/.dorisdump.yaml`

---

### Monitoring Dump/Replay Process

`--log-level debug/trace`

`debug` outputs a brief process, while `trace` shows a detailed process, such as SQL execution time and details during replay.

---

### Multi-FE Replay

Each FE's audit log is separate. When dumping, they must be dump separately. When replaying, they must also be replayed separately and simultaneously. For example, for a 2 FE cluster:

```sh
# Dump audit logs for fe1 and fe2 separately
dorisdump dump --dump-query --audit-logs fe1.audlt.log -O fe1
dorisdump dump --dump-query --audit-logs fe2.audlt.log -O fe2

# Replay audit logs for fe1 and fe2 simultaneously
nohup dorisdump replay -H <fe1.ip> -f fe1/sql/q0.sql -O fe1 &
nohup dorisdump replay -H <fe2.ip> -f fe2/sql/q0.sql -O fe2 &
```

---

### Large-scale Batch Replay

When the volume of SQL to be replayed is too large, for example, replaying logs for a whole month (31 days), it's best to replay in hourly batches. Use `--from` and `--to` during dump to batch (or batch manually after dump). Example:

```sh
dump YEAR_MONTH="2025-03" # <-- Change this line
dump DORIS_YES=1

for day in {1..31} ; do
  day=$(printf "%02d" $day)

  for hour in {0..23} ; do
      hour=$(printf "%02d" $hour)
      output=output/$day/$hour
      sql=$output/q0.sql

      echo "dumping and replaying at $day-$hour"

      # Dump
      dorisdump dump --dump-query --from "$YEAR_MONTH-$day $hour:00:00" --to "$YEAR_MONTH-$day $hour:59:59" --audit-log-table __internal_schema.audit_log --output "$output"

      # Replay, clear previous replay results, 50 clients concurrently, run continuously
      dorisdump replay -f "$sql" --result-dir result --clean --client-count 50 --speed 999999

      # View replay results
      dorisdump diff --min-duration-diff 1s --original-sqls $sql result -Ldebug 2>&1 | tee -a "result-$day.txt"
  done
done
```

---

### Find SQLs with Replay Duration Exceeding 1s

Search the replay result directory directly. It is recommended to use [ripgrep](https://github.com/BurntSushi/ripgrep); using `grep` is similar:

```sh
# Find SQLs with execution time exceeding 1s
rg '"durationMs":\d{4}' output/replay

# Find SQLs with execution time exceeding 6s
rg -e '"durationMs":[6-9]\d{3}' -e '"durationMs":\d{5}' output/replay
```

---

### Automation

For example, when writing scripts to dump/replay multiple files, it's inconvenient to manually input `y` for confirmation. You can set the environment variable `DORIS_YES=1` or `DORIS_YES=0` to automatically confirm or deny.

---

## Anonymization

`dorisdump anonymize --help`

For basic usage, see [README.md](./README.md#anonymize).

Anonymization uses the Go version of Doris Antlr4 Parser, which is currently case-insensitive. For example, `table1` and `TABLE1` will produce the same result.

### Parameters

- `-f, --file`: Read SQL from a file. If '-', read from standard input.
- `--anonymize-reserve-ids`: Reserve ID fields, do not anonymize them.
- `--anonymize-id-min-length`: ID fields with a length less than this value will not be anonymized. Default is `3`.
- `--anonymize-method`: Hash method, `hash` or `minihash`. The latter generates a concise dictionary based on the former, making anonymized IDs shorter. Default is `minihash`.
- `--anonymize-minihash-dict`: When the hash method is `minihash`, specify the concise dictionary file. Default is `./dorisdump_hashdict.yaml`.

## FAQ

### How to provide the tool to customers, and is there any impact on the production environment?

If customers cannot access the internet, download the [latest binary](https://github.com/Thearas/dorisdump/releases) and provide it directly. The Linux version has no dependencies and can run directly on the machine. By default, the tool will not perform any write operations on the cluster.

If you are concerned about resource consumption during dump, you can set `--parallel=1`. Memory consumption will be at most tens of megabytes, and execution time is generally in seconds.

When replaying, please execute in batches and ensure sufficient resources.

---

### The number of dump SQLs is less than in the audit log

First, check the [Dump Queries](#dump-queries) parameters to see if any SQLs were filtered out.

Then, enable `--log-level=debug` to see if any of the following situations occurred:

- `ignore sql with duplicated query_id`: Duplicate `query_id`s will be ignored. This is a bug in Doris itself.
- `query has been truncated`: SQL that is too long will be truncated. Please check Doris's [`audit_plugin_max_sql_length`](https://doris.apache.org/docs/admin-manual/audit-plugin#audit-log-configuration) configuration.

---

### Dump SQL has syntax errors

It is recommended to add the `-s, --strict` parameter during dump to validate SQL syntax correctness. Then, check the `query_id` output by the tool and find it in the audit log. It could be one of the following two situations:

1. The tool unescapes `\r`, `\n`, and `\t` in log SQL. However, if the original SQL itself contains these characters, it may cause syntax errors.
2. There is an issue with the audit log itself.

Generally, errors are infrequent. You can manually modify the dump SQL.

---

### Dump statistics do not match the actual ones

Check the `method` field of columns in the dump `stats.yaml`. If it is `SAMPLE` (i.e., sampling), then there might be a significant deviation from the actual values.

```yaml
columns:
  - name: col_int
    ndv: 10
    null_count: 4969
    data_size: 800000
    avg_size_byte: 8
    min: "2022"
    max: "2030"
    method: SAMPLE # <-- here
```

It is recommended to specify `--analyze` during dump, or first manually execute `ANALYZE DATABASE WITH SYNC`/`ANALYZE TABLE WITH SYNC`, and then dump.

---

### Replay error: too many connections

There are two situations:

1. Because audit logs do not contain a `connection id`, the tool replays based on clients (`ip:port`) rather than connections. However, a client may consist of multiple serial connections. In this case, the tool does not know when to disconnect, leading to connections not being released promptly.

    By default, replay connections are automatically released after 5s of inactivity. However, sometimes too many connections and loss of session variables can still occur. You can adjust `--max-conn-idle-time`.
2. If `--speed` is set too high, too many SQLs are squeezed into a short period for execution. Reducing the `--speed` value can solve this. Refer to [Custom Speed and Concurrency](#custom-speed-and-concurrency).

Additionally, there is a general solution: increase the user's maximum number of connections [`max_user_connections`](https://doris.apache.org/docs/admin-manual/config/user-property#max_user_connections).
