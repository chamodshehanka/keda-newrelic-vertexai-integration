import os
import requests
from flask import Flask, jsonify
from google.cloud import aiplatform
from flask_cors import CORS

# Initialize the Flask app
app = Flask(__name__)
CORS(app)  # Enable Cross-Origin Resource Sharing (CORS) if needed

# New Relic API Configuration
NEW_RELIC_API_KEY = os.getenv("NEW_RELIC_API_KEY")
NEW_RELIC_ACCOUNT_ID = os.getenv("NEW_RELIC_ACCOUNT_ID")
APP_NAME = os.getenv("APP_NAME")
NEW_RELIC_QUERY_URL = (
    f"https://insights-api.newrelic.com/v1/accounts/{NEW_RELIC_ACCOUNT_ID}/query"
)

# Vertex AI Configuration
aiplatform.init(project="your-project-id", location="your-location")
ENDPOINT_ID = "your-endpoint-id"  # Replace with your Vertex AI endpoint ID
endpoint = aiplatform.Endpoint(endpoint_id=ENDPOINT_ID)


# Function to fetch metrics from New Relic
def fetch_new_relic_metrics():
    """
    Fetch metrics from New Relic.
    We're querying the rate of transactions for the past 5 minutes.
    """
    nrql_query = "SELECT rate(count(*), 1 minute) FROM Transaction WHERE appName = '{}' SINCE 5 minutes ago".format(
        APP_NAME
    )
    headers = {
        "Accept": "application/json",
        "X-Query-Key": NEW_RELIC_API_KEY,
    }
    params = {"nrql": nrql_query}
    response = requests.get(NEW_RELIC_QUERY_URL, headers=headers, params=params)

    if response.status_code == 200:
        data = response.json()
        # Assuming the rate is under 'results' -> 'events' -> 'count'
        try:
            transaction_rate = data["results"][0]["count"]
            return transaction_rate
        except KeyError:
            raise Exception("Invalid response format from New Relic")
    else:
        raise Exception(
            f"Error fetching data from New Relic: {response.status_code} {response.text}"
        )


# Function to make prediction using Vertex AI
def predict_traffic(input_data):
    """
    Use Vertex AI model to predict traffic based on New Relic metrics.
    """
    try:
        instances = [{"metric": input_data}]
        prediction = endpoint.predict(instances=instances)
        predicted_traffic = prediction.predictions[
            0
        ]  # Assuming prediction is a scalar value
        return predicted_traffic
    except Exception as e:
        raise Exception(f"Error predicting traffic: {str(e)}")


# Endpoint for fetching the prediction (for KEDA)
@app.route("/get-prediction", methods=["GET"])
def get_prediction():
    """
    This endpoint will be used by KEDA to get the predicted traffic value.
    It fetches New Relic metrics, generates a prediction via Vertex AI, and returns it.
    """
    try:
        # Step 1: Fetch New Relic Metrics
        transaction_rate = fetch_new_relic_metrics()

        # Step 2: Generate Prediction using Vertex AI
        predicted_traffic = predict_traffic(transaction_rate)

        # Step 3: Return prediction as a JSON response
        return jsonify({"predicted_traffic": predicted_traffic}), 200

    except Exception as e:
        return jsonify({"error": str(e)}), 500


if __name__ == "__main__":
    # Run the Flask app
    app.run(host="0.0.0.0", port=8080)
