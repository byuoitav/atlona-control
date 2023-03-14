package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func PostHTTPWithContext(ctx context.Context, address string, service string, payload []byte) ([]byte, error) {

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", address, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	switch {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, errors.New(string(body))
	case body == nil:
		return nil, errors.New("response from device was blank")
	}

	return body, nil
}

func PostHTTP(address string, payload []byte, service string) ([]byte, error) {

	return PostHTTPWithContext(context.TODO(), address, service, payload)
}
