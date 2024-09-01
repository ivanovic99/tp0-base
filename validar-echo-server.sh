#!/bin/bash

# Nombre del contenedor del servidor
SERVER_CONTAINER_NAME="server"

# Nombre de la red de Docker
NETWORK_NAME="tp0_testing_net"

# Mensaje de prueba
TEST_MESSAGE="Hello, Echo Server! I am working fine yeyyyyy"

# Puerto del servidor (puede ser parametrizado si es necesario)
SERVER_PORT=12345

# Crear un contenedor temporal para realizar la prueba
docker run --rm --network $NETWORK_NAME alpine sh -c "
    # Instalar netcat
    apk add --no-cache netcat-openbsd

    # Obtener la IP del servidor
    SERVER_IP=\$(getent hosts $SERVER_CONTAINER_NAME | awk '{ print \$1 }')

    # Verificar si se obtuvo la IP del servidor
    if [ -z \"\$SERVER_IP\" ]; then
        echo 'action: test_echo_server | result: fail | error: could not resolve server IP'
        exit 1
    fi

    # Enviar el mensaje de prueba al servidor y recibir la respuesta
    RESPONSE=\$(echo $TEST_MESSAGE | nc \$SERVER_IP $SERVER_PORT)

    # Verificar si la respuesta es igual al mensaje de prueba
    if [ \"\$RESPONSE\" = \"$TEST_MESSAGE\" ]; then
        echo 'action: test_echo_server | result: success'
    else
        echo 'action: test_echo_server | result: fail'
    fi
"
