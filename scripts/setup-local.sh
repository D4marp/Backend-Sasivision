#!/usr/bin/env bash
# Setup & test API lokal tanpa Docker (butuh MySQL native di port 3306)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "==> Cek MySQL native..."
if ! mysql -h127.0.0.1 -P3306 -uroot -e "SELECT 1" >/dev/null 2>&1; then
  echo "MySQL tidak bisa diakses (127.0.0.1:3306, user root tanpa password)."
  echo "Pastikan MySQL/MariaDB lokal sudah jalan."
  exit 1
fi

echo "==> Buat database sasivision..."
mysql -h127.0.0.1 -P3306 -uroot -e "CREATE DATABASE IF NOT EXISTS sasivision CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

echo "==> Seed migrations..."
go run ./cmd/seed

echo "==> Unit tests..."
go test ./... -count=1

echo ""
echo "==> Jalankan server di terminal lain:"
echo "    cd $ROOT && go run ./cmd/server"
echo ""
echo "==> Lalu test API:"
echo "    ./scripts/test-api-local.sh"
echo ""
echo "==> Integration test (server harus sudah jalan):"
echo "    API_BASE_URL=http://127.0.0.1:8080 go test -tags=integration ./internal/integration/... -v"
