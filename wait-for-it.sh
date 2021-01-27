#!/bin/ash
# wait-for-it.sh

set -e

until nc -zn elasticsearch 9200; do
  >&2 echo "Elasticsearch is unavailable - sleeping"
  sleep 3
done
>&2 echo "Elasticsearch is up - executing command"

exec $@
