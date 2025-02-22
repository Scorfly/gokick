## Get Public Key

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	response, _ := client.GetPublicKey(context.Background())

	spew.Dump("response", response)
```
output
```
(string) (len=8) "response"
(gokick.PublicKeyResponseWrapper) {
 Result: (gokick.PublicKeyResponse) {
  PublicKey: (string) (len=450) "-----BEGIN PUBLIC KEY-----\nMIIBIxxxxDAQAB\n-----END PUBLIC KEY-----"
 }
}
```