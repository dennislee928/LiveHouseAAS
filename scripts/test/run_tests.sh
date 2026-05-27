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

# psycopg2-binary 2.9.10 does not publish wheels for Python 3.14 yet. On
# Homebrew systems where python3 points at 3.14, pip falls back to a source
# build and can fail while linking OpenSSL/libpq.
MIN_PYTHON="3.9"
MAX_PYTHON_EXCLUSIVE="3.14"

python_supported() {
    "$1" -c 'import sys; raise SystemExit(not ((3, 9) <= sys.version_info[:2] < (3, 14)))' 2>/dev/null
}

find_python() {
    if [[ -n "${PYTHON:-}" ]]; then
        if python_supported "$PYTHON"; then
            echo "$PYTHON"
            return 0
        fi
        echo "ERROR: PYTHON=$PYTHON is not supported. Use Python >= $MIN_PYTHON and < $MAX_PYTHON_EXCLUSIVE." >&2
        return 1
    fi

    local candidate
    for candidate in python3.13 python3.12 python3.11 python3; do
        if command -v "$candidate" >/dev/null 2>&1 && python_supported "$candidate"; then
            command -v "$candidate"
            return 0
        fi
    done

    echo "ERROR: No supported Python found. Install Python >= $MIN_PYTHON and < $MAX_PYTHON_EXCLUSIVE, or set PYTHON=/path/to/python." >&2
    return 1
}

PYTHON_BIN="$(find_python)"

# ── Virtual environment setup ─────────────────────────────────────────────────
if [[ -d "$VENV_DIR" ]] && ! "$VENV_DIR/bin/python" -c 'import sys; raise SystemExit(not ((3, 9) <= sys.version_info[:2] < (3, 14)))' 2>/dev/null; then
    echo "Existing venv uses unsupported Python; recreating $VENV_DIR ..."
    rm -rf "$VENV_DIR"
fi

if [[ ! -d "$VENV_DIR" ]]; then
    echo "Creating Python virtual environment at $VENV_DIR ..."
    "$PYTHON_BIN" -m venv "$VENV_DIR"
fi

source "$VENV_DIR/bin/activate"

echo "Installing / verifying dependencies in venv..."
pip install -q -r requirements.txt

is_placeholder() {
    [[ "$1" == *"<"*">"* ]]
}

# ── Env-var validation ────────────────────────────────────────────────────────
WARN=0
if [[ -z "${CHOREO_API_KEY}" ]]; then
    echo ""
    echo "WARNING: CHOREO_API_KEY is not set — all backend API tests will return 401"
    echo "         The key expires 10 minutes after generation. Get a fresh one from:"
    echo "         Choreo UI → backend.prd → Test tab → copy 'Security Header' value"
    echo "         Then: export CHOREO_API_KEY='eyJraWQi...'  (use single quotes!)"
    WARN=1
else
    echo "INFO:    CHOREO_API_KEY is set; sending it as header '${CHOREO_API_KEY_HEADER:-api-key}'"
fi
if [[ -z "${POSTGRES_PASSWORD}" ]]; then
    echo ""
    echo "WARNING: POSTGRES_PASSWORD is not set — PostgreSQL tests will fail"
    echo "         Use single quotes: export POSTGRES_PASSWORD='Popoman1217\!1996'"
    echo "         Or escape the !:   export POSTGRES_PASSWORD=Popoman1217\!1996"
    WARN=1
fi
if [[ -z "${NATS_MONITOR_URL}" ]] || is_placeholder "${NATS_MONITOR_URL}"; then
    echo "INFO:    NATS_MONITOR_URL not configured — NATS tests will fail unless this service is exposed"
fi
if [[ -z "${MINIO_URL}" ]] || is_placeholder "${MINIO_URL}"; then
    echo "INFO:    MINIO_URL not configured — MinIO tests will fail unless this service is exposed"
fi
if [[ -z "${REDIS_HOST}" ]] || is_placeholder "${REDIS_HOST}"; then
    echo "INFO:    REDIS_HOST not configured — Redis tests will fail unless this service is exposed"
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
echo "Results: open $SCRIPT_DIR/results/report.html"
