import os
import requests_unixsocket
from dotenv import load_dotenv
import json
from urllib.parse import quote_plus

class GetMeClient:
    """
    A client for interacting with the getMe daemon over a Unix Domain Socket.
    """
    def __init__(self):
        load_dotenv()
        self.unix_session = requests_unixsocket.Session()
        self.sock_path = quote_plus(os.getenv("GETME_SOCK_PATH", "/tmp/getMeStore/getMe.sock"))

    def put(self, key, value):
        """
        Puts a key-value pair into the store. The value is sent in the request body.

        :param key: The key for the entry.
        :param value: The value to store.
        :raises Exception: If the server returns a non-200 status code.
        """

        url = f"http+unix://{self.sock_path}/put"
        
        
        payload = {
            "key": key,
            "value": value
        }
        headers = {
            'Content-Type': 'application/json'
        }

        response = self.unix_session.post(url, data=json.dumps(payload), headers=headers)
        if response.status_code != 200:
            raise Exception(f"Failed to put key-value pair: {response.text}")

        return response.json()

    def get(self, key) -> str:
        """
        Gets a value for a given key from the store.

        :param key: The key to retrieve.
        :return: The value associated with the key as a string.
        :raises Exception: If the server returns a non-200 status code.
        """
        resp = self.unix_session.get(f"http+unix://{self.sock_path}/get",
                                     params={"key": key})
        
        if resp.status_code != 200:
            raise Exception(f"Failed to get value for key '{key}': {resp.text}")
        
        print( resp.json())
        return resp.json()

    def delete(self, key):
        """
        Deletes a key from the store.

        :param key: The key to delete.
        :raises Exception: If the server returns a non-200 status code.
        """
        url = f"http+unix://{self.sock_path}/delete"
        response = self.unix_session.delete(url, params={"key": key})
        if response.status_code != 200:
            raise Exception(f"Failed to delete key-value pair: {response.text}")