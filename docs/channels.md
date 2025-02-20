## Get Channels

```go
	client, err := gokick.NewClient(&http.Client{}, "", "xxx")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

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
	client, err := gokick.NewClient(&http.Client{}, "", "xxxx")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

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
	client, err := gokick.NewClient(&http.Client{}, "", "xxx")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	response, err := client.UpdateStreamCategory(context.Background(), 9569)
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