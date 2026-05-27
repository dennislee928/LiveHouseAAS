*** Settings ***
Documentation     NATS component tests via HTTP monitoring API (port 8222)
...               Set NATS_MONITOR_URL env var to the exposed monitoring endpoint.
Resource          ../resources/common.resource
Library           CustomLibrary

*** Test Cases ***
NATS Health Endpoint Responds
    [Tags]    nats    smoke
    NATS Monitoring Healthz    ${NATS_MON_URL}

NATS Server Info Contains Version
    [Tags]    nats
    ${data}=    NATS Monitoring Varz    ${NATS_MON_URL}
    Dictionary Should Contain Key    ${data}    version
    Log    NATS server version: ${data}[version]
    Should Not Be Empty    ${data}[version]

NATS JetStream Is Enabled
    [Tags]    nats    jetstream
    ${data}=    NATS Monitoring Varz    ${NATS_MON_URL}
    # -js flag was passed in CMD — jetstream section should be present
    Dictionary Should Contain Key    ${data}    jetstream
    ${js}=    Get From Dictionary    ${data}    jetstream
    Should Not Be Equal    ${js}    ${None}
    Log    JetStream config: ${js}

NATS Connections Endpoint Accessible
    [Tags]    nats
    ${lib}=    Get Library Instance    CustomLibrary
    Create Session    nats_mon    ${NATS_MON_URL}    verify=False
    ${resp}=    GET On Session    nats_mon    /connz
    Should Be Equal As Integers    ${resp.status_code}    200
    Dictionary Should Contain Key    ${resp.json()}    num_connections

NATS Routes Endpoint Accessible
    [Tags]    nats
    Create Session    nats_mon    ${NATS_MON_URL}    verify=False
    ${resp}=    GET On Session    nats_mon    /routez
    Should Be Equal As Integers    ${resp.status_code}    200

NATS Subscriptions Endpoint Accessible
    [Tags]    nats
    Create Session    nats_mon    ${NATS_MON_URL}    verify=False
    ${resp}=    GET On Session    nats_mon    /subsz
    Should Be Equal As Integers    ${resp.status_code}    200
