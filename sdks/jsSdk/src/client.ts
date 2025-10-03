import axios, { AxiosInstance } from "axios";

class GetMeClient {
  private axiosClient: AxiosInstance;

  constructor() {
    this.axiosClient = axios.create({
      baseURL: "http://unix",
      socketPath : process.env.GETME_SOCKET_PATH || '/tmp/getMeStore/getMe.sock',
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

  put = async (key: string, value: string): Promise<void> => {
    try {
      await this.axiosClient.post(`/put`, { Key: key, Value: value });
    } catch (error) {
      console.error("Error storing value:", error);
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
}

export { GetMeClient };
