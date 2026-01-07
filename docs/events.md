## Get Events Subscriptions

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, _ := client.GetSubscriptions(context.Background())

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.EventsResponseWrapper) {
 Result: ([]gokick.EventResponse) (len=2 cap=2) {
  (gokick.EventResponse) {
   AppID: (string) (len=26) "01JMEFN25GFCxxxxxx",
   BroadcasterUserID: (int) 721956,
   CreatedAt: (string) (len=20) "2025-02-20T23:33:10Z",
   Event: (string) (len=17) "chat.message.sent",
   ID: (string) (len=26) "01JMJVAGE9JQS9xxxxxx",
   Method: (string) (len=7) "webhook",
   UpdatedAt: (string) (len=20) "2025-02-20T23:34:14Z",
   Version: (int) 1
  },
  (gokick.EventResponse) {
   AppID: (string) (len=26) "01JMEFN25GFCxxxxx",
   BroadcasterUserID: (int) 721956,
   CreatedAt: (string) (len=20) "2025-02-20T23:33:10Z",
   Event: (string) (len=16) "channel.followed",
   ID: (string) (len=26) "01JMJVAGF7Rxxxxxx",
   Method: (string) (len=7) "webhook",
   UpdatedAt: (string) (len=20) "2025-02-20T23:34:14Z",
   Version: (int) 1
  }
 }
}
```

## Post Events Subscriptions
```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	subscriptions := []gokick.SubscriptionRequest{
		{
			Name:    gokick.SubscriptionNameChatMessage,
			Version: 1,
		},
		{
			Name:    gokick.SubscriptionNameChannelFollow,
			Version: 1,
		},
	}
	// broadcasterUserID is optional - if nil, uses the authenticated user's channel
	response, _ := client.CreateSubscriptions(context.Background(), gokick.SubscriptionMethodWebhook, subscriptions, nil)

	spew.Dump("response", response)
```

### Post Events Subscriptions with broadcasterUserID

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	subscriptions := []gokick.SubscriptionRequest{
		{
			Name:    gokick.SubscriptionNameChatMessage,
			Version: 1,
		},
	}
	broadcasterUserID := 721956
	response, _ := client.CreateSubscriptions(context.Background(), gokick.SubscriptionMethodWebhook, subscriptions, &broadcasterUserID)

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.CreateSubscriptionsResponseWrapper) {
 Result: ([]gokick.CreateSubscriptionResponse) (len=2 cap=2) {
  (gokick.CreateSubscriptionResponse) {
   Error: (string) "",
   Name: (string) (len=17) "chat.message.sent",
   SubscriptionID: (string) (len=26) "01JMJVAGE9JQS97VRBxxxx",
   Version: (int) 1
  },
  (gokick.CreateSubscriptionResponse) {
   Error: (string) "",
   Name: (string) (len=16) "channel.followed",
   SubscriptionID: (string) (len=26) "01JMJVAGF7RET51N1Gyyyyyy",
   Version: (int) 1
  }
 }
}
```

## Delete Events Subscriptions

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, _ := client.DeleteSubscriptions(context.Background(), gokick.NewSubscriptionToDeleteFilter().SetIDs([]string{"01JMMNxxxx"}))

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.EmptyResponse) {
}
```