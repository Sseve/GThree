#!/bin/bash

# 开发测试的游戏区服主机
IP="172.16.9.129"

if [ "$UID" == "0" ];then
    echo "请使用普通用户运行此脚本!"
    exit 0
fi

# 构建gtmaster
gtmaster() {
    if [ "$1" == "gtmaster" ];then
        cd cmd/gtmaster && go build -ldflags="-s -w"
        if [ -f gtmaster ];then
            cp -f gtmaster ../../release/gtmaster && \
            rm gtmaster && cd ../../release/gtmaster && \
            cp ../../config/gtmaster.yaml . && \
            echo "build gtmaster success"
        else
            echo "build gtmaster failed"
        fi
    fi
}

# 构建gtservant
gtservant() {
    if [ "$1" == "gtservant" ];then
        cd cmd/gtservant && go build -ldflags="-s -w"
        if [ -f gtservant ];then
            cp -f gtservant ../../release/gtservant && \
            rm gtservant && cd ../../release/gtservant && \
            cp ../../config/gtservant.yaml . && \
            echo "build gtservant success"
            if [[ "$IP" ]];then
                if  ping -c 1 ${IP} >/dev/null 2>&1;then 
                    scp -r ./* gamecpp@${IP}:/home/gamecpp/gtservant
                fi
            fi
        else
            echo "build gtservant failed!"
        fi
    fi
}

case "$1" in
    gtmaster)
        gtmaster "$1"
        ;;
    gtservant)
        gtservant "$1"
        ;;
    *)
    echo "Usage: $0 [server|client]"
    echo "    gtmaster     构建gtmaster服务"
    echo "    gtservant    构建gtservant服务"    
    exit 1
esac