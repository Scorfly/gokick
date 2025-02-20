## Token Introspect

```go
	client, err := gokick.NewClient(&http.Client{}, "", "xxx")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	token, err := client.TokenIntrospect(context.Background())
	if err != nil {
		log.Fatalf("Failed to fetch token: %v", err)
	}

	spew.Dump("token", token)
```
output
```
(string) (len=5) "token"
(gokick.TokenIntrospectResponseWrapper) {
 Result: (gokick.TokenIntrospectResponse) {
  Active: (bool) true,
  ClientID: (string) (len=26) "zzzzzzzzz",
  Exp: (int) 1740079154,
  Scope: (string) (len=9) "user:read"
 }
}
```

## Get Users

```go
	client, err := gokick.NewClient(&http.Client{}, "", "xxx")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	users, err := client.GetUsers(context.Background(), gokick.NewUserListFilter())
	if err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}

	spew.Dump("users", users)
```
output
```
(string) (len=5) "users"
(gokick.UsersResponseWrapper) {
 Result: ([]gokick.UserResponse) (len=1 cap=1) {
  (gokick.UserResponse) {
   Email: (string) (len=17) "user@domain.tld",
   Name: (string) (len=7) "Scorfly",
   ProfilePicture: (string) (len=117) "https://files.kick.com/images/user/721956/profile_image/conversion/44a9f1fb-0498-47b5-820e-ef9399fd23d4-fullsize.webp",
   UserID: (int) 123456
  }
 }
}
```