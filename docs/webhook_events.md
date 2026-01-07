## Validate event (without parsing)

```go
	req := http.Request{} // your request here
	body, _ := io.ReadAll(req.Body)
	
	isValid := gokick.ValidateEvent(req.Header, body)
	
	if !isValid {
		log.Fatal("Event signature is invalid")
	}
```

## Validate and parse event

Headers
```json
"User-Agent":  "Go-http-client/1.1"
"Content-Type": "application/json"
"Accept":  "application/json"
"Kick-Event-Signature":  "EINDkB8ZBed…bCdBLuguc8yfAjXKEvtvVNfhQ=="
"Kick-Event-Message-Timestamp": "2025-02-21T23:23:36Z"
"Kick-Event-Version": "1"
"Kick-Event-Type":  "chat.message.sent"
"Kick-Event-Subscription-Id": "01JMN13xxxxxx"
"Kick-Event-Message-Id": "01JMND5PSxxxxxx"
```

```go
    subscriptionName, _ := gokick.NewSubscriptionName("chat.message.sent")  // value from "Kick-Event-Type" header
	response, _ := gokick.ValidateAndParseEvent(
		subscriptionName,
		"1",                                                                // value from "Kick-Event-Version" header
		"EINDkB8ZBed…bCdBLuguc8yfAjXKEvtvVNfhQ==",                          // value from "Kick-Event-Version" header
		"01JMND5PSxxxxxx",
		"2025-02-21T23:23:36Z",                                             // value from "Kick-Event-Message-Timestamp" header
		`{"message_id":"bb9832e4-e865-48f4…"content":"Test [emote:39261:kkHuh] test[emote:39265:EDMusiC]","emotes":null}`,
	)

	event := response.(*gokick.ChatMessageEvent) // need to cast the type depending of the subscriptionName

	spew.Dump("event", event)
```
output
```
(string) (len=5) "event"
(*gokick.ChatMessageEvent)(0xc000025110)({
 MessageID: (string) (len=36) "bb9832e4-e865-48f4-a0c3-392f78bf3b1a",
 Broadcaster: (gokick.BroadcasterEvent) {
  IsAnonymous: (bool) false,
  UserID: (int) 721956,
  Username: (string) (len=7) "Scorfly",
  IsVerified: (bool) false,
  ProfilePicture: (string) (len=117) "https://files.kick.com/images/user/721956/profile_image/conversion/44a9f1fb-0498-47b5-820e-ef9399fd23d4-fullsize.webp",
  ChannelSlug: (string) (len=7) "scorfly"
 },
 Sender: (gokick.BroadcasterEvent) {
  IsAnonymous: (bool) false,
  UserID: (int) 721956,
  Username: (string) (len=7) "Scorfly",
  IsVerified: (bool) false,
  ProfilePicture: (string) (len=117) "https://files.kick.com/images/user/721956/profile_image/conversion/44a9f1fb-0498-47b5-820e-ef9399fd23d4-fullsize.webp",
  ChannelSlug: (string) (len=7) "scorfly"
 },
 Content: (string) (len=6) "Test [emote:39261:kkHuh] test[emote:39265:EDMusiC]",
 Emotes: ([]gokick.ChatMessageEmotesEvent) <nil>,
 CreatedAt: (string) (len=20) "2025-02-21T23:23:36Z",
 RepliesTo: (struct {
  MessageID string "json:\"message_id\"";
  Sender gokick.UserEvent "json:\"sender\"";
  Content string "json:\"content\""
 }) {
  MessageID: (string) "",
  Sender: (gokick.UserEvent) {
   IsAnonymous: (bool) false,
   UserID: (int) 0,
   Username: (string) "",
   IsVerified: (bool) false,
   ProfilePicture: (string) "",
   ChannelSlug: (string) ""
  },
  Content: (string) ""
 }
})
```

## Get event from HTTP Request

```go
    req := http.Request{} // your request here
	response, _ := gokick.GetEventFromRequest(req)

	event := response.(*gokick.ChatMessageEvent) // need to cast the type depending of the subscriptionName

	spew.Dump("event", event)
```

## Kicks Gifted Event

```go
    subscriptionName, _ := gokick.NewSubscriptionName("kicks.gifted")
	response, _ := gokick.ValidateAndParseEvent(
		subscriptionName,
		"1",
		"EINDkB8ZBed…bCdBLuguc8yfAjXKEvtvVNfhQ==",
		"01JMND5PSxxxxxx",
		"2025-02-21T23:23:36Z",
		`{"broadcaster":{...},"sender":{...},"gift":{"amount":100,"name":"Kick","type":"kick","tier":"1","message":"Thanks!"},"created_at":"2025-02-21T23:23:36Z"}`,
	)

	event := response.(*gokick.KicksGiftedEvent)

	spew.Dump("event", event)
```

## Channel Reward Redemption Updated Event

```go
    subscriptionName, _ := gokick.NewSubscriptionName("channel.reward.redemption.updated")
	response, _ := gokick.ValidateAndParseEvent(
		subscriptionName,
		"1",
		"EINDkB8ZBed…bCdBLuguc8yfAjXKEvtvVNfhQ==",
		"01JMND5PSxxxxxx",
		"2025-02-21T23:23:36Z",
		`{"id":"01JMxxxxx","user_input":"test","status":"fulfilled","redeemed_at":"2025-02-21T23:23:36Z","reward":{"id":"01JMxxxxx","title":"Test Reward","cost":100,"description":"A test reward"},"redeemer":{...},"broadcaster":{...}}`,
	)

	event := response.(*gokick.ChannelRewardRedemptionUpdatedEvent)

	spew.Dump("event", event)
```