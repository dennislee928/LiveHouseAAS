*** Settings ***
Documentation     Core API endpoints — venues, events, bookings, tickets
Resource          ../resources/common.resource
Library           String
Suite Setup       Run Keywords    Create API Session    AND    Suite Login

*** Variables ***
${AUTH_TOKEN}      ${EMPTY}
${VENUE_ID}        ${EMPTY}
${EVENT_ID}        ${EMPTY}

*** Keywords ***
Suite Login
    ${suffix}=    Generate Random String    6    [NUMBERS]
    Set Suite Variable    ${SUITE_EMAIL}    api_test_${suffix}@livehouseaas.test
    ${headers}=    API Headers
    # Register
    ${body}=    Create Dictionary
    ...    email=${SUITE_EMAIL}    password=${TEST_PASSWORD}
    ...    name=API Test User    role=venue_owner
    POST On Session    api    /auth/register    json=${body}    headers=${headers}    expected_status=any
    # Login
    ${resp}=    POST On Session    api    /auth/login    json=${body}    headers=${headers}
    Set Suite Variable    ${AUTH_TOKEN}    ${resp.json()}[token]

*** Test Cases ***
# ── Venues ───────────────────────────────────────────────────────────────────

List Venues Returns 200
    [Tags]    venues
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/venues    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

Create Venue
    [Tags]    venues
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${body}=    Create Dictionary
    ...    name=Robot Test Venue
    ...    address=123 Test St, Taipei
    ...    capacity=200
    ...    description=Created by Robot Framework tests
    ${resp}=    POST On Session    api    /api/v1/venues
    ...    json=${body}    headers=${headers}    expected_status=any
    Log Response    ${resp}
    Should Be True    ${resp.status_code} in [200, 201]
    Set Suite Variable    ${VENUE_ID}    ${resp.json()}[id]

Get Venue By ID
    [Tags]    venues
    Skip If    '${VENUE_ID}' == ''    Create Venue must pass first
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/venues/${VENUE_ID}    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200
    Should Be Equal As Strings    ${resp.json()}[id]    ${VENUE_ID}

Get My Venues
    [Tags]    venues
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/venues/my    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

# ── Events ────────────────────────────────────────────────────────────────────

List Events Returns 200
    [Tags]    events
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/events    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

Create Event
    [Tags]    events
    Skip If    '${VENUE_ID}' == ''    Venue must exist first
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${body}=    Create Dictionary
    ...    title=Robot Test Event
    ...    venue_id=${VENUE_ID}
    ...    start_time=2026-12-01T18:00:00Z
    ...    end_time=2026-12-01T21:00:00Z
    ...    description=Automated test event
    ...    genre=Rock
    ${resp}=    POST On Session    api    /api/v1/events
    ...    json=${body}    headers=${headers}    expected_status=any
    Log Response    ${resp}
    Should Be True    ${resp.status_code} in [200, 201]
    Set Suite Variable    ${EVENT_ID}    ${resp.json()}[id]

Get Event By ID
    [Tags]    events
    Skip If    '${EVENT_ID}' == ''    Create Event must pass first
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/events/${EVENT_ID}    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

Search Events Returns 200
    [Tags]    events    search
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /search/events?q=Robot    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

# ── Bookings ──────────────────────────────────────────────────────────────────

List My Bookings
    [Tags]    bookings
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/bookings    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

# ── Tickets ───────────────────────────────────────────────────────────────────

List My Tickets
    [Tags]    tickets
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/tickets    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

# ── Analytics ─────────────────────────────────────────────────────────────────

Analytics Summary
    [Tags]    analytics
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/analytics/summary    headers=${headers}    expected_status=any
    Should Be True    ${resp.status_code} in [200, 403]

# ── Notifications ─────────────────────────────────────────────────────────────

List Notifications
    [Tags]    notifications
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/notifications    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

Unread Notification Count
    [Tags]    notifications
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/notifications/unread-count    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

# ── KYB ───────────────────────────────────────────────────────────────────────

KYB Status
    [Tags]    kyb
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/kyb/status    headers=${headers}    expected_status=any
    Should Be True    ${resp.status_code} in [200, 404]

# ── Dashboard ─────────────────────────────────────────────────────────────────

Dashboard Stats
    [Tags]    dashboard
    ${headers}=    Authenticated Headers    ${AUTH_TOKEN}
    ${resp}=    GET On Session    api    /api/v1/dashboard/stats    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200
