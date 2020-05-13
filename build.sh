#!/bin/bash
set -xe

# 设置https代理
#export https_proxy=http://10.13.40.145:8118

# 安装gitlab依赖库
#cat go.mod | grep "git.intra" | grep -v "module" | awk '{print $1"@"$2}' | xargs go get -insecure

arg="$1"
if [[ ! -f "$arg" ]]
then
    rm -f service_core
    go build
else
    case $arg in
    linux32)
        rm -f service_core_linux32
        env GOOS=linux GOARCH=386 go build -o service_core_linux32
        ;;
    linux64)
        rm -f service_core_linux64
        env GOOS=linux GOARCH=amd64 go build -o service_core_linux64
        ;;
    mac)
        rm -f service_core_mac64
        env GOOS=darwin GOARCH=amd64 go build -o service_core_mac64
        ;;
    local)
        rm -f service_core
        go build
        ;;
    *)
    echo "Usage: build.sh [ linux32 | linux64 | mac | local ]"
    exit 0
    ;;
    esac
fi
rm -rf target
mkdir -p target/bin
#cp -r docs target/
cp -r configs target/conf
cp service_core target/bin/
