package raftd

import (
	"io"
	"net/http"
)

func httpPost(url, contentType string, body io.Reader, header http.Header) (resp *http.Response, err error) {
	c := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	for k, v := range header {
		for _, v1 := range v {
			req.Header.Add(k, v1)
		}
	}
	return c.Do(req)
}
