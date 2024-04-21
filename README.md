# TwitchGO

TwitchGO is an easy to use library for Twitch API and IRC Chat implemented in GO.

You can create a connection as simple as

```go
bot := twitchgo.New("myClientID", "myClientSecret", "myIRCToken")

err := bot.Connect()
```

Remember to never use your bots secret and token in plain text!

### Development

It is still under development. Errors and unexpected behavior may occour. And a large potion of the code is not documented well
