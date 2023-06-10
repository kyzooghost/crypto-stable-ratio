#!/bin/bash
# Need language-agnostic way of defining and sharing environment variables
# Note that this file is not for storing private credentials like an API key
set -a

# These variables are injected into shell environment variables
STABLE_MCAP_ENDPOINT=https://api.btctools.io/api/marketcap-stable-chart?period=1y
TOTAL_MCAP_ENDPOINT=https://api.btctools.io/api/marketcap-total-chart?period=1y
CSV_S3_PATH=crypto-stable-ratio-graph/crypto-stable-ratios

set +a