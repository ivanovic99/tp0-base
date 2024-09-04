import struct
from .serializer import deserialize_bets, serialize_winners

AMOUNT_OF_BYTES = 4

class Protocol:
    def __init__(self, conn):
        self.conn = conn

    def receive_case(self):
        """Receive a case ID from the client."""
        data = self._recv_n_bytes(1)
        if len(data) < 1:
            raise ValueError("Incomplete case ID received")
        case_id = data[0]
        return case_id

    def receive_agency_id(self):
        """Receive an agency ID from the client."""
        data = self._recv_n_bytes(AMOUNT_OF_BYTES)
        if len(data) < AMOUNT_OF_BYTES:
            raise ValueError("Incomplete agency ID received")
        agency_id = struct.unpack('!I', data)[0]
        return agency_id
    
    def send_winners(self, documents_of_winners):
        """Send the documents of the winners to the client."""
        data = serialize_winners(documents_of_winners)
        length = len(data).to_bytes(AMOUNT_OF_BYTES, byteorder='big')
        self._send_all(length)
        self._send_all(data)

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

    def _send_all(self, data):
        """Helper method to send all bytes or raise an error."""
        total_sent = 0
        while total_sent < len(data):
            sent = self.conn.send(data[total_sent:])
            if sent == 0:
                raise ValueError("Connection closed before sending all data")
            total_sent += sent

    def receive_total_bets(self):
        total_bets_bytes = self._recv_n_bytes(AMOUNT_OF_BYTES)
        if len(total_bets_bytes) < AMOUNT_OF_BYTES:
            raise ValueError("Failed to receive total bets")
        total_bets = int.from_bytes(total_bets_bytes, byteorder='big')
        return total_bets
    
    def send_ok(self, ok):
        data = b'\x02' if ok else b'\x00'
        self._send_all(data)
