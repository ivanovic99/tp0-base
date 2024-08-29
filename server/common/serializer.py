import json

class Bet:
    def __init__(self, nombre, apellido, dni, nacimiento, numero):
        self.nombre = nombre
        self.apellido = apellido
        self.dni = dni
        self.nacimiento = nacimiento
        self.numero = numero

    def to_dict(self):
        return {
            "nombre": self.nombre,
            "apellido": self.apellido,
            "dni": self.dni,
            "nacimiento": self.nacimiento,
            "numero": self.numero
        }

    @staticmethod
    def from_dict(data):
        return Bet(
            nombre=data["nombre"],
            apellido=data["apellido"],
            dni=data["dni"],
            nacimiento=data["nacimiento"],
            numero=data["numero"]
        )

def serialize_bet(bet):
    return json.dumps(bet.to_dict()).encode('utf-8')

def deserialize_bet(data):
    return Bet.from_dict(json.loads(data.decode('utf-8')))
