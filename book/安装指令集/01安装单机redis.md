## redis安装

前置条件：

1. 配置好`yum`源
2. 下载[redis](https://download.redis.io/releases/redis-5.0.14.tar.gz) 安装包
3. 安装[easyctl](../-安装文档/README.md)

版本支持：[v0.7.10-alpha以上](https://github.com/weiliang-ms/easyctl/releases/tag/v0.7.10-alpha)

### 安装

> 1.生成配置文件

```shell
$ easyctl install redis
INFO[0000] 生成配置文件样例, 请携带 -c 参数重新执行 -> config.yaml
```

> 2.调整配置

`vi config.yaml`，调整以下参数

- `server`主机信息（用于安装`redis`主机）
- `redis:` `redis`配置
    - `password`: `redis`密码
    - `port`: `redis`监听端口
    - `package`: `redis`安装包路径

```yaml
server:
  - host: 10.10.10.1
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
redis:
  password: "dddd"
  port: 26379
  package: "redis-5.0.14.tar.gz"
```

> 3.安装

```shell
$ easyctl install redis -c config.yaml --debug
```

结果输出如下：

```
[easyctl] DESKTOP-O8QG6I5 | 2021-10-10T14:36:01+08:00 | info | redis安装完毕,相关信息如下：
1.节点列表: 10.10.10.1:26379
2.密码: dddd
3.日志目录: /var/log/redis
4.数据目录: /var/data/redis
5.启动命令/节点: 
service redis-33333 start
6.二进制目录：/usr/local/bin/redis-*--- PASS: TestInstallRedis (31.00s)
```