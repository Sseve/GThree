#!/bin/bash

## gtservant程序管理脚本

set -e

# 使用普通用户运行该脚本
if [ "$UID" == "0" ];then
  echo "请使用普通用户运行该脚本"
  exit 1
fi

# 当前程序目录
GT_PATH=$(cd $(dirname "$0") && pwd)

# 启动gtservant程序
gtservant_start(){
  ${GT_PATH}/gtservant >/dev/null 2>&1 &
  if [ $? -eq 0 ];then
    sleep 5
    echo "gtservant start success"
    pgrep -l gtservant | awk '{print $1}' >${GT_PATH}/pid.txt
  else
    echo "gtservant start failed"
  fi
}

# 停止gtservant程序
gtservant_stop(){
  if [[ "$(gtservant_check)" == "0" ]];then
    local pid=$(cat ${GT_PATH}/pid.txt)
    kill $pid && >${GT_PATH}/pid.txt
  else
    echo "gtservant is stop"
  fi
}

# 检查gtservant程序
gtservant_check(){
  local pid=$(cat ${GT_PATH}/pid.txt)
  if [[ "$pid" ]];then
    echo 0
  else
    echo 1
  fi 
}

case $1 in
  start)
    gtservant_start
    ;;
  stop)
    gtservant_stop
    ;;
  check)
    if [ "$(gtservant_check)" == "0" ];then
        echo "gtservant is running..."
    else
        echo "gtservant is stopped..."
    fi
    ;;
  *)
  echo "Usage: $0 start|stop|check"
  echo "  start  -- start gtservant"
  echo "  stop   -- stop gtservant"
  echo "  check  -- check gtservant"
  exit 1
esac