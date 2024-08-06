package ai21

type Role string

const (
	// RoleSystem is a Role for providing initial instructions
	// to provide general guidance on the tone and voice of the generated message.
	RoleSystem Role = "system"

	// RoleAssistant is a Role for AI model responses.
	RoleAssistant Role = "assistant"

	// RoleUser is a Role for user requests.
	RoleUser = "user"
)
