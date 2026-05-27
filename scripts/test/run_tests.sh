#!/usr/bin/env bash
# Run the full LiveHouseAAS Robot Framework test suite
# Usage: ./run_tests.sh [--suite <name>] [--tag <tag>]
#
# Required env vars:
#   CHOREO_API_KEY     — Choreo internal key (from component → Test → Security Header)
#   POSTGRES_PASSWORD  — AlwaysData PostgreSQL password
#
# Optional env vars (for direct component tests):
#   NATS_MONITOR_URL   — e.g. http://<nats-choreo-host>:8222
#   MINIO_URL          — e.g. http://<minio-choreo-host>:9000
#   MINIO_ACCESS_KEY   — MinIO root user
#   MINIO_SECRET_KEY   — MinIO root password
#   REDIS_HOST         — Redis hostname
#   REDIS_PORT         — Redis port (default 6379)
#   REDIS_PASSWORD     — Redis password (if set)

set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# Install dependencies if not already installed
if ! python3 -c "import robot" &>/dev/null; then
    echo "Installing Robot Framework dependencies..."
    pip3 install -r requirements.txt
fi

# Build robot arguments
ROBOT_ARGS=(
    --outputdir results
    --log log.html
    --report report.html
    --timestampoutputs
)

# Optional: filter by suite or tag
while [[ $# -gt 0 ]]; do
    case $1 in
        --suite) ROBOT_ARGS+=(--suite "$2"); shift 2 ;;
        --tag)   ROBOT_ARGS+=(--include "$2"); shift 2 ;;
        *)       shift ;;
    esac
done

mkdir -p results

echo "Running tests..."
python3 -m robot "${ROBOT_ARGS[@]}" tests/

echo ""
echo "Results: $SCRIPT_DIR/results/report.html"
