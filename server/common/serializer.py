import datetime
from .utils import Bet

def deserialize_bets(data: bytes) -> list[Bet]:
    bets = []
    offset = 0
    while offset < len(data):
        # Deserialize each Bet object
        agency = int.from_bytes(data[offset:offset + 4], byteorder='big')
        offset += 4

        first_name_length = int.from_bytes(data[offset:offset + 4], byteorder='big')
        offset += 4
        first_name = data[offset:offset + first_name_length].decode('utf-8')
        offset += first_name_length

        last_name_length = int.from_bytes(data[offset:offset + 4], byteorder='big')
        offset += 4
        last_name = data[offset:offset + last_name_length].decode('utf-8')
        offset += last_name_length

        document_length = int.from_bytes(data[offset:offset + 4], byteorder='big')
        offset += 4
        document = data[offset:offset + document_length].decode('utf-8')
        offset += document_length

        birthdate_length = int.from_bytes(data[offset:offset + 4], byteorder='big')
        offset += 4
        birthdate = data[offset:offset + birthdate_length].decode('utf-8')
        birthdate = datetime.date.fromisoformat(birthdate)
        offset += birthdate_length

        number = int.from_bytes(data[offset:offset + 4], byteorder='big')
        offset += 4

        bet = Bet(
            agency=agency,
            first_name=first_name,
            last_name=last_name,
            document=document,
            birthdate=birthdate.isoformat(),
            number=number
        )
        bets.append(bet)
    return bets

def serialize_winners(winners: list[str]) -> bytes:
    buf = bytearray()
    for winner in winners:
        winner_bytes = winner.encode('utf-8')
        winner_length = len(winner_bytes)
        buf.extend(winner_length.to_bytes(4, byteorder='big'))
        buf.extend(winner_bytes)
    return bytes(buf)
