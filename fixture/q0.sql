/*dorisdump{"ts":"2024-08-06 23:44:11.041","client":"192.168.48.119:51970","user":"root","db":"__internal_schema","queryId":"8cb2e4f433e74463-a0ededde7b648b35","durationMs":10}*/ SELECT
	CONCAT('90465', '-', '-1', '-', 'col_varchar_1024__undef_signed_not_null_index_inverted_parser_english_support_phrase_true') AS `id`, 0 AS `catalog_id`, 90443 AS `db_id`, 90465 AS `tbl_id`, -1 AS `idx_id`, 'col_varchar_1024__undef_signed_not_null_index_inverted_parser_english_support_phrase_true' AS `col_id`, NULL AS `part_id`, 30 AS `row_count`, SUM(`t1`.`count`) * COUNT(1) / (SUM(`t1`.`count`) - SUM(IF(`t1`.`count` = 1, 1, 0)) + SUM(IF(`t1`.`count` = 1, 1, 0)) * SUM(`t1`.`count`) / 30) as `ndv`, IFNULL(SUM(IF(`t1`.`column_key` IS NULL, `t1`.`count`, 0)), 0) * 1.0 as `null_count`, SUBSTRING(CAST('' AS STRING), 1, 1024) AS `min`, SUBSTRING(CAST('都发生交到写出来英特有你们上午有车吗发到曲子' AS STRING), 1, 1024) AS `max`, SUM(LENGTH(`column_key`) * count) * 1.0 AS `data_size`, NOW() FROM (     SELECT t0.`colValue` as `column_key`, COUNT(1) as `count`     FROM     (SELECT SUBSTRING(CAST(`col_varchar_1024__undef_signed_not_null_index_inverted_parser_english_support_phrase_true` AS STRING), 1, 1024) AS `colValue`          FROM `internal`.`example`.`table_30_undef_partitions2_keys3_properties4_distributed_by5`       limit 30) as `t0`     GROUP BY `t0`.`colValue` ) as `t1`;
/*dorisdump{"ts":"2024-08-06 23:44:11.042","client":"192.168.48.119:51970","user":"root","db":"__internal_schema","queryId":"f194c444ed2a4e72-87cf9c7091cda06e","durationMs":10}*/ SELECT CONCAT('90465', '-', '-1', '-', 'col_varchar_1024__undef_signed_index_inverted_parser_english_support_phrase_true') AS `id`, 0 AS `catalog_id`, 90443 AS `db_id`, 90465 AS `tbl_id`, -1 AS `idx_id`, 'col_varchar_1024__undef_signed_index_inverted_parser_english_support_phrase_true' AS `col_id`, NULL AS `part_id`, 30 AS `row_count`, SUM(`t1`.`count`) * COUNT(1) / (SUM(`t1`.`count`) - SUM(IF(`t1`.`count` = 1, 1, 0)) + SUM(IF(`t1`.`count` = 1, 1, 0)) * SUM(`t1`.`count`) / 30) as `ndv`, IFNULL(SUM(IF(`t1`.`column_key` IS NULL, `t1`.`count`, 0)), 0) * 1.0 as `null_count`, SUBSTRING(CAST('--' AS STRING), 1, 1024) AS `min`, SUBSTRING(CAST('西欧啊浇水去一次北站拉美徐州配合驱动中科' AS STRING), 1, 1024) AS `max`, SUM(LENGTH(`column_key`) * count) * 1.0 AS `data_size`, NOW() FROM (     SELECT t0.`colValue` as `column_key`, COUNT(1) as `count`     FROM     (SELECT SUBSTRING(CAST(`col_varchar_1024__undef_signed_index_inverted_parser_english_support_phrase_true` AS STRING), 1, 1024) AS `colValue`          FROM `internal`.`example`.`table_30_undef_partitions2_keys3_properties4_distributed_by5`       limit 30) as `t0`     GROUP BY `t0`.`colValue` ) as `t1`;
/*dorisdump{"ts":"2024-08-06 23:44:11.043","client":"192.168.48.118:51970","user":"root","db":"__internal_schema","queryId":"8eaf2c126a249c7-8d48a95bd8501cc9","durationMs":10}*/ SELECT CONCAT('90465', '-', '-1', '-', 'col_varchar_1024__undef_signed_index_inverted') AS `id`, 0 AS `catalog_id`, 90443 AS `db_id`, 90465 AS `tbl_id`, -1 AS `idx_id`, 'col_varchar_1024__undef_signed_index_inverted' AS `col_id`, NULL AS `part_id`, 30 AS `row_count`, SUM(`t1`.`count`) * COUNT(1) / (SUM(`t1`.`count`) - SUM(IF(`t1`.`count` = 1, 1, 0)) + SUM(IF(`t1`.`count` = 1, 1, 0)) * SUM(`t1`.`count`) / 30) as `ndv`, IFNULL(SUM(IF(`t1`.`column_key` IS NULL, `t1`.`count`, 0)), 0) * 1.0 as `null_count`, SUBSTRING(CAST('' AS STRING), 1, 1024) AS `min`, SUBSTRING(CAST('而是' AS STRING), 1, 1024) AS `max`, SUM(LENGTH(`column_key`) * count) * 1.0 AS `data_size`, NOW() FROM (     SELECT t0.`colValue` as `column_key`, COUNT(1) as `count`     FROM     (SELECT SUBSTRING(CAST(`col_varchar_1024__undef_signed_index_inverted` AS STRING), 1, 1024) AS `colValue`          FROM `internal`.`example`.`table_30_undef_partitions2_keys3_properties4_distributed_by5`       limit 30) as `t0`     GROUP BY `t0`.`colValue` ) as `t1`;
/*dorisdump{"ts":"2024-08-06 23:44:12.044","client":"192.168.48.119:51970","user":"root","db":"__internal_schema","queryId":"60b885f02d014194-b225555e4ed26d7e","durationMs":10}*/ SELECT 
    CONCAT('90465', '-', '-1', '-', 'col_varchar_1024__undef_signed_not_null_index_inverted_parser_unicode_support_phrase_true') AS `id`,
    0 AS `catalog_id`,
    90443 AS `db_id`,
    90465 AS `tbl_id`,
    -1 AS `idx_id`,
    'col_varchar_1024__undef_signed_not_null_index_inverted_parser_unicode_support_phrase_true' AS `col_id`,
    NULL AS `part_id`,
    30 AS `row_count`, SUM(`t1`.`count`) * COUNT(1) / (SUM(`t1`.`count`) - SUM(IF(`t1`.`count` = 1, 1, 0)) + SUM(IF(`t1`.`count` = 1, 1, 0)) * SUM(`t1`.`count`) / 30) as `ndv`,
    IFNULL(SUM(IF(`t1`.`column_key` IS NULL, `t1`.`count`, 0)), 0) * 1.0 as `null_count`,
    SUBSTRING(CAST('--' AS STRING), 1, 1024) AS `min`, SUBSTRING(CAST('难道' AS STRING), 1, 1024) AS `max`,
    SUM(LENGTH(`column_key`) * count) * 1.0 AS `data_size`, NOW() 
FROM (
    SELECT
        t0.`colValue` as `column_key`,
        COUNT(1) as `count`
    FROM (
        SELECT SUBSTRING(CAST(`col_varchar_1024__undef_signed_not_null_index_inverted_parser_unicode_support_phrase_true` AS STRING), 1, 1024) AS `colValue`
        FROM `internal`.`example`.`table_30_undef_partitions2_keys3_properties4_distributed_by5`
        limit 30
    ) as `t0`
    GROUP BY `t0`.`colValue`
) as `t1`;
/*dorisdump{"ts":"2024-08-06 23:44:11.045","client":"192.168.48.118:51970","user":"root","db":"__internal_schema","queryId":"ffb1d743a9eb4394-9b48a38bcc0b8b19","durationMs":10}*/ SELECT CONCAT('90465', '-', '-1', '-', 'col_varchar_1024__undef_signed_not_null') AS `id`, 0 AS `catalog_id`, 90443 AS `db_id`, 90465 AS `tbl_id`, -1 AS `idx_id`, 'col_varchar_1024__undef_signed_not_null' AS `col_id`, NULL AS `part_id`, 30 AS `row_count`, SUM(`t1`.`count`) * COUNT(1) / (SUM(`t1`.`count`) - SUM(IF(`t1`.`count` = 1, 1, 0)) + SUM(IF(`t1`.`count` = 1, 1, 0)) * SUM(`t1`.`count`) / 30) as `ndv`, IFNULL(SUM(IF(`t1`.`column_key` IS NULL, `t1`.`count`, 0)), 0) * 1.0 as `null_count`, SUBSTRING(CAST('' AS STRING), 1, 1024) AS `min`, SUBSTRING(CAST('说明心理学而来丰富提货本科退货热插拔' AS STRING), 1, 1024) AS `max`, SUM(LENGTH(`column_key`) * count) * 1.0 AS `data_size`, NOW() FROM (     SELECT t0.`colValue` as `column_key`, COUNT(1) as `count`     FROM     (SELECT SUBSTRING(CAST(`col_varchar_1024__undef_signed_not_null` AS STRING), 1, 1024) AS `colValue`          FROM `internal`.`example`.`table_30_undef_partitions2_keys3_properties4_distributed_by5`       limit 30) as `t0`     GROUP BY `t0`.`colValue` ) as `t1`;
/*dorisdump{"ts":"2024-08-06 23:44:13.046","client":"192.168.48.119:51970","user":"root","db":"__internal_schema","queryId":"41dc7120df0040c0-a677b93ad1a28d27","durationMs":10}*/ SELECT CONCAT('90465', '-', '-1', '-', 'col_varchar_1024__undef_signed_not_null_index_inverted_parser_english_support_phrase_false') AS `id`, 0 AS `catalog_id`, 90443 AS `db_id`, 90465 AS `tbl_id`, -1 AS `idx_id`, 'col_varchar_1024__undef_signed_not_null_index_inverted_parser_english_support_phrase_false' AS `col_id`, NULL AS `part_id`, 30 AS `row_count`, SUM(`t1`.`count`) * COUNT(1) / (SUM(`t1`.`count`) - SUM(IF(`t1`.`count` = 1, 1, 0)) + SUM(IF(`t1`.`count` = 1, 1, 0)) * SUM(`t1`.`count`) / 30) as `ndv`, IFNULL(SUM(IF(`t1`.`column_key` IS NULL, `t1`.`count`, 0)), 0) * 1.0 as `null_count`, SUBSTRING(CAST('' AS STRING), 1, 1024) AS `min`, SUBSTRING(CAST('附注利落王峰亚太时机编号硬盘服务器站用得上太厚' AS STRING), 1, 1024) AS `max`, SUM(LENGTH(`column_key`) * count) * 1.0 AS `data_size`, NOW() FROM (     SELECT t0.`colValue` as `column_key`, COUNT(1) as `count`     FROM     (SELECT SUBSTRING(CAST(`col_varchar_1024__undef_signed_not_null_index_inverted_parser_english_support_phrase_false` AS STRING), 1, 1024) AS `colValue`          FROM `internal`.`example`.`table_30_undef_partitions2_keys3_properties4_distributed_by5`       limit 30) as `t0`     GROUP BY `t0`.`colValue` ) as `t1`;
/*dorisdump{"ts":"2024-08-06 23:44:12.047","client":"192.168.48.118:51970","user":"root","db":"__internal_schema","queryId":"18799575029447f9-a6a3fc65c8eda3f1","durationMs":10}*/ SELECT CONCAT('90465', '-', '-1', '-', 'col_varchar_1024__undef_signed_index_inverted_parser_unicode_support_phrase_false') AS `id`, 0 AS `catalog_id`, 90443 AS `db_id`, 90465 AS `tbl_id`, -1 AS `idx_id`, 'col_varchar_1024__undef_signed_index_inverted_parser_unicode_support_phrase_false' AS `col_id`, NULL AS `part_id`, 30 AS `row_count`, SUM(`t1`.`count`) * COUNT(1) / (SUM(`t1`.`count`) - SUM(IF(`t1`.`count` = 1, 1, 0)) + SUM(IF(`t1`.`count` = 1, 1, 0)) * SUM(`t1`.`count`) / 30) as `ndv`, IFNULL(SUM(IF(`t1`.`column_key` IS NULL, `t1`.`count`, 0)), 0) * 1.0 as `null_count`, SUBSTRING(CAST('' AS STRING), 1, 1024) AS `min`, SUBSTRING(CAST('董事长' AS STRING), 1, 1024) AS `max`, SUM(LENGTH(`column_key`) * count) * 1.0 AS `data_size`, NOW() FROM (     SELECT t0.`colValue` as `column_key`, COUNT(1) as `count`     FROM     (SELECT SUBSTRING(CAST(`col_varchar_1024__undef_signed_index_inverted_parser_unicode_support_phrase_false` AS STRING), 1, 1024) AS `colValue`          FROM `internal`.`example`.`table_30_undef_partitions2_keys3_properties4_distributed_by5`  where a = 1     limit 30) as `t0`     GROUP BY `t0`.`colValue` ) as `t1`;
/*dorisdump{"ts":"2024-08-06 23:44:13.048","client":"192.168.48.118:51970","user":"root","db":"__internal_schema","queryId":"18799575029447f9-a6a3fc65c8eda3f2","durationMs":10}*/ SELECT CONCAT('90465', '-', '-1', '-', 'col_varchar_1024__undef_signed_index_inverted_parser_unicode_support_phrase_false') AS `id`, 0 AS `catalog_id`, 90443 AS `db_id`, 90465 AS `tbl_id`, -1 AS `idx_id`, 'col_varchar_1024__undef_signed_index_inverted_parser_unicode_support_phrase_false' AS `col_id`, NULL AS `part_id`, 30 AS `row_count`, SUM(`t1`.`count`) * COUNT(1) / (SUM(`t1`.`count`) - SUM(IF(`t1`.`count` = 1, 1, 0)) + SUM(IF(`t1`.`count` = 1, 1, 0)) * SUM(`t1`.`count`) / 30) as `ndv`, IFNULL(SUM(IF(`t1`.`column_key` IS NULL, `t1`.`count`, 0)), 0) * 1.0 as `null_count`, SUBSTRING(CAST('' AS STRING), 1, 1024) AS `min`, SUBSTRING(CAST('董事长' AS STRING), 1, 1024) AS `max`, SUM(LENGTH(`column_key`) * count) * 1.0 AS `data_size`, NOW() FROM (     SELECT t0.`colValue` as `column_key`, COUNT(1) as `count`     FROM     (SELECT SUBSTRING(CAST(`col_varchar_1024__undef_signed_index_inverted_parser_unicode_support_phrase_false` AS STRING), 1, 1024) AS `colValue`          FROM `internal`.`example`.`table_30_undef_partitions2_keys3_properties4_distributed_by5`  where a = 2     limit 30) as `t0`     GROUP BY `t0`.`colValue` ) as `t1`;
/*dorisdump{"ts":"2024-08-06 23:44:13.049","client":"192.168.48.118:51970","user":"root","db":"__internal_schema","queryId":"18799575029447f9-a6a3fc65c8eda3f3","durationMs":10}*/ set @a=1;
/*dorisdump{"ts":"2024-08-06 23:44:13.049","client":"192.168.48.119:51970","user":"root","db":"__internal_schema","queryId":"18799575029447f9-a6a3fc65c8eda314","durationMs":10}*/ create database if not exists test;
/*dorisdump{"ts":"2024-08-06 23:44:13.050","client":"192.168.48.119:51970","user":"root","db":"__internal_schema","queryId":"18799575029447f9-a6a3fc65c8eda3f4","durationMs":59000}*/ create table if not exists test.t1(a int) DISTRIBUTED BY HASH (a) PROPERTIES ("replication_allocation" = "tag.location.default: 1");
/*dorisdump{"ts":"2024-08-06 23:45:13.050","client":"192.168.48.119:51970","user":"root","db":"__internal_schema","queryId":"18799575029447f9-a6a3fc65c8eda3f5","durationMs":10}*/ insert into test.t1 values (1);
