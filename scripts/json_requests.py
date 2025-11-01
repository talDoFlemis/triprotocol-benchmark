#!/usr/bin/env python3

import socket
import json

host = "54.174.195.77"
port = 8081
aluno_id = "538349"


def tcp_request(message: str) -> str:
    print(f"Sending request: {message}")
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect((host, port))
    client_socket.sendall(bytes(message + "\n", "utf-8"))
    data = client_socket.recv(64 * 1024)
    string = data.decode("utf-8").strip()
    print(f"Received response: {string}")
    return string


def do_operation(op: str, token: str, **params):
    args = {
        "tipo": "operacao",
        "operacao": op,
        "token": token,
        "parametros": params,
    }

    result_string = json.dumps(args)
    response = tcp_request(result_string)

    return response


auth_request = json.dumps({"tipo": "autenticar", "aluno_id": aluno_id})
auth_response = json.loads(tcp_request(auth_request))
token = auth_response["token"]

do_operation("echo", token, mensagem="ola mundo")
do_operation("soma", token, numeros=[1, 2, 3])
do_operation("timestamp", token)
do_operation("status", token, detalhado=True)
do_operation("historico", token, limite=2)

logout_request = json.dumps({"tipo": "logout", "token": token})
tcp_request(logout_request)
