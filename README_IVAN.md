# Project Overview

This project implements a server-client architecture for managing bets and performing a draw. The server handles multiple clients concurrently using multiprocessing, and the communication between the client and server is done using custom serialization methods.

## How to Run the Project

To run the project, you can use the following `make` commands:

1. **Start the Docker Compose environment**:
   ```sh
   make docker-compose-up
    ```
2. **View the logs**:
   ```sh
   make docker-compose-logs
    ```
3. **Stop the Docker Compose environment**:
    ```sh
    make docker-compose-down
    ```

## Protocol and Serialization

### Protocol

The communication protocol between the client and server is defined as follows:

Case ID: The client sends a case ID to the server to indicate the type of request.

1: Send bets to the server.
2: Indicate that the client has finished sending bets.
3: Request the list of winners for a specific agency.

### Serialization

The serialization and deserialization of data are handled manually without using libraries like JSON. The custom methods ensure that the data is correctly serialized and deserialized for communication over sockets.

#### Serialize Bets

The serialize_bets function converts a list of Bet objects into a byte stream. Each field of the Bet object is serialized in a specific order, with length-prefixed strings for variable-length fields.

#### Deserialize Bets

The deserialize_bets function converts a byte stream back into a list of Bet objects. It reads each field in the same order as they were serialized, reconstructing the Bet objects.

#### Serialize Winners

The serialize_winners function converts a list of winner document numbers into a byte stream. Each document number is serialized as a length-prefixed string.

#### Deserialize Winners

The deserialize_winners function converts a byte stream back into a list of winner document numbers. It reads each length-prefixed string to reconstruct the list of winners.


### Concurrency Management in the Server

The server handles multiple clients concurrently using the multiprocessing module. Here are the key components:

Server Socket: The server listens for incoming connections using a server socket.

Client Processes: For each new client connection, the server spawns a new process using multiprocessing.Process. This allows the server to handle multiple clients concurrently.

Shared State: The server uses multiprocessing.Value and multiprocessing.Manager().dict() to manage shared state across processes. This includes:

self._finished_clients: A counter to track the number of clients that have finished sending bets.
self._draw_finished: A flag to indicate whether the draw has been performed.
self._winners: A dictionary to store the list of winners for each agency.
Barrier: The server uses a multiprocessing.Barrier to synchronize the clients. This ensures that all clients wait at the barrier until the draw is performed and the winners are determined.

Lock: The server uses a multiprocessing.Lock to ensure thread-safe access to shared resources, such as storing bets and updating the list of winners.

By using these concurrency management techniques, the server can efficiently handle multiple clients and ensure correct synchronization and data integrity.
