## Get Channels

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, err := client.GetChannels(context.Background(), gokick.NewChannelListFilter())
	if err != nil {
		log.Fatalf("Failed to fetch response: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.ChannelsResponseWrapper) {
 Result: ([]gokick.ChannelResponse) (len=1 cap=1) {
  (gokick.ChannelResponse) {
   BannerPicture: (string) (len=78) "https://files.kick.com/images/channel/700014/banner_image/default-banner-2.jpg",
   BroadcasterUserID: (int) 721956,
   Category: (gokick.CategoryResponse) {
    ID: (int) 0,
    Name: (string) "",
    Thumbnail: (string) ""
   },
   ChannelDescription: (string) "",
   Slug: (string) (len=7) "scorfly",
   Stream: (struct { Key string "json:\"key\""; URL string "json:\"url\"" }) {
    Key: (string) (len=56) "sk_us-west-2_OB1yyyyy",
    URL: (string) (len=53) "rtmps://uuuu.global-contribute.live-video.net"
   },
   StreamTitle: (string) ""
  }
 }
}
```

## Patch Channels

### Update Stream title

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, err := client.UpdateStreamTitle(context.Background(), "Test KICK API")
	if err != nil {
		log.Fatalf("Failed to fetch response: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.EmptyResponse) {
}
```

### Update Stream category

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, _ := client.UpdateStreamCategory(context.Background(), 9569)


	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.EmptyResponse) {
}
```

### Update Stream tags

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, _ := client.UpdateStreamTags(context.Background(), []string{"tag1", "tag2"})

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.EmptyResponse) {
}
```

## Get Channel Rewards

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, err := client.GetChannelRewards(context.Background())
	if err != nil {
		log.Fatalf("Failed to fetch response: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.ChannelRewardsResponseWrapper) {
 Result: ([]gokick.ChannelRewardResponse) (len=1 cap=1) {
  (gokick.ChannelRewardResponse) {
   BackgroundColor: (*string)(<nil>),
   Cost: (int) 100,
   Description: (*string)(<nil>),
   ID: (string) (len=26) "01JMxxxxx",
   IsEnabled: (*bool)(<nil>),
   IsPaused: (*bool)(<nil>),
   IsUserInputRequired: (*bool)(<nil>),
   ShouldRedemptionsSkipRequestQueue: (*bool)(<nil>),
   Title: (string) (len=10) "Test Reward"
  }
 }
}
```

## Create Channel Reward

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	req := gokick.CreateChannelRewardRequest{
		Title: "My Reward",
		Cost:  100,
	}

	response, err := client.CreateChannelReward(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to create reward: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.Response[gokick.ChannelRewardResponse]) {
 Result: (gokick.ChannelRewardResponse) {
  BackgroundColor: (*string)(<nil>),
  Cost: (int) 100,
  Description: (*string)(<nil>),
  ID: (string) (len=26) "01JMxxxxx",
  IsEnabled: (*bool)(<nil>),
  IsPaused: (*bool)(<nil>),
  IsUserInputRequired: (*bool)(<nil>),
  ShouldRedemptionsSkipRequestQueue: (*bool)(<nil>),
  Title: (string) (len=10) "My Reward"
 }
}
```

## Update Channel Reward

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	title := "Updated Reward"
	cost := 200
	req := gokick.UpdateChannelRewardRequest{
		Title: &title,
		Cost:  &cost,
	}

	response, err := client.UpdateChannelReward(context.Background(), "01JMxxxxx", req)
	if err != nil {
		log.Fatalf("Failed to update reward: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.Response[gokick.ChannelRewardResponse]) {
 Result: (gokick.ChannelRewardResponse) {
  BackgroundColor: (*string)(<nil>),
  Cost: (int) 200,
  Description: (*string)(<nil>),
  ID: (string) (len=26) "01JMxxxxx",
  IsEnabled: (*bool)(<nil>),
  IsPaused: (*bool)(<nil>),
  IsUserInputRequired: (*bool)(<nil>),
  ShouldRedemptionsSkipRequestQueue: (*bool)(<nil>),
  Title: (string) (len=14) "Updated Reward"
 }
}
```

## Delete Channel Reward

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, err := client.DeleteChannelReward(context.Background(), "01JMxxxxx")
	if err != nil {
		log.Fatalf("Failed to delete reward: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.EmptyResponse) {
}
```