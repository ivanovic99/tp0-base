#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <output_file> <num_clients>"
    exit 1
fi

# Assign arguments to variables
output_file=$1
num_clients=$2

# Print the provided arguments
echo "Nombre del archivo de salida: $output_file"
echo "Cantidad de clientes: $num_clients"

# Run the Python script to generate the docker-compose file
python3 mi-generador.py "$output_file" "$num_clients"
