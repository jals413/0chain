#!/bin/bash

keys=$1
values=$2

./zwalletcli/zwallet --wallet testing.json faucet \
  --methodName pour --input "{Pay day}" --tokens 1

./zwalletcli/zwallet --wallet testing.json sc-update-config \
    --keys $1 --values $2