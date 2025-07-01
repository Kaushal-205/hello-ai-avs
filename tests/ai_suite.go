package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Layr-Labs/hourglass-avs-template/internal/ai"
	"github.com/Layr-Labs/hourglass-avs-template/internal/config"
)

func main() {
	fmt.Println("🔍 HelloAI-AVS Comprehensive Test Suite")
	fmt.Println("=======================================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	if cfg.TogetherAPIKey == "" {
		log.Fatalf("❌ TOGETHER_API_KEY is not set")
	}

	fmt.Printf("✅ Config loaded - Model: %s\n", cfg.DefaultModel)
	aiClient := ai.NewClient(cfg)

	// Test Suite
	tests := []struct {
		name    string
		request *ai.TaskRequest
		expectError bool
	}{
		{
			name: "Simple Math Question",
			request: &ai.TaskRequest{
				TaskType: "text_completion",
				Prompt:   "What is 15 + 27?",
			},
			expectError: false,
		},
		{
			name: "Blockchain Explanation",
			request: &ai.TaskRequest{
				TaskType:    "text_completion",
				Prompt:      "Explain blockchain technology in 2 sentences.",
				MaxTokens:   80,
				Temperature: 0.5,
			},
			expectError: false,
		},
		{
			name: "Smart Contract Best Practices",
			request: &ai.TaskRequest{
				TaskType:    "text_completion",
				Prompt:      "List 3 key principles for writing secure smart contracts.",
				MaxTokens:   150,
				Temperature: 0.7,
			},
			expectError: false,
		},
		{
			name: "Invalid Request (Empty Prompt)",
			request: &ai.TaskRequest{
				TaskType: "text_completion",
				Prompt:   "",
			},
			expectError: true,
		},
	}

	var passedTests, failedTests int

	for i, test := range tests {
		fmt.Printf("\n📝 Test %d: %s\n", i+1, test.name)

		// Validate request first
		if err := aiClient.ValidateRequest(test.request); err != nil {
			if test.expectError {
				fmt.Printf("✅ Validation correctly rejected: %v\n", err)
				passedTests++
				continue
			} else {
				fmt.Printf("❌ Unexpected validation error: %v\n", err)
				failedTests++
				continue
			}
		}

		if test.expectError {
			fmt.Printf("❌ Expected validation error but got none\n")
			failedTests++
			continue
		}

		// Process the request
		result, err := aiClient.ProcessTask(context.Background(), test.request)
		if err != nil {
			fmt.Printf("❌ Processing failed: %v\n", err)
			failedTests++
			continue
		}

		fmt.Printf("✅ Response: %s\n", result.Result)
		fmt.Printf("   Stats: %d tokens, %dms, Model: %s\n", 
			result.TokensUsed, result.ProcessingTimeMs, result.ModelUsed)

		// Test JSON serialization
		if _, err := json.Marshal(result); err != nil {
			fmt.Printf("❌ JSON serialization failed: %v\n", err)
			failedTests++
		} else {
			fmt.Printf("   ✅ JSON serialization: OK\n")
			passedTests++
		}
	}

	// Summary
	fmt.Printf("\n🎯 Test Results Summary\n")
	fmt.Printf("=======================\n")
	fmt.Printf("✅ Passed: %d\n", passedTests)
	fmt.Printf("❌ Failed: %d\n", failedTests)
	fmt.Printf("📊 Total: %d\n", passedTests+failedTests)

	if failedTests == 0 {
		fmt.Println("\n🎉 All tests passed! HelloAI-AVS is ready for deployment!")
	} else {
		fmt.Printf("\n⚠️  %d tests failed. Please review the issues above.\n", failedTests)
	}
} 