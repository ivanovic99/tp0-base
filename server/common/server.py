import socket
import logging
import threading
from .protocol import Protocol
from .utils import store_bets, load_bets, has_won

TOTAL_CLIENTS = 2

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._clients = []
        self._is_running = True
        self._finished_clients = 0
        self._total_clients = TOTAL_CLIENTS
        self._winners = {1: [], 2: []}
        self._barrier = threading.Barrier(self._total_clients)


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
            has_finished = False
            while not has_finished:
                try:
                    case_id = protocol.receive_case()
                    if case_id == 1:
                        bets = protocol.receive_bets()
                        if not bets:
                            has_finished = True
                        try:
                            store_bets(bets)
                            logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
                        except Exception as e:
                            logging.info(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')
                    elif case_id == 2:
                        has_finished = True
                except Exception as e:
                    logging.error("action: receive_message | result: fail | error: {e}")
                    break
            if not has_finished:
                logging.error("action: receive_message | result: fail | error: has not finished")
                client_sock.close()
                return
            logging.info("action: receive_all_bets | result: success")
            protocol.send_ok(True)

            self._finished_clients += 1
            if self._finished_clients == self._total_clients:
                logging.info("action: sorteo | result: success")
                self.__perform_draw()

            case_id = protocol.receive_case()
            if case_id == 3:
                logging.info(f'action: receive_case | result: success | case_id: {case_id}')
                agency_id = protocol.receive_agency_id()
                logging.info(f'action: receive_agency_id | result: success | agency_id: {agency_id}')
                # Wait at the barrier
                logging.info(f'action: winners_status_BEFORE_barrier | result: in_progress | winners: {self._winners}')
                self._barrier.wait()
                logging.info(f'action: winners_status_AFTER_barrier | result: in_progress | winners: {self._winners}')
                winners = []
                if agency_id in self._winners:
                    winners = self._winners[agency_id]
                    logging.info(f'action: send_winners | result: in_progress | winners: {winners}')
                    protocol.send_winners(winners)
                    logging.info(f'action: winners_sent | result: success')
            else:
                logging.error("action: receive_message | result: fail | error: invalid case_id")
                client_sock.close()
                return
            
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

    def __perform_draw(self):
        """
        Perform a draw of the bets
        """
        bets = list(load_bets())
        logging.info(f'action: sorteo | result: in_progress | cantidad: {len(bets)}')
        for bet in bets:
            if has_won(bet):
                logging.info(f'action: sorteo | result: winner | winner_id: {bet.document}')
                self._winners[bet.agency].append(bet.document)
