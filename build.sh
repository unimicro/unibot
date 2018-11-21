#! /bin/bash
set -e

OUT_FILE="unibot.new"

echo '≫ Buidling "'${OUT_FILE}'" executable...'
GOOS=linux go build -o "$OUT_FILE"
