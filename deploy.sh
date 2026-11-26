#!/bin/bash
set -e  # exit on any error

trap 'echo "ERROR: Script failed at line $LINENO with exit code $?"' ERR

ENV_FILE=$(mktemp)
trap 'rm -f "$ENV_FILE"' EXIT
chmod 600 "$ENV_FILE"

SECRET_NAME="net-watch"

echo "Fetching secrets from Secrets Manager..."
NETWATCH_SECRET=$(aws secretsmanager get-secret-value \
  --secret-id "${SECRET_NAME}" \
  --query SecretString \
  --output text)

DB_USER=$(echo "${NETWATCH_SECRET}" | grep -oP '"DB_USER"\s*:\s*"\K[^"]+')
DB_PWD=$(echo "${NETWATCH_SECRET}" | grep -oP '"DB_PWD"\s*:\s*"\K[^"]+')

cat > "$ENV_FILE" <<EOF
DB_NAME=netwatch
DB_USER=${DB_USER}
DB_PWD=${DB_PWD}
EOF

unset NETWATCH_SECRET DB_USER DB_PWD

PRIMARY_DATA_PATH="D:/Projects/.data"
FALLBACK_DATA_PATH="C:/ProgramData/docker/volumes"

if [[ -d "${PRIMARY_DATA_PATH}" ]]; then
  echo "Using primary data path: ${PRIMARY_DATA_PATH}"
  export DATA_PATH="${PRIMARY_DATA_PATH}/net-watch"
else
  echo "Primary data path not found, using fallback: ${FALLBACK_DATA_PATH}"
  export DATA_PATH="${FALLBACK_DATA_PATH}/net-watch"
fi

echo "Starting services..."
docker compose --env-file "$ENV_FILE" -p net-watch up -d

DB_BINDING=$(docker compose --env-file "$ENV_FILE" -p net-watch port netwatch-db 5432)
printf 'PostgreSQL listening on %s\n' "${DB_BINDING}"

echo "Done."
