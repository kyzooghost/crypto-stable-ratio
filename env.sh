#!/usr/bin/env bash
# Use 'source .env.sh', not './env.sh' or else environment variables set here are not available - https://askubuntu.com/a/53179
# Need language-agnostic way of defining and sharing environment variables
# Note that this file is not for storing private credentials like an API key
set -a

# These variables are injected into shell environment variables
BUCKET_DB_NAME=crypto-stable-ratio-db
STABLE_MCAP_ENDPOINT=https://api.btctools.io/api/marketcap-stable-chart?period=1y
TOTAL_MCAP_ENDPOINT=https://api.btctools.io/api/marketcap-total-chart?period=1y

set +a

echo "ENV VARIABLES SET"