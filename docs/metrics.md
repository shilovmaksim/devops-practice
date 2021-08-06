Metrics
=======
**Description:** Optimization service can be deployed with Prometheus to gather metrics.

## Deployment

Use 'docker-compose.metrics.yml' file to deploy locally.

```
docker-compose -f docker-compose.yml -f docker-compose.metrics.yml up --remove-orphans
```
or
```
make docker_metrics_run
```

Docker volumes are used to store data in local deployment.

## Stack

Prometheus is uses to gather metrics. Prometheus configuration file is located in the `.data/prometheus/config` folder. It contains a list of targets that provide metrics.


Grafana is used for visualization. Admin's password and other environment variables are located in the `.data/grafana/config.monitoring` file. Several predefined dashboards and data sources configuration files are located in the `.data/grafana/provisioning` folder.

'node-exporter', 'nginx-exporter', 'cadvisor' are used to gather infrastructure data.


## Usage

`http://localhost:3000` leads to Grafana's main page with several predefined dashboards. Use username `admin` and admin's password provided in the configuration file.

`http://localhost:9090` leads to Prometheus UI which can be used for debugging and observing a target's status.

There is a dashboard used to provide information about containers and host available resources:
![Grafana Containers Info](/assets/GrafanaLoad.png)

## Metrics

API Service and Optimization Service provide general http request metrics presented as HistogramVec, which is bundles a set of Histograms that all share the same Desc, but have different values for their variable labels.
The name of the metric is 'http_request_duration_seconds', the labels are '"handler", "method", "code"'.

Use `{job="job_name"}` to select the metric of the desired service. The *job_name* corresponds to the job name provided in the configuration file for Prometheus.

Example:
```
http_request_duration_seconds_bucket{job="optimization_service"}
```
Histogram example:
```
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{job="optimization_service"}[5m])) by (le))
```

Rate example:
```
rate(ui_upload_duration_time_seconds_sum[5m])/rate(ui_upload_duration_time_seconds_count[5m])
```

API Service also provides a number of additional histograms. Default bucket values are used, because execution times should fall into common http request times:

* Files upload time - `ui_upload_duration_time_seconds`
* Files compression time - `compression_time_seconds`
* Files upload to the storage time - `storage_upload_duration_time_seconds`
* Optimization request time - `optimization_request_duration_time_seconds`

Optimization service and API service dashboards:
![Grafana Optimization Service](/assets/GrafanaOptimization.png)

![Grafana API Service Info](/assets/GrafanaAPIService.png)

## Logs
Logs from all the running containers are collected using Loki. In order to access them use Grafana and a predefined dashboard *Logs*. 

The picture below shows a sample view of logs representation in grafana:
![Grafana Loki Dashboard](/assets/LokiLogsGrafana.png)