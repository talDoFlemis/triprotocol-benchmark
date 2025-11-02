#!/usr/bin/env python3

import socket
import triprotocol_pb2


host = "54.174.195.77"
port = 8082
aluno_id = "538349"


def tcp_request(message: bytes) -> str:
    print(f"Sending request: {message}")
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect((host, port))
    client_socket.sendall(message)
    data = client_socket.recv(64 * 1024)
    string = data.decode("utf-8").strip()
    print(f"Received response: {string}")
    return string

# def do_operation(req: triprotocol_pb2.Requisicao, resp: proto) -> triprotocol_pb2.RespostaOk:
#     byteString = req.SerializeToString()
#     response = tcp_request(byteString)
#     resposta = triprotocol_pb2.Resposta()
#     resposta.ParseFromString(response.encode("utf-8"))

#     if not resposta.HasField("ok"):
#         raise Exception(f"Operation failed: {resposta}")

#     return resposta.ok



auth_request = triprotocol_pb2.Requisicao(auth=triprotocol_pb2.ComandoAuth(
    aluno_id=aluno_id,
))
byteString = auth_request.SerializeToString()
response = tcp_request(byteString)
auth_response = triprotocol_pb2.Resposta()
auth_response.ParseFromString(response.encode("utf-8"))

token = auth_response.ok
print(f"Auth token: {token}")

