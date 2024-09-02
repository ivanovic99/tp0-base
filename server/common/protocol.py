from .serializer import deserialize_bet
import struct

AMMOUNT_OF_BYTES = 4

class Protocol:
    def __init__(self, conn):
        self.conn = conn


    def receive_bet(self):
        length_data = self.conn.recv(AMMOUNT_OF_BYTES)
        if len(length_data) < AMMOUNT_OF_BYTES:
            raise ValueError("Incomplete length data received")

        length = struct.unpack('!I', length_data)[0]

        data = self.conn.recv(length)
        if len(data) < length:
            raise ValueError("Incomplete bet data received")

        bet = deserialize_bet(data)
        return bet
