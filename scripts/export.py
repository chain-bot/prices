"""
Export InfluxDb bucket w/ pagination
"""

import os
import time
from datetime import datetime
from influxdb_client import InfluxDBClient, Point, WritePrecision
import pandas as pd
from dotenv import load_dotenv

load_dotenv(dotenv_path='../.env')

if __name__ == "__main__":
    token = os.getenv('INFLUXDB_ADMIN_USER_TOKEN')
    org = os.getenv('INFLUXDB_ORG')
    bucket = os.getenv('INFLUXDB_BUCKET_CANDLE')
    host = os.getenv('INFLUXDB_HOST')
    port = os.getenv('INFLUXDB_PORT')
    url = f"http://{host}:{port}"

    client = InfluxDBClient(url=url, token=token)
    query_api = client.query_api()

    end = int(time.time()) + 60  # Add extra minute to get last time
    # Fri Jan 01 2016 05:00:00 GMT+0000
    start = 1451624400
    # Get 30 days at a time
    page_size = 30 * 24 * 60 * 60
    res = None
    while start < end:
        new_end = min(start + page_size, end)
        print(f"Querying {bucket} from {start} to {new_end}")

        query = f'''
            from(bucket: "candle")
          |> range(start: {start}, stop: {new_end})'''
        df = query_api.query_data_frame(query, org=org)
        if df.shape[0] > 0:
            df.drop(columns=['result', 'table', '_start', '_stop'], inplace=True)
            df.rename(columns={"_time": "time", "_field": "field", "_value": "value", "_measurement": "base"},
                      inplace=True)
        print(list(df.columns))
        if res is None and df.shape[0] > 0:
            res = df
        else:
            res = pd.concat([res, df])
        start = new_end
    print(f"rows: {res.shape[0]}")
    print(f"columns: {res.shape[1]}")

    compression_opts = dict(method='zip',
                            archive_name='out.csv')
    res.to_csv('out.zip', index=False,
               compression=compression_opts)
