## redis集群安装

前置条件：

1. 配置好`yum`源，目前支持三种集群规模
2. 下载[redis](https://download.redis.io/releases/redis-5.0.14.tar.gz) 安装包
3. 安装[easyctl](../-安装文档/README.md)

支持的集群类型：

- [x] 单机伪集群
- [x] 三节点集群（三主三从，每个节点共运行两个服务）
- [x] 六节点集群（三主三从，每个节点共运行一个服务）

版本支持：[v0.7.9-alpha以上](https://github.com/weiliang-ms/easyctl/releases/tag/v0.7.9-alpha)

### 单机伪集群

> 1.生成配置文件

```shell
$ easyctl install redis-cluster
INFO[0000] 生成配置文件样例, 请携带 -c 参数重新执行 -> config.yaml
```

> 2.调整配置

`vi config.yaml`，调整以下参数

- `server`主机信息（用于安装redis集群主机）
- `redis-cluster:` 集群配置
    - `password`: `redis`密码
    - `cluster-type`: 部署集群类型
    - `package`: `redis`安装包路径

```yaml
server:
  - host: 10.10.10.1
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
redis-cluster:
  password: ""
  cluster-type: 1 # [0] 本地伪集群 ; [1] 三节点3分片2副本 ; [2] 6节点3分片2副本
  package: "redis-5.0.14.tar.gz
```

> 3.安装

```shell
$ easyctl install redis-cluster -c config.yaml --debug
```

结果输出如下：

```
1.节点列表: 10.10.10.1:26379,10.10.10.1:26380,10.10.10.1:26381,10.10.10.1:26382,10.10.10.1:26383,10.10.10.1:26384
2.密码:
3.日志目录: /var/log/redis
4.数据目录: /var/data/redis
5.启动命令/节点:
service redis-26379 start
service redis-26380 start
service redis-26381 start
service redis-26382 start
service redis-26383 start
service redis-26384 start

```

### 三节点集群

> 1.生成配置文件

```shell
$ easyctl install redis-cluster
INFO[0000] 生成配置文件样例, 请携带 -c 参数重新执行 -> config.yaml
```

> 2.调整配置

`vi config.yaml`，调整以下参数

- `server`主机信息（用于安装redis集群主机）
- `redis-cluster:` 集群配置
    - `password`: `redis`密码
    - `cluster-type`: 部署集群类型
    - `package`: `redis`安装包路径

```yaml
server:
  - host: 10.10.10.[1:3]
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
redis-cluster:
  password: ""
  cluster-type: 1 # [0] 本地伪集群 ; [1] 三节点3分片2副本 ; [2] 6节点3分片2副本
  package: "redis-5.0.14.tar.gz
```

> 3.安装

```shell
$ easyctl install redis-cluster -c config.yaml --debug
```

结果输出如下：

```
1.节点列表: 10.10.10.1:26379,10.10.10.2:26379,10.10.10.3:26379,10.10.10.1:26380,10.10.10.2:26380,10.10.10.3:26380
2.密码: redis@ddd
3.日志目录: /var/log/redis
4.数据目录: /var/data/redis
5.启动命令/节点: 
service redis-26379 start
service redis-26380 start
```

### 三节点集群

> 1.生成配置文件

```shell
$ easyctl install redis-cluster
INFO[0000] 生成配置文件样例, 请携带 -c 参数重新执行 -> config.yaml
```

> 2.调整配置

`vi config.yaml`，调整以下参数

- `server`主机信息（用于安装redis集群主机）
- `redis-cluster:` 集群配置
    - `password`: `redis`密码
    - `cluster-type`: 部署集群类型
    - `package`: `redis`安装包路径

```yaml
server:
  - host: 10.10.10.[1:6]
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
redis-cluster:
  password: "redis@ddd"
  cluster-type: 2 # [0] 本地伪集群 ; [1] 三节点3分片2副本 ; [2] 6节点3分片2副本
  package: "redis-5.0.14.tar.gz
```

> 3.安装

```shell
$ easyctl install redis-cluster -c config.yaml --debug
```

结果输出如下：

```
1.节点列表: 10.10.10.1:26379,10.10.10.2:26379,10.10.10.3:26379,10.10.10.4:26379,10.10.10.5:26379,10.10.10.6:26379
2.密码: redis@ddd
3.日志目录: /var/log/redis
4.数据目录: /var/data/redis
5.启动命令/节点: 
service redis-26379 start
```