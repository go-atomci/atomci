#!/bin/bash

# func: atomci init
# author: colynn

# prepare verify
DOCKER_NUM=$(docker-compose ps | grep 'Up' | wc -l)

if [ $DOCKER_NUM -ne 3 ]
then
    NOT_UP=$(docker-compose ps | grep -v "Name\|-----\|Up" | awk '{print $1}')
    for dockerName in $NOT_UP
    do
        if [ $dockerName == "atomci" ]
        then
            docker-compose restart mysql
            echo "sleep 10 seconds, wait for mysql init ready"
            sleep 10
        fi
    echo "restart $dockerName..."
    docker-compose restart $dockerName;
    done
fi

DOCKER_NUM=$(docker-compose ps | grep 'Up' | wc -l)

if [ $DOCKER_NUM -ne 3 ]
then
    NOT_UP=$(docker-compose ps | grep -v "Name\|-----\|Up" | awk '{print $1}')
    echo "$NOT_UP restart failed, please check docker log use 'docker-compose logs [docker-name]'"
fi

# init result verify
[ $? -eq 0 ] && echo -e "AtomCI 初始化成功(:\n\n访问atomci: http://localhost:8090 \n" || echo -e "AtomCI 初始化失败, 请确认atomci 容器日志，\n或是 https://github.com/go-atomci/atomci-press/issues/new 反馈你的问题(:"
