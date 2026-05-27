*** Settings ***
Documentation     MinIO component tests — HTTP health + S3 API round-trip
...               Set MINIO_URL, MINIO_ACCESS_KEY, MINIO_SECRET_KEY env vars.
Resource          ../resources/common.resource
Library           ../resources/CustomLibrary.py
Suite Setup       Run Keywords    Create Session    minio    ${MINIO_URL}    verify=False
...               AND    Connect To MinIO    ${MINIO_URL}    ${MINIO_ACCESS_KEY}    ${MINIO_SECRET_KEY}
Suite Teardown    MinIO Delete Object    robot-test-bucket    robot-test-object.txt

*** Test Cases ***
MinIO Live Health Endpoint
    [Tags]    minio    smoke
    ${resp}=    GET On Session    minio    /minio/health/live
    Should Be Equal As Integers    ${resp.status_code}    200

MinIO Ready Health Endpoint
    [Tags]    minio    smoke
    ${resp}=    GET On Session    minio    /minio/health/ready    expected_status=any
    Should Be True    ${resp.status_code} in [200, 503]
    Log    MinIO ready status: ${resp.status_code}

MinIO List Buckets
    [Tags]    minio
    ${buckets}=    MinIO List Buckets
    Log    Existing buckets: ${buckets}

MinIO Create Bucket
    [Tags]    minio
    MinIO Create Bucket If Not Exists    robot-test-bucket

MinIO Upload Object
    [Tags]    minio
    MinIO Upload And Download
    ...    robot-test-bucket
    ...    robot-test-object.txt
    ...    Hello from Robot Framework tests!

MinIO Console Port Responds
    [Tags]    minio    console
    ${console_url}=    Evaluate
    ...    '${MINIO_URL}'.replace(':9000', ':9001')
    Create Session    minio_console    ${console_url}    verify=False
    ${resp}=    GET On Session    minio_console    /    expected_status=any
    Should Be True    ${resp.status_code} in [200, 302, 404]
    Log    MinIO console responded: ${resp.status_code}
