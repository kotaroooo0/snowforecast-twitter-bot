curl -H "Content-Type: application/json" -XPUT "localhost:9200/snow_resorts_v1" --data-binary "@settings.json"
curl -H "Content-Type: application/json" -XPOST "localhost:9200/snow_resorts_v1/_bulk" --data-binary "@snow_resorts.json"
