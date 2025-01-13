import urllib.parse
import multipart
import multipart.multipart
import python_multipart
import io
import json


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
                 context: dict[str, str], cookies: dict[str, tuple[str, int]]):
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
            fields[field.field_name] = field

        def on_file(file: multipart.multipart.File):
            files[file.field_name] = file

        python_multipart.parse_form(self.headers, io.BytesIO(self.raw_body),
                                    on_field, on_file)
        return (fields, files)

    def parse_body_json(self):
        return json.loads(self.raw_body)


class Response:

    def __init__(self, status: int):
        self.status = status
        self.headers = []
        self.data = io.BytesIO()
