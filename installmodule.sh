#!/usr/bin/env bash
readonly REPO=$1

rm -rf temp/* \
&& git clone https://github.com/$REPO.git temp \
&& rm -rf vendor/github.com/$REPO/* \
&& mkdir -p vendor/github.com/$REPO \
&& mv temp/* vendor/github.com/$REPO/ && rm -rf temp/*