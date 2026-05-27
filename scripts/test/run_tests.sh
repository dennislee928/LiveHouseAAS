#!/usr/bin/env bash
# Run the full LiveHouseAAS Robot Framework test suite
#
# Usage:
#   ./run_tests.sh                   # all tests
#   ./run_tests.sh --suite 07        # single suite (prefix match)
#   ./run_tests.sh --tag smoke       # filter by tag
#   ./run_tests.sh --tag postgres    # only postgres tests
#
# Required env vars:
#   CHOREO_API_KEY     Choreo internal key (from component → Test → Security Header)
#   POSTGRES_PASSWORD  AlwaysData PostgreSQL password
#
# Optional env vars (leave unset to skip direct component tests):
#   NATS_MONITOR_URL   http://<host>:8222
#   MINIO_URL          http://<host>:9000
#   MINIO_ACCESS_KEY   MinIO root user   (default: minioadmin)
#   MINIO_SECRET_KEY   MinIO root secret (default: minioadmin)
#   REDIS_HOST         Redis hostname
#   REDIS_PORT         Redis port        (default: 6379)
#   REDIS_PASSWORD     Redis password    (if set)

set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# ── Dependency check ──────────────────────────────────────────────────────────
echo "Installing / verifying Python dependencies..."
pip3 install -q -r requirements.txt

# ── Env-var validation ────────────────────────────────────────────────────────
WARN=0
if [[ -z "${CHOREO_API_KEY}" ]]; then
    echo "WARNING: CHOREO_API_KEY is not set — all backend API tests will return 401"
    WARN=1
fi
if [[ -z "${POSTGRES_PASSWORD}" ]]; then
    echo "WARNING: POSTGRES_PASSWORD is not set — PostgreSQL tests will fail"
    WARN=1
fi
if [[ -z "${NATS_MONITOR_URL}" ]]; then
    echo "INFO:    NATS_MONITOR_URL not set — NATS tests will use localhost:8222 (likely fail)"
fi
if [[ -z "${MINIO_URL}" ]]; then
    echo "INFO:    MINIO_URL not set — MinIO tests will use localhost:9000 (likely fail)"
fi
if [[ -z "${REDIS_HOST}" ]]; then
    echo "INFO:    REDIS_HOST not set — Redis tests will use localhost (likely fail)"
fi
[[ $WARN -eq 1 ]] && echo ""

# ── Robot arguments ───────────────────────────────────────────────────────────
ROBOT_ARGS=(
    --outputdir results
    --log log.html
    --report report.html
    --timestampoutputs
    --pythonpath resources
)

while [[ $# -gt 0 ]]; do
    case $1 in
        --suite) ROBOT_ARGS+=(--suite "$2"); shift 2 ;;
        --tag)   ROBOT_ARGS+=(--include "$2"); shift 2 ;;
        --skip)  ROBOT_ARGS+=(--exclude "$2"); shift 2 ;;
        *)       shift ;;
    esac
done

mkdir -p results

echo "Running Robot Framework tests..."
echo "────────────────────────────────────────"
python3 -m robot "${ROBOT_ARGS[@]}" tests/ || true

echo ""
echo "────────────────────────────────────────"
echo "Results written to: $SCRIPT_DIR/results/"
echo "  Open: $SCRIPT_DIR/results/report.html"
