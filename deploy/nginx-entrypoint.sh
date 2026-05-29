#!/bin/sh
set -e

# Generate self-signed SSL certificate if it doesn't exist
if [ ! -f /etc/nginx/ssl/selfsigned.crt ]; then
    echo "Generating self-signed SSL certificate..."
    apk add --no-cache openssl > /dev/null 2>&1
    mkdir -p /etc/nginx/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout /etc/nginx/ssl/selfsigned.key \
        -out /etc/nginx/ssl/selfsigned.crt \
        -subj "/CN=localhost/O=GMAO/C=FR"
    echo "Self-signed SSL certificate generated."
fi

# Start Nginx
exec nginx -g 'daemon off;'
