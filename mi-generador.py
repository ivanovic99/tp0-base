"""
Generates a docker-compose file with a server and N clients
"""

import sys
import yaml

def generate_docker_compose(output_file: str, num_clients: int):
    """
    Generate a docker-compose file with a server and N clients

    Parameters:
        output_file (str): The output file path
        num_clients (int): The number of clients to generate
    
    Returns:
        None
    """

    docker_compose = {
        'name': 'tp0',
        'services': {
            'server': {
                'container_name': 'server',
                'image': 'server:latest',
                'entrypoint': 'python3 /main.py',
                'environment': [
                    'PYTHONUNBUFFERED=1',
                    'LOGGING_LEVEL=DEBUG'
                ],
                'volumes': [
                    './server/config.ini:/config.ini'
                ],
                'networks': ['testing_net']
            }
        },
        'networks': {
            'testing_net': {
                'ipam': {
                    'driver': 'default',
                    'config': [{'subnet': '172.25.125.0/24'}]
                }
            }
        }
    }

    for client_N in range(1, num_clients + 1):
        client_name = f'client{client_N}'
        docker_compose['services'][client_name] = {
            'container_name': client_name,
            'image': 'client:latest',
            'entrypoint': '/client',
            'environment': [
                f'CLI_ID={client_N}',
                'CLI_LOG_LEVEL=DEBUG'
            ],
            'volumes': [
                './client/config.yaml:/config.yaml'
            ],
            'networks': ['testing_net'],
            'depends_on': ['server']
        }

    with open(output_file, 'w') as file:
        yaml.dump(docker_compose, file, default_flow_style=False)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 mi-generador.py <output_file> <num_clients>")
        sys.exit(1)

    output_file = sys.argv[1]
    num_clients = int(sys.argv[2])
    generate_docker_compose(output_file, num_clients)
