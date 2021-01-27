#!/bin/bash

chmod 755 ipfs
cp ipfs /usr/local/bin

ipfspath="$HOME/.ipfs/"

if [ ! -d "$ipfspath" ];then
	mkdir $ipfspath
	echo "创建文件夹 $ipfspath 成功"
fi

rm -rf $ipfspsth/swarm.key

if [ ! -d "swarm.key" ];then
        cp swarm.key $HOME/.ipfs/
        echo "密钥完成"
fi

ipfs init
ipfs bootstrap rm --all
ipfs config --json API.HTTPHeaders.Access-ContrZZol-Allow-Origin '["*"]'
ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'
ipfs config --json API.HTTPHeaders.Access-Control-Allow-Credentials '["true"]'
ipfs config --json Swarm.EnableAutoRelay 'true'

echo "===============删除默认中继节点=============="
for line in `cat bootstrap.txt`
do
 echo "加入可信节点:$line"
 ipfs bootstrap add $line
done
