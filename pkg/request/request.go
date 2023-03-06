package request

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

func HttpRequest(method, path string, reqBody interface{}, respBody interface{}, header map[string]string) error {

	var (
		req     *http.Request
		res     *http.Response
		err     error
		content []byte
	)

	// 表单不为空
	d, _ := json.Marshal(reqBody)

	req, _ = http.NewRequest(method, path, bytes.NewReader(d))
	req.Header.Set("Content-Type", "application/json")
	if header != nil {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}
	client := &http.Client{Timeout: 10 * time.Second}

	if res, err = client.Do(req); err != nil {
		return err
	}

	defer func(b io.ReadCloser) { _ = b.Close() }(res.Body)

	if res.StatusCode != 200 {
		return errors.New("Request error: " + res.Status)
	}

	if content, err = io.ReadAll(res.Body); err != nil {
		return err
	}
	if err = json.Unmarshal(content, respBody); err != nil {
		return err
	}
	return nil
}
