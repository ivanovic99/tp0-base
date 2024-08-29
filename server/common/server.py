import socket
import logging
import threading
from .protocol import Protocol
from .utils import store_bets

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._clients = []
        self._is_running = True

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # the server
        while self._is_running:
            try:
                client_sock = self.__accept_new_connection()
                client_thread = threading.Thread(target=self.__handle_client_connection, args=(client_sock,))
                client_thread.start()
                self._clients.append(client_thread)
            except socket.error as e:
                if self._running:
                    logging.error(f"action: accept_new_connection | result: fail | error: {e}")


    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            protocol = Protocol(client_sock)
            bet = protocol.receive_bet()
            logging.info(f'action: receive_message | result: success | dni: {bet.document} | numero: {bet.number}')
            store_bets([bet])
            logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}')
            protocol.send_bet(bet)
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        connection, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return connection

    def shutdown(self):
        """
        Gracefully shutdown the server
        """
        logging.info('action: shutdown | result: in_progress | reason: graceful_shutdown')
        self._is_running = False
        self._server_socket.close()
        for client in self._clients:
            client.join()
        logging.info('action: shutdown | result: all_clients_joined')


