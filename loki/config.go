package loki

import (
	"net/url"
	"time"
)

const (
	GrpcRequestProto LokiRequestProto = "grpc"
	HttpRequestProto LokiRequestProto = "http"
)

type LokiRequestProto string

const (
	MAX_ERR_MESSAGE_LEN     = 1024
	DEFAULT_BATCH_MAX_SIZE  = 10
	DEFAULT_BATCH_MAX_WAIT  = 10 * time.Second
	DEFAULT_REQUEST_TIMEOUT = 30 * time.Second
	DEFAULT_REQUEST_PROTO   = GrpcRequestProto
)

type NoLogger struct{}

func (l *NoLogger) Error(e string) {}

type LogHandler func()

type Requesthandler func()

type Config struct {
	Url             string           // Push Url of the loki server including http:// or https://, ex: https://example.com/loki/api/v1/push
	BatchMaxSize    int              // BatchMaxSize is the maximum number of log lines that are sent in one request. Default: 5
	BatchMaxWait    time.Duration    // BatchMaxWait is the maximum time to wait before sending a request. Default: 10s
	Timeout         time.Duration    // Request timeout. Default: 30s
	RequestProtocol LokiRequestProto // Http or grpc. Default: grpc

	Logger Logger

	// Static labels that are added to all log lines.
	Labels map[string]string

	// Dynamic labels that are extracted from log fields, and add to log line as Structure Metadata.
	// Useful when using with logger, ex: log.Info("message", "log_field1", "val1", "log_field2", "val2")
	// The Key is required, but the Value is fallback when extracted field is empty and it is optional.
	StructuredMetadata map[string]string

	// Basic Auth
	Username string
	Password string
}

func (c *Config) Apply(conf Config) (*Config, error) {
	return nil, nil
}

func defaultConfig() *Config {
	return &Config{
		Url:                "",
		BatchMaxSize:       DEFAULT_BATCH_MAX_SIZE,
		BatchMaxWait:       DEFAULT_BATCH_MAX_WAIT,
		Timeout:            DEFAULT_REQUEST_TIMEOUT,
		RequestProtocol:    DEFAULT_REQUEST_PROTO,
		Logger:             &NoLogger{},
		Labels:             nil,
		StructuredMetadata: nil,
		Username:           "",
		Password:           "",
	}
}

func NewWithDefaultConfig(cfg Config) (*Config, error) {
	_, err := url.ParseRequestURI(cfg.Url)
	if err != nil {
		return nil, err
	}
	config := defaultConfig()
	config.Url = cfg.Url

	if cfg.BatchMaxSize > 0 {
		config.BatchMaxSize = cfg.BatchMaxSize
	}

	if cfg.BatchMaxWait > 0 {
		config.BatchMaxWait = cfg.BatchMaxWait
	}

	if cfg.Timeout > 0 {
		config.Timeout = cfg.Timeout
	}

	if cfg.Logger != nil {
		config.Logger = cfg.Logger
	}

	if cfg.RequestProtocol == HttpRequestProto {
		config.RequestProtocol = HttpRequestProto
	}

	return config, nil
}

func (c *Config) Validate() bool {
	_, err := url.ParseRequestURI(c.Url)
	if err != nil {
		return false
	}

	if c.BatchMaxSize <= 0 || c.BatchMaxWait <= 0 || c.Timeout == 0 || c.Logger == nil {
		return false
	}

	if c.RequestProtocol != HttpRequestProto && c.RequestProtocol != GrpcRequestProto {
		return false
	}

	return true
}
