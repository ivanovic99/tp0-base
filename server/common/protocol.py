import struct
from .serializer import deserialize_bets

AMOUNT_OF_BYTES = 4

class Protocol:
    def __init__(self, conn):
        self.conn = conn

    def receive_bets(self):
        """Receive a list of bets from the client."""
        length_data = self._recv_n_bytes(AMOUNT_OF_BYTES)
        if len(length_data) < AMOUNT_OF_BYTES:
            raise ValueError("Incomplete length data received")

        length = struct.unpack('!I', length_data)[0]

        data = self._recv_n_bytes(length)
        if len(data) < length:
            raise ValueError("Incomplete bets data received")

        bets = deserialize_bets(data)
        return bets

    def _recv_n_bytes(self, n):
        """Helper method to receive exactly n bytes or raise an error."""
        data = bytearray()
        while len(data) < n:
            packet = self.conn.recv(n - len(data))
            if not packet:
                raise ValueError("Connection closed before receiving all data")
            data.extend(packet)
        return bytes(data)

    def receive_total_bets(self):
        total_bets_bytes = self.conn.recv(AMOUNT_OF_BYTES)
        if len(total_bets_bytes) < AMOUNT_OF_BYTES:
            raise ValueError("Failed to receive total bets")
        total_bets = int.from_bytes(total_bets_bytes, byteorder='big')
        return total_bets
    
    def send_ok(self, ok):
        data = b'\x01' if ok else b'\x00'
        self.conn.sendall(data)
