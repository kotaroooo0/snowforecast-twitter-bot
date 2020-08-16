curl -H "Content-Type: application/json" -XPOST "localhost:9200/snow_resorts/_bulk?pretty&refresh" --data-binary "@snow_resorts.json"
