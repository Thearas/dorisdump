package src

import (
	"fmt"

	"github.com/samber/lo"
)

var BuiltinFunctionHashs = lo.SliceToMap(BuiltinFunctions, func(s string) (string, string) {
	b := hashstr(hasher, s)
	return fmt.Sprintf(AnonymizeHashFmt, b[:AnonymizeHashBytes]), s
})

// Generate from `SHOW BUILTIN FUNCTIONS`.
var BuiltinFunctions = []string{
	"%element_extract%",
	"%element_slice%",
	"abs",
	"acos",
	"add",
	"add_months",
	"adddate",
	"aes_decrypt",
	"aes_encrypt",
	"and",
	"any",
	"any_value",
	"append_trailing_char_if_absent",
	"approx_count_distinct",
	"array",
	"array_agg",
	"array_apply",
	"array_avg",
	"array_compact",
	"array_concat",
	"array_contains",
	"array_contains_all",
	"array_count",
	"array_cum_sum",
	"array_difference",
	"array_distinct",
	"array_enumerate",
	"array_enumerate_uniq",
	"array_except",
	"array_exists",
	"array_filter",
	"array_first_index",
	"array_intersect",
	"array_join",
	"array_last_index",
	"array_max",
	"array_min",
	"array_popback",
	"array_popfront",
	"array_position",
	"array_product",
	"array_pushback",
	"array_pushfront",
	"array_range",
	"array_remove",
	"array_repeat",
	"array_reverse_sort",
	"array_reverse_split",
	"array_shuffle",
	"array_size",
	"array_slice",
	"array_sort",
	"array_sortby",
	"array_split",
	"array_sum",
	"array_union",
	"array_with_constant",
	"array_zip",
	"arrays_overlap",
	"ascii",
	"asin",
	"atan",
	"atan2",
	"auto_partition_name",
	"avg",
	"avg_weighted",
	"bin",
	"bit_count",
	"bit_length",
	"bit_shift_left",
	"bit_shift_right",
	"bitand",
	"bitmap_agg",
	"bitmap_and",
	"bitmap_and_count",
	"bitmap_and_not",
	"bitmap_and_not_count",
	"bitmap_andnot",
	"bitmap_andnot_count",
	"bitmap_contains",
	"bitmap_count",
	"bitmap_empty",
	"bitmap_from_array",
	"bitmap_from_base64",
	"bitmap_from_string",
	"bitmap_has_all",
	"bitmap_has_any",
	"bitmap_hash",
	"bitmap_hash64",
	"bitmap_intersect",
	"bitmap_max",
	"bitmap_min",
	"bitmap_not",
	"bitmap_or",
	"bitmap_or_count",
	"bitmap_remove",
	"bitmap_subset_in_range",
	"bitmap_subset_limit",
	"bitmap_to_array",
	"bitmap_to_base64",
	"bitmap_to_string",
	"bitmap_union",
	"bitmap_union_count",
	"bitmap_union_int",
	"bitmap_xor",
	"bitmap_xor_count",
	"bitnot",
	"bitor",
	"bitxor",
	"cardinality",
	"casttobigint",
	"casttoboolean",
	"casttochar",
	"casttodate",
	"casttodatetime",
	"casttodatetimev2",
	"casttodatev2",
	"casttodecimal128",
	"casttodecimal256",
	"casttodecimal32",
	"casttodecimal64",
	"casttodecimalv2",
	"casttodouble",
	"casttofloat",
	"casttoint",
	"casttoipv4",
	"casttoipv6",
	"casttojson",
	"casttolargeint",
	"casttosmallint",
	"casttostring",
	"casttotime",
	"casttotimev2",
	"casttotinyint",
	"casttovarchar",
	"casttovariant",
	"cbrt",
	"ceil",
	"ceiling",
	"char",
	"char_length",
	"character_length",
	"coalesce",
	"collect_list",
	"collect_set",
	"concat",
	"concat_ws",
	"conv",
	"convert_to",
	"convert_tz",
	"corr",
	"cos",
	"cosh",
	"cosine_distance",
	"count",
	"count_by_enum",
	"countequal",
	"covar",
	"covar_pop",
	"covar_samp",
	"crc32",
	"cume_dist",
	"curdate",
	"current_date",
	"current_time",
	"current_timestamp",
	"curtime",
	"cut_ipv6",
	"date",
	"date_add",
	"date_ceil",
	"date_floor",
	"date_format",
	"date_sub",
	"date_trunc",
	"datediff",
	"datev2",
	"day",
	"day_ceil",
	"day_floor",
	"dayname",
	"dayofmonth",
	"dayofweek",
	"dayofyear",
	"days_add",
	"days_diff",
	"days_sub",
	"dceil",
	"decode_as_varchar",
	"degrees",
	"dense_rank",
	"dexp",
	"dfloor",
	"digital_masking",
	"divide",
	"dlog1",
	"dlog10",
	"domain",
	"domain_without_www",
	"dpow",
	"dround",
	"dsqrt",
	"e",
	"element_at",
	"elt",
	"encode_as_bigint",
	"encode_as_int",
	"encode_as_largeint",
	"encode_as_smallint",
	"ends_with",
	"eq",
	"eq_for_null",
	"esquery",
	"exp",
	"extract_url_parameter",
	"field",
	"find_in_set",
	"first_value",
	"first_value_rewrite",
	"floor",
	"fmod",
	"fpow",
	"from_base64",
	"from_days",
	"from_microsecond",
	"from_millisecond",
	"from_second",
	"from_unixtime",
	"g",
	"ge",
	"get_json_bigint",
	"get_json_double",
	"get_json_int",
	"get_json_string",
	"greatest",
	"group_array",
	"group_array_intersect",
	"group_bit_and",
	"group_bit_or",
	"group_bit_xor",
	"group_bitmap_xor",
	"group_concat",
	"group_uniq_array",
	"grouping",
	"grouping_id",
	"gt",
	"hex",
	"hist",
	"histogram",
	"hll_cardinality",
	"hll_empty",
	"hll_from_base64",
	"hll_hash",
	"hll_raw_agg",
	"hll_to_base64",
	"hll_union",
	"hll_union_agg",
	"hour",
	"hour_ceil",
	"hour_floor",
	"hours_add",
	"hours_diff",
	"hours_sub",
	"if",
	"ifnull",
	"ignore",
	"in_iterate",
	"in_set_lookup",
	"inet6_aton",
	"inet6_ntoa",
	"inet_aton",
	"inet_ntoa",
	"initcap",
	"inner_product",
	"instr",
	"int_divide",
	"int_to_uuid",
	"intersect_count",
	"ipv4_cidr_to_range",
	"ipv4_num_to_string",
	"ipv4_string_to_num",
	"ipv4_string_to_num_or_default",
	"ipv4_string_to_num_or_null",
	"ipv4_to_ipv6",
	"ipv6_cidr_to_range",
	"ipv6_num_to_string",
	"ipv6_string_to_num",
	"ipv6_string_to_num_or_default",
	"ipv6_string_to_num_or_null",
	"is_ip_address_in_range",
	"is_ipv4_compat",
	"is_ipv4_mapped",
	"is_ipv4_string",
	"is_ipv6_string",
	"is_not_null_pred",
	"is_null_pred",
	"json_array",
	"json_contains",
	"json_exists_path",
	"json_extract",
	"json_extract_bigint",
	"json_extract_bool",
	"json_extract_double",
	"json_extract_int",
	"json_extract_isnull",
	"json_extract_largeint",
	"json_extract_string",
	"json_insert",
	"json_keys",
	"json_length",
	"json_object",
	"json_parse",
	"json_parse_error_to_invalid",
	"json_parse_error_to_null",
	"json_parse_error_to_value",
	"json_parse_notnull",
	"json_parse_notnull_error_to_invalid",
	"json_parse_notnull_error_to_value",
	"json_parse_nullable",
	"json_parse_nullable_error_to_invalid",
	"json_parse_nullable_error_to_null",
	"json_parse_nullable_error_to_value",
	"json_quote",
	"json_replace",
	"json_set",
	"json_type",
	"json_unquote",
	"json_valid",
	"jsonb_exists_path",
	"jsonb_extract",
	"jsonb_extract_bigint",
	"jsonb_extract_bool",
	"jsonb_extract_double",
	"jsonb_extract_int",
	"jsonb_extract_isnull",
	"jsonb_extract_largeint",
	"jsonb_extract_string",
	"jsonb_keys",
	"jsonb_parse",
	"jsonb_parse_error_to_invalid",
	"jsonb_parse_error_to_null",
	"jsonb_parse_error_to_value",
	"jsonb_parse_notnull",
	"jsonb_parse_notnull_error_to_invalid",
	"jsonb_parse_notnull_error_to_value",
	"jsonb_parse_nullable",
	"jsonb_parse_nullable_error_to_invalid",
	"jsonb_parse_nullable_error_to_null",
	"jsonb_parse_nullable_error_to_value",
	"jsonb_type",
	"l1_distance",
	"l2_distance",
	"lag",
	"last_day",
	"last_value",
	"lcase",
	"le",
	"lead",
	"least",
	"left",
	"length",
	"like",
	"ln",
	"localtime",
	"localtimestamp",
	"locate",
	"log",
	"log10",
	"log2",
	"lower",
	"lpad",
	"lt",
	"ltrim",
	"makedate",
	"map",
	"map_agg",
	"map_contains_key",
	"map_contains_value",
	"map_keys",
	"map_size",
	"map_values",
	"mask",
	"mask_first_n",
	"mask_last_n",
	"match_all",
	"match_any",
	"match_phrase",
	"match_phrase_edge",
	"match_phrase_prefix",
	"match_regexp",
	"max",
	"max_by",
	"md5",
	"md5sum",
	"microsecond",
	"microsecond_timestamp",
	"microseconds_add",
	"microseconds_diff",
	"microseconds_sub",
	"millisecond_timestamp",
	"milliseconds_add",
	"milliseconds_diff",
	"milliseconds_sub",
	"min",
	"min_by",
	"minute",
	"minute_ceil",
	"minute_floor",
	"minutes_add",
	"minutes_diff",
	"minutes_sub",
	"mod",
	"money_format",
	"month",
	"month_ceil",
	"month_floor",
	"monthname",
	"months_add",
	"months_diff",
	"months_sub",
	"multi_distinct_count",
	"multi_distinct_group_concat",
	"multi_distinct_sum",
	"multi_distinct_sum0",
	"multi_match",
	"multi_match_any",
	"multi_search_all_positions",
	"multiply",
	"murmur_hash3_32",
	"murmur_hash3_64",
	"named_struct",
	"ndv",
	"ne",
	"negative",
	"non_nullable",
	"not",
	"not_in_iterate",
	"not_in_set_lookup",
	"not_null_or_empty",
	"now",
	"ntile",
	"null_or_empty",
	"nullable",
	"nullif",
	"nvl",
	"or",
	"orthogonal_bitmap_expr_calculate",
	"orthogonal_bitmap_expr_calculate_count",
	"orthogonal_bitmap_intersect",
	"orthogonal_bitmap_intersect_count",
	"orthogonal_bitmap_union_count",
	"overlay",
	"parse_url",
	"percent_rank",
	"percentile",
	"percentile_approx",
	"percentile_array",
	"pi",
	"pmod",
	"positive",
	"pow",
	"power",
	"protocol",
	"quantile_percent",
	"quantile_state_empty",
	"quantile_union",
	"quarter",
	"quote",
	"radians",
	"rand",
	"random",
	"random_bytes",
	"rank",
	"regexp",
	"regexp_extract",
	"regexp_extract_all",
	"regexp_replace",
	"regexp_replace_one",
	"repeat",
	"replace",
	"retention",
	"reverse",
	"right",
	"rlike",
	"round",
	"round_bankers",
	"row_number",
	"rpad",
	"rtrim",
	"sec_to_time",
	"second",
	"second_ceil",
	"second_floor",
	"second_timestamp",
	"seconds_add",
	"seconds_diff",
	"seconds_sub",
	"sequence",
	"sequence_count",
	"sequence_match",
	"sha",
	"sha1",
	"sha2",
	"shuffle",
	"sign",
	"sin",
	"size",
	"sleep",
	"sm3",
	"sm3sum",
	"sm4_decrypt",
	"sm4_encrypt",
	"space",
	"split_by_regexp",
	"split_by_string",
	"split_part",
	"sqrt",
	"st_angle",
	"st_angle_sphere",
	"st_area_square_km",
	"st_area_square_meters",
	"st_asbinary",
	"st_astext",
	"st_aswkt",
	"st_azimuth",
	"st_circle",
	"st_contains",
	"st_distance_sphere",
	"st_geometryfromtext",
	"st_geometryfromwkb",
	"st_geomfromtext",
	"st_geomfromwkb",
	"st_linefromtext",
	"st_linestringfromtext",
	"st_point",
	"st_polyfromtext",
	"st_polygon",
	"st_polygonfromtext",
	"st_x",
	"st_y",
	"starts_with",
	"stddev",
	"stddev_pop",
	"stddev_samp",
	"str_to_date",
	"strcmp",
	"strleft",
	"strright",
	"struct",
	"struct_element",
	"sub_bitmap",
	"sub_replace",
	"subdate",
	"substr",
	"substring",
	"substring_index",
	"subtract",
	"sum",
	"sum0",
	"sum_distinct",
	"tan",
	"tanh",
	"time_to_sec",
	"timediff",
	"timestamp",
	"to_base64",
	"to_bitmap",
	"to_bitmap_with_check",
	"to_date",
	"to_datev2",
	"to_days",
	"to_ipv4",
	"to_ipv4_or_default",
	"to_ipv4_or_null",
	"to_ipv6",
	"to_ipv6_or_default",
	"to_ipv6_or_null",
	"to_monday",
	"to_quantile_state",
	"tokenize",
	"topn",
	"topn_array",
	"topn_weighted",
	"trim",
	"truncate",
	"ucase",
	"unhex",
	"unix_timestamp",
	"upper",
	"url_decode",
	"utc_timestamp",
	"uuid",
	"uuid_numeric",
	"uuid_to_int",
	"var_pop",
	"var_samp",
	"variance",
	"variance_pop",
	"variance_samp",
	"version",
	"week",
	"week_ceil",
	"week_floor",
	"weekday",
	"weekofyear",
	"weeks_add",
	"weeks_diff",
	"weeks_sub",
	"width_bucket",
	"window_funnel",
	"xxhash_32",
	"xxhash_64",
	"year",
	"year_ceil",
	"year_floor",
	"years_add",
	"years_diff",
	"years_sub",
	"yearweek",
}
