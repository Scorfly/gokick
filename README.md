![GoKICK logo](assets/gokick-logo.png)

# GoKICK

A comprehensive Go client library for the [Kick.com](https://kick.com) API. GoKICK provides a type-safe, idiomatic Go interface to interact with Kick's streaming platform, including authentication, channels, chat, moderation, livestreams, and webhooks.

[![Coverage Status](https://coveralls.io/repos/github/Scorfly/gokick/badge.svg)](https://coveralls.io/github/Scorfly/gokick)
[![Release](https://img.shields.io/github/release/Scorfly/gokick.svg?color=%23007ec6)](https://github.com/Scorfly/gokick/releases/latest)

## Features

- 🔐 **Complete OAuth2 Authentication** - User and app access token management with automatic token refresh
- 📺 **Channel Management** - Get channels, update stream metadata, manage channel rewards
- 💬 **Chat Operations** - Send and delete chat messages
- 🛡️ **Moderation Tools** - Ban and unban users
- 📊 **Livestream Data** - Get livestreams and statistics
- 🎁 **Kicks & Rewards** - Access leaderboards and manage channel rewards
- 🔔 **Webhook Events** - Subscribe to and handle webhook events
- 🏷️ **Categories & Users** - Browse categories and user information
- 🔄 **Auto Token Refresh** - Automatic user token refresh with callback support
- 🧪 **Well Tested** - Comprehensive test coverage

## Installation

```bash
go get github.com/scorfly/gokick
```

## Quick Start

### Basic Client Setup

```go
package main

import (
    "context"
    "fmt"
    "github.com/scorfly/gokick"
)

func main() {
    // Create a client with app access token
    client, err := gokick.NewClient(&gokick.ClientOptions{
        AppAccessToken: "your-app-access-token",
    })
    if err != nil {
        panic(err)
    }

    // Get categories
    categories, err := client.GetCategories(context.Background())
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d categories\n", len(categories.Result))
}
```

### OAuth2 Authentication Flow

```go
// 1. Get authorization URL
client, _ := gokick.NewClient(&gokick.ClientOptions{
    ClientID: "your-client-id",
})

authURL, _ := client.GetAuthorize(
    "http://localhost:3000/oauth/kick/callback",
    "state-value",
    "code-challenge",
    []gokick.Scope{gokick.ScopeUserRead, gokick.ScopeChannelRead},
)

// 2. Redirect user to authURL, then exchange code for token
client, _ = gokick.NewClient(&gokick.ClientOptions{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
})

token, _ := client.GetToken(
    context.Background(),
    "http://localhost:3000/oauth/kick/callback",
    "authorization-code",
    "code-verifier",
)

// 3. Use the token
client.SetUserAccessToken(token.AccessToken)
```

### Automatic Token Refresh

```go
client, _ := gokick.NewClient(&gokick.ClientOptions{
    ClientID:        "your-client-id",
    ClientSecret:    "your-client-secret",
    UserAccessToken: "current-access-token",
    UserRefreshToken: "refresh-token",
})

// Set callback to save new tokens when refreshed
client.OnUserAccessTokenRefreshed(func(accessToken, refreshToken string) {
    // Save tokens to your storage
    fmt.Printf("Tokens refreshed: %s\n", accessToken)
})

// Client will automatically refresh expired tokens
channels, err := client.GetChannels(context.Background(), gokick.NewChannelListFilter())
```

## Documentation

### API Reference

See the [documentation directory](docs/README.md) for detailed examples and supported endpoints:

- [Authentication](docs/authentication.md) - OAuth2 flows, token management
- [Channels](docs/channels.md) - Channel operations and rewards
- [Chat](docs/chat.md) - Send and manage chat messages
- [Moderation](docs/moderation.md) - Ban/unban users
- [Livestreams](docs/livestreams.md) - Get livestream data and stats
- [Users](docs/users.md) - User information and token introspection
- [Categories](docs/categories.md) - Browse categories
- [Kicks](docs/kicks.md) - Kicks leaderboard
- [Events](docs/events.md) - Webhook event subscriptions
- [Webhook Events](docs/webhook_events.md) - Webhook payload structures

### Supported Endpoints

GoKICK supports a comprehensive set of Kick API endpoints:

- ✅ **Authentication** - Authorization, token exchange, refresh, revoke, app tokens
- ✅ **Categories** - List and get category details
- ✅ **Users** - Get users, token introspection
- ✅ **Channels** - Get channels, update stream metadata, manage rewards
- ✅ **Chat** - Send and delete messages
- ✅ **Moderation** - Ban and unban users
- ✅ **Livestreams** - Get livestreams and statistics
- ✅ **Public Key** - Get public key for webhook verification
- ✅ **Kicks** - Get kicks leaderboard
- ✅ **Events** - Subscribe to webhook events

See the [full feature list](docs/README.md) for complete details.

## Development

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Requirements

- Go 1.25.5 or later

## Kick API Documentation

For official Kick API documentation, see the [KICK Developer website](https://dev.kick.com/).

## License

This package is distributed under the terms of the [MIT](LICENSE) license.