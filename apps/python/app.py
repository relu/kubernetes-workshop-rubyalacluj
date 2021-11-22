import os
from http.server import HTTPServer, BaseHTTPRequestHandler

PORT = 3000
name = os.getenv('NAME', __file__)

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type','text/html')
        self.end_headers()
        self.wfile.write(bytes(('Hello from %s' % name).encode('utf-8')))

httpd = HTTPServer(('0.0.0.0', PORT), Handler)
print('serving at port', PORT)
httpd.serve_forever()
