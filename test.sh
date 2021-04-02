#!/usr/bin/env bash

set -e

for testPath in "./colors" "./pgxs" "./useragent" "./ipaddr"; do
  go test $testPath
done
