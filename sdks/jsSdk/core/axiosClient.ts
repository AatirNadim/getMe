import axios from "axios";


const axiosClient = axios.create({
  baseURL: 'http://unix',
  socketPath: process.env.GETME_SOCKET_PATH || '/tmp/getMeStore/getMe.sock',
  headers: {
    'Content-Type': 'application/json',
  },
});

export { axiosClient }