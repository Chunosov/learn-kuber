import os
from flask import Flask
server = Flask(__name__)

_HOST = '0.0.0.0'
_PORT = int(os.environ.get('PORT', 5000))

# NB: when knative:
# Taking a port number from the PORT variable is important
# when service is run under knative
# because this is knative who provides a port number.
# If we just hardcode some port, the service will be unavailable.
# knative will probe the port several times
# and if we don't listen to it, the service we be marked as invalid.

@server.route("/")
def hello():
    return f"Hello World! I'm on internal port {_PORT}\n"

if __name__ == "__main__":
    print(f'Running server on {_HOST}:{_PORT}')
    server.run(debug=True, host=_HOST, port=_PORT)
