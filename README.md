# SnowforecastTwitterBot

## Twitter

- [@snowfall_bot](https://twitter.com/snowfall_bot)

  - 毎日主要スキー場の降雪予報をツイートする

    <img width="400" src="https://user-images.githubusercontent.com/31947384/81565326-1ac9bd00-93d4-11ea-8ebf-bb3499d7566c.png">

  - リプライに反応して対応したスキー場の降雪予報をリプライする(漢字でもひらがなでもローマ字でも対応)

    <img width="400" src="https://user-images.githubusercontent.com/31947384/81564307-80b54500-93d2-11ea-82c7-ea5a3adc2f46.png">

    <img width="400" src="https://user-images.githubusercontent.com/31947384/81564354-8f036100-93d2-11ea-96a4-235108bbde9e.png">

## Setup

```
// Install git
$ apt-get update
$ apt-get install git-all

// Install docker-compose
$ sudo curl -L "https://github.com/docker/compose/releases/download/1.25.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose

// Clone this repository
$ git clone https://github.com/kotaroooo0/snowforecast-twitter-bot.git
$ cd snowforecast-twitter-bot

// Set .env
$ vi app/.env

// Boot!
$ docker-compose up -d
```

A example of `.env`

```
CONSUMER_KEY=<your twitter app's consumer key>
CONSUMER_SECRET=<your twitter app's consumer secret>
ACCESS_TOKEN_KEY=<your account's access token>
ACCESS_TOKEN_SECRET=<your account's access secret>
REDIS_HOST=redis
```

## Command for Develop

You can use Makefile.

```
# Hot Reload
# ref: https://github.com/oxequa/realize
$ realize start

# Test
# ref: https://github.com/eaburns/Watch
# ref: https://github.com/izumin5210/cgt
$ Watch -t make test | cgt

# Redis
# Setup Server(Backup: /usr/local/var/db/redis/dump.rdb)
$ redis-server /usr/local/etc/redis.conf

# Initial Data
# 1 is Test DB
$ cat data.txt | redis-cli --pipe
$ cat data.txt | redis-cli -n 1 --pipe

# Delete Selected DB
redis-cli) flushdb
```

## Author

kotaroooo0

## LICENSE

Apache License
