#!/bin/bash
#
# This is an example of how to interact with the MCP server.
#
# Usage:
# 1. Make sure the MCP server is running (docker compose up).
# 2. Run this script from your terminal: ./example_usage.sh

echo "Sending request to the 'write_essay' pattern..."
echo ""

curl -X POST -H "Content-Type: application/json" \
-d '{"input": "Write a short essay about the future of space exploration."}' \
http://localhost:3333/write_essay

echo ""
echo ""
echo "Sending request to the 'summarize' pattern..."
echo ""

curl -X POST -H "Content-Type: application/json" \
-d '{"input": "The solar system consists of the Sun and the astronomical objects bound to it by gravity. Of the eight planets, the four smaller, inner planets, Mercury, Venus, Earth and Mars, are terrestrial planets, being primarily composed of rock and metal. The four outer planets are giant planets, being substantially more massive than the terrestrials."}' \
http://localhost:3333/summarize

echo ""
