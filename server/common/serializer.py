import json
from .utils import Bet

def deserialize_bets(data: bytes) -> list[Bet]:
    return [Bet.from_dict(bet) for bet in json.loads(data.decode('utf-8'))]

def serialize_winners(winners: list[str]) -> bytes:
    return json.dumps(winners).encode('utf-8')
