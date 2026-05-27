*** Settings ***
Documentation     PostgreSQL tests — AlwaysData hosted instance
...               Host: postgresql-dennislee928.alwaysdata.net:5432
...               DB:   dennislee928_2 / User: dennislee928
...               Set POSTGRES_PASSWORD env var before running.
Resource          ../resources/common.resource
Library           DatabaseLibrary
Suite Setup       Connect To AlwaysData Postgres
Suite Teardown    Disconnect From Database

*** Keywords ***
Connect To AlwaysData Postgres
    Connect To Database    psycopg2
    ...    dbName=${DB_NAME}    dbUsername=${DB_USER}    dbPassword=${DB_PASSWORD}
    ...    dbHost=${DB_HOST}    dbPort=${DB_PORT}

*** Test Cases ***
Database Connection Is Live
    [Tags]    postgres    smoke
    ${result}=    Query    SELECT 1 AS alive
    Should Be Equal As Integers    ${result[0][0]}    1

PostgreSQL Version Logged
    [Tags]    postgres    smoke
    ${result}=    Query    SELECT version()
    Log    PostgreSQL: ${result[0][0]}
    Should Contain    ${result[0][0]}    PostgreSQL

List Application Tables
    [Tags]    postgres
    ${result}=    Query
    ...    SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY tablename
    Log    Tables: ${result}
    Should Not Be Empty    ${result}

Users Table Exists
    [Tags]    postgres    schema
    Table Must Exist    users

Events Table Exists
    [Tags]    postgres    schema
    Table Must Exist    events

Venues Table Exists
    [Tags]    postgres    schema
    Table Must Exist    venues

Bookings Table Exists
    [Tags]    postgres    schema
    Table Must Exist    bookings

Orders Table Exists
    [Tags]    postgres    schema
    Table Must Exist    orders

Tickets Table Exists
    [Tags]    postgres    schema
    Table Must Exist    tickets

Row Count Queries Return Without Error
    [Tags]    postgres
    ${users}=     Query    SELECT COUNT(*) FROM users
    ${venues}=    Query    SELECT COUNT(*) FROM venues
    ${events}=    Query    SELECT COUNT(*) FROM events
    Log    users=${users[0][0]}  venues=${venues[0][0]}  events=${events[0][0]}

Connection Pool Under Load
    [Tags]    postgres    load
    FOR    ${i}    IN RANGE    10
        ${r}=    Query    SELECT ${i} * 2 AS val
        Should Be Equal As Integers    ${r[0][0]}    ${i * 2}
    END
