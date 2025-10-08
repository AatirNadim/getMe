package gosdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"getMeClientForGo/core"
	"getMeClientForGo/core/constants"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type GetMeClient struct {
	httpClient *http.Client
}

type PutRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (client *GetMeClient) Init() error {

	err := godotenv.Load()
	if err != nil {
		return err
	}
	var socketPath string
	socketPath = os.Getenv("SOCKET_PATH")
	if socketPath == "" {
		socketPath = constants.SocketPath
	}

	httpClient, err := core.CreateHttpClient(socketPath)
	if err != nil {
		return err
	}
	client.httpClient = httpClient
	return nil
}

func (client *GetMeClient) Get(key string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/get", constants.BaseUrl), nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", key)
	req.URL.RawQuery = q.Encode()

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get key: %s, status code: %d", key, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (client *GetMeClient) Put(key, value string) error {

	fmt.Println("Preparing JSON payload for PUT request with key:", key, " and value:", value)
	jsonPayload, err := json.Marshal(PutRequestBody{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	fmt.Println("preparing io reader payload with:", jsonPayload)
	readerPayload := bytes.NewReader(jsonPayload)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/put", constants.BaseUrl), readerPayload)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("key", key)
	q.Add("value", value)
	req.URL.RawQuery = q.Encode()

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to put key: %s, status code: %d", key, resp.StatusCode)
	}

	return nil
}

func (client *GetMeClient) Delete(key string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete", constants.BaseUrl), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("key", key)
	req.URL.RawQuery = q.Encode()

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete key: %s, status code: %d", key, resp.StatusCode)
	}

	return nil
}

func (client *GetMeClient) ClearStore() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/clearStore", constants.BaseUrl), nil)
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to clear store, status code: %d", resp.StatusCode)
	}

	return nil
}
