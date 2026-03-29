#!/bin/sh

API_URL="${API_BASE:-http://ticketfair-app:8000/api/v1}"

sed -i "s|http://localhost:8000/api/v1|${API_URL}|g" /usr/share/nginx/html/index.html

echo "Dashboard starting — API: ${API_URL}"

nginx -g "daemon off;"