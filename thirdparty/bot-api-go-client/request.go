package bot

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 5 * time.Second}

func Request(ctx context.Context, method, path string, body []byte, accessToken string) ([]byte, error) {
	req, err := http.NewRequest(method, "https://mixin-api.zeromesh.net"+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, ServerError(ctx, nil)
	}
	return ioutil.ReadAll(resp.Body)
}
