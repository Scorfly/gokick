## Get Public Key

```go
	client, err := gokick.NewClient(&http.Client{}, "", "xxxx")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	response, err := client.GetPublicKey(context.Background())
	if err != nil {
		log.Fatalf("Failed to fetch response: %v", err)
	}

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