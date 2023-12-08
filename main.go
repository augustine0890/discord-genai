package main

import (
	"context"
	"discord-genai/util"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	systemInstruction = "Generate a response that can be displayed in Discord"
	modelID           = "amazon.titan-text-lite-v1" //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
	// modelID = "amazon.titan-text-express-v1"
)

const prompt = "Who invented the airplane?"

// TextGenerationConfig holds the configuration for text generation.
type TextGenerationConfig struct {
	MaxTokenCount int      `json:"maxTokenCount"`
	StopSequences []string `json:"stopSequences"`
	Temperature   float64  `json:"temperature"`
	TopP          float64  `json:"topP"`
}

// Input represents the payload for model invocation.
type Input struct {
	InputText            string               `json:"inputText"`
	TextGenerationConfig TextGenerationConfig `json:"textGenerationConfig"`
}

type Result struct {
	TokenCount       int    `json:"tokenCount"`
	OutputText       string `json:"outputText"`
	CompletionReason string `json:"completionReason"`
}

type Output struct {
	InputTextTokenCount int      `json:"inputTextTokenCount"`
	Results             []Result `json:"results"`
}

func main() {
	stage := flag.String("stage", "prod", "The enviroment running")
	flag.Parse()

	// Loading enivroment variables
	err := util.LoadEnv(*stage)
	if err != nil {
		fmt.Printf("Error loading environment variables: %v\n", err)
	}
	log.Printf("Running with %v environment\n", *stage)

	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Error gettingthe aws config: %v\n", err)
	}

	bedrockClient := bedrockruntime.NewFromConfig(cfg)
	textGenerationConfig := TextGenerationConfig{
		Temperature:   0.8,
		TopP:          1.0,
		MaxTokenCount: 256,
		StopSequences: []string{},
	}

	// Append the system instruction to the prompt
	fullPrompt := prompt + " {{" + systemInstruction + " }}"

	payload := Input{
		InputText:            fullPrompt,
		TextGenerationConfig: textGenerationConfig,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	output, err := bedrockClient.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		log.Fatalf("Failed to invoke model: %v\n", err)
	}

	var resp Output
	err = json.Unmarshal(output.Body, &resp)
	if err != nil {
		log.Fatalf("Failed to unmarshal response: %v\n", err)
	}

	fmt.Println("Response from LLM: \n", resp)
}
