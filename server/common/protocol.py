from .serializer import deserialize_bet
import struct

AMMOUNT_OF_BYTES = 4

class Protocol:
    def __init__(self, conn):
        self.conn = conn

    def receive_bet(self):
        length_data = self._recv_all(AMMOUNT_OF_BYTES)
        length = struct.unpack('!I', length_data)[0]
        data = self._recv_all(length)
        bet = deserialize_bet(data)
        return bet

    def _recv_all(self, num_bytes):
        data = b''
        while len(data) < num_bytes:
            packet = self.conn.recv(num_bytes - len(data))
            if not packet:
                raise ValueError("Incomplete data received")
            data += packet
        return data
