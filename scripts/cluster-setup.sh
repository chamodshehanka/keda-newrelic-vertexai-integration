#!/bin/bash

# Check if yq is installed
if ! command -v yq &> /dev/null
then
    echo "yq could not be found, please install it."
    exit 1
fi

# Load environment variables from configs/config.yaml file
echo "Loading environment variables from configs/config.yaml file..."
NEW_RELIC_LICENSE_KEY=$(yq e '.newrelic.licenseKey' configs/config.yaml)
CLUSTER_NAME=$(yq e '.newrelic.clusterName' configs/config.yaml)


# Check if the environment variables are set
if [ -z "$NEW_RELIC_LICENSE_KEY" ]; then
    echo "NEW_RELIC_LICENSE_KEY is not set. Please check your configs/config.yaml file."
    exit 1
fi
echo "NEW_RELIC_LICENSE_KEY is set."

if [ -z "$CLUSTER_NAME" ]; then
    echo "CLUSTER_NAME is not set. Please check your configs/config.yaml file."
    exit 1
fi
echo "CLUSTER_NAME is set."

# Install KEDA
# echo "Installing KEDA..."
# helm repo add kedacore https://kedacore.github.io/charts
# helm repo update
# helm install keda kedacore/keda --namespace keda --create-namespace

# Install New Relic
echo "Installing New Relic..."
# helm repo add newrelic https://helm-charts.newrelic.com
# helm repo update
helm install newrelic newrelic/nri-bundle \
    --set global.licenseKey=$NEW_RELIC_LICENSE_KEY \
    --set global.cluster=$CLUSTER_NAME \
    --namespace newrelic \
    --create-namespace