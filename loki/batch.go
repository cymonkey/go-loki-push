package loki

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/grafana/loki/pkg/push"
)

const (
	ERROR_STREAM_LIMIT_EXCEED = "streams: %d exceeds limit: %d, stream: '%s'"
)

type BatchRequest interface {
	Add(push.Entry) error
	Compose(context context.Context, url string) (*http.Request, int, error)
	Size() int
	Reset()
	setBatch(*batch)
}

type batch struct {
	streams map[string]*push.Stream
	bytes   int

	labels     map[string]string
	maxStreams int
}

func NewBatch[T BatchRequest](maxStreamConfig int, labels map[string]string, entries ...push.Entry) T {
	t := new(T)
	var b = &batch{
		streams:    make(map[string]*push.Stream),
		bytes:      0,
		labels:     labels,
		maxStreams: maxStreamConfig,
	}
	(*t).setBatch(b)

	for _, entry := range entries {
		_ = (*t).Add(entry)
	}

	return *t
}

func (b *batch) Size() int {
	return len(b.streams)
}

// add an entry to the batch
func (b *batch) Add(entry push.Entry) error {
	b.bytes += entrySize(entry)

	// Append the entry to an already existing stream (if any)
	labels := labelsToString(b.labels)
	if stream, ok := b.streams[labels]; ok {
		stream.Entries = append(stream.Entries, entry)
		return nil
	}

	streams := len(b.streams)
	if b.maxStreams > 0 && streams >= b.maxStreams {
		return fmt.Errorf(ERROR_STREAM_LIMIT_EXCEED, streams, b.maxStreams, labels)
	}

	// Add the entry as a new stream
	b.streams[labels] = &push.Stream{
		Labels:  labels,
		Entries: []push.Entry{entry},
	}
	return nil
}

func entrySize(entry push.Entry) int {
	structuredMetadataSize := 0
	for _, label := range entry.StructuredMetadata {
		structuredMetadataSize += label.Size()
	}
	return len(entry.Line) + structuredMetadataSize
}

func labelsToString(labels map[string]string) string {
	var b strings.Builder
	totalSize := 2
	lnameArr := make([]string, 0, len(labels))

	for lkey, lval := range labels {
		lnameArr = append(lnameArr, lkey)
		// guess size increase: 2 for `, ` between labels and 3 for the `=` and quotes around label value
		totalSize += len(lkey) + 2 + len(lval) + 3
	}

	b.Grow(totalSize)
	b.WriteByte('{')
	slices.Sort(lnameArr)
	for i, l := range lnameArr {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString(l)
		b.WriteString(`=`)
		b.WriteString(strconv.Quote(labels[l]))
	}
	b.WriteByte('}')

	return b.String()
}

func (b *batch) createPushRequest() (*push.PushRequest, int) {
	req := push.PushRequest{
		Streams: make([]push.Stream, 0, len(b.streams)),
	}

	entriesCount := 0
	for _, stream := range b.streams {
		req.Streams = append(req.Streams, *stream)
		entriesCount += len(stream.Entries)
	}
	return &req, entriesCount
}

func (b *batch) Reset() {
	b.streams = make(map[string]*push.Stream)
	b.bytes = 0
}
