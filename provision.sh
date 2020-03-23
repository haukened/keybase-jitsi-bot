#!/usr/bin/env bash
keybase --no-auto-fork \
    oneshot \
    -u $KEYBASE_USERNAME \
    --paperkey "$(cat /run/secrets/$KEYBASE_USERNAME-paperkey)"
./app