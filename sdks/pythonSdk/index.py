
from core.client import GetMeClient


from flask import Flask, request

server = Flask(__name__)




@server.route('/get', methods=['GET'], )
def get_value():
  key = request.args.get('key')
  if not key:
    return "Key parameter is required", 400
  return getMeClient.get(key), 200


if __name__ == '__main__':
  getMeClient = GetMeClient()
  
  server.run(port = 3001, debug=True)