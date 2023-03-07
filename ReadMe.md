# WorkStation Client

**[English](ReadMe_en.md)**

WorkStation 是一个服务器、路由器等网络设备的集中管理工具，包含了一些基本信息集中式管理，配置信息的远程查看和修改。
适合家庭用户和小型企业用户，ddns集中管理功能方便更好的管理您在不通区域的设备，加入ipv6的支持。支持设备包括树莓派、Openwrt、Windows、Linux、FreeBSD甚至是常用的NAS。

## 目录

## 编译

**编译Linux版本**

```shell
# x86_64 at windows
C:\wsc> SET CGO_ENABLED=0
C:\wsc> SET GOOS=linux
C:\wsc> SET GOARCH=amd64
C:\wsc> go build main.go
# at linux
~/wsc$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

# arm64 at windows
C:\wsc> SET CGO_ENABLED=0
C:\wsc> SET GOOS=linux
C:\wsc> SET GOARCH=arm64
# at linux
~/wsc$ CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build main.go
```

**编译MacOS版本**

```shell
# x86_64 at windows
C:\> SET CGO_ENABLED=0
C:\> SET GOOS=darwin
C:\> SET GOARCH=amd64
# at linux
~/wsc$ CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go

# arm64 at windows
C:\wsc> SET CGO_ENABLED=0
C:\wsc> SET GOOS=darwin
C:\wsc> SET GOARCH=arm64
# at linux
~/wsc$ CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build main.go
```

**编译FreeBSD版本**

```shell
# x86_64 at windows
C:\> SET CGO_ENABLED=0
C:\> SET GOOS=freebsd
C:\> SET GOARCH=amd64
# at linux
~/wsc$ CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build main.go

# arm64 at windows
C:\wsc> SET CGO_ENABLED=0
C:\wsc> SET GOOS=freebsd
C:\wsc> SET GOARCH=arm64
# at linux
~/wsc$ CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build main.go

```

## 开源声明

WorkStation 基于 GPL V3 协议开源。