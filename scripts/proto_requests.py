#!/usr/bin/env python3

import socket
import triprotocol_pb2
import struct
from datetime import datetime


host = "3.88.99.255"
port = 8082
aluno_id = "538349"


def tcp_request(message: bytes) -> bytes:
    tamanho = len(message)
    header = struct.pack("!I", tamanho)
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect((host, port))
    client_socket.sendall(header + message)
    data = client_socket.recv(64 * 1024)
    header_size = struct.unpack("!I", data[:4])[0]
    data = data[4 : 4 + header_size]
    print(f"Header size: {header_size}\n")
    return data


def do_request(request: triprotocol_pb2.Requisicao) -> triprotocol_pb2.Resposta:
    print(f"Request: {request}\n")
    byteString = request.SerializeToString()
    data = tcp_request(byteString)
    response = triprotocol_pb2.Resposta()
    response.ParseFromString(data)
    print(f"Response: {response}\n")
    return response


auth_response = do_request(
    triprotocol_pb2.Requisicao(
        auth=triprotocol_pb2.ComandoAuth(
            aluno_id=aluno_id, timestamp_cliente=datetime.now().isoformat()
        )
    )
)
token = auth_response.ok.dados["token"]

print(f"Received token: {token}\n")

do_request(
    triprotocol_pb2.Requisicao(
        operacao=triprotocol_pb2.ComandoOperacao(
            operacao="echo",
            parametros={"mensagem": "Hello, World!"},
            token=token,
        )
    )
)

do_request(
    triprotocol_pb2.Requisicao(
        operacao=triprotocol_pb2.ComandoOperacao(
            operacao="soma",
            parametros={
                "numeros": "1,2,3",
            },
            token=token,
        )
    )
)

do_request(
    triprotocol_pb2.Requisicao(
        operacao=triprotocol_pb2.ComandoOperacao(
            operacao="timestamp",
            parametros={},
            token=token,
        )
    )
)

do_request(
    triprotocol_pb2.Requisicao(
        operacao=triprotocol_pb2.ComandoOperacao(
            operacao="status",
            parametros={
                "detalhado": "true",
            },
            token=token,
        )
    )
)

do_request(
    triprotocol_pb2.Requisicao(
        operacao=triprotocol_pb2.ComandoOperacao(
            operacao="historico",
            parametros={
                "limite": "1",
            },
            token=token,
        )
    )
)

do_request(
    triprotocol_pb2.Requisicao(
        logout=triprotocol_pb2.ComandoLogout(
            token=token,
        )
    )
)
