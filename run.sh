#!/bin/bash

git add .
git commit -a -m 'Updated'
git push
# kill $(ps aux | grep '[/]home/mkc/go/.*tinode.conf' | awk '{print $2}')
# kill $(ps aux | grep '[/]home/mkc/go/.*tinode.conf' | awk '{print $2}')
go get -tags rethinkdb github.com/mkc188/chat/server && go install -tags rethinkdb github.com/mkc188/chat/server && $GOPATH/bin/server -config=$GOPATH/src/github.com/mkc188/chat/server/tinode.conf -static_data=$HOME/webapp
