## 降雪量予報 bot

https://twitter.com/snowfall_sb

GCE(debian9) で本番稼働

```
$ wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
$ tar -xvf go1.13.3.linux-amd64.tar.gz
$ sudo mv go /usr/local
$ echo 'export GOROOT=/usr/local/go' >> ~/.profile
$ echo 'export GOPATH=$HOME' >> ~/.profile
$ echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH'  >> ~/.profile

$ apt-get update
$ apt-get install git-all

ソースを持ってきて
$ go run main.go &
```
