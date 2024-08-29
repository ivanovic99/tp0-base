import socket
from .serializer import serialize_bet, deserialize_bet

class Protocol:
    def __init__(self, conn):
        self.conn = conn

    def send_bet(self, bet):
        data = serialize_bet(bet)
        self.conn.sendall(data + b'\n')

    def receive_bet(self):
        data = self.conn.recv(1024).rstrip()
        return deserialize_bet(data)
