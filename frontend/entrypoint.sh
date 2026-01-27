#!/bin/sh
set -euo pipefail

HTTP_CONF="/etc/nginx/nginx.http.conf"
TLS_TEMPLATE="/etc/nginx/nginx.ssl.conf.template"
TLS_CONF="/etc/nginx/nginx.tls.conf"
CERT_PATH_BASE="/etc/letsencrypt/live"

if [ "${ENABLE_TLS:-false}" = "true" ]; then
  : "${SERVER_NAME:?SERVER_NAME must be set when ENABLE_TLS=true}"
  : "${CERTBOT_EMAIL:?CERTBOT_EMAIL must be set when ENABLE_TLS=true}"

  export SERVER_NAME
  envsubst '${SERVER_NAME}' < "${TLS_TEMPLATE}" > "${TLS_CONF}"

  mkdir -p /var/www/certbot /etc/letsencrypt

  cert_dir="${CERT_PATH_BASE}/${SERVER_NAME}"
  if [ ! -f "${cert_dir}/fullchain.pem" ]; then
    echo "[entrypoint] Bootstrapping TLS certificates for ${SERVER_NAME}"
    cp "${HTTP_CONF}" /etc/nginx/nginx.conf
    nginx -c "${HTTP_CONF}"
    if certbot certonly --non-interactive --agree-tos --email "${CERTBOT_EMAIL}" --webroot -w /var/www/certbot -d "${SERVER_NAME}"; then
      echo "[entrypoint] Certificate request completed"
    else
      echo "[entrypoint] WARNING: Initial certificate request failed; serving HTTP only and will retry later."
    fi
    nginx -c "${HTTP_CONF}" -s stop || true
  fi

  if [ -f "${cert_dir}/fullchain.pem" ]; then
    cp "${TLS_CONF}" /etc/nginx/nginx.conf
  else
    cp "${HTTP_CONF}" /etc/nginx/nginx.conf
  fi

  renew_hours=${CERTBOT_RENEW_INTERVAL_HOURS:-12}
  (
    while true; do
      sleep $((renew_hours * 3600))
      if certbot renew --webroot -w /var/www/certbot --quiet; then
        cp "${TLS_CONF}" /etc/nginx/nginx.conf
        nginx -s reload || true
      elif [ ! -f "${cert_dir}/fullchain.pem" ]; then
        if certbot certonly --non-interactive --agree-tos --email "${CERTBOT_EMAIL}" --webroot -w /var/www/certbot -d "${SERVER_NAME}"; then
          cp "${TLS_CONF}" /etc/nginx/nginx.conf
          nginx -s reload || true
        fi
      fi
    done
  ) &
else
  echo "[entrypoint] TLS disabled, using default nginx configuration"
  cp "${HTTP_CONF}" /etc/nginx/nginx.conf
fi

exec nginx -g "daemon off;"
