#!/bin/bash
if [ $# -eq 0 ];then
	echo "enter one of windows,linux,mac"
	exit
elif [ $# -eq 1 ];then
    platform=$1
	if [ $platform == "linux" ];then
		GOOS=linux GOARCH=amd64 go build
		cd webroot
		npm run build
		cd ..
		rm -rf build/linux
		mkdir build/linux
		mv sipe-client build/linux/
		mkdir -p build/linux/webroot
		mv webroot/dist build/linux/webroot
		cp -r etc build/linux/

	elif [ $platform == "widows" ];then
		GOOS=widows GOARCH=amd64 go build
		cd webroot
		npm run build
		cd ..
		rm -rf build/windows
		mkdir build/windows
		mv sipe-client build/windows/
		mkdir -p build/windows/webroot
		mv webroot/dist build/windows/webroot
		cp -r etc build/windows/

	elif [ $platform == "mac" ];then	
		GOOS=darwin GOARCH=amd64 go build
		cd webroot
		npm run build
		cd ..
		rm -rf build/mac
		mkdir build/mac
		mv sipe-client build/mac/
		mkdir -p build/mac/webroot
		mv webroot/dist build/mac/webroot
		cp -r etc build/mac/
	fi
fi



