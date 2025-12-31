import axios, { AxiosInstance } from "axios";
import { baseUrl, defaultSocketPath } from "./constants";

class GetMeClient {
  private axiosClient: AxiosInstance;

  constructor() {
    this.axiosClient = axios.create({
      baseURL: baseUrl,
      socketPath : process.env.GETME_SOCKET_PATH || defaultSocketPath,
      headers: {
        "Content-Type": "application/json",
      },
    });
    console.log("GetMeClient initialized");
  }

  get = async (key: string): Promise<string> => {
    try {
      const response = await this.axiosClient.get(`/get`, {
        params: { key },
      });

      console.log("Fetched value:", response.data);

      return response.data;
    } catch (error) {
      console.error("Error fetching value:", error);
      throw error;
    }
  };

  getJson = async <T = any>(key: string): Promise<T> => {
    const raw = await this.get(key);
    try {
      return JSON.parse(raw) as T;
    } catch (error) {
      console.error("Error parsing JSON value:", error);
      throw error;
    }
  };

  put = async (key: string, value: string): Promise<void> => {
    try {
      await this.axiosClient.post(`/put`, { Key: key, Value: value });
    } catch (error) {
      console.error("Error storing value:", error);
      throw error;
    }
  };

  putJson = async (key: string, value: unknown): Promise<void> => {
    try {
      const jsonString = JSON.stringify(value);
      await this.put(key, jsonString);
    } catch (error) {
      console.error("Error serializing JSON value:", error);
      throw error;
    }
  };

  delete = async (key: string): Promise<void> => {
    try {
      await this.axiosClient.delete(`/delete`, {
        params: { key },
      });
    } catch (error) {
      console.error("Error deleting value:", error);
      throw error;
    }
  };

  clearStore = async () => {
    try {
      await this.axiosClient.delete(`/clearStore`);
    } catch (error) {
      console.error("Error clearing store:", error);
      throw error;
    }
  }
}

export { GetMeClient };
