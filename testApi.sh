#!/bin/bash
if [ "$#" -eq 3 ];
then
    KEY=$URLBOXKEY
    SECRET=$URLBOXSECRET

    QUERY_STRING="url=$1&width=$2&height=$3"
    TOKEN=$(echo -n "$QUERY_STRING" | openssl sha1 -hmac "$SECRET")
    URL="https://api.urlbox.io/v1/$KEY/$TOKEN/png?$QUERY_STRING"

    echo "[+] To request: $URL"
    curl "$URL" >> "$1-$2x$3.png"
else
    echo "Usage: ./testApi.sh URL WIDTH HEIGHT"
fi
