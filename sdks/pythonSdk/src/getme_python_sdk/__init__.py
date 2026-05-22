from .core.client import GetMeClient

__all__ = ["GetMeClient"]

# ==========================================
# Example Usage (Flask Server integration):
# ==========================================
# from flask import Flask, request
# from getme import GetMeClient
#
# server = Flask(__name__)
# getMeClient = GetMeClient()
#
# @server.route('/get', methods=['GET'])
# def get_value():
#     key = request.args.get('key')
#     if not key:
#         return "Key parameter is required", 400
#     return getMeClient.get(key), 200
#
# if __name__ == '__main__':
#     server.run(port=3001, debug=True)
