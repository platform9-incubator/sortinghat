
# Sorting Hat

Sorting Hat is an application which turns logs into useful alerts. It does so by categorizing logs statements into different buckets. Kind of like fingerprinting a log statement.

An operator (typically a support engineer) will mark each bucket as useful or not-useful. 

Each incoming log statement is categorized and put into an existing bucket or a new bucket is created. A dashboard helps the operator perform these categorizations and also show him/her the alerts.

The service is in alpha-state and looking for maintainers.

## Services

There are two major services:
* __logingest__: A log ingestion service which listens on a webserver. Performs the pattern matching as per the buckets.
* __sortinghat-ui__: Another service that powers the web interface and helps querying and categorization. The UI is heavily influenced by Rollbar's design.

## Deployment

### Kuberentes
The service(s) can be deployed on any Kubernetes cluster. The current deployment assumes use of service proxy. 

__TODO__: One major asks has been to start using servicetype==loadbalancer for k8s deployments.


## Architecture

The architecture is very simple pipeline (and there are many many ways to improve it, including many requests to start using Spark and perhaps ML libraries for categorization).

### Pipeline

LogAggregation Service --> filter by keywords (like ERROR) --> SortingHat-logingest (parse and sort into buckets) --> Monog DB

SortingHat UI --> SortingHat UI backend --> Mongo DB

## How to customize

The project is in early stages and needs a lot of customization in order to be useful for someone. Here are two that are specific to each log aggregator and format (we wish there was time to fix this to provide a better out of the box example).

* log_parse: Parsing of the log file is based on regular expression and as you can imagine each log format needs code change.
* parse_xxxx_payload: This is the code that can act as a webhook sink for filtered log statement that you want to analyze. Needless to say that each webhook imoplementation is different and hence the transformation.

For more information please contact 

opensource at Platform9.com