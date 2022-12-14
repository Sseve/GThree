#!/bin/bash

## gtmaster程序管理脚本

set -e

# 使用普通用户运行该脚本
if [ "$UID" == "0" ];then
  echo "请使用普通用户运行该脚本"
  exit 1
fi

# 当前程序目录
GT_PATH=$(cd $(dirname "$0") && pwd)

# 启动gtmaster程序
gtmaster_start(){
  ${GT_PATH}/gtmaster
  if [ $? -eq 0 ];then
    sleep 3
    echo "gtmaster start success"
    pgrep gtmaster >${GT_PATH}/pid.txt
  else
    echo "gtmaster start failed"
  fi
}

# 停止gtmaster程序
gtmaster_stop(){
  if [[ "$(gtmaster_check)" == "0" ]];then
    local pid=$(cat ${GT_PATH}/pid.txt)
    kill $pid && >${GT_PATH}/pid.txt
  else
    echo "gtmaster is stop"
  fi
}

# 检查gtmaster程序
gtmaster_check(){
  local pid=$(cat ${GT_PATH}/pid.txt)
  if [[ "$pid" ]];then
    echo 0
  else
    echo 1
  fi 
}

case $1 in
  start)
    gtmaster_start
    ;;
  stop)
    gtmaster_stop
    ;;
  check)
    if [ "$(gtmaster_check)" == "0" ];then
        echo "gtmaster is running..."
    else
        echo "gtmaster is stopped..."
    fi
    ;;
  *)
  echo "Usage: $0 start|stop|check"
  echo "  start  -- start gtmaster"
  echo "  stop   -- stop gtmaster"
  echo "  check  -- check gtmaster"
  exit 1
esac