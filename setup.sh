#!/bin/bash

wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz

tar -xvf go1.13.3.linux-amd64.tar.gz

sudo mv go /usr/local

echo 'export GOROOT=/usr/local/go' >> ~/.profile

echo 'export GOPATH=$HOME' >> ~/.profile

echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH'  >> ~/.profile

apt-get update

apt-get install git-all

git pull https://github.com/kotaroooo0/snowforecast-twitter-bot.git
