# Supported Endpoints & Features

## APIs

**Authentication:**

- [x] Authorization Endpoint
- [x] Token Endpoint
- [x] App Access Token Endpoint
- [x] Refresh Token Endpoint
- [x] Revoke Token Endpoint

**Categories:**

- [x] Get Categories (`GET /public/v2/categories`)
- [x] Get Category

**Users:**

- [x] Token Introspect (`POST /oauth/token/introspect`)
- [x] Get Users

**Channels:**

- [x] Get Channels
- [x] Patch Channels
  - [x] Update Stream title
  - [x] Update Stream category
  - [x] Update Stream tags
- [x] Get Channel Rewards
- [x] Create Channel Reward
- [x] Update Channel Reward
- [x] Delete Channel Reward
- [x] Get Channel Reward Redemptions
- [x] Accept Channel Reward Redemptions
- [x] Reject Channel Reward Redemptions

**Chat:**

- [x] Post Chat Message
- [x] Delete Chat Message

**Moderation:**

- [x] Post Moderation Bans
- [x] Delete Moderation Bans

**Livestreams:**

- [x] Get Livestreams V2 (`GET /public/v2/livestreams`) — preferred
- [x] Get Users Livestreams (`GET /public/v1/users/livestreams`)
- [x] Get Livestreams (`GET /public/v1/livestreams`) — **deprecated** (prefer V2 / by users)
- [x] Get Livestreams Stats

**Drops** (organization-linked OAuth apps only):

- [x] Get Drops Claims
- [x] Update Drops Claims

**Public Key:**

- [x] Get Public Key

**Kicks:**

- [x] Get Kicks Leaderboard

## Events

**Subscribe to Events:**

- [x] Get Events Subscriptions
- [x] Post Events Subscriptions
- [x] Delete Events Subscriptions

**Webhook Payloads:**

- [x] Chat Message
- [x] Channel Follow
- [x] Channel Subscription Renewal
- [x] Channel Subscription Gifts
- [x] Channel Subscription Created
- [x] Livestream Status Updated
- [x] Livestream Metadata Updated
- [x] Moderation Banned
- [x] Kicks Gifted
- [x] Channel Reward Redemption Updated
