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
	promptFormat    = "\n\nHuman:%s\n\nAssistant:"
	claudeV2ModelID = "anthropic.claude-v2" //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
)

const prompt = `<paragraph> 
"In 1758, the Swedish botanist and zoologist Carl Linnaeus published in his Systema Naturae, the two-word naming of species (binomial nomenclature). Canis is the Latin word meaning "dog", and under this genus, he listed the domestic dog, the wolf, and the golden jackal."
</paragraph>

Please rewrite the above paragraph to make it understandable to a 5th grader.

Please output your rewrite in <rewrite></rewrite> tags.`

type Request struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	Temperature       float64  `json:"temperature,omitempty"`
	TopP              float64  `json:"top_p,omitempty"`
	TopK              float64  `json:"top_k,omitempty"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
}

type Response struct {
	Completion string `json:"completion"`
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

	payload := Request{
		Prompt:            fmt.Sprintf(promptFormat, prompt),
		MaxTokensToSample: 2048,
		Temperature:       0.7,
		TopK:              250,
		TopP:              1,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	output, err := bedrockClient.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String(claudeV2ModelID),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		log.Fatalf("Failed to invoke model: %v\n", err)
	}

	var resp Response
	err = json.Unmarshal(output.Body, &resp)
	if err != nil {
		log.Fatalf("Failed to unmarshal response: %v\n", err)
	}

	fmt.Println("Response from LLM: \n", resp.Completion)
}
