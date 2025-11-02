from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Requisicao(_message.Message):
    __slots__ = ("auth", "operacao", "info", "logout")
    AUTH_FIELD_NUMBER: _ClassVar[int]
    OPERACAO_FIELD_NUMBER: _ClassVar[int]
    INFO_FIELD_NUMBER: _ClassVar[int]
    LOGOUT_FIELD_NUMBER: _ClassVar[int]
    auth: ComandoAuth
    operacao: ComandoOperacao
    info: ComandoInfo
    logout: ComandoLogout
    def __init__(self, auth: _Optional[_Union[ComandoAuth, _Mapping]] = ..., operacao: _Optional[_Union[ComandoOperacao, _Mapping]] = ..., info: _Optional[_Union[ComandoInfo, _Mapping]] = ..., logout: _Optional[_Union[ComandoLogout, _Mapping]] = ...) -> None: ...

class Resposta(_message.Message):
    __slots__ = ("ok", "erro")
    OK_FIELD_NUMBER: _ClassVar[int]
    ERRO_FIELD_NUMBER: _ClassVar[int]
    ok: RespostaOk
    erro: RespostaErro
    def __init__(self, ok: _Optional[_Union[RespostaOk, _Mapping]] = ..., erro: _Optional[_Union[RespostaErro, _Mapping]] = ...) -> None: ...

class ComandoAuth(_message.Message):
    __slots__ = ("aluno_id", "timestamp_cliente")
    ALUNO_ID_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_CLIENTE_FIELD_NUMBER: _ClassVar[int]
    aluno_id: str
    timestamp_cliente: str
    def __init__(self, aluno_id: _Optional[str] = ..., timestamp_cliente: _Optional[str] = ...) -> None: ...

class ComandoOperacao(_message.Message):
    __slots__ = ("token", "operacao", "parametros")
    class ParametrosEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    TOKEN_FIELD_NUMBER: _ClassVar[int]
    OPERACAO_FIELD_NUMBER: _ClassVar[int]
    PARAMETROS_FIELD_NUMBER: _ClassVar[int]
    token: str
    operacao: str
    parametros: _containers.ScalarMap[str, str]
    def __init__(self, token: _Optional[str] = ..., operacao: _Optional[str] = ..., parametros: _Optional[_Mapping[str, str]] = ...) -> None: ...

class ComandoInfo(_message.Message):
    __slots__ = ("tipo",)
    TIPO_FIELD_NUMBER: _ClassVar[int]
    tipo: str
    def __init__(self, tipo: _Optional[str] = ...) -> None: ...

class ComandoLogout(_message.Message):
    __slots__ = ("token",)
    TOKEN_FIELD_NUMBER: _ClassVar[int]
    token: str
    def __init__(self, token: _Optional[str] = ...) -> None: ...

class RespostaOk(_message.Message):
    __slots__ = ("comando", "dados", "timestamp")
    class DadosEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    COMANDO_FIELD_NUMBER: _ClassVar[int]
    DADOS_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    comando: str
    dados: _containers.ScalarMap[str, str]
    timestamp: str
    def __init__(self, comando: _Optional[str] = ..., dados: _Optional[_Mapping[str, str]] = ..., timestamp: _Optional[str] = ...) -> None: ...

class RespostaErro(_message.Message):
    __slots__ = ("comando", "mensagem", "timestamp", "detalhes")
    class DetalhesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    COMANDO_FIELD_NUMBER: _ClassVar[int]
    MENSAGEM_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    DETALHES_FIELD_NUMBER: _ClassVar[int]
    comando: str
    mensagem: str
    timestamp: str
    detalhes: _containers.ScalarMap[str, str]
    def __init__(self, comando: _Optional[str] = ..., mensagem: _Optional[str] = ..., timestamp: _Optional[str] = ..., detalhes: _Optional[_Mapping[str, str]] = ...) -> None: ...

class DadosAuth(_message.Message):
    __slots__ = ("token", "nome", "matricula", "timestamp", "timeout_segundos")
    TOKEN_FIELD_NUMBER: _ClassVar[int]
    NOME_FIELD_NUMBER: _ClassVar[int]
    MATRICULA_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_SEGUNDOS_FIELD_NUMBER: _ClassVar[int]
    token: str
    nome: str
    matricula: str
    timestamp: str
    timeout_segundos: int
    def __init__(self, token: _Optional[str] = ..., nome: _Optional[str] = ..., matricula: _Optional[str] = ..., timestamp: _Optional[str] = ..., timeout_segundos: _Optional[int] = ...) -> None: ...

class ResultadoEcho(_message.Message):
    __slots__ = ("mensagem_original", "mensagem_eco", "hash_md5", "tamanho_mensagem", "timestamp_servidor")
    MENSAGEM_ORIGINAL_FIELD_NUMBER: _ClassVar[int]
    MENSAGEM_ECO_FIELD_NUMBER: _ClassVar[int]
    HASH_MD5_FIELD_NUMBER: _ClassVar[int]
    TAMANHO_MENSAGEM_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_SERVIDOR_FIELD_NUMBER: _ClassVar[int]
    mensagem_original: str
    mensagem_eco: str
    hash_md5: str
    tamanho_mensagem: int
    timestamp_servidor: str
    def __init__(self, mensagem_original: _Optional[str] = ..., mensagem_eco: _Optional[str] = ..., hash_md5: _Optional[str] = ..., tamanho_mensagem: _Optional[int] = ..., timestamp_servidor: _Optional[str] = ...) -> None: ...

class ResultadoSoma(_message.Message):
    __slots__ = ("numeros_originais", "quantidade", "soma", "media", "maximo", "minimo", "timestamp_calculo")
    NUMEROS_ORIGINAIS_FIELD_NUMBER: _ClassVar[int]
    QUANTIDADE_FIELD_NUMBER: _ClassVar[int]
    SOMA_FIELD_NUMBER: _ClassVar[int]
    MEDIA_FIELD_NUMBER: _ClassVar[int]
    MAXIMO_FIELD_NUMBER: _ClassVar[int]
    MINIMO_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_CALCULO_FIELD_NUMBER: _ClassVar[int]
    numeros_originais: _containers.RepeatedScalarFieldContainer[float]
    quantidade: int
    soma: float
    media: float
    maximo: float
    minimo: float
    timestamp_calculo: str
    def __init__(self, numeros_originais: _Optional[_Iterable[float]] = ..., quantidade: _Optional[int] = ..., soma: _Optional[float] = ..., media: _Optional[float] = ..., maximo: _Optional[float] = ..., minimo: _Optional[float] = ..., timestamp_calculo: _Optional[str] = ...) -> None: ...

class ResultadoTimestamp(_message.Message):
    __slots__ = ("timestamp_unix", "timestamp_iso", "timestamp_formatado", "ano", "mes", "dia", "hora", "minuto", "segundo", "microsegundo")
    TIMESTAMP_UNIX_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_ISO_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FORMATADO_FIELD_NUMBER: _ClassVar[int]
    ANO_FIELD_NUMBER: _ClassVar[int]
    MES_FIELD_NUMBER: _ClassVar[int]
    DIA_FIELD_NUMBER: _ClassVar[int]
    HORA_FIELD_NUMBER: _ClassVar[int]
    MINUTO_FIELD_NUMBER: _ClassVar[int]
    SEGUNDO_FIELD_NUMBER: _ClassVar[int]
    MICROSEGUNDO_FIELD_NUMBER: _ClassVar[int]
    timestamp_unix: float
    timestamp_iso: str
    timestamp_formatado: str
    ano: int
    mes: int
    dia: int
    hora: int
    minuto: int
    segundo: int
    microsegundo: int
    def __init__(self, timestamp_unix: _Optional[float] = ..., timestamp_iso: _Optional[str] = ..., timestamp_formatado: _Optional[str] = ..., ano: _Optional[int] = ..., mes: _Optional[int] = ..., dia: _Optional[int] = ..., hora: _Optional[int] = ..., minuto: _Optional[int] = ..., segundo: _Optional[int] = ..., microsegundo: _Optional[int] = ...) -> None: ...

class StatusServidor(_message.Message):
    __slots__ = ("status", "operacoes_processadas", "sessoes_ativas", "tempo_ativo", "versao", "estatisticas_banco", "sessoes_detalhes", "metricas")
    class EstatisticasBancoEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    class SessoesDetalhesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    class MetricasEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: float
        def __init__(self, key: _Optional[str] = ..., value: _Optional[float] = ...) -> None: ...
    STATUS_FIELD_NUMBER: _ClassVar[int]
    OPERACOES_PROCESSADAS_FIELD_NUMBER: _ClassVar[int]
    SESSOES_ATIVAS_FIELD_NUMBER: _ClassVar[int]
    TEMPO_ATIVO_FIELD_NUMBER: _ClassVar[int]
    VERSAO_FIELD_NUMBER: _ClassVar[int]
    ESTATISTICAS_BANCO_FIELD_NUMBER: _ClassVar[int]
    SESSOES_DETALHES_FIELD_NUMBER: _ClassVar[int]
    METRICAS_FIELD_NUMBER: _ClassVar[int]
    status: str
    operacoes_processadas: int
    sessoes_ativas: int
    tempo_ativo: float
    versao: str
    estatisticas_banco: _containers.ScalarMap[str, str]
    sessoes_detalhes: _containers.ScalarMap[str, str]
    metricas: _containers.ScalarMap[str, float]
    def __init__(self, status: _Optional[str] = ..., operacoes_processadas: _Optional[int] = ..., sessoes_ativas: _Optional[int] = ..., tempo_ativo: _Optional[float] = ..., versao: _Optional[str] = ..., estatisticas_banco: _Optional[_Mapping[str, str]] = ..., sessoes_detalhes: _Optional[_Mapping[str, str]] = ..., metricas: _Optional[_Mapping[str, float]] = ...) -> None: ...

class InfoServidor(_message.Message):
    __slots__ = ("nome", "versao", "host", "port", "protocolo", "formato", "operacoes_disponiveis", "total_operacoes")
    NOME_FIELD_NUMBER: _ClassVar[int]
    VERSAO_FIELD_NUMBER: _ClassVar[int]
    HOST_FIELD_NUMBER: _ClassVar[int]
    PORT_FIELD_NUMBER: _ClassVar[int]
    PROTOCOLO_FIELD_NUMBER: _ClassVar[int]
    FORMATO_FIELD_NUMBER: _ClassVar[int]
    OPERACOES_DISPONIVEIS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_OPERACOES_FIELD_NUMBER: _ClassVar[int]
    nome: str
    versao: str
    host: str
    port: int
    protocolo: str
    formato: str
    operacoes_disponiveis: _containers.RepeatedScalarFieldContainer[str]
    total_operacoes: int
    def __init__(self, nome: _Optional[str] = ..., versao: _Optional[str] = ..., host: _Optional[str] = ..., port: _Optional[int] = ..., protocolo: _Optional[str] = ..., formato: _Optional[str] = ..., operacoes_disponiveis: _Optional[_Iterable[str]] = ..., total_operacoes: _Optional[int] = ...) -> None: ...

class HistoricoOperacao(_message.Message):
    __slots__ = ("operacao", "parametros", "resultado", "timestamp", "sucesso")
    class ParametrosEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    class ResultadoEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    OPERACAO_FIELD_NUMBER: _ClassVar[int]
    PARAMETROS_FIELD_NUMBER: _ClassVar[int]
    RESULTADO_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    SUCESSO_FIELD_NUMBER: _ClassVar[int]
    operacao: str
    parametros: _containers.ScalarMap[str, str]
    resultado: _containers.ScalarMap[str, str]
    timestamp: str
    sucesso: bool
    def __init__(self, operacao: _Optional[str] = ..., parametros: _Optional[_Mapping[str, str]] = ..., resultado: _Optional[_Mapping[str, str]] = ..., timestamp: _Optional[str] = ..., sucesso: _Optional[bool] = ...) -> None: ...

class HistoricoAluno(_message.Message):
    __slots__ = ("aluno_id", "operacoes", "total")
    ALUNO_ID_FIELD_NUMBER: _ClassVar[int]
    OPERACOES_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    aluno_id: str
    operacoes: _containers.RepeatedCompositeFieldContainer[HistoricoOperacao]
    total: int
    def __init__(self, aluno_id: _Optional[str] = ..., operacoes: _Optional[_Iterable[_Union[HistoricoOperacao, _Mapping]]] = ..., total: _Optional[int] = ...) -> None: ...
