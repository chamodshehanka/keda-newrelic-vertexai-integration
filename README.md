# Vertex AI Traffic Predictor

## Prerequisites

- NewRelic Account
- NewRelic API Key
- GKE Cluster


## Setup

1. Create a new project in the Google Cloud Console.
2. Enable the Vertex AI API.
3. Create a service account and download the JSON key.
4. Create a GKE cluster.
5. Create a new Relic account.
6. Create a new Relic API key.
7. Create a new Relic application.
8. Install the New Relic Infrastructure agent on the GKE cluster.
9. Deploy the New Relic Infrastructure agent on the GKE cluster.
10. Configure New Relic for your application.


### New Relic Queries 

```
SELECT average(cpuRequestedCores) AS 'CPU Requests (Cores)', average(cpuLimitCores) AS 'CPU Limits (Cores)',
       average(memoryRequestedBytes) AS 'Memory Requests (Bytes)', average(memoryLimitBytes) AS 'Memory Limits (Bytes)'
FROM K8sContainerSample
WHERE clusterName = 'gke_t-pointer-442615-u8_us-central1-a_devfest-2024-demo' AND namespace = 'default' AND label.app = 'devfest-app-2024'
FACET podName
SINCE 1 hour ago
```

Request Rate in last 5 minutes
```
SELECT rate(count(*), 1 minute) 
FROM Transaction 
WHERE appName = 'devfest-app-2024' 
SINCE 5 minute ago
```