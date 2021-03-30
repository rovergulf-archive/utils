#!/usr/bin/env bash

set -e

for testPath in "./colors" "./pgxs" "./useragent"; do
  go test $testPath
done
