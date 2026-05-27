#!/usr/bin/env bash
# Run the LiveHouseAAS Robot Framework test suite
#
# Usage:
#   ./run_tests.sh                   # all tests
#   ./run_tests.sh --suite 07        # single suite by prefix
#   ./run_tests.sh --tag smoke       # filter by tag
#   ./run_tests.sh --tag backend
#   ./run_tests.sh --tag postgres
#
# Set env vars BEFORE running (use single quotes to avoid zsh special chars):
#   export CHOREO_API_KEY='eyJraWQi...'    # fresh key from Choreo UI (10 min TTL)
#   export POSTGRES_PASSWORD='your_pass'   # single quotes handle ! $ etc safely
#
# Optional (direct component tests):
#   export NATS_MONITOR_URL='http://<host>:8222'
#   export MINIO_URL='http://<host>:9000'
#   export MINIO_ACCESS_KEY='minioadmin'
#   export MINIO_SECRET_KEY='minioadmin'
#   export REDIS_HOST='<host>'
#   export REDIS_PORT='6379'
#   export REDIS_PASSWORD='<pass>'

set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

VENV_DIR="$SCRIPT_DIR/.venv"

# ── Virtual environment setup ─────────────────────────────────────────────────
if [[ ! -d "$VENV_DIR" ]]; then
    echo "Creating Python virtual environment at $VENV_DIR ..."
    python3 -m venv "$VENV_DIR"
fi

source "$VENV_DIR/bin/activate"

echo "Installing / verifying dependencies in venv..."
pip install -q -r requirements.txt

# ── Env-var validation ────────────────────────────────────────────────────────
WARN=0
if [[ -z "${CHOREO_API_KEY}" ]]; then
    echo ""
    echo "WARNING: CHOREO_API_KEY is not set — all backend API tests will return 401"
    echo "         The key expires 10 minutes after generation. Get a fresh one from:"
    echo "         Choreo UI → backend.prd → Test tab → copy 'Security Header' value"
    echo "         Then: export CHOREO_API_KEY='eyJraWQi...'  (use single quotes!)"
    WARN=1
fi
if [[ -z "${POSTGRES_PASSWORD}" ]]; then
    echo ""
    echo "WARNING: POSTGRES_PASSWORD is not set — PostgreSQL tests will fail"
    echo "         Use single quotes: export POSTGRES_PASSWORD='Popoman1217\!1996'"
    echo "         Or escape the !:   export POSTGRES_PASSWORD=Popoman1217\!1996"
    WARN=1
fi
[[ -z "${NATS_MONITOR_URL}" ]] && echo "INFO:    NATS_MONITOR_URL not set — NATS tests will fail (expected if not exposed)"
[[ -z "${MINIO_URL}" ]]        && echo "INFO:    MINIO_URL not set — MinIO tests will fail (expected if not exposed)"
[[ -z "${REDIS_HOST}" ]]       && echo "INFO:    REDIS_HOST not set — Redis tests will fail (expected if not exposed)"
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
echo "Results: open $SCRIPT_DIR/results/report.html"
