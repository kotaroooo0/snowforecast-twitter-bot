# SnowforecastTwitterBot

## Twitter

- [@snowfall_bot](https://twitter.com/snowfall_bot)

<img width="400" alt="スクリーンショット 2020-04-24 19 13 29" src="https://user-images.githubusercontent.com/31947384/80201900-d2826f00-865f-11ea-95bb-5e3d475ba5d4.png">

## Setup

```
$ wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
$ tar -xvf go1.13.3.linux-amd64.tar.gz
$ sudo mv go /usr/local
$ echo 'export GOROOT=/usr/local/go' >> ~/.profile
$ echo 'export GOPATH=$HOME' >> ~/.profile
$ echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH'  >> ~/.profile

$ apt-get update
$ apt-get install git-all

# Get source code, then
$ go run main.go &
```

A example of `.env`

```
CONSUMER_KEY=<your twitter app's consumer key>
CONSUMER_SECRET=<your twitter app's consumer secret>
ACCESS_TOKEN_KEY=<your account's access token>
ACCESS_TOKEN_SECRET=<your account's access secret>
```

## Author

kotaroooo0

## LICENSE

Apache License
