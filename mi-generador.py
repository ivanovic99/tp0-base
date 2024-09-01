"""
Generates a docker-compose file with a server and N clients.
"""

import sys
import yaml
from typing import Dict, Any

# Constants
DOCKER_COMPOSE_VERSION = '3.8'
NETWORK_SUBNET = '172.25.125.0/24'
SERVER_IMAGE = 'server:latest'
CLIENT_IMAGE = 'client:latest'
SERVER_ENTRYPOINT = 'python3 /main.py'
CLIENT_ENTRYPOINT = '/client'
LOGGING_LEVEL = 'DEBUG'

def generate_docker_compose(output_file: str, num_clients: int) -> None:
    """
    Generate a docker-compose file with a server and N clients.

    Parameters:
        output_file (str): The output file path.
        num_clients (int): The number of clients to generate.
    
    Returns:
        None
    """
    docker_compose: Dict[str, Any] = {
        'version': DOCKER_COMPOSE_VERSION,
        'services': {
            'server': {
                'container_name': 'server',
                'image': SERVER_IMAGE,
                'entrypoint': SERVER_ENTRYPOINT,
                'environment': [
                    'PYTHONUNBUFFERED=1',
                    f'LOGGING_LEVEL={LOGGING_LEVEL}'
                ],
                'networks': ['testing_net']
            }
        },
        'networks': {
            'testing_net': {
                'ipam': {
                    'driver': 'default',
                    'config': [{'subnet': NETWORK_SUBNET}]
                }
            }
        }
    }

    for client_N in range(1, num_clients + 1):
        client_name = f'client{client_N}'
        docker_compose['services'][client_name] = {
            'container_name': client_name,
            'image': CLIENT_IMAGE,
            'entrypoint': CLIENT_ENTRYPOINT,
            'environment': [
                f'CLI_ID={client_N}',
                f'CLI_LOG_LEVEL={LOGGING_LEVEL}'
            ],
            'networks': ['testing_net'],
            'depends_on': ['server']
        }

    try:
        with open(output_file, 'w') as file:
            yaml.dump(docker_compose, file, default_flow_style=False)
        print(f"Docker-compose file generated successfully: {output_file}")
    except IOError as e:
        print(f"Failed to write docker-compose file: {e}")
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 mi-generador.py <output_file> <num_clients>")
        sys.exit(1)

    output_file = sys.argv[1]
    try:
        num_clients = int(sys.argv[2])
    except ValueError:
        print("The number of clients must be an integer.")
        sys.exit(1)

    generate_docker_compose(output_file, num_clients)
