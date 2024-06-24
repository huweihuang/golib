package httplib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/huweihuang/golib/logger/zap"
)

func RequestURL(method, url, path string, header map[string]string, request interface{}, response interface{}) (
	statusCode int, body []byte, err error) {

	log.Logger().With(
		"method", method,
		"url", url,
		"path", path,
		"header", header,
		"request", request,
	).Debug("request url info")

	params, err := encodeData(request)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to encode request, err: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", url, path), params)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to new request, %v", err)
	}

	// set header
	req.Header.Set("Content-Type", "application/json")
	for key, value := range header {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to send http request, err: %v", err)
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read body, err: %v", err)
	}

	if response != nil {
		err = json.Unmarshal(body, response)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to unmarshal body, err: %v", err)
		}
	}

	return resp.StatusCode, body, nil
}

func encodeData(data interface{}) (*bytes.Buffer, error) {
	params := bytes.NewBuffer(nil)
	if data != nil {
		buf, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		if _, err := params.Write(buf); err != nil {
			return nil, err
		}
	}
	return params, nil
}
