package river

import (
	"context"
	"testing"

	"log/slog"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLogger(t *testing.T) {
	tests := []struct {
		name   string
		config Logger
	}{
		{
			name: "Default Logger",
			config: Logger{
				Debug:  false,
				Pretty: false,
			},
		},
		{
			name: "Debug Logger",
			config: Logger{
				Debug:  true,
				Pretty: false,
			},
		},
		{
			name: "Pretty Logger",
			config: Logger{
				Debug:  false,
				Pretty: true,
			},
		},
		{
			name: "Debug Pretty Logger",
			config: Logger{
				Debug:  true,
				Pretty: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			logger := createLogger(tt.config)

			require.NotNil(t, logger, "Logger should not be nil")

			handler := logger.Handler()
			require.NotNil(t, handler, "Logger handler should not be nil")

			assert.Equal(t, handler.Enabled(ctx, slog.LevelDebug), tt.config.Debug)

			if tt.config.Pretty {
				assert.IsType(t, &slog.TextHandler{}, handler, "Logger handler should be TextHandler")
			} else {
				assert.IsType(t, &slog.JSONHandler{}, handler, "Logger handler should be JSONHandler")
			}
		})
	}
}

func TestCreateQueueConfig(t *testing.T) {
	tests := []struct {
		name   string
		queues []Queue
		expect map[string]river.QueueConfig
	}{
		{
			name:   "Default Queue",
			queues: []Queue{},
			expect: map[string]river.QueueConfig{
				river.QueueDefault: {MaxWorkers: defaultMaxWorkers},
			},
		},
		{
			name: "Custom Queues",
			queues: []Queue{
				{Name: "queue1", MaxWorkers: 10},
				{Name: "queue2", MaxWorkers: 20},
			},
			expect: map[string]river.QueueConfig{
				river.QueueDefault: {MaxWorkers: defaultMaxWorkers},
				"queue1":           {MaxWorkers: 10},
				"queue2":           {MaxWorkers: 20},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qc := createQueueConfig(tt.queues)
			assert.Equal(t, tt.expect, qc)
		})
	}
}
