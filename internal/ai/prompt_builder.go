package ai

func BuildJSONPrompt(task string, context string) GenerateInput {
	return GenerateInput{
		SystemPrompt: "You are an audit log investigator assistant. Use only the provided logs. Do not invent facts. Return valid JSON only.",
		UserPrompt:   task + "\n\n" + context,
		Temperature:  0.2,
		MaxTokens:    700,
	}
}
