#!/bin/bash

# Install Latest ATOMCI Release

os=`uname -a`
server_url=""
git_urls=('github.com' 'hub.fastgit.org')

echo '正在检测可用服务器'
for git_url in ${git_urls[*]}
do
	success="true"
	for i in {1..3}
	do
		echo -ne "检测 ${git_url} ... ${i} "
	    curl -m 5 -kIs https://${git_url} >/dev/null
		if [ $? != 0 ];then
			echo "failed"
			success="false"
			break
		else
			echo "ok"
            sleep 1s
		fi
	done
	if [ ${success} == "true" ];then
		server_url=${git_url}
		break
	fi
done

if [ "x${server_url}" == "x" ];then
    echo "没有找到稳定的下载服务器，请稍候重试"
    exit 1
fi


echo "使用下载服务器 ${server_url}"

# 支持MacOS
if [[ $os =~ 'Darwin' ]];then
    ATOMCIVERSION=$(curl -s https://${server_url}/go-atomci/atomci/releases/latest |grep -Eo '[0-9]+.[0-9]+.[0-9]+')
else
	ATOMCIVERSION=$(curl -s https://${server_url}/go-atomci/atomci/releases/latest/download 2>&1 | grep -Po '[0-9]+\.[0-9]+\.[0-9]+.*(?=")')
fi

echo '当前最新版本为：${ATOMCIVERSION} 开始下载...'

wget --no-check-certificate https://${server_url}/go-atomci/atomci/releases/download/v${ATOMCIVERSION}/atomci-${ATOMCIVERSION}-docker-compose.tgz

tar zxf atomci-${ATOMCIVERSION}-docker-compose.tgz
cd atomci-${ATOMCIVERSION}
docker-compose down && docker-compose up -d
/bin/bash init.sh