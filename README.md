# TwitchGO

TwitchGO is an easy to use library for Twitch IRC Chat implemented in GO.

You can create a connection as simple as

```go
bot := twitchgo.New("mybot", "mytoken")

err := bot.Connect()
```

Remember to never use your bots token in plain text!

### Development

It is still under development. Errors and unexpected behavior may occour. And a large potion of the code is not documented well
