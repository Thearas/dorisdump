# Introduction

- [导出](#导出)
  - [导出表和视图](#导出表和视图)
  - [导出查询](#导出查询)
  - [其他导出参数](#其他导出参数)
- [回放](#回放)
  - [回放速度和并发](#回放速度和并发)
  - [其他回放参数](#其他回放参数)
- [对比回放结果](#对比回放结果)
- [最佳实践](#最佳实践)
  - [命令行提示与自动补全](#命令行提示与自动补全)
  - [环境变量和配置文件](#环境变量和配置文件)
  - [监看导出/回放过程](#监看导出回放过程)
  - [分批回放](#分批回放)
  - [找出回放时长超过 1s 的 SQL](#找出回放时长超过-1s-的-sql)
  - [自动化](#自动化)
- [脱敏](#脱敏)
- [FAQ](#faq)
  - [怎么把工具给客户，对生产环境有没有影响](#怎么把工具给客户对生产环境有没有影响)
  - [导出的 SQL 数量比审计日志里的少](#导出的-sql-数量比审计日志里的少)
  - [导出的 SQL 有语法错误](#导出的-sql-有语法错误)
  - [导出的统计信息与实际不符](#导出的统计信息与实际不符)
  - [回放报错连接数超了](#回放报错连接数超了)

## 导出

`dorisdump dump --help`

分为两部分，「导出表」和「导出查询」。这两部分可以合入一条 `dorisdump` 命令。

### 导出表和视图

`dorisdump dump --dump-schema`

需要连上 Doris 数据库，导出表和视图的 `CREATE` 语句，默认也会导出表的统计信息。

```sh
# 导出 db1 和 db2 的所有表和视图
dorisdump dump --dump-schema --host xxx --port xxx --user xxx --password xxx --dbs db1,db2

# 默认导出到 output/ddl 下：
output
└── ddl
    ├── db1.t1.table.sql
    ├── db1.stats.yaml
    ├── db2.t2.table.sql
    └── db2.stats.yaml

```

### 导出查询

`dorisdump dump --dump-query`

可以从审计日志表或文件导出，默认只导出 `SELECT` 语句，可以加上 `--only-select=false` 一并导出其他语句。

```sh
# 从审计日志表导出，表名一般是 __internal_schema.audit_log
dorisdump dump --dump-query --audit-log-table <db.table> --from '2024-11-14 17:00:00' --to '2024-11-14 18:00:00'

# 从审计日志文件导出，'*' 代表匹配多个文件（注意一定要用引号括起来）
dorisdump dump --dump-query --audit-logs 'fe.audit.log,fe.audit.log.20240802*'

# 默认导出到 output/sql 下：
output
└── sql
    ├── q0.sql
    └── q1.sql
```

注意：

- 从日志文件导出时，`q0.sql` 对应第一个日志文件、`q1.sql` 对应第二个、以此类推，但从日志表导出时，只会写入到 `q0.sql`
- 多次执行导出命令，结果会追加到前一次导出的 SQL 文件中，除非指定 `--clean`，删除之前的 `output/ddl` 和 `output/sql` 目录

### 其他导出参数

- `--parallel` 控制导出并发量，调大导出更快，调小占用资源更少，默认 `min(机器核数, 10)`
- `--dump-stats` 导出表时也导出统计信息，导出在 `output/ddl/db.stats.yaml` 文件，默认开启
- `--only-select` 是否从只导出 `SELECT` 语句，默认开启
- `--from` 和 `--to` 导出时间范围内的 SQL
- `--query-min-duration` 导出 SQL 的最小执行时长
- `--query-states` 导出 SQL 的状态，可以是 `ok`、`eof` 和 `err`
- `-s, --strict` 从审计日志导出时校验 SQL 语法正确性
- `--audit-log-encoding` 审计日志文件编码，默认自动检测
- `--anonymize` 导出时脱敏，比如 `select * from table1` 变为 `select * from a`
- `--anonymize-xxx` 其他脱敏参数，见 [脱敏](#脱敏)

## 回放

`dorisdump replay --help`

需要先[导出查询](#导出查询)，然后基于导出的 SQL 文件回放。

```sh
# 导出
dorisdump dump --dump-query --audit-logs fe.audit.log

# 回放，工具会自动推荐并发数，我们确认就行
dorisdump replay -f output/q0.sql
```

结果默认放在 `output/replay` 目录下，每个文件代表一个客户端，文件中每行代表一条 SQL 的结果。

> 注意：执行多次回放命令，结果会追加到前一次的结果中，除非指定 `--clean`，删除之前的 `output/replay` 目录

---

### 回放速度和并发

回放的原理是不同客户端的 SQL 并发，同一客户端的 SQL 串行且有间隔时长，严格按照审计日志计算：

```sh
# sql1 和 sql2 是同一个客户端相邻执行的两条 SQL
间隔时长 = sql1 开始时间 - sql2 开始时间 - sql1 执行时长
```

#### 自定义速度和并发

由以下参数控制：

- `--speed` 控制回放速度，作用于上面提到的「间隔时长」，比如 `--speed 0.5` 代表放慢一倍，而 `--speed 2` 代表加快一倍。原理是按比例增加或减少同一客户端的相邻 SQL 的间隔时长，注意如果 SQL 本身的执行时间过长，则 `--speed` 效果不佳

- `--client-count int` 重新设置客户端数目，并且所有 SQL 都将被均衡地分散到各个客户端并行跑，**！！！设置此值无法预料回放效果！！！**。默认跟日志里的客户端数一样，以达到跟线上相同的效果

示例：如果只想以 50 并发无间隔回放，且每条 SQL 都独立无依赖，可以设置 `--speed 999999 --client-count 50`。

---

### 其他回放参数

- `-c, --cluster` 回放的集群，仅在 Cloud 模式下有用
- `--result-dir` 回放结果目录，默认 `output/replay`
- `--users` 只回放这些用户发起的 SQL，默认回放全部用户的
- `--from` 和 `--to` 回放时间范围内的 SQL
- `--count` 回放 SQL 数量，默认回放全部
- `--max-hash-rows` 回放时记录的最大 hash 结果行数，用于对比两次回放结果是否一致，默认不 hash
- `--max-conn-idle-time` 客户端连接的最大空闲时间，同一客户端的相邻 SQL 的间隔时长超出此值时，连接会被回收，默认 `10s`

## 对比回放结果

`dorisdump diff --help`

按照场景的不同有两种方式：

1. 对比两次回放结果，比如升降级场景：

    ```sh
    dorisdump diff output/replay1 output/replay2
    ```

2. 对比导出的 SQL 和它的回放结果：

    ```sh
    dorisdump diff --min-duration-diff 2s --original-sqls 'output/sql/*.sql' output/replay
    ```

> `--min-duration-diff` 表示打印执行时长差异超过此值的 SQL，默认 `100ms`

## 最佳实践

### 命令行提示与自动补全

`dorisdump completion --help`

安装完成或执行上面的命令时，会给出启用自动补全的方法。

---

### 环境变量和配置文件

`dorisdump --help`

除了命令行传参，还有两种方式：

1. 通过前缀为 `DORIS_xxx` 的大写环境变量传参，比如 `DORIS_HOST=xxx` 等价于  `--host xxx`
2. 通过配置文件传参，比如 `dorisdump --config-file xxx.yaml`，见 [example](./example/example.dorisdump.yaml)

参数优先级由高到低：

1. 命令行参数
2. 环境变量
3. `--config` 指定的配置文件
4. 默认配置文件 `~/.dorisdump.yaml`

---

### 监看导出/回放过程

`--log-level debug/trace`

`debug` 输出简略过程，而 `trace` 可以看到详细过程，比如回放时 SQL 的执行时间和详情等。

---

### 分批回放

回放的 SQL 量太大时，比如回放一个月 31 天的日志，最好以小时为单位分批回放，在导出时用 `--from` 和 `--to` 分批（或导出后手动分批），示例：

```sh
export YEAR_MONTH="2025-03" # <-- 改这一行
export DORIS_YES=1

for day in {1..31} ; do
  day=$(printf "%02d" $day)

  for hour in {0..23} ; do
      hour=$(printf "%02d" $hour)
      output=output/$day/$hour
      sql=$output/q0.sql

      echo "dumping and replaying at $day-$hour"

      # 导出
      dorisdump dump --dump-query --from "$YEAR_MONTH-$day $hour:00:00" --to "$YEAR_MONTH-$day $hour:59:59" --audit-log-table __internal_schema.audit_log --output "$output"

      # 回放，并清除前一次回放结果，50 个客户端并发，不间断跑
      dorisdump replay -f "$sql" --result-dir result --clean --client-count 50 --speed 999999

      # 查看回放结果
      dorisdump diff --min-duration-diff 1s --original-sqls $sql result -Ldebug 2>&1 | tee -a "result-$day.txt"
  done
done
```

---

### 找出回放时长超过 1s 的 SQL

直接搜索回放结果目录，建议用 [ripgrep](https://github.com/BurntSushi/ripgrep)，用 `grep` 也类似：

```sh
# 找出执行时间超过 1s 的
rg '"durationMs":\d{4}' output/replay

# 找出执行时间超过 6s 的
rg -e '"durationMs":[6-9]\d{3}' -e '"durationMs":\d{5}' output/replay
```

---

### 自动化

比如写脚本导出/回放多个文件时，不方便手动输入 `y` 确认，可以设置环境变量 `DORIS_YES=1` 或 `DORIS_YES=0` 自动确认或否认。

---

## 脱敏

`dorisdump anonymize --help`

基础使用见 [README.md](./README.md#anonymize)。

脱敏使用 Go 版本的 Doris Anltr4 Parser，目前是大小写不敏感的，比如 `table1` 和 `TABLE1` 会有相同的结果。

### 参数

- `-f, --file` 从文件读取 SQL，如果为 '-' 则从标准输入读取
- `--anonymize-reserve-ids` 保留 ID 字段，不做脱敏
- `--anonymize-id-min-length` 长度小于此值的 ID 字段不做脱敏，默认 `3`
- `--anonymize-method` hash 方法，`hash` 或 `minihash`，后者在前者的基础上生成简要字典，让脱敏后的 ID 变短，默认是 `minhash`
- `--anonymize-minihash-dict` 当 hash 方法为 `minihash` 时，指定简要字典文件，默认 `./dorisdump_hashdict.yaml`

## FAQ

### 怎么把工具给客户，对生产环境有没有影响

客户不能科学上网的话，把[最新二进制](https://github.com/Thearas/dorisdump/releases)下载下来直接给，Linux 版是无依赖的，放机器上就能跑，工具对集群不会有任何写入操作。

导出时担心消耗资源的话，可以设置 `--parallel=1`，内存消耗最多几十兆，一般执行时间在秒级。

回放时请分批执行并保证资源充足。

---

### 导出的 SQL 数量比审计日志里的少

首先检查[导出参数](#导出查询)，看看是不是过滤掉了。

然后打开 `--log-level=debug`，看看是不是以下情况：

- `ignore sql with duplicated query_id`：重复的 `query_id` 会被忽略，这是 Doris 本身的 bug
- `query has been truncated`：SQL 过长会被截断，请检查 Doris 的 [`audit_plugin_max_sql_length`](https://doris.apache.org/docs/admin-manual/audit-plugin#audit-log-configuration) 配置

---

### 导出的 SQL 有语法错误

建议在导出时加上 `-s, --strict` 参数，校验 SQL 语法正确性。然后看工具输出的 `query_id`，去审计日志里找，有可能是以下两种情况：

1. 工具默认会反转义日志 SQL 中的 `\\r`、`\\n` 和 `\\t`，但如果原始 SQL 本身包含这些字符，就可能导致语法错误
2. 审计日志本身有问题

一般出错的情况不多，手动修改导出后的 SQL 即可。

---

### 导出的统计信息与实际不符

检查导出的 `stats.yaml` 中列的 `method` 字段，如果是 `SAMPLE`（即采样），那么与实际可能偏差较大。

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

推荐先执行 `ANALYZE DATABASE WITH SYNC` 或 `ANALYZE TABLE WITH SYNC`，然后再导出。

---

### 回放报错连接数超了

有两种情况：

1. 因为审计日志中没有 `connection id`，所以工具以客户端（`ip:port`）而不是连接为单位进行回放，但一个客户端可能由多个串行的连接组成，这样工具不知道何时断开连接，从而导致连接不能及时释放

    回放连接默认 10s 无活动自动释放，但有时还是会出现连接过多和 session 变量丢失的情况，可以调整 `--max-conn-idle-time`
2. `--speed` 设置过大，过多的 SQL 被挤压到一小段时间执行，减小 `--speed` 值即可解决，参考[自定义速度和并发](#自定义速度和并发)

另外有一个通解：调大用户最大连接数 [`max_user_connections`](https://doris.apache.org/docs/admin-manual/config/user-property#max_user_connections)。
