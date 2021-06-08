# API service of the optimization engine

REST API service capable of receiving incoming files for optimization script, storing them in the S3 or local bucket 
and performing calls to the optimization service. <br/>

## Using the API
| url | method | params | response code | response body | description |  
|-----------|-----------|-----------|-----------|-----------|-----------|
| /api/v1/health | GET | - |  200 |```{"health": true}```| Success health check |
| /api/v1/upload | POST | files |  200 |```{"filename":"file.csv","location":"http://s3_location/file.csv","etag":"md5_like_s3_etag"}```| Success optimize run |
| /api/v1/upload | POST | files |  400 |```{"text":"script execution error"}```| Failed optimize run |

----

## Testing with `curl`
Health check:
```
  curl -X GET  localhost:8080/api/v1/health
```
Upload and get optimization result:
```
  curl -F 'file=@/path/file1.csv' -F 'file=@/path/file2.csv' http://localhost:8090/api/v1/upload

```

## Logging
Incoming requests are logged in the Apache [Common Log Format](http://httpd.apache.org/docs/2.2/logs.html#common) and can be grepped in `{server_name}/log` folder.
