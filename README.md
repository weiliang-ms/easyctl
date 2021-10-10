![GitHub Workflow Status (event)](https://img.shields.io/github/workflow/status/weiliang-ms/easyctl/Go?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/weiliang-ms/easyctl)](https://goreportcard.com/report/github.com/weiliang-ms/easyctl)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/weiliang-ms/easyctl?filename=go.mod&style=flat-square)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/weiliang-ms/easyctl?style=flat-square)
![GitHub all releases](https://img.shields.io/github/downloads/weiliang-ms/easyctl/total?style=flat-square)
[![codecov](https://codecov.io/gh/weiliang-ms/easyctl/branch/master/graph/badge.svg?token=7RGD5V5L9Y)](https://codecov.io/gh/weiliang-ms/easyctl)
![GitHub](https://img.shields.io/github/license/weiliang-ms/easyctl?style=flat-square)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fweiliang-ms%2Feasyctl.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fweiliang-ms%2Feasyctl?ref=badge_shield)
# easyctl 

`easyctl`是一款基于`golang`轻量级运维工具集，详情使用请移步[使用文档](https://weiliang-ms.github.io/easyctl/)

兼容性：

- [x] 全部功能兼容`CentOS7`
- [ ] 部分功能兼容`CentOS6`

## 安装使用

### 编译安装最新版

```shell
git clone https://github.com/weiliang-ms/easyctl.git
cd easyctl
go build -ldflags "-w -s" -o /usr/local/bin/easyctl
```

### 下载release版本

> [latest release](https://github.com/weiliang-ms/easyctl/releases/tag/latest)

- [Mac OS](https://github.com/weiliang-ms/easyctl/releases/download/latest/easyctl-latest-darwin-amd64.tar.gz)
```shell
sudo tar zxvf easyctl-latest-darwin-amd64.tar.gz
sudo cp easyctl /usr/local/bin
```

- [linux-amd64](https://github.com/weiliang-ms/easyctl/releases/download/latest/easyctl-latest-linux-amd64.tar.gz)
```shell
sudo tar zxvf easyctl-latest-linux-amd64.tar.gz
sudo cp easyctl /usr/local/bin
```

- [Windows](https://github.com/weiliang-ms/easyctl/releases/download/latest/easyctl-latest-windows-amd64.zip)

> 下载[tag版本](https://github.com/weiliang-ms/easyctl/tags)

## 迭代计划

> 里程碑

- [x] [文档站点](https://weiliang-ms.github.io/easyctl/)
- [x] 集成`github workflow`
- `v0.x.y-alpha`: 开发
- `v1.0.0-beta`: `bug`修复
- `v1.0.0-release`: 正式版本

> `v1`功能列表

- `add`
- `clean`
  - [x] [清理多机redis](https://weiliang-ms.github.io/easyctl/%E6%B8%85%E7%90%86%E6%8C%87%E4%BB%A4%E9%9B%86/01%E5%A4%9A%E4%B8%BB%E6%9C%BAredis%E6%B8%85%E7%90%86.html)
  - [x] [清理主机dns配置](https://weiliang-ms.github.io/easyctl/%E6%B8%85%E7%90%86%E6%8C%87%E4%BB%A4%E9%9B%86/01%E5%A4%9A%E4%B8%BB%E6%9C%BAdns%E9%85%8D%E7%BD%AE%E6%B8%85%E7%90%86.html)
- `deny`
  - [x] [防火墙](https://weiliang-ms.github.io/easyctl/%E7%A6%81%E7%94%A8%E6%8C%87%E4%BB%A4%E9%9B%86/03%E4%B8%BB%E6%9C%BA%E7%A6%81%E7%94%A8%E9%98%B2%E7%81%AB%E5%A2%99.html)
  - [x] [ping](https://weiliang-ms.github.io/easyctl/%E7%A6%81%E7%94%A8%E6%8C%87%E4%BB%A4%E9%9B%86/01%E4%B8%BB%E6%9C%BA%E7%A6%81Ping.html)
  - [x] [selinux](https://weiliang-ms.github.io/easyctl/%E7%A6%81%E7%94%A8%E6%8C%87%E4%BB%A4%E9%9B%86/02%E4%B8%BB%E6%9C%BA%E7%A6%81%E7%94%A8selinux.html)
- `export`
  - [x] [从harbor批量下载chart](https://weiliang-ms.github.io/easyctl/%E5%AF%BC%E5%87%BA%E6%8C%87%E4%BB%A4%E9%9B%86/01%E4%BB%8Eharbor%E5%BA%93%E4%B8%AD%E5%AF%BC%E5%87%BAchart.html)
  - [x] [从harbor导出镜像tag列表](https://weiliang-ms.github.io/easyctl/%E5%AF%BC%E5%87%BA%E6%8C%87%E4%BB%A4%E9%9B%86/02%E4%BB%8Eharbor%E5%BA%93%E4%B8%AD%E5%AF%BC%E5%87%BA%E9%95%9C%E5%83%8Ftag%E5%88%97%E8%A1%A8.html)
- `install`
  - [x] [单机redis](https://weiliang-ms.github.io/easyctl/%E5%AE%89%E8%A3%85%E6%8C%87%E4%BB%A4%E9%9B%86/01%E5%AE%89%E8%A3%85%E5%8D%95%E6%9C%BAredis.html)
  - [x] [redis集群](https://weiliang-ms.github.io/easyctl/%E5%AE%89%E8%A3%85%E6%8C%87%E4%BB%A4%E9%9B%86/02%E5%AE%89%E8%A3%85redis%E9%9B%86%E7%BE%A4.html)
- `set`
  - [x] [时区](https://weiliang-ms.github.io/easyctl/%E8%AE%BE%E7%BD%AE%E6%8C%87%E4%BB%A4%E9%9B%86/05%E5%A4%9A%E4%B8%BB%E6%9C%BA%E8%AE%BE%E7%BD%AE%E6%97%B6%E5%8C%BA.html)
  - [x] [dns](https://weiliang-ms.github.io/easyctl/%E8%AE%BE%E7%BD%AE%E6%8C%87%E4%BB%A4%E9%9B%86/06%E5%A4%9A%E4%B8%BB%E6%9C%BA%E9%85%8D%E7%BD%AEdns.html)
  - [x] [主机互信](https://weiliang-ms.github.io/easyctl/%E8%AE%BE%E7%BD%AE%E6%8C%87%E4%BB%A4%E9%9B%86/03%E5%A4%9A%E4%B8%BB%E6%9C%BA%E5%85%8D%E5%AF%86%E7%99%BB%E5%BD%95.html)
  - [x] [主机host解析](https://weiliang-ms.github.io/easyctl/%E8%AE%BE%E7%BD%AE%E6%8C%87%E4%BB%A4%E9%9B%86/01%E5%A4%9A%E4%B8%BB%E6%9C%BAhost%E8%A7%A3%E6%9E%90.html)
  - [x] [文件描述符数](https://weiliang-ms.github.io/easyctl/%E8%AE%BE%E7%BD%AE%E6%8C%87%E4%BB%A4%E9%9B%86/04%E5%A4%9A%E4%B8%BB%E6%9C%BA%E8%AE%BE%E7%BD%AE%E6%96%87%E4%BB%B6%E6%8F%8F%E8%BF%B0%E7%AC%A6.html)
  - [x] [修改root口令](https://weiliang-ms.github.io/easyctl/%E8%AE%BE%E7%BD%AE%E6%8C%87%E4%BB%A4%E9%9B%86/02%E5%A4%9A%E4%B8%BB%E6%9C%BA%E4%BF%AE%E6%94%B9root%E5%8F%A3%E4%BB%A4.html)
- `track`
  - [x] [日志tail](https://weiliang-ms.github.io/easyctl/%E8%BF%BD%E8%B8%AA%E6%8C%87%E4%BB%A4%E9%9B%86/01%E5%A4%9A%E4%B8%BB%E6%9C%BA%E6%97%A5%E5%BF%97%E5%AE%9E%E6%97%B6%E8%BF%BD%E8%B8%AA.html)

> `v2`功能预览

- `windows GUI`

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fweiliang-ms%2Feasyctl.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fweiliang-ms%2Feasyctl?ref=badge_large)