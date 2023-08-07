# Mongo logs analyser

> A tool that takes a file with mongo logs and analyses, giving data like most used collections or slow queries.

---

## Install / build

```bash
make build
```

## Stats

Having an input file, calculates and retrieves statistics of that log and saves them on an output file.
Output file is optional. In the case there isn't an output, it will log to the stdout.

```bash
./bin/mongo_logs_analyser stats -i "<file_path>" -o "<output_file_path>"
```
