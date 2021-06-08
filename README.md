Rack planner tool
=======
**Description:** Rack planner optimization service aims to deliver greater experience to run optimization algorithm against custom input.

## High level overview

Optimization engine consists of three microservices:

* UI service consists ot NGINX and a simple React-based application
* API service is responsible for backend
* Optimization service is responsible for running optimization script

Workflow:
1. React-based frontend enables the user to upload input files and inspect optimization result.
2. API service compresses input files and uploads the archive to the S3 bucket with `input_` prefix and timestamp and generates a REST request to the optimization service.
3. Optimization service downloads the archive, decompresses it, runs optimization script, compresses the files produced by the script, uploads the archive to the S3 bucket with `output_` prefix and timestamp, and sends a response to the API service.
4. API service sends a successful response to the UI application with the path and location of the compressed file.

Each service can be scaled independently, one possible structure is depicted on the image.

![Internal Project Dataflow](/assets/ProjectStructure.png)

### Prerequisites

System requires a file storage which can be presented as an S3 bucket or a local folder. The `type` parameter in the `storage` category can take `s3` or `local` value. The `bucket` parameter must point to a valid already created folder or s3 bucket.


Environment variables must be exported or edited in docker-compose.yml file. 

Running with local storage:

```
export STORAGE_TYPE=local
export STORAGE_BUCKET=/bucket
```

In order to run with s3 storage one must provide a valid bucket S3 bucket name and AWS credentials:

```
export STORAGE_TYPE=s3
export STORAGE_BUCKET=s3-bucket-name
export AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
export AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
```

Do not store AWS credentials in configuration files! Set environment variables with essential information before running the containers.

### CI/CD

### Running locally with docker

### Setting logging level

Valid options for log level: trace, debug, info, warning, error, fatal and panic.
```
export LOG_LEVEL=debug
```

### Running locally

```
$ make run_optimization
$ make run_api
```

### Running tests and linters

```
$ make test
```

### Vendoring

All dependencies are vendorred and managed via Go modules
