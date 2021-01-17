#!/bin/sh

set -e

until curl -s localhost:9200 > /dev/null; do
  >&2 echo "Elasticsearch is unavailable - sleeping"
  sleep 3
done
>&2 echo "Elasticsearch is up - executing command"

curl -H "Content-Type: application/json" -XPUT "localhost:9200/snow_resorts_v1" --data-binary "@docker-compose/elasticsearch/settings.json"
curl -X POST "localhost:9200/_aliases?pretty" -H 'Content-Type: application/json' --data-binary "@docker-compose/elasticsearch/alias.json"
curl -H "Content-Type: application/json" -XPOST "localhost:9200/snow_resorts_v1/_bulk" --data-binary "@docker-compose/elasticsearch/snow_resorts.json"
