## Get Kicks Leaderboard

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, err := client.GetKicksLeaderboard(context.Background(), gokick.NewKicksLeaderboardFilter().SetTop(10))
	if err != nil {
		log.Fatalf("Failed to fetch leaderboard: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.KicksLeaderboardResponseWrapper) {
 Result: (gokick.KicksLeaderboardResponse) {
  Lifetime: ([]gokick.KicksLeaderboardEntry) (len=10 cap=10) {
   (gokick.KicksLeaderboardEntry) {
    GiftedAmount: (int) 5000,
    Rank: (int) 1,
    UserID: (int) 123456,
    Username: (string) (len=7) "user1"
   },
   ...
  },
  Month: ([]gokick.KicksLeaderboardEntry) (len=10 cap=10) {
   (gokick.KicksLeaderboardEntry) {
    GiftedAmount: (int) 2000,
    Rank: (int) 1,
    UserID: (int) 789012,
    Username: (string) (len=7) "user2"
   },
   ...
  },
  Week: ([]gokick.KicksLeaderboardEntry) (len=10 cap=10) {
   (gokick.KicksLeaderboardEntry) {
    GiftedAmount: (int) 1000,
    Rank: (int) 1,
    UserID: (int) 345678,
    Username: (string) (len=7) "user3"
   },
   ...
  }
 }
}
```

