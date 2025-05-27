## Ban a user

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, _ := client.BanUser(context.Background(), 721956, 34242, nil, nil)

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.BanUserResponseWrapper) {
 Result: (gokick.BanUserResponse) {
 }
}
```

## Unban a user

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, _ := client.UnbanUser(context.Background(), 721956, 34242)

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.BanUserResponseWrapper) {
 Result: (gokick.BanUserResponse) {
 }
}
```
