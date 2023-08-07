# Mongo logs analyser

> A tool that takes a file with mongo logs and analyses, giving data like most used collections or slow queries.

---

## Install / build

```bash
make build
```

## Stats

Having an input file, calculates and retrieves statistics of that log.

```bash
# with stdout
./bin/mongo_logs_analyser stats -i "<file_path>"

# with output file
./bin/mongo_logs_analyser stats -i "<file_path>" -o "<output_file_path>"

# minimum ms duration per request
./bin/mongo_logs_analyser stats -i "<file_path>" -t 1000

# filters out a regex pattern (from valid json array) out of the command log
./bin/mongo_logs_analyser stats -i "<file_path>" -f '["fooId", "find"]'

```
