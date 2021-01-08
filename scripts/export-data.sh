# Not an Actual Script
docker exec -it coinprice-scraper_influxdb_1 bash
influx config create --config-name admin \
  --host-url http://localhost:8086 \
  --org coinprice \
  --token <INFLUXDB TOKEN> \
  --active
influx query 'from(bucket: "candle") |> range(start: -10y) |> filter(fn: (r) => r["_measurement"] == "BTC") |> filter(fn: (r) => r["_field"] == "close" or r["_field"] == "high" or r["_field"] == "open" or r["_field"] == "low" or r["_field"] == "volume") |> filter(fn: (r) => r["exchange"] == "BINANCE") |> filter(fn: (r) => r["quote"] == "USDT") |> aggregateWindow(every: 1m, fn: mean, createEmpty: false) |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")' --raw > btc.csv
influx query 'from(bucket: "candle") |> range(start: -10y) |> filter(fn: (r) => r["_measurement"] == "ETH") |> filter(fn: (r) => r["_field"] == "close" or r["_field"] == "high" or r["_field"] == "open" or r["_field"] == "low" or r["_field"] == "volume") |> filter(fn: (r) => r["exchange"] == "BINANCE") |> filter(fn: (r) => r["quote"] == "USDT") |> aggregateWindow(every: 1m, fn: mean, createEmpty: false) |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")' --raw > eth.csv
influx query 'from(bucket: "candle") |> range(start: -10y) |> filter(fn: (r) => r["_measurement"] == "LTC") |> filter(fn: (r) => r["_field"] == "close" or r["_field"] == "high" or r["_field"] == "open" or r["_field"] == "low" or r["_field"] == "volume") |> filter(fn: (r) => r["exchange"] == "BINANCE") |> filter(fn: (r) => r["quote"] == "USDT") |> aggregateWindow(every: 1m, fn: mean, createEmpty: false) |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")' --raw > ltc.csv
influx query 'from(bucket: "candle") |> range(start: -10y) |> filter(fn: (r) => r["_measurement"] == "XRP") |> filter(fn: (r) => r["_field"] == "close" or r["_field"] == "high" or r["_field"] == "open" or r["_field"] == "low" or r["_field"] == "volume") |> filter(fn: (r) => r["exchange"] == "BINANCE") |> filter(fn: (r) => r["quote"] == "USDT") |> aggregateWindow(every: 1m, fn: mean, createEmpty: false) |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")' --raw > xrp.csv
docker cp d08131acffa5:/btc.csv /home/zahin/Dev/mochahub/coinprice-scraper/influx_export/
docker cp d08131acffa5:/eth.csv /home/zahin/Dev/mochahub/coinprice-scraper/influx_export/
docker cp d08131acffa5:/ltc.csv /home/zahin/Dev/mochahub/coinprice-scraper/influx_export/
docker cp d08131acffa5:/xrp.csv /home/zahin/Dev/mochahub/coinprice-scraper/influx_export/