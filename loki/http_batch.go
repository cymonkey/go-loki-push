package loki

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpBatchRequest struct {
	*batch
}

var (
	_ BatchRequest = (*HttpBatchRequest)(nil)
)

func (r *HttpBatchRequest) Compose(context context.Context, url string) (req *http.Request, entriesCount int, err error) {
	buf := bytes.NewBuffer([]byte{})
	gz := gzip.NewWriter(buf)
	pushRequest, entriesCount := r.createPushRequest()

	if err := json.NewEncoder(gz).Encode(pushRequest); err != nil {
		return nil, -1, err
	}

	if err := gz.Close(); err != nil {
		return nil, -1, err
	}

	req, err = http.NewRequestWithContext(context, http.MethodPost, url, buf)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	return req, entriesCount, nil
}

func (r *HttpBatchRequest) setBatch(ba *batch) {
	r.batch = ba
}
