#!/bin/bash
if [ $# -eq 0 ];then
	echo "enter one of windows,linux,mac"
	exit
elif [ $# -eq 1 ];then
    platform=$1
	if [ $platform == "linux" ];then
		GOOS=linux GOARCH=amd64 go build
#		cd webroot
#		npm run build
#		cd ..
		rm -rf build/linux
		mkdir build/linux
		mv sipemanager build/linux/
#		mkdir -p build/linux/webroot
#		mv webroot/dist build/linux/webroot
		cp -r etc build/linux/
		cp -r docs build/linux/
		cp start.sh build/linux/

	elif [ $platform == "widows" ];then
		GOOS=widows GOARCH=amd64 go build
#		cd webroot
#		npm run build
#		cd ..
		rm -rf build/windows
		mkdir build/windows
		mv sipemanager build/windows/
#		mkdir -p build/windows/webroot
#		mv webroot/dist build/windows/webroot
		cp -r etc build/windows/
		cp -r docs build/windows/
		cp start.sh build/windows/

	elif [ $platform == "mac" ];then	
		GOOS=darwin GOARCH=amd64 go build
		cd webroot
	  npm run build
#   npm install
		cd ..
		rm -rf build/mac
		mkdir -p build/mac
		mv sipemanager build/mac/
		mkdir -p build/mac/webroot
		mv webroot/dist build/mac/webroot
		cp -r etc build/mac/
		cp -r docs build/mac/
		cp start.sh build/mac/
	fi
fi



