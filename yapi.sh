#!/bin/bash
which yapi-cli
if [ "$?" -ne 0 ]; then
    npm install -g yapi-cli
fi
#提交到服务器
yapi import