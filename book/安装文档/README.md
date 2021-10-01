### 安装easyctl

1. 编译安装最新版

```shell
git clone https://github.com/weiliang-ms/easyctl.git
cd easyctl
go build -ldflags "-w -s" -o /usr/local/bin/easyctl
```

2. 下载编译好的文件(建议最新版本)

[easyctl release](https://github.com/weiliang-ms/easyctl/releases)

```
chmod +x easyctl
mv easyctl /usr/local/bin
```