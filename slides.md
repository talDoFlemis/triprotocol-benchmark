---
author: Said Cavalcante Rodrigues
date: MMMM DD, YYYY
paging: Slide %d / %d
---

# TriProtocol Benchmark

Um benchmark de serialização de protocolo e servidor de validação que suporta três formatos diferentes:
**String**, **JSON**, e **Protocol Buffers**.

**Aluno:** Said Cavalcante Rodrigues
**Matrícula:** 538349

---

## 🚀 Features

- **Suporte Multi-Protocolo**: Testa e compara três protocolos de serialização.
  - Serialização baseada em String
  - Serialização JSON
  - Serialização binária Protocol Buffers (protobuf)
- **Arquitetura em Camadas**: Separação clara entre as camadas de Transporte e Apresentação.
- **TUI Interativa**: Interface de terminal (TUI) construída com [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- **Modelo de Domínio Universal**: Um único conjunto de entidades de domínio funciona com todos os serializadores.
- **Observabilidade**: Instrumentação OpenTelemetry para tracing e métricas.
- **Scripts de Cliente Python**: Scripts prontos para testar os diferentes protocolos.

---

## 📁 Estrutura do Projeto

```

.
├── main.go                 # Ponto de entrada
├── tui.go                  # Implementação da TUI
├── app_layer.go            # Camada de aplicação com cliente genérico
├── round_tripper.go        # Abstração da camada de transporte (TCP)
├── domain.go               # Modelos e tipos de domínio
├── dto.go                  # Data transfer objects
├── serde.go                # Interface de serialização (Apresentação)
├── string_serde.go         # Implementação do protocolo String
├── json_serde.go           # Implementação do protocolo JSON
├── protobuf_serde.go       # Implementação do protocolo Protobuf
├── proto/
│   └── triprotocol.proto   # Definições do Protocol Buffer
├── protogenerated/         # Código protobuf gerado
├── scripts/                # Scripts de cliente Python
│   ├── proto_requests.py   # Cliente Protobuf
│   ├── json_requests.py    # Cliente JSON
│   └── string_request.py   # Cliente String
├── Dockerfile              # Definição do contêiner
├── Taskfile.yaml           # Automação de tarefas
└── base.yaml               # Configuração base

```

---

## 🏗️ Visão Geral da Arquitetura

Este projeto implementa uma separação de interesses abstraindo duas camadas críticas da pilha de rede: Transporte e Apresentação.

```
+---------------------+
|  Application Layer  |
|  Domain -> AppLayer |
+---------------------+
|
+---------+---------+
|                   |
v                   v
+---------------------+   +---------------------+
| Presentation Layer  |   |   Transport Layer   |
|---------------------|   |---------------------|
|  Serde Interface    |   |  RoundTripper I.    |
|   |                 |   |    |                |
|   v                 |   |    v                |
| [Str, JSON, Proto]  |   |  [TCP, UDP]         |
|   |                 |   |    ^                |
|   +-----------------|---+----| (AppLayer -> RT)
+---------------------+   +---------------------+

```

---

## 1. Abstração da Camada de Transporte

A interface **RoundTripper** abstrai o protocolo de transporte (TCP/UDP).

```
+--------------------------+
|  <<interface>>           |
|      RoundTripper        |
+--------------------------+
| +RequestReply(...)       |
+--------------------------+
^
|
+---------+---------------------------------+
|                                           |
v                                           v
+----------------------+  +---------------------+
|   TCPRoundTripper    |  |   UDPRoundTripper   |
+----------------------+  +---------------------+
| -DialTimeout         |  | +RequestReply(...)  |
| -WriteTimeout        |  +---------------------+
| -ReadTimeout         |
| +RequestReply(...)   |
+----------------------+

```

**Benefícios:**
- Trocar entre TCP e UDP sem alterar a lógica da aplicação.
- Gerenciamento centralizado de timeout e conexão.

---

## 2. Abstração da Camada de Apresentação

A interface **Serde** (Serializer/Deserializer) lida com a tradução entre entidades de domínio e formatos de protocolo.

```

+--------------------------+
|  <<interface>>           |
|         Serde            |
+--------------------------+
| +Marshal(...)            |
| +Unmarshal(...)          |
+--------------------------+
^
|
+-------------------+----------------+
|                   |                |
v                   v                v
+-------------+ +-----------+ +---------------+
| StringSerde | | JSONSerde | | ProtobufSerde |
+-------------+ +-----------+ +---------------+
| +Marshal    | | +Marshal  | | +Marshal      |
| +Unmarshal  | | +Unmarshal| | +Unmarshal    |
+-------------+ +-----------+ +---------------+

```
**Benefícios:**
- Um único modelo de domínio funciona com todos os formatos.
- Código da aplicação agnóstico ao protocolo.

---

## 3. Fluxo de Requisição

```

Client      AppLayerClient      Serde Impl.       RoundTripper      Remote Server
|               |                   |                 |                 |
| Do(...)       |                   |                 |                 |
|-------------->|                   |                 |                 |
|               | Wrap Request      |                 |                 |
|               |------------------>|                 |                 |
|               | Marshal(req)      |                 |                 |
|               |------------------>|                 |                 |
|               |                   | Domain -> Proto |                 |
|               |     []byte        |                 |                 |
|               |<------------------|                 |                 |
|               | RequestReply(...) |                 |                 |
|               |------------------------------------>|                 |
|               |                   |                 | Send TCP        |
|               |                   |                 |---------------->|
|               |                   |                 | Response bytes  |
|               |                   |                 |<----------------|
|               |     []byte        |                 |                 |
|               |<------------------------------------|                 |
|               | Unmarshal(bytes)  |                 |                 |
|               |------------------>|                 |                 |
|               |                   | Proto -> Domain |                 |
|               | Populated resp    |                 |                 |
|               |<------------------|                 |                 |
| Response      |                   |                 |                 |
|<--------------|                   |                 |                 |

```

---

## 4. Padrão de Serializador Universal

Todos os três serializadores funcionam com as **mesmas entidades de domínio**.

```

+----------------+
|  Domain Layer  |
|----------------|   +--------------------+
| [AuthRequest]  |-->| Presentation Layer |
| [OpRequest]    |-->| (PLR)              |
| [LogoutRequest]|-->| Token + Body       |
+----------------+   +--------------------+
|
+-------------------------------------+------------------+
|                                     |                  |
v                                     v                  v
+-----------------------+ +--------------------+ +---------------------+
| String Format (K=V)   | | JSON Format (JSON) | | Protobuf Fmt (Bin)  |
+-----------------------+ +--------------------+ +---------------------+
|                                     |                  |
+-------------------------------------+------------------+
|
v
[Wire Protocol]

````
---

## 🔑 Componentes Chave

### Formatos de Serialização

1.  **Protocolo String**: Formato simples baseado em `CHAVE=VALOR` para comunicação leve.
2.  **JSON**: Formato JSON legível por humanos com codificação UTF-8.
3.  **Protocol Buffers**: Serialização binária para transmissão de rede eficiente.

### Features Principais

- **Validação**: Validação de requisições usando `go-playground/validator`.
- **Observabilidade**: Integração OpenTelemetry para tracing distribuído.
- **Configuração**: Configuração flexível com Viper.
- **Type Safety**: `AppLayerClient` genérico com verificação de tipo em tempo de compilação.

---

## 🐳 Quick Start com Docker

Execute a TUI interativa diretamente usando Docker:

```bash
docker run --rm --pull always -it 
  ghcr.io/taldoflemis/triprotocol-benchmark/tui:latest
````

Este comando irá:

  - Baixar a imagem mais recente do GitHub Container Registry.
  - Executar a TUI (Interface de Usuário do Terminal) interativa.
  - Remover automaticamente o contêiner ao sair.

---

## 🐍 Clientes Python

O projeto inclui scripts Python para testar cada protocolo:

**Protocol Buffers Client:**

```bash
cd scripts
python3 proto_requests.py
```

**JSON Client:**

```bash
python3 json_requests.py
```

**String Protocol Client:**

```bash
python3 string_request.py
```

---

# Obrigado!