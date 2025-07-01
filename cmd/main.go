package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Layr-Labs/hourglass-avs-template/internal/ai"
	"github.com/Layr-Labs/hourglass-avs-template/internal/config"
	"github.com/Layr-Labs/hourglass-monorepo/ponos/pkg/performer/server"
	performerV1 "github.com/Layr-Labs/protocol-apis/gen/protos/eigenlayer/hourglass/v1/performer"
	"go.uber.org/zap"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// This offchain binary is run by Operators running the Hourglass Executor. It contains
// the business logic of the AVS and performs worked based on the tasked sent to it.
// The Hourglass Aggregator ingests tasks from the TaskMailbox and distributes work
// to Executors configured to run the AVS Performer. Performers execute the work and
// return the result to the Executor where the result is signed and return to the
// Aggregator to place in the outbox once the signing threshold is met.

type TaskWorker struct {
	logger   *zap.Logger
	aiClient *ai.Client
	config   *config.Config
}

func NewTaskWorker(logger *zap.Logger, aiClient *ai.Client, config *config.Config) *TaskWorker {
	return &TaskWorker{
		logger:   logger,
		aiClient: aiClient,
		config:   config,
	}
}

func (tw *TaskWorker) ValidateTask(t *performerV1.TaskRequest) error {
	tw.logger.Sugar().Infow("Validating task",
		zap.String("task_id", string(t.TaskId)),
		zap.Int("payload_size", len(t.Payload)),
	)

	// Parse the task payload as AI request
	var aiRequest ai.TaskRequest
	if err := json.Unmarshal(t.Payload, &aiRequest); err != nil {
		tw.logger.Error("Failed to parse task payload", zap.Error(err))
		return fmt.Errorf("invalid task payload: %w", err)
	}

	// Validate AI request
	if err := tw.aiClient.ValidateRequest(&aiRequest); err != nil {
		tw.logger.Error("Invalid AI request", zap.Error(err))
		return fmt.Errorf("invalid AI request: %w", err)
	}

	tw.logger.Sugar().Infow("Task validation successful",
		zap.String("task_id", string(t.TaskId)),
		zap.String("task_type", aiRequest.TaskType),
		zap.String("prompt_preview", aiRequest.Prompt[:min(50, len(aiRequest.Prompt))]),
	)

	return nil
}

func (tw *TaskWorker) HandleTask(t *performerV1.TaskRequest) (*performerV1.TaskResponse, error) {
	tw.logger.Sugar().Infow("Handling task",
		zap.String("task_id", string(t.TaskId)),
	)

	// Parse the task payload as AI request
	var aiRequest ai.TaskRequest
	if err := json.Unmarshal(t.Payload, &aiRequest); err != nil {
		tw.logger.Error("Failed to parse task payload", zap.Error(err))
		return nil, fmt.Errorf("invalid task payload: %w", err)
	}

	// Process AI task
	ctx := context.Background()
	aiResponse, err := tw.aiClient.ProcessTask(ctx, &aiRequest)
	if err != nil {
		tw.logger.Error("Failed to process AI task", zap.Error(err))
		return nil, fmt.Errorf("AI processing failed: %w", err)
	}

	// Convert AI response to bytes
	resultBytes, err := json.Marshal(aiResponse)
	if err != nil {
		tw.logger.Error("Failed to marshal AI response", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	tw.logger.Sugar().Infow("Task completed successfully",
		zap.String("task_id", string(t.TaskId)),
		zap.String("model_used", aiResponse.ModelUsed),
		zap.Int("tokens_used", aiResponse.TokensUsed),
		zap.Int64("processing_time_ms", aiResponse.ProcessingTimeMs),
	)

	return &performerV1.TaskResponse{
		TaskId: t.TaskId,
		Result: resultBytes,
	}, nil
}

func main() {
	ctx := context.Background()
	l, _ := zap.NewProduction()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}

	// Validate API key
	if cfg.TogetherAPIKey == "" {
		panic("TOGETHER_API_KEY is required but not set")
	}

	// Create AI client
	aiClient := ai.NewClient(cfg)

	// Create task worker
	w := NewTaskWorker(l, aiClient, cfg)

	l.Info("Starting HelloAI-AVS Performer",
		zap.String("default_model", cfg.DefaultModel),
		zap.Int("max_tokens", cfg.MaxTokens),
		zap.Float64("temperature", cfg.Temperature),
	)

	pp, err := server.NewPonosPerformerWithRpcServer(&server.PonosPerformerConfig{
		Port:    8080,
		Timeout: 5 * time.Second,
	}, w, l)
	if err != nil {
		panic(fmt.Errorf("failed to create performer: %w", err))
	}

	if err := pp.Start(ctx); err != nil {
		panic(err)
	}
}
