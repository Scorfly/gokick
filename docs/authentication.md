## Refresh token

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		ClientID:        "01JMFMAxxxx",
		ClientSecret:    "894b8190xxxxxx",
	})

	response, _ := client.RefreshToken(context.Background(), "your-refresh-token")

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.TokenResponse) {
 AccessToken: (string) (len=48) "ZJYZM2QYZMIxxxxx",
 TokenType: (string) (len=6) "Bearer",
 ExpiresIn: (int) 7200,
 Scope: (string) (len=79) "user:read chat:write channel:read channel:write streamkey:read events:subscribe",
 RefreshToken: (string) (len=48) "MJJIYWU2ZMMTNZxxxx"
}
```