## Get Livestreams

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "ZMY2OTHxxxxxxxx",
	})

	response, err := client.GetLivestreams(context.Background(), gokick.NewLivestreamListFilter().SetLimit(1))
	if err != nil {
		log.Fatalf("Failed to fetch response: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.LivestreamsResponseWrapper) {
 Result: ([]gokick.LivestreamResponse) (len=1 cap=1) {
  (gokick.LivestreamResponse) {
   BroadcasterUserID: (int) 30111,
   Category: (gokick.CategoryResponseV2) {
    ID: (int) 15,
    Name: (string) (len=13) "Just Chatting",
    ImageURL: (string) (len=90) "https://files.kick.com/images/subcategories/15/banner/b697a8a3-62db-4779-aa76-e4e47662af97"
   },
   ChannelID: (int) 291111,
   HasMatureContent: (bool) false,
   Language: (string) (len=2) "en",
   Slug: (string) (len=14) "inxxxxx",
   StartedAt: (string) (len=20) "2025-04-01T14:38:29Z",
   StreamTitle: (string) (len=20) "Super first stream",
   ThumbnailURL: (string) (len=75) "https://images.kick.com/video_thumbnails/xxxx/yyy/480.webp",
   ViewerCount: (int) 18081,
   CustomTags: ([]string) (len=2 cap=2) {
    (string) (len=4) "tag1",
    (string) (len=4) "tag2"
   },
   ProfilePicture: (string) (len=117) "https://files.kick.com/images/user/30111/profile_image/conversion/44a9f1fb-0498-47b5-820e-ef9399fd23d4-fullsize.webp"
  }
 }
}
```

## Get Livestreams Stats

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, err := client.GetLivestreamsStats(context.Background())
	if err != nil {
		log.Fatalf("Failed to fetch stats: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.LivestreamStatsResponseWrapper) {
 Result: (gokick.LivestreamStatsResponse) {
  TotalCount: (int) 12345
 }
}
```