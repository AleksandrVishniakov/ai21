package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AleksandrVishniakov/ai21"
)

func main() {
	ctx := context.Background()
	client := ai21.NewClient("your_api_token")
	conversation := ai21.NewConversation(
		client,
		ai21.WithInitialMessage(
			"You are experienced senior Go developer. You should solve problems for the junior developer. Explain with simple words",
		),
	)

	fmt.Printf("--- Chat with %s ---\n", ai21.ModelJambaInstructPreview)
	fmt.Print(">>> ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		resp, err := conversation.CompletionRequest(ctx, scanner.Text())
		if err != nil {
			log.Fatalf("response error: %s\n", err.Error())
		}

		fmt.Println(resp.Content())
		fmt.Print(">>> ")
	}
	fmt.Printf("\n Conversation end.\n")
}
