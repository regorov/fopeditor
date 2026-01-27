#!/bin/sh
set -euo pipefail

if [ "${ENABLE_TLS:-false}" = "true" ]; then
  : "${SERVER_NAME:?SERVER_NAME must be set when ENABLE_TLS=true}"
  : "${CERTBOT_EMAIL:?CERTBOT_EMAIL must be set when ENABLE_TLS=true}"

  export SERVER_NAME
  envsubst '${SERVER_NAME}' < /etc/nginx/nginx.ssl.conf.template > /etc/nginx/nginx.conf

  mkdir -p /var/www/certbot /etc/letsencrypt

  if [ ! -f "/etc/letsencrypt/live/${SERVER_NAME}/fullchain.pem" ]; then
    certbot certonly --non-interactive --agree-tos --email "${CERTBOT_EMAIL}" --webroot -w /var/www/certbot -d "${SERVER_NAME}"
  fi

  renew_hours=${CERTBOT_RENEW_INTERVAL_HOURS:-12}
  (
    while true; do
      sleep $((renew_hours * 3600))
      certbot renew --webroot -w /var/www/certbot --quiet && nginx -s reload || true
    done
  ) &
else
  echo "[entrypoint] TLS disabled, using default nginx configuration"
fi

exec nginx -g "daemon off;"
