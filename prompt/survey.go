package prompt

//go:generate mockgen -destination ../internal/mock/prompt/survey.go -package mock -source survey.go

import (
	"github.com/AlecAivazis/survey/v2"
	"log"
)

type Prompt interface {
	Confirm(string) bool
}

func NewPrompt() Prompt {
	return prompt{}
}

type prompt struct{}

func NewPromptContext(interactive bool, prompt Prompt) PromptContext{
	return PromptContext{
		Interactive: interactive,
		Prompt: prompt,
	}
}

type PromptContext struct {
	Interactive bool
	Prompt Prompt
}

func (p prompt) Confirm(message string) bool {
	if message == "" {
		message = "Confirm?"
	}

	confirmPrompt := &survey.Confirm{
		Default: false,
		Message: message,
	}

	confirmation := false
	err := survey.AskOne(confirmPrompt, &confirmation)
	if err != nil {
		log.Printf("error occurred when getting input from user: %s", err)
		return false
	}

	return confirmation
}
