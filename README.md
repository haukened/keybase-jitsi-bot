[![Build](https://github.com/haukened/keybase-jitsi-bot/workflows/Build/badge.svg)](https://github.com/haukened/keybase-jitsi-bot/actions)

# keybase-jitsi-bot
A bot for Keybase that start Jitsi meetings

This package requires the keybase binary installed on your system, and works on linux, macOS, and Windows 10

#### Tested on:
 - Ubuntu Latest
 - macOS Latest
 - Windows Latest

## Running on the command line:
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

## Running in the docker container:
#### Pulling the container:

`docker pull haukeness/keybase-jitsi-bot`

#### Running the container:
You need to set ENV vars instead of passing command line flags:

Required by keybase: (Must set all of these)
 - `KEYBASE_USERNAME=foo`
 - `KEYBASE_PAPERKEY="bar baz ..."`
 - `KEYBASE_SERVICE=1`
 
Required by this package: (Set the values you feel like, if you don't set them they won't be used)
 - `BOT_DEBUG=true`
 - `BOT_LOG_CONVID=<your keybase conversation id>`
 - `BOT_FEEDBACK_CONVID=<your keybase conversation id>`
 - `BOT_FEEDBACK_TEAM_ADVERT="@team#channel"`

#### Example:
`docker run --name myJitsi --rm -d -e KEYBASE_USERNAME=FOO -e KEYBASE_PAPERKEY="bar baz ..." -e KEYBASE_SERVICE=1 -e BOT_DEBUG=true haukeness/keybase-jitsi-bot`
