## Get authorize endpoint

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		ClientID: "client-id",
	})

	response, _ := client.GetAuthorizeEndpoint(
		"http://localhost:3000/oauth/kick/callback",
		"custom-state",
		"custome-code-challenge",
		[]gokick.Scope{gokick.ScopeUserRead, gokick.ScopeChannelRead},
	)
	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(string) (len=347) "https://id.kick.com/oauth/authorize?client_id=client-id&response_type=code&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Foauth%2Fkick%2Fcallback&state=custom-state&scope=user%3Aread+channel%3Aread&code_challenge=custom-code-challenge&code_challenge_method=S256"
```

## Get token

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		ClientID:     "01JMFMxxx",
		ClientSecret: "894b8xxx",
	})

	response, _ := client.GetToken(
		context.Background(),
		"http://localhost:3000/oauth/kick/callback",
		"code",
		"code-verifier",
	)

	spew.Dump("response", response)
```
output
```
(string) (len=10) "response"
(gokick.TokenResponse) {
 AccessToken: (string) (len=48) "MDJMMWNMxxxxx",
 TokenType: (string) (len=6) "Bearer",
 ExpiresIn: (int) 7200,
 Scope: (string) (len=79) "user:read chat:write channel:read channel:write streamkey:read events:subscribe",
 RefreshToken: (string) (len=48) "MGNMxxxx"
}
```

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

## Revoke token

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		ClientID:     "01JMFMAxxxxx",
		ClientSecret: "894b81xxxxx",
	})

	err := client.RevokeToken(
		context.Background(),
		gokick.TokenTypeAccess,
		"MGNMZJxxxx",
	)
```