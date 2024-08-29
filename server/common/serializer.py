import json
from .utils import Bet

def serialize_bet(bet: Bet) -> bytes:
    return json.dumps(bet.to_dict()).encode('utf-8')

def deserialize_bet(data: bytes) -> Bet:
    return Bet.from_dict(json.loads(data.decode('utf-8')))
