package loki

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/grafana/loki/pkg/push"
)

type Logger interface {
	Error(string)
}

// Client pushes entries to Loki and can be stopped
type Client interface {
	WithLabels(map[string]string) Client
	Name() string
	Chan() chan<- push.Entry
	Stop()
}

type client struct {
	Client
	name       string
	config     *Config
	httpClient *http.Client
	ctx        context.Context
	cancel     context.CancelFunc
	once       *sync.Once
	entries    chan push.Entry
	labels     map[string]string
	logsBatch  BatchRequest
	waitGroup  *sync.WaitGroup
	logger     Logger
	auth       Authentication
	lock       *sync.Cond
}

func NewClient(cfg Config, ctx context.Context) (Client, error) {
	config, err := NewWithDefaultConfig(cfg)
	if err != nil {
		return nil, err
	}

	if !config.Validate() {
		panic("Invalid config")
	}

	var batch BatchRequest
	if config.RequestProtocol == HttpRequestProto {
		batch = NewBatch[*HttpBatchRequest](config.BatchMaxSize, config.Labels)
	} else {
		batch = NewBatch[*GrpcBatchRequest](config.BatchMaxSize, config.Labels)
	}

	var auth Authentication
	if config.Username != "" && config.Password != "" {
		auth = NewAuth(config.Username, config.Password)
	} else {
		auth = &NoAuth{}
	}

	var labels map[string]string
	if len(config.Labels) > 0 {
		labels = config.Labels
	} else {
		labels = make(map[string]string)
	}

	ctx, cancel := context.WithCancel(ctx)
	client := &client{
		config:     &cfg,
		httpClient: &http.Client{},
		ctx:        ctx,
		cancel:     cancel,
		entries:    make(chan push.Entry),
		labels:     labels,
		logsBatch:  batch,
		logger:     config.Logger,
		auth:       auth,
	}
	client.lock = sync.NewCond(&sync.Mutex{})

	client.waitGroup.Add(1)
	go client.run()
	return client, nil
}

func (c *client) SetLabels(map[string]string) {

}

func (c *client) Name() string {
	return c.name
}

func (c *client) Chan() chan<- push.Entry {
	return c.entries
}

// Stop stops the loki pusher/client
func (c *client) Stop() {
	c.waitGroup.Wait()
	c.once.Do(func() { close(c.entries) })
	c.cancel()
}

func (c *client) WithLabels(labels map[string]string) Client {
	clone := c.clone()
	for k, v := range labels {
		clone.labels[k] = v
	}
	return clone
}

func (c *client) clone() *client {
	log := *c
	return &log
}

func (c *client) run() {
	ticker := time.NewTimer(c.config.BatchMaxWait)
	defer ticker.Stop()

	defer func() {
		if c.logsBatch.Size() > 0 {
			// c.sendBatch()
			c.Flush()
		}

		c.waitGroup.Done()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		case entry, ok := <-c.entries:
			if !ok {
				continue
			}

			c.logsBatch.Add(entry)
			if c.logsBatch.Size() >= c.config.BatchMaxSize {
				c.sendBatch()
			}
		case <-ticker.C:
			if c.logsBatch.Size() > 0 {
				c.sendBatch()
			}
		}
	}
}

func (c *client) sendBatch() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()
	req, _, err := c.logsBatch.Compose(ctx, c.config.Url)
	if err != nil {
		return err
	}

	c.auth.Apply(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(resp.Body, MAX_ERR_MESSAGE_LEN))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = fmt.Errorf("server returned HTTP status %s (%d): %s", resp.Status, resp.StatusCode, line)
		c.logger.Error(err.Error())
		// Retry if needed
	}
	c.logsBatch.Reset()
	return err
}

func (c *client) Flush() error {
	for len(c.entries) > 0 {
		if entry, ok := <-c.entries; ok {
			c.logsBatch.Add(entry)
		}
	}
	if c.logsBatch.Size() > 0 {
		c.sendBatch()
	}
	return nil
}
