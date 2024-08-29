import socket
from .serializer import serialize_bet, deserialize_bet

class Protocol:
    def __init__(self, conn):
        self.conn = conn

    def send_bet(self, bet):
        data = serialize_bet(bet)
        self.conn.sendall(data + b'\n')

    def receive_bet(self):
        # Change the way we receive data, we need to remove the trailing newline character. 1024 is the maximum size of the data we can receive but it is fixed, we should change it to a variable that can be received.
        data = self.conn.recv(1024).rstrip()
        return deserialize_bet(data)
