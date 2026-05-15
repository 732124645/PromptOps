"""PromptOps runtime client: fetch prompts, render templates, hot-reload."""

import base64
import json
import os
import re
import socket
import ssl
import threading
import urllib.request
from urllib.parse import quote, urlparse

DEFAULT_TOKEN = "promptops-dev-token"
_VAR_RE = re.compile(r"{{\s*([\w.]+)\s*}}")


def render_template(content, variables=None):
    """Replace {{variable}} placeholders; unknown variables are left as-is."""
    variables = variables or {}

    def repl(match):
        key = match.group(1)
        return str(variables[key]) if key in variables else match.group(0)

    return _VAR_RE.sub(repl, "" if content is None else str(content))


class PromptOpsClient:
    """Fetches prompts by key, caches them, and (via watch) keeps them fresh."""

    def __init__(self, server, namespace="prod", token=DEFAULT_TOKEN):
        if not server:
            raise ValueError('PromptOps: "server" is required')
        self.server = server.rstrip("/")
        self.namespace = namespace
        self.token = token
        self._cache = {}
        self._on_update = None
        self._stop = threading.Event()
        self._sock = None
        self._thread = None

    def get_prompt(self, key, refresh=False):
        if not refresh and key in self._cache:
            return self._cache[key]
        url = "{}/api/sdk/prompts/{}?env={}".format(
            self.server, quote(key, safe=""), quote(self.namespace, safe="")
        )
        req = urllib.request.Request(url, headers={"Authorization": "Bearer " + self.token})
        with urllib.request.urlopen(req, timeout=10) as resp:
            prompt = json.loads(resp.read().decode("utf-8"))
        self._cache[key] = prompt
        return prompt

    def render(self, key, variables=None):
        return render_template(self.get_prompt(key)["content"], variables)

    def watch(self, on_update=None):
        """Open a WebSocket and refresh cached prompts on hot-reload events."""
        self._on_update = on_update
        self._stop.clear()
        self._thread = threading.Thread(target=self._ws_loop, daemon=True)
        self._thread.start()

    def close(self):
        self._stop.set()
        if self._sock is not None:
            try:
                self._sock.shutdown(socket.SHUT_RDWR)
            except OSError:
                pass

    # --- minimal WebSocket client (RFC 6455, text frames only) ---

    def _ws_loop(self):
        parsed = urlparse(self.server)
        host = parsed.hostname
        port = parsed.port or (443 if parsed.scheme == "https" else 80)
        try:
            sock = socket.create_connection((host, port), timeout=10)
            if parsed.scheme == "https":
                sock = ssl.create_default_context().wrap_socket(sock, server_hostname=host)
        except OSError:
            return
        self._sock = sock
        key = base64.b64encode(os.urandom(16)).decode("ascii")
        handshake = (
            "GET /ws HTTP/1.1\r\n"
            "Host: {}:{}\r\n"
            "Upgrade: websocket\r\n"
            "Connection: Upgrade\r\n"
            "Sec-WebSocket-Key: {}\r\n"
            "Sec-WebSocket-Version: 13\r\n\r\n"
        ).format(host, port, key)
        try:
            sock.sendall(handshake.encode("ascii"))
            response = self._read_until(sock, b"\r\n\r\n")
            if b" 101 " not in response.split(b"\r\n", 1)[0]:
                return
            while not self._stop.is_set():
                message = self._read_frame(sock)
                if message is None:
                    break
                self._handle(message)
        except OSError:
            pass
        finally:
            try:
                sock.close()
            except OSError:
                pass
            self._sock = None

    @staticmethod
    def _read_until(sock, delimiter):
        buf = b""
        while delimiter not in buf:
            chunk = sock.recv(1024)
            if not chunk:
                break
            buf += chunk
        return buf

    @staticmethod
    def _recv_exact(sock, count):
        buf = b""
        while len(buf) < count:
            chunk = sock.recv(count - len(buf))
            if not chunk:
                return None
            buf += chunk
        return buf

    def _read_frame(self, sock):
        header = self._recv_exact(sock, 2)
        if header is None:
            return None
        opcode = header[0] & 0x0F
        masked = header[1] & 0x80
        length = header[1] & 0x7F
        if length == 126:
            ext = self._recv_exact(sock, 2)
            length = int.from_bytes(ext, "big") if ext else 0
        elif length == 127:
            ext = self._recv_exact(sock, 8)
            length = int.from_bytes(ext, "big") if ext else 0
        mask = self._recv_exact(sock, 4) if masked else b""
        payload = self._recv_exact(sock, length) if length else b""
        if payload is None:
            return None
        if masked and mask:
            payload = bytes(b ^ mask[i % 4] for i, b in enumerate(payload))
        if opcode == 0x8:  # close
            return None
        if opcode == 0x1:  # text
            return payload.decode("utf-8", "replace")
        return ""  # ping / pong / binary are ignored

    def _handle(self, raw):
        if not raw:
            return
        try:
            event = json.loads(raw)
        except json.JSONDecodeError:
            return
        key = event.get("key")
        if key and key in self._cache:
            updated = self.get_prompt(key, refresh=True)
            if self._on_update:
                self._on_update({"key": key, "prompt": updated, "event": event.get("type")})
