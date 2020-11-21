go-chat
=======

I'm on my bullshit about chat-ops again.

Examples
--------

The examples are an attempt to keep track of how the heck I think
this library should work.

### everyx

Drops a message into the channel(s) every X minutes.

If you run it without any special environment variables set, it
runs in terminal mode, using stdout as "the channel" and reading
standard input for messages.

For Slack mode, set the `BOT_SLACK_TOKEN` to your bot's API token.
You can also set `BOT_SLACK_CHANNEL` to a single channel to join.
Otherwise, the channel "testing" is assumed.
