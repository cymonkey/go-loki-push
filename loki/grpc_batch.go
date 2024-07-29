package loki

import (
	"bytes"
	"context"
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
)

var (
	_ BatchRequest = (*GrpcBatchRequest)(nil)
)

type GrpcBatchRequest struct {
	*batch
}

func (r *GrpcBatchRequest) Compose(context context.Context, url string) (req *http.Request, entriesCount int, err error) {
	pushRequest, entriesCount := r.createPushRequest()
	buf, err := proto.Marshal(pushRequest)
	if err != nil {
		return nil, 0, err
	}
	buf = snappy.Encode(nil, buf)
	req, err = http.NewRequestWithContext(context, "POST", url, bytes.NewReader(buf))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/x-protobuf")

	return req, entriesCount, nil
}

func (r *GrpcBatchRequest) setBatch(ba *batch) {
	r.batch = ba
}
