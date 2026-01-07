## Post Chat Message

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	// broadcasterUserID is optional - if nil, uses the authenticated user's channel
	broadcasterUserID := 721956
	response, _ := client.SendChatMessage(context.Background(), &broadcasterUserID, "my message", nil, gokick.MessageTypeUser)

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

### Reply to a message

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	broadcasterUserID := 721956
	replyToMessageID := "5138d04d-68f8-4eca-aa65-93123f6f97fe"
	response, _ := client.SendChatMessage(context.Background(), &broadcasterUserID, "my reply", &replyToMessageID, gokick.MessageTypeUser)

	spew.Dump("response", response)
```

## Delete Chat Message

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, err := client.DeleteChatMessage(context.Background(), "5138d04d-68f8-4eca-aa65-93123f6f97fe")
	if err != nil {
		log.Fatalf("Failed to delete message: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.EmptyResponse) {
}
```