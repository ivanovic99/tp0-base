import datetime
from .utils import Bet

def deserialize_bet(data: bytes) -> Bet:
    offset = 0
    
    # Extract the agency (first 4 bytes)
    agency = int.from_bytes(data[offset:offset + 4], byteorder='big')
    offset += 4
    
    # Extract the first name length (next 4 bytes)
    first_name_length = int.from_bytes(data[offset:offset + 4], byteorder='big')
    offset += 4
    # Extract the first name
    first_name = data[offset:offset + first_name_length].decode('utf-8')
    offset += first_name_length
    
    # Extract the last name length (next 4 bytes)
    last_name_length = int.from_bytes(data[offset:offset + 4], byteorder='big')
    offset += 4
    # Extract the last name
    last_name = data[offset:offset + last_name_length].decode('utf-8')
    offset += last_name_length
    
    # Extract the document length (next 4 bytes)
    document_length = int.from_bytes(data[offset:offset + 4], byteorder='big')
    offset += 4
    # Extract the document
    document = data[offset:offset + document_length].decode('utf-8')
    offset += document_length
    
    # Extract the birthdate length (next 4 bytes)
    birthdate_length = int.from_bytes(data[offset:offset + 4], byteorder='big')
    offset += 4
    # Extract the birthdate
    birthdate = data[offset:offset + birthdate_length].decode('utf-8')
    birthdate = datetime.date.fromisoformat(birthdate)
    offset += birthdate_length
    
    # Extract the number (last 4 bytes)
    number = int.from_bytes(data[offset:offset + 4], byteorder='big')
    
    # Create and return the Bet object
    return Bet(
        agency=agency,
        first_name=first_name,
        last_name=last_name,
        document=document,
        birthdate=birthdate.isoformat(),
        number=number
    )
