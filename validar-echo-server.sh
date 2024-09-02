#!/bin/bash

SERVER_CONTAINER_NAME="server"

NETWORK_NAME="tp0_testing_net"

TEST_MESSAGE="Hello, Echo Server! I am working fine yeyyyyy"

SERVER_PORT=12345

docker run --rm --network $NETWORK_NAME alpine sh -c "
    apk add --no-cache netcat-openbsd

    SERVER_IP=\$(getent hosts $SERVER_CONTAINER_NAME | awk '{ print \$1 }')

    if [ -z \"\$SERVER_IP\" ]; then
        echo 'action: test_echo_server | result: fail | error: could not resolve server IP'
        exit 1
    fi

    RESPONSE=\$(echo $TEST_MESSAGE | nc -w 1 \$SERVER_IP $SERVER_PORT)

    if [ \"\$RESPONSE\" = \"$TEST_MESSAGE\" ]; then
        echo 'action: test_echo_server | result: success'
    else
        echo 'action: test_echo_server | result: fail'
    fi
"
