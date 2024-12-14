#!/bin/bash

echo "Building subnet..."

echo "Please adding the following subnet info..."
echo "Network name: LIFENetwork"
echo "Chain ID: 9999"
echo "Currency Symbol: LIFE"

avalanche subnet deploy LIFEnetwork --avalanchego-version v1.11.13
