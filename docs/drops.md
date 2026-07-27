## Get Drops Claims

Only OAuth apps associated with a Kick organization can access this endpoint. Use an app access token.

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		AppAccessToken: "xxxx",
	})

	response, err := client.GetDropsClaims(
		context.Background(),
		gokick.NewDropsClaimsFilter().
			SetCampaignID("01JAXK8N4QWRTY5PM7ZEBVJDS8S").
			SetLimit(10),
	)
	if err != nil {
		log.Fatalf("Failed to fetch claims: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.DropsClaimsResponseWrapper) {
 Result: (gokick.DropsClaimsData) {
  Claims: ([]gokick.DropsClaimResponse) (len=1 cap=1) {
   (gokick.DropsClaimResponse) {
    CampaignID: (string) (len=26) "01JAXK8N4QWRTY5PM7ZEBVJDS8S",
    ClaimID: (string) (len=26) "01JAXK8N4QWRTY5PM7ZEBVGH2S",
    CreatedAt: (string) (len=20) "2026-01-15T12:00:00Z",
    ExternalID: (string) (len=7) "ext-123",
    ExternalStatus: (string) (len=9) "initiated",
    RewardID: (string) (len=26) "01JAXK8N4QWRTY5PM7ZEBVJDS9T",
    UpdatedAt: (string) (len=20) "2026-01-15T12:00:00Z",
    UserID: (int) 1234
   }
  },
  Cursor: (string) (len=26) "01K0TDDR08Q5ZNWXK92SDH7SDH"
 }
}
```

### Filter by claim ID

```go
	response, err := client.GetDropsClaims(
		context.Background(),
		gokick.NewDropsClaimsFilter().SetClaimID("01JAXK8N4QWRTY5PM7ZEBVGH2S"),
	)
```

## Update Drops Claims

Update `external_status` for up to 100 claims. On success the API returns 204 No Content.

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		AppAccessToken: "xxxx",
	})

	response, err := client.UpdateDropsClaims(context.Background(), []gokick.DropsClaimUpdate{
		{ClaimID: "01KAAFHJ2PNXS48NG8XXPGWKCZ", ExternalStatus: "processed"},
		{ClaimID: "01KAAFHJ2PNXS48NG8XXPGWKCA", ExternalStatus: "initiated"},
	})
	if err != nil {
		log.Fatalf("Failed to update claims: %v", err)
	}

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.EmptyResponse) {
}
```
