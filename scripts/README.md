## generate_user_access_token.go

This script is designed to generate a user access token for a specific user in a GoKICK application.

This is an example of how to use the OAuth 2.1 flow to generate a user access token.
It’s not design to be used in production, but rather as a tool to generate a user access token for testing purposes.

### How to run it

```sh
$ KICK_CLIENT_ID=01JMFMARZ9GN12JNCTEZWWWGRE \
    KICK_CLIENT_SECRET=894b819067849af1a50defdd47f62a0204a7b17c149d028fe31a1d64d49ff9b0 \
    go run scripts/generate_user_access_token.go 
```

It will launch an HTTP server on port 3000

### Create access token

#### Step 1

Open your browser and go to `http://localhost:3000/oauth/kick/` .
This will redirect you to the KICK authorization page.

#### Step 2

On the KICK page, validate the permissions and click on "Authorize".

#### Step 3

You will be redirected to `http://localhost:3000/oauth/kick/callback` .

This will display the access token in the browser.

Example:

```json
{
  "access_token": "NWU2YJNKxxxxx",
  "expires_in": 7200,
  "refresh_token": "NDBMZMJJZxxxxx",
  "scope": "user:read chat:write channel:read channel:write streamkey:read events:subscribe",
  "token_type": "Bearer"
}
```
