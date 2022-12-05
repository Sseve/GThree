#!/bin/bash

if [ "$UID" == "0" ];then
    echo "请使用普通用户运行此脚本!"
    exit 0
fi

# build gtmaster and start it
gtmaster() {
    if [ "$1" == "server" ];then
        cd cmd/gtmaster && go build -ldflags="-s -w"
        if [ -f gtmaster ];then
            cp -f gtmaster ../../release/gtmaster && \
            rm gtmaster && cd ../../release/gtmaster && \
            cp ../../config/gtmaster.yaml .
            
        else
            echo "build gtmaster failed"
        fi
    fi
}

# build gtservant and start it
gtservant() {
    if [ "$1" == "client" ];then
        cd cmd/gtservant && go build -ldflags="-s -w"
        if [ -f gtservant ];then
            cp -f gtservant ../../release/gtservant && \
            rm gtservant && cd ../../release/gtservant && \
            cp ../../config/gtservant.yaml . && \
            scp -r ./* gamecpp@172.16.9.128:/home/gamecpp/gtservant
        else
            echo "build gtservant failed!"
        fi
    fi
}

case "$1" in
    master)
        gtmaster server
        ;;
    servant)
        gtservant client
        ;;
    *)
    echo "Usage: $0 [server|client]"
    echo "    master     构建gtmaster服务"
    echo "    servant    构建gtservant服务"    
    exit 1
esac