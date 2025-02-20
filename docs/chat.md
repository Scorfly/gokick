## Post Chat Message

```go
	client, err := gokick.NewClient(&http.Client{}, "", "xxx")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	response, err := client.SendChatMessage(context.Background(), 721956, "my message", gokick.MessageTypeUser)
	if err != nil {
		log.Fatalf("Failed to fetch response: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.ChatResponseWrapper) {
 Result: (gokick.ChatResponse) {
  IsSent: (bool) true,
  MessageID: (string) (len=36) "5138d04d-68f8-4eca-aa65-93123f6f97fe"
 }
}
```