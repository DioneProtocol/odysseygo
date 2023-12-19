#!/bin/bash

go_ver=go1.21.4.linux-amd64
coreth_repo=https://git.sfxdx.com/green-energy1/go-ethereum

#check root user
if [[ $EUID > 0 ]]
  then echo -e "\nRun this script as root! Exiting...\n"
  exit
fi

#welcome
clear
echo -e "\nWelcome to the node-installation script!\n"
echo "               What are you planning to do?"
echo "--------------------------------------------------------"
echo "1: Prepare server, compile, configure and run bootnodes;"
echo "2: Prepare server, compile, configure and run node;"
echo "--------------------------------------------------------"
echo ""
read -p "Choice: " n < /dev/tty
if [ -z "$n" ] || [[ $n != 1 ]] && [[ $n != 2 ]] ; then echo -e '\nSelect an existing option! Exiting...\n' && exit 1; fi
if [ $n == 2 ]
then
  read -p "Enter node name from run_config file 'node' section: " node_name < /dev/tty
  if [ -z "$node_name" ]; then echo -e '\nNode name is blank! Exiting...\n' && exit 1; fi
fi 
echo ""

#pre-config
export DEBIAN_FRONTEND=noninteractive

#install packages
{ apt update && apt install -y curl wget gcc g++ git jq; } || { echo -e '\nComponents installation failed: Ubuntu servers are not available or dpkg is misconfigured or frozen. Exiting...\n' && exit 1; }

#install go
{ wget https://go.dev/dl/$go_ver.tar.gz && rm -rf /usr/local/go && tar -D /usr/local -xzf $go_ver.tar.gz && rm $go_ver.tar.gz && ln -sf /usr/local/go/bin/go /usr/local/bin/go; } || { echo -e '\nComponents installation failed: Go servers are not available. Exiting...\n' && exit 1; }

#clone and prepare repos
cd /usr/local/go/bin/ && mkdir -p src/github.com/DioneProtocol && cd src/github.com/DioneProtocol 
{ git clone $coreth_repo coreth && ln -sf /usr/local/go/bin/src/github.com/DioneProtocol/coreth /root/coreth; } || { git clone $coreth_repo coreth && ln -sf /usr/local/go/bin/src/github.com/DioneProtocol/coreth /root/coreth; } || { echo -e '\nRepo cloning failed: insufficient permissions. Exiting...\n' && exit 1; }

#execute script
if [ $n == 1 ]
then
  bash /root/odyssey-avax-fork/scripts/run_blockchain.sh
else 
  bash /root/odyssey-avax-fork/scripts/run_node.sh -n $node_name
fi

#echo fin
echo -e "\nScript's completed the selected stage!\n"