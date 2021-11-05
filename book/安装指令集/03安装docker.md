## docker安装

二进制方式安装配置`docker`

前置条件：

1. 下载[docker-ce](https://download.docker.com/linux/static/stable/x86_64/) 安装包
2. 安装[easyctl](../-安装文档/README.md)
3. `CentOS7`系统

版本支持：[v0.7.15-alpha以上](https://github.com/weiliang-ms/easyctl/releases/tag/v0.7.15-alpha)

## 本地安装

> 1.生成配置文件

```shell
$ easyctl install docker-ce
INFO[0000] 生成配置文件样例, 请携带 -c 参数重新执行 -> config.yaml
```

> 2.调整配置

`vi config.yaml`，调整以下参数

- `server`主机信息（用于安装`docker`主机），如果为空表示本地安装
- `excludes`排除`server.host`声明地址段内的主机
- `docker:` `docker`配置
    - `package`: `docker`二进制包
    - `preserveDir`: `docker`持久化目录（默认`/var/lib/docker`）
    - `insecureRegistries`: 非`https`仓库列表
    - `registryMirrors`: # 镜像源列表

```yaml
server:
  - host: 10.10.10.1-3
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
docker:
  package: docker-19.03.15.tgz   # 二进制安装包目录
  preserveDir: /data/lib/docker  # docker数据持久化目录
  insecureRegistries: # 非https仓库列表
    - gcr.azk8s.cn
    - quay.azk8s.cn
  registryMirrors:               # 镜像源
```

`config.yaml`修改后样例:

```yaml
docker:
  package: docker-19.03.15.tgz   # 二进制安装包目录
  preserveDir: /data/lib/docker  # docker数据持久化目录
  insecureRegistries: # 非https仓库列表
    - gcr.azk8s.cn
    - quay.azk8s.cn
  registryMirrors:               # 镜像源
```

> 3.安装样例

可以去掉`--debug`参数减少内容输出

```shell
$ easyctl install docker-ce -c config.yaml --debug
```

结果输出如下：

```
...
[easyctl] localhost.localdomain | 2021-11-04T23:20:16-04:00 | info | 启动docker
[easyctl] localhost.localdomain | 2021-11-04T23:20:16-04:00 | debug | 执行指令:
setenforce 0
groupadd docker
useradd docker -g docker
systemctl daemon-reload
systemctl restart docker

[easyctl] localhost.localdomain | 2021-11-04T23:20:17-04:00 | debug |
[easyctl] localhost.localdomain | 2021-11-04T23:20:17-04:00 | info | docker安装完
```

## 远程安装

通过`ssh`方式远程安装，需要指定`servers`列表

> 1.生成配置文件

```shell
$ easyctl install docker-ce
INFO[0000] 生成配置文件样例, 请携带 -c 参数重新执行 -> config.yaml
```

> 2.调整配置

`vi config.yaml`，调整以下参数

- `server`主机信息（用于安装`docker`主机），如果为空表示本地安装
- `excludes`排除`server.host`声明地址段内的主机
- `docker:` `docker`配置
    - `package`: `docker`二进制包
    - `preserveDir`: `docker`持久化目录（默认`/var/lib/docker`）
    - `insecureRegistries`: 非`https`仓库列表
    - `registryMirrors`: # 镜像源列表

```yaml
server:
  - host: 10.10.10.1-3
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
docker:
  package: docker-19.03.15.tgz   # 二进制安装包目录
  preserveDir: /data/lib/docker  # docker数据持久化目录
  insecureRegistries: # 非https仓库列表
    - gcr.azk8s.cn
    - quay.azk8s.cn
  registryMirrors:               # 镜像源
```

`config.yaml`修改后样例:

```yaml
server:
  - host: 192.168.109.143
    username: root
    password: 1
    port: 22
docker:
  package: docker-19.03.15.tgz   # 二进制安装包目录
  preserveDir: /data/lib/docker  # docker数据持久化目录
  insecureRegistries: # 非https仓库列表
    - gcr.azk8s.cn
    - quay.azk8s.cn
  registryMirrors:               # 镜像源
```

> 3.安装样例

可以去掉`--debug`参数减少内容输出

```shell
$ easyctl install docker-ce -c config.yaml --debug
```

结果输出如下：

```
...
WantedBy=sockets.target
[easyctl] localhost.localdomain | 2021-11-04T23:22:18-04:00 | info | 启动docker
[easyctl] localhost.localdomain | 2021-11-04T23:22:18-04:00 | info | 开始并行执行命令...
[easyctl] localhost.localdomain | 2021-11-04T23:22:18-04:00 | info | [192.168.109.143] 开始执行指令 ->
setenforce 0
groupadd docker
useradd docker -g docker
systemctl daemon-reload
systemctl restart docker

[easyctl] localhost.localdomain | 2021-11-04T23:22:20-04:00 | info | docker安装完毕```
```