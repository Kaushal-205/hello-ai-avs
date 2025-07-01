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
	fmt.Println("🔍 Quick AI Functionality Test")
	fmt.Println("==============================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	if cfg.TogetherAPIKey == "" {
		log.Fatalf("❌ TOGETHER_API_KEY is not set")
	}

	fmt.Printf("✅ Config loaded - Model: %s\n", cfg.DefaultModel)

	// Create AI client
	aiClient := ai.NewClient(cfg)

	// Test 1: Simple prompt
	fmt.Println("\n📝 Test 1: Simple AI prompt")
	testPrompt1 := &ai.TaskRequest{
		TaskType: "text_completion",
		Prompt:   "What is 2+2?",
	}

	result1, err := aiClient.ProcessTask(context.Background(), testPrompt1)
	if err != nil {
		log.Fatalf("❌ Test 1 failed: %v", err)
	}
	fmt.Printf("✅ Response: %s\n", result1.Result)
	fmt.Printf("   Tokens: %d, Time: %dms\n", result1.TokensUsed, result1.ProcessingTimeMs)

	// Test 2: AVS-related prompt
	fmt.Println("\n📝 Test 2: AVS explanation")
	testPrompt2 := &ai.TaskRequest{
		TaskType:    "text_completion",
		Prompt:      "In one sentence, what is an AVS in the context of EigenLayer?",
		MaxTokens:   50,
		Temperature: 0.5,
	}

	result2, err := aiClient.ProcessTask(context.Background(), testPrompt2)
	if err != nil {
		log.Fatalf("❌ Test 2 failed: %v", err)
	}
	fmt.Printf("✅ Response: %s\n", result2.Result)
	fmt.Printf("   Tokens: %d, Time: %dms\n", result2.TokensUsed, result2.ProcessingTimeMs)

	// Test 3: JSON serialization (what the AVS framework expects)
	fmt.Println("\n📝 Test 3: JSON serialization")
	jsonBytes, err := json.Marshal(result2)
	if err != nil {
		log.Fatalf("❌ JSON serialization failed: %v", err)
	}
	fmt.Printf("✅ JSON Output: %s\n", string(jsonBytes))

	// Test 4: Validation
	fmt.Println("\n📝 Test 4: Request validation")
	invalidRequest := &ai.TaskRequest{
		TaskType:    "text_completion",
		Prompt:      "", // Empty prompt should fail
		Temperature: 3.0, // Invalid temperature should fail
	}

	if err := aiClient.ValidateRequest(invalidRequest); err != nil {
		fmt.Printf("✅ Validation correctly rejected invalid request: %v\n", err)
	} else {
		fmt.Printf("❌ Validation should have failed but didn't\n")
	}

	fmt.Println("\n🎉 All tests completed successfully!")
	fmt.Println("💡 HelloAI-AVS is ready for full deployment testing!")
} 