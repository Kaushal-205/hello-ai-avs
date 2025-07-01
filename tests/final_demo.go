package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Layr-Labs/hourglass-avs-template/internal/ai"
	"github.com/Layr-Labs/hourglass-avs-template/internal/config"
)

func main() {
	fmt.Println("🚀 HelloAI-AVS Final Demo")
	fmt.Println("=========================")
	
	cfg, _ := config.Load()
	aiClient := ai.NewClient(cfg)
	
	// Demo with a clear, well-defined prompt
	request := &ai.TaskRequest{
		TaskType: "text_completion",
		Prompt:   "Write a short explanation of what makes a good smart contract. Keep it under 100 words.",
		MaxTokens: 120,
		Temperature: 0.7,
	}
	
	fmt.Printf("📝 Prompt: %s\n\n", request.Prompt)
	fmt.Println("🤖 Processing...")
	
	result, err := aiClient.ProcessTask(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("✅ AI Response:\n%s\n\n", result.Result)
	fmt.Printf("📊 Performance:\n")
	fmt.Printf("   Model: %s\n", result.ModelUsed)
	fmt.Printf("   Tokens Used: %d\n", result.TokensUsed)
	fmt.Printf("   Processing Time: %dms\n", result.ProcessingTimeMs)
	
	fmt.Println("\n🎉 HelloAI-AVS is working perfectly!")
} 

