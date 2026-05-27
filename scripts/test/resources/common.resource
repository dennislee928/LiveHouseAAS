*** Settings ***
Library    RequestsLibrary
Library    Collections
Library    String
Library    OperatingSystem

*** Variables ***
# --- Backend API (Choreo) ---
${API_BASE_URL}         https://960b78ad-4106-4b00-a989-39ae3209f7d4-dev.e1-us-east-azure.choreoapis.dev/default/backend-ky/v1.0
${API_KEY_HEADER}       %{CHOREO_API_KEY_HEADER=api-key}
${API_KEY}              %{CHOREO_API_KEY=}

# --- PostgreSQL (AlwaysData) ---
${DB_HOST}              postgresql-dennislee928.alwaysdata.net
${DB_PORT}              5432
${DB_NAME}              dennislee928_2
${DB_USER}              dennislee928
${DB_PASSWORD}          %{POSTGRES_PASSWORD=}

# --- Internal Choreo service URLs (set via env if exposed) ---
${NATS_MON_URL}         %{NATS_MONITOR_URL=http://localhost:8222}
${MINIO_URL}            %{MINIO_URL=http://localhost:9000}
${MINIO_ACCESS_KEY}     %{MINIO_ACCESS_KEY=minioadmin}
${MINIO_SECRET_KEY}     %{MINIO_SECRET_KEY=minioadmin}
${REDIS_HOST}           %{REDIS_HOST=localhost}
${REDIS_PORT}           %{REDIS_PORT=6379}
${REDIS_PASSWORD}       %{REDIS_PASSWORD=}

# --- Test data ---
${TEST_EMAIL}           robot_test_%{RANDOM_SUFFIX=001}@livehouseaas.test
${TEST_PASSWORD}        RobotTest@2026!
${TEST_VENUE_NAME}      Robot Test Venue

*** Keywords ***
API Headers
    [Documentation]    Return standard headers for Choreo API calls
    ${headers}=    Create Dictionary
    ...    Content-Type=application/json
    ...    Accept=application/json
    Run Keyword If    '${API_KEY}' != ''    Set To Dictionary    ${headers}    ${API_KEY_HEADER}    ${API_KEY}
    RETURN    ${headers}

Authenticated Headers
    [Arguments]    ${token}
    [Documentation]    Return headers with both Choreo API-Key and user JWT
    ${headers}=    API Headers
    Set To Dictionary    ${headers}    Authorization    Bearer ${token}
    RETURN    ${headers}

Create API Session
    Create Session    api    ${API_BASE_URL}    verify=True

Log Response
    [Arguments]    ${response}
    Log    Status: ${response.status_code}
    Log    Body: ${response.text}
