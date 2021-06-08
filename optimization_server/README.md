# Optimization service

REST API service capable to execute optimization logic via exposed endpoints. <br/>
Under the hood it executes secret _python_&#8471; script which is doing all the magic.

## Using the API
| url | method | params | response code | response body | description |  
|-----------|-----------|-----------|-----------|-----------|-----------|
| /api/v1/health | GET | - |  200 |```{"health": true}```| Success health check |
| /api/v1/optimize | GET | ```{"args":["file1.csv","file2.csv"]}```|  200 |```{"exitCode": 0,"shellOutput": "","scriptOutput": "","executionTime": 253}```| Success optimize run |
| /api/v1/optimize | POST | ```{"args":[""]}``` |  400 |```{"text":"error validating json body"}```| Failed optimize run |

----

## Testing with `curl`
Health check:
```
  curl -X GET  localhost:8080/api/v1/health
```
Optimize run:
```
  curl -X GET -H "Content-Type: application/json" -d '{"args":["file1.csv","file2.csv"]}' localhost:8080/api/v1/optimize
```

## Logging
Incoming requests are logged in the Apache [Common Log Format](http://httpd.apache.org/docs/2.2/logs.html#common) and can be grepped in `{server_name}/log` folder.