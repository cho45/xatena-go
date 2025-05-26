#!/bin/sh

set -e

echo

perl bench/bench.pl sample.txt
go run bench/bench.go sample.txt
perl bench/bench.pl sample.txt
go run bench/bench.go sample.txt
perl bench/bench.pl sample.txt
go run bench/bench.go sample.txt
