## Post Chat Message

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, _ := client.SendChatMessage(context.Background(), 721956, "my message", gokick.MessageTypeUser)

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