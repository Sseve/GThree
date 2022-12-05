#!/bin/bash

## 管理游戏进程的脚本

GAME_PATH=$(cd "$(dirname $0)"; pwd)
PID_FILE=${GAME_PATH}/tmp/gameserv.pid

# 检查游戏进程状态
gs_ck() {
    local pid=$(cat ${PID_FILE})
    if [[ -L /proc/${pid}/exe ]]; then
        echo "{\"code\": \"0\", \"msg\": \"starting\"}"
        exit 1
    else
        echo "{\"code\": \"1\", \"msg\": \"stopped\"}"
    fi
}

# 启动游戏进程
gs_st() {
    [[ ! -d ${GAME_PATH}/tmp ]] && mkdir -p ${GAME_PATH}/tmp && touch $PID_FILE
    result=$(gs_ck)
    ret=$(echo ${result}|grep -o 'starting')
    if [[ "${ret}" ]]; then
        echo "{\"code\": \"0\", \"msg\": \"starting\"}"
        exit 1
    fi
    if [[ -f gameserv-beta ]];then
        mv gameserv-beta gameserv
    fi
    chmod +x gameserv
    ${GAME_PATH}/gameserv
    if [[ $? -eq 0 ]];then
        echo "{\"code\": \"0\", \"msg\": \"startSucess\"}"
    else
        echo "{\"code\": \"1\", \"msg\": \"startFailed\"}"
    fi
}

# 关闭游戏进程
gs_sp() {
    local pid=$(cat ${PID_FILE})
    local apid=$(ps aux|grep "${GAME_PATH}"|grep -v 'grep'|awk '{print $2}')
    if [[ "${apid}" ]] && [[ "${pid}" ]];then
        kill -15 ${pid} && >${PID_FILE}
        if [[ $? -eq 0 ]];then
            echo "{\"code\": \"0\", \"msg\": \"stopSucess\"}"
        else
            echo "{\"code\": \"1\", \"msg\": \"stopFailed\"}"
        fi
    else
        echo "{\"code\": \"1\", \"msg\": \"stopped\"}"
    fi
}

# 重载游戏进程
gs_reload() {
    local pid=$(cat ${PID_FILE})
    local apid=$(ps aux|grep "${GAME_PATH}"|grep -v 'grep'|awk '{print $2}')
    if [[ ${pid} = ${apid} ]];then
        kill -10 ${pid} #&& >${PID_FILE}
        if [[ $? -eq 0 ]];then
            echo "{\"code\": \"0\", \"msg\": \"reloadSucess\"}"
        else
            echo "{\"code\": \"1\", \"msg\": \"reloadFailed\"}"
        fi
    else
        echo "{\"code\": \"1\", \"msg\": \"stopped\"}"
    fi
}


# �����Ϸ״̬ 
gs_ss() {
    gs_ck
}


case "$1" in
    start)
        gs_st
        ;;
    stop)
        gs_sp
        ;;
    reload)
	gs_reload
   	 ;;
    check)
        gs_ss
        ;;
    *)
    echo $"Usage: $0 {start|stop|check}"
    exit 1
esac