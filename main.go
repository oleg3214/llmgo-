package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

type LLMSearchTool struct {
	llm llms.Model
}

func (t LLMSearchTool) Name() string {
	return "LLMSearch"
}

func (t LLMSearchTool) Description() string {
	return "Поиск информации "
}

func (t LLMSearchTool) Call(ctx context.Context, input string) (string, error) {

	prompt := fmt.Sprintf("Ты интеллектуальный поисковик.Дай краткий ответ без лишний воды: %s", input)

	resp, err := llms.GenerateFromSinglePrompt(ctx, t.llm, prompt)
	if err != nil {
		return "", err
	}

	return resp, nil
}

type UpLoadBD struct {
	llm llms.Model
}

func (t UpLoadBD) Name() string {
	return "LLMSearch"
}

func (t UpLoadBD) Description() string {
	return "Поиск информации "
}

func (t UpLoadBD) Call(ctx context.Context, input string) (string, error) {

	prompt := fmt.Sprintf("Ты интеллектуальный поисковик.Дай краткий ответ без лишний воды: %s", input)

	resp, err := llms.GenerateFromSinglePrompt(ctx, t.llm, prompt)
	if err != nil {
		return "", err
	}

	return resp, nil
}

type VerboseTool struct {
	Tool tools.Tool
}

func (v VerboseTool) Name() string {
	return v.Tool.Name()
}

func (v VerboseTool) Description() string {
	return v.Tool.Description()
}

func (v VerboseTool) Call(ctx context.Context, input string) (string, error) {
	fmt.Printf("\n🛠 Action: %s\n📥 Input: %s\n", v.Name(), input)
	output, err := v.Tool.Call(ctx, input)
	if err != nil {
		fmt.Printf("⛔️ Error: %v\n", err)
		return "", err
	}

	fmt.Printf("📤 Observation: %s\n", output)
	return output, nil
}

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Set OPENROUTER_API_KEY env variable with your OpenRouter API key.")
	}
	llm, err := openai.New(
		openai.WithBaseURL("https://openrouter.ai/api/v1"),
		openai.WithModel("deepseek/deepseek-r1-0528:free"),
		openai.WithToken(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	// searchTool := LLMSearchTool{llm: llm}

	//

	// agentTools := []tools.Tool{
	// 	searchTool,
	// 	tools.Calculator{},
	// }
	wrappedTools := []tools.Tool{
		VerboseTool{Tool: LLMSearchTool{llm: llm}},
		VerboseTool{Tool: tools.Calculator{}},
	}

	agent := agents.NewOneShotAgent(llm, wrappedTools, agents.WithMaxIterations(10))
	executor := agents.NewExecutor(agent)

	question := "Узнай возраст Тома Холанда у LLMSearchTool, обработай ответ и вытяни возраст и умнож его на 2, используя Calculator"
	answer, err := chains.Run(context.Background(), executor, question)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(answer)
}
