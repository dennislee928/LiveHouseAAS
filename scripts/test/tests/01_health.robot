*** Settings ***
Documentation     Smoke test — verify all 4 Choreo components respond
Resource          ../resources/common.resource
Suite Setup       Create API Session

*** Test Cases ***
Backend API Health Check
    [Tags]    smoke    backend
    ${headers}=    API Headers
    ${resp}=    GET On Session    api    /health    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200
    ${body}=    Set Variable    ${resp.json()}
    Log    ${body}
    Dictionary Should Contain Key    ${body}    status

NATS Monitoring Health Check
    [Tags]    smoke    nats
    ${lib}=    Get Library Instance    CustomLibrary
    Run Keyword And Continue On Failure
    ...    NATS Monitoring Healthz    ${NATS_MON_URL}

NATS Server Info
    [Tags]    smoke    nats
    ${lib}=    Get Library Instance    CustomLibrary
    ${data}=    Run Keyword And Continue On Failure
    ...    NATS Monitoring Varz    ${NATS_MON_URL}
    Run Keyword If    ${data}    Log    NATS version: ${data}[version]

MinIO Health Check
    [Tags]    smoke    minio
    ${headers}=    Create Dictionary    Accept=application/json
    Create Session    minio    ${MINIO_URL}    verify=False
    ${resp}=    GET On Session    minio    /minio/health/live    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

PostgreSQL Connectivity Check
    [Tags]    smoke    postgres
    Connect To Database Using Custom Params
    ...    psycopg2
    ...    database='${DB_NAME}', user='${DB_USER}', password='${DB_PASSWORD}', host='${DB_HOST}', port=${DB_PORT}, connect_timeout=10
    ${result}=    Query    SELECT version()
    Log    PostgreSQL version: ${result[0][0]}
    Disconnect From Database
