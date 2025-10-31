#!/usr/bin/env python3

import socket
import re

host = "54.174.195.77"
port = 8080
aluno_id = 538349


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
    args = ["OP", "operacao=" + op, "token=" + token]
    for key, value in params.items():
        args.append(f"{key}={value}")

    args.append("FIM")

    result_string = "|".join(args)
    response = tcp_request(result_string)

    return response


auth_request = f"AUTH|aluno_id={aluno_id}|FIM"

auth_response = tcp_request(auth_request)
match = re.search(r"token=([^|]*)", auth_response)

if not match:
    raise ValueError("Token not found in AUTH response")

token = match.group(1)


do_operation("echo", token, mensagem="ola mundo")
do_operation("soma", token, numeros=[1, 2, 3])
do_operation("timestamp", token)
do_operation("status", token, detalhado=True)
do_operation("historico", token, limite=10)

logout_request = f"LOGOUT|token={token}|FIM"
tcp_request(logout_request)
