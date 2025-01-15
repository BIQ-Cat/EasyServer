import urllib.parse
import multipart
import multipart.multipart
import python_multipart
import io
import json

import typing


class URL:

    def __init__(self, scheme: str, host: str, port: str, path: str,
                 query: dict[str, list[str]]):
        self.scheme = scheme
        self.host = host
        self.port = port
        self.path = path
        self.query = query


class Request:

    def __init__(self, method: str, url: URL, protocol: str,
                 headers: dict[str, bytes], body: bytes,
                 context: dict[str, bytes], cookies: dict[str, tuple[str, int]]):
        self.method = method
        self.url = url
        self.protocol = protocol
        self.headers = headers
        self.raw_body = body
        self.context = context
        self.cookies = cookies

    def parse_body_urlencoded(self):
        return urllib.parse.parse_qs(self.raw_body)

    def parse_body_multipart(self):
        fields = {}  # type: dict[bytes, multipart.multipart.Field]
        files = {}  # type: dict[bytes, multipart.multipart.File]

        def on_field(field: multipart.multipart.Field):
            
            fields[typing.cast(bytes, field.field_name)] = field

        def on_file(file: multipart.multipart.File):
            files[typing.cast(bytes, file.field_name)] = file

        python_multipart.parse_form(self.headers, io.BytesIO(self.raw_body),
                                    on_field, on_file)  # type: ignore
        return (fields, files)

    def parse_body_json(self):
        return json.loads(self.raw_body)


class Response:

    def __init__(self, status = 200):
        self.status = status
        self.headers = {}  # type: dict[str, bytes]
        self.data = b""

type Handler = typing.Callable[[Request], Response]
type MiddlewareFunc = typing.Callable[[Handler], Handler]

class Controller:

    def __init__(self,
                 handler: Handler,
                 methods: typing.Optional[list[bytes]] = None,
                 headers: typing.Optional[dict[str, bytes]] = None,
                 schemas: typing.Optional[list[bytes]] = None,
                 data: typing.Optional[dict[str, typing.Any]] = None):
        self.handler = handler
        self.methods = methods
        self.headers = headers
        self.schemas = schemas
        self.data = data

    def __call__(self, req: Request):
        return self.handler(req)


class Module:
    def __init__(self, route: dict[str, Controller], middlewares: list[MiddlewareFunc] = [], is_explict = True):
        self.route = route
        self.middlewares = middlewares
        self.is_expict = is_explict