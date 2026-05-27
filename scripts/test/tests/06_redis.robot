*** Settings ***
Documentation     Redis component tests — PING, SET/GET, TTL, data type round-trips
...               Set REDIS_HOST, REDIS_PORT, REDIS_PASSWORD env vars.
Resource          ../resources/common.resource
Library           CustomLibrary
Suite Setup       Connect To Redis    ${REDIS_HOST}    ${REDIS_PORT}    ${REDIS_PASSWORD}
Suite Teardown    Disconnect From Redis

*** Test Cases ***
Redis PING
    [Tags]    redis    smoke
    Redis Ping

Redis String Set And Get
    [Tags]    redis
    Redis Set And Get    robot:test:string    RobotFrameworkValue

Redis Overwrite Value
    [Tags]    redis
    Redis Set And Get    robot:test:string    UpdatedValue

Redis Key Expiry Works
    [Tags]    redis
    Redis Set And Get    robot:test:ttl    TemporaryValue
    Redis Delete Key    robot:test:ttl

Redis Multiple Keys
    [Tags]    redis
    Redis Set And Get    robot:test:key1    value1
    Redis Set And Get    robot:test:key2    value2
    Redis Set And Get    robot:test:key3    value3
    [Teardown]    Run Keywords
    ...    Redis Delete Key    robot:test:key1    AND
    ...    Redis Delete Key    robot:test:key2    AND
    ...    Redis Delete Key    robot:test:key3

Redis Clean Up Test Keys
    [Tags]    redis    teardown
    Redis Delete Key    robot:test:string
