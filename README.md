# Automating Microservice Scaling with KEDA, Google Vertex AI, and New Relic

This repository demonstrates how to set up automated scaling of Kubernetes microservices using KEDA, powered by Google Vertex AI predictions and New Relic metrics.

---

## Prerequisites

Ensure the following are set up before proceeding:

1. **Kubernetes Cluster**
   - A Kubernetes cluster with at least 3 worker nodes (e.g., [EKS](https://aws.amazon.com/eks/), [GKE](https://cloud.google.com/kubernetes-engine/), or [Minikube](https://minikube.sigs.k8s.io/docs/)).
   - `kubectl` configured to access your cluster.
2. **KEDA**
   - Install KEDA in your cluster:
     ```bash
        helm repo add kedacore https://kedacore.github.io/charts
        helm repo update
        helm install keda kedacore/keda --namespace keda --create-namespace
     ```
3. **Google Cloud Account**
   - Access to Google Cloud with Vertex AI APIs enabled.
   - Service account key for authentication.
4. **New Relic Account**
   - A New Relic account with access to NRQL queries.
   - Query Key and Account ID from your New Relic account.
   - New Relic Installed in your Kubernetes cluster.
5. **Docker**
   - Docker installed for containerizing and deploying your app.
6. **Helm**
   - Helm installed if deploying additional tools (like monitoring solutions).

---

## Cluster Setup

### 1. Set up Kubernetes Cluster

- Create a Kubernetes cluster using your preferred platform (EKS/GKE/Minikube).
- Ensure your cluster has sufficient resources for testing scaling.

### 2. Install KEDA

- Deploy KEDA in the cluster:
  ```bash
  kubectl apply -f https://github.com/kedacore/keda/releases/download/v2.9.0/keda-2.9.0.yaml
  ```

### 3. Deploy Application

- Build and deploy your application:
  ```bash
  docker build -t your-app-name .
  kubectl apply -f k8s/deployment.yaml
  ```

### 4. Verify Deployment

- Check that your app is running:
  ```bash
  kubectl get pods
  ```

---

## KEDA Configuration

### 1. Create a New Relic Secret

- Encode your New Relic Account ID and Query Key as base64:
  ```bash
  echo -n "YOUR_ACCOUNT_ID" | base64
  echo -n "YOUR_QUERY_KEY" | base64
  ```
- Create the secret:
  ```yaml
  apiVersion: v1
  data:
    accountId: <base64_encoded_account_id>
    queryKey: <base64_encoded_query_key>
  kind: Secret
  metadata:
    name: newrelic-secret
    namespace: default
  ```

### 2. Configure TriggerAuthentication

```yaml
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: newrelic-auth
  namespace: default
spec:
  secretTargetRef:
    - parameter: account
      name: newrelic-secret
      key: accountId
    - parameter: queryKey
      name: newrelic-secret
      key: queryKey
```

### 3. Create the ScaledObject

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: devfest-app-scaledobject
  namespace: default
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: devfest-app-2024
  minReplicaCount: 1
  maxReplicaCount: 10
  pollingInterval: 30
  cooldownPeriod: 300
  triggers:
    - type: new-relic
      metadata:
        account: "6225327"
        queryKey: "NRAK-XXXXXX"
        region: "US"
        nrql: |
          SELECT rate(count(*), 1 minute) 
          FROM Transaction 
          WHERE appName = 'devfest-app-2024' 
          SINCE 5 minute ago
        threshold: "5"
      authenticationRef:
        name: newrelic-auth
```

---

## Testing the Setup

1. **Simulate Traffic**

   - Use tools like `k6` or `wrk` to send traffic to your application.
   - Verify if scaling happens as per the defined thresholds.

2. **Monitor Metrics**

   - Check the logs for KEDA:
     ```bash
     kubectl logs -l app=keda-operator -n keda
     ```
   - Monitor New Relic to ensure the query returns expected results.

3. **Verify Scaling**
   - Confirm that the replicas of your deployment scale up/down:
     ```bash
     kubectl get hpa
     kubectl get pods
     ```

---

## Cleanup

To clean up all resources:

```bash
kubectl delete -f k8s/
kubectl delete -f https://github.com/kedacore/keda/releases/download/v2.9.0/keda-2.9.0.yaml
```

---

## References

- [KEDA Documentation](https://keda.sh/docs/)
- [New Relic NRQL](https://docs.newrelic.com/docs/query-data/nrql-new-relic-query-language/)
- [Google Vertex AI](https://cloud.google.com/vertex-ai/)

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
