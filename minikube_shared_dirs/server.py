import os
from flask import Flask
server = Flask(__name__)

@server.route("/")
def hello():
    return 'Content of hosthome: ' + ','.join(os.listdir('/opt/hosthome')) + '\n'

if __name__ == "__main__":
    server.run(host='0.0.0.0')
