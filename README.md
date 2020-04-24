[![Build](https://github.com/haukened/keybase-jitsi-bot/workflows/Build/badge.svg)](https://github.com/haukened/keybase-jitsi-bot/actions)

# keybase-jitsi-bot
A bot for Keybase that start Jitsi meetings

## Running in the docker container:
#### Pulling the container:

`docker pull haukeness/keybase-jitsi-bot`

#### Running the container:
Firstly, this container has been updated to work with a docker swarm service that can pass encrypt secrets and pass the decrpyted secret to the running container.  This helps in multi-tenancy applications where you don't want the paperkey sitting in an environment variable for everyone to see.  
The container expects docker swarm to mount the secret file at `/run/secrets/$BOT_NAME-paperkey` and is read from the file when the application starts the keybase service.  This is done by `provision.sh` so if you still want to pass the `KEYBASE_PAPERKEY` env var you'll need to change the entrypoint to launch `./app` directly.

##### You need to set ENV vars instead of passing command line flags:

Required by keybase: (Must set all of these)

 - `KEYBASE_USERNAME=foo`
 - `KEYBASE_SERVICE=1`

Optional for keybase: (Your mileage may vary!)

 - `KEYBASE_SERVICE_ARGS="-enable-bot-lite-mode=1"`
 
Used by this package: (Set the values you feel like, if you don't set them they won't be used.)
 - `BOT_DEBUG=true` 
   - This enables debugging to the console and logging chat.
 - `BOT_LOG_CONVID=<your keybase conversation id>` 
   - The keybase ConversationID you want the bot to send logs to (if you want the bot to log to keybase team chat).  If this value is not set, the bot will only log to the docker console.
 - `BOT_FEEDBACK_CONVID=<your keybase conversation id>` 
   - The keybase ConversationID where you want the bot to send feedback to. If this value is not set, the bot won't be able to send feedback, and will instead tell users you didn't configure feedback.
 - `BOT_FEEDBACK_TEAM_ADVERT="@team#channel"` 
   - The human readable name of the keybase team where the bot will send feedback.  This informs users where the feedback will be sent to in the bot command help menu when sending feedback.  If you don't set this, it will still work but users won't know where feedback goes to.
 - `BOT_KVSTORE_TEAM="<your team without the @symbol>"` 
   - This is the name of the team where the bot with store all settings in the keybase kvstore.  Its important that this is a private team that only you have access to, and the bot must be at least a writer in the team.

This package requires the keybase binary installed on your system, and works on linux, macOS, and Windows 10

#### Tested on:
 - Ubuntu Latest
 - macOS Latest
 - Windows Latest

## Running on the command line:
Really this is just for testing - don't do this.

#### Installation:
 - `git clone https://github.com/haukened/keybase-jitsi-bot.git`
 - `cd keybase-jitsi-bot`
 - `go get -u ./...`
 - `go build`
 - `go install`
 
#### Running:
```
  -debug
        enables command debugging to stdout
  -feedback-convid string
        sets the keybase chat1.ConvIDStr to send feedback to.
  -feedback-team-advert string
        sets the keybase team/channel to advertise feedback. @team#channel
  -log-convid string
        sets the keybase chat1.ConvIDStr to log debugging to keybase chat.
```

#### Example: 
`jitsi-bot --debug`
