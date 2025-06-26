package src

import (
	"context"
	"fmt"
	"strings"

	"github.com/Thearas/dodo/src/prompt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/sirupsen/logrus"
)

const (
	LLMOutputPrefix = "```yaml\n"
)

// Use Deepseek by default, but you can use OpenAI by setting the apiKey and baseURL
func LLMGendataConfig(
	ctx context.Context,
	apiKey, baseURL, model, prompt_ string,
	tables, columnStats, sqls []string,
) (string, error) {
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/beta"
	}
	if model == "" {
		model = "deepseek-coder"
	}

	userPrompt := fmt.Sprintf(`
<tables>
%s
</tables>


<column-stats>
%s
</column-stats>


<queries>
%s
</queries>
			`,
		strings.Join(tables, "\n"),
		strings.Join(columnStats, "\n---\n"),
		strings.Join(sqls, "\n"))
	if prompt_ != "" {
		userPrompt = fmt.Sprintf(`%s

<additional-user-prompt>
%s
</additional-user-prompt>
		`, userPrompt, prompt_)
	}

	logrus.Debugln("LLM user prompt:", userPrompt)

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)
	stop := openai.ChatCompletionNewParamsStopUnion{
		OfString: openai.String("\n```"),
	}
	stop.SetExtraFields(map[string]any{"prefix": true})
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       model,
		Temperature: openai.Float(0.3),
		Stop:        stop,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt.Gendata),
			openai.UserMessage(userPrompt),
			openai.AssistantMessage(LLMOutputPrefix),
		},
	})
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(chatCompletion.Choices[0].Message.Content, LLMOutputPrefix), nil
}
