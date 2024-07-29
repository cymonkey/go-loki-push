package loki

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	cfg, err := NewWithDefaultConfig(Config{Url: "http://loki.com"})
	assert.Nil(t, err)
	assert.Equal(t, "http://loki.com", cfg.Url)
	assert.Equal(t, DEFAULT_BATCH_MAX_SIZE, cfg.BatchMaxSize)
	assert.Equal(t, DEFAULT_BATCH_MAX_WAIT, cfg.BatchMaxWait)
	assert.Equal(t, DEFAULT_REQUEST_TIMEOUT, cfg.Timeout)
	assert.Equal(t, DEFAULT_REQUEST_PROTO, cfg.RequestProtocol)
	assert.NotNil(t, cfg.Logger)
	assert.Zero(t, cfg.Labels)
	assert.Zero(t, cfg.StructuredMetadata)
	assert.Equal(t, "", cfg.Username)
	assert.Equal(t, "", cfg.Password)
}

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		Url                string
		BatchMaxSize       int
		BatchMaxWait       time.Duration
		Timeout            time.Duration
		RequestProtocol    LokiRequestProto
		Logger             Logger
		Labels             map[string]string
		StructuredMetadata map[string]string
		Username           string
		Password           string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Url:                tt.fields.Url,
				BatchMaxSize:       tt.fields.BatchMaxSize,
				BatchMaxWait:       tt.fields.BatchMaxWait,
				Timeout:            tt.fields.Timeout,
				RequestProtocol:    tt.fields.RequestProtocol,
				Logger:             tt.fields.Logger,
				Labels:             tt.fields.Labels,
				StructuredMetadata: tt.fields.StructuredMetadata,
				Username:           tt.fields.Username,
				Password:           tt.fields.Password,
			}
			if got := c.Validate(); got != tt.want {
				t.Errorf("Config.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWithDefaultConfig(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "New default config with only url",
			args: args{Config{
				Url: "http://localhost:3100",
			}},
			want: &Config{
				Url:                "http://localhost:3100",
				BatchMaxSize:       10,
				BatchMaxWait:       10 * time.Second,
				Timeout:            30 * time.Second,
				RequestProtocol:    GrpcRequestProto,
				Labels:             *new(map[string]string),
				StructuredMetadata: *new(map[string]string),
				Logger:             &NoLogger{},
				Username:           "",
				Password:           "",
			},
			wantErr: false,
		},
		{
			name: "New config with custom one",
			args: args{Config{
				Url:             "http://testing.com",
				BatchMaxSize:    3,
				BatchMaxWait:    5 * time.Second,
				RequestProtocol: HttpRequestProto,
			}},
			want: &Config{
				Url:                "http://testing.com",
				BatchMaxSize:       3,
				BatchMaxWait:       5 * time.Second,
				RequestProtocol:    HttpRequestProto,
				Timeout:            30 * time.Second,
				Labels:             *new(map[string]string),
				StructuredMetadata: *new(map[string]string),
				Username:           "",
				Password:           "",
				Logger:             &NoLogger{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWithDefaultConfig(tt.args.cfg)
			assert.NoError(t, err, "NewWithDefaultConfig() error = %v", err)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithDefaultConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualExportedValues(t, tt.want, got, "NewWithDefaultConfig() = %v, want %v", got, tt.want)
		})
	}

	errTests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "New default config with invalid url",
			args: args{Config{
				Url: "1203aaa",
			}},
			want: &Config{
				Url:                "http://localhost:3100",
				BatchMaxSize:       10,
				BatchMaxWait:       10 * time.Second,
				Timeout:            30 * time.Second,
				RequestProtocol:    GrpcRequestProto,
				Labels:             *new(map[string]string),
				StructuredMetadata: *new(map[string]string),
				Logger:             &NoLogger{},
				Username:           "",
				Password:           "",
			},
			wantErr: false,
		},
	}
	for _, tt := range errTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWithDefaultConfig(tt.args.cfg)
			assert.Error(t, err)
		})
	}
}
