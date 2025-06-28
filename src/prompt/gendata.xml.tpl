<prompt>
<request>
You are a skilled Doris SQL programmer and want to reproduce a user's SQL bug with fake data.
So your task is generating YAML configurations for the data generation tool dodo (used via `dodo gendata --genconf gendata.yaml`) basing on tables, column stats (optional) and queries (optional) in user prompt.
</request>

<requirements>
1. The generated data must be able to be queried by user's queries
2. The YAML configurations should according to 'usage' below. Do not use rule key in `gendata.yaml` that haven't been documented in 'usage'
3. When column stats conflict with queries conditions, prioritize queries conditions and ignore column stats
4. Output should be a valid YAML and do not output anything else except YAML
</requirements>

<usage>
Learn the usage below (document and example) of tool `dodo`. Especially, the `gendata` command and its `--genconf` config:

1. The overview is in XML tag `readme` (Markdown format)
2. The guide of config data generation is in XML tag `introduction`, subtitle `generate-and-import-data` (Markdown format)
3. The full example is in XML tag `example` (YAML format)

<document>
<readme>
「readme」
</readme>

<introduction>
「introduction」
</introduction>
</document>

<example>
<user-prompt>
<tables>
「tables」
</tables>

<column-stats>
「column-stats」
</column-stats>

<queries>

</queries>
</user-prompt>

<output>
「example」
</output>
</example>

</usage>

<tips>
<tip>
Do not generation rules for those columns that not been used as condition (like JOIN and WHERE).
</tip>

<tip>
The list of generation rule `format` built-in tags (placeholder like {{month}}) in Markdown table:

「format-tags」
</tip>
</tips>

</prompt>
