# easyctl 

`easyctl`是一款基于`golang`轻量级运维工具集，详情使用请移步[使用文档](https://weiliang-ms.github.io/easyctl/)

** 适用平台：** `CentOS7`

## 安装使用

### 编译安装最新版

```shell
git clone https://github.com/weiliang-ms/easyctl.git
cd easyctl
go build -ldflags "-w -s" -o /usr/local/bin/easyctl
```

### 下载release版本

[release](https://github.com/weiliang-ms/easyctl/releases)

## 迭代计划

> 里程碑

- `v0.x.y-alpha`: 添加常用指令集功能（`x`为一级指令集，如`set`;`y`为二级指令集，如`set`指令集中的`dns`子指令集）
- `v1.0.0-beta`: `bug`修复、文档站点、集成`github workflow`
- `v1.0.0-release`: 正式版本

> `v1`已实现功能

- `deny`
  - 防火墙
  - `ping`
  - `selinux`
- `set`
  - 时区
  - `dns`
  - 主机互信
  - 主机`host`解析
  - 文件描述符数
  - 修改`root`口令

> `v2`功能预览

- `windows GUI`