*** Settings ***
Documentation     Auth flow: register → login → /auth/me → change password
Resource          ../resources/common.resource
Library           String
Suite Setup       Create API Session

*** Variables ***
${REGISTERED_TOKEN}    ${EMPTY}

*** Test Cases ***
Register New User
    [Tags]    auth
    ${suffix}=    Generate Random String    6    [NUMBERS]
    Set Suite Variable    ${UNIQUE_EMAIL}    robot_test_${suffix}@livehouseaas.test
    ${headers}=    API Headers
    ${body}=    Create Dictionary
    ...    email=${UNIQUE_EMAIL}
    ...    password=${TEST_PASSWORD}
    ...    name=Robot Test User
    ...    role=user
    ${resp}=    POST On Session    api    /auth/register
    ...    json=${body}    headers=${headers}    expected_status=any
    Log Response    ${resp}
    Should Be Equal As Integers    ${resp.status_code}    201

Login With Valid Credentials
    [Tags]    auth
    ${headers}=    API Headers
    ${body}=    Create Dictionary
    ...    email=${UNIQUE_EMAIL}
    ...    password=${TEST_PASSWORD}
    ${resp}=    POST On Session    api    /auth/login
    ...    json=${body}    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200
    ${data}=    Set Variable    ${resp.json()}
    Dictionary Should Contain Key    ${data}    token
    Set Suite Variable    ${REGISTERED_TOKEN}    ${data}[token]
    Log    Login token acquired

Login With Wrong Password Returns 401
    [Tags]    auth    negative
    ${headers}=    API Headers
    ${body}=    Create Dictionary
    ...    email=${UNIQUE_EMAIL}
    ...    password=WrongPassword!
    ${resp}=    POST On Session    api    /auth/login
    ...    json=${body}    headers=${headers}    expected_status=any
    Should Be Equal As Integers    ${resp.status_code}    401

Get Current User Profile
    [Tags]    auth
    Skip If    '${REGISTERED_TOKEN}' == ''    Login test must pass first
    ${headers}=    Authenticated Headers    ${REGISTERED_TOKEN}
    ${resp}=    GET On Session    api    /auth/me    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200
    ${data}=    Set Variable    ${resp.json()}
    Should Be Equal As Strings    ${data}[email]    ${UNIQUE_EMAIL}

Unauthenticated Request Returns 401
    [Tags]    auth    negative
    ${headers}=    API Headers
    ${resp}=    GET On Session    api    /auth/me
    ...    headers=${headers}    expected_status=any
    Should Be Equal As Integers    ${resp.status_code}    401
