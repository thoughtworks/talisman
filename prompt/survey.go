package prompt

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
		log.Printf("error occured when getting input from user: %s", err)
		return false
	}

	return confirmation
}


