### 修改主机root口令

> 生成默认配置文件

```shell
easyctl set new-password
```

> 修改配置文件

`config.yaml`

- 调整主机信息，新`root`口令的值

```yaml
server:
  - host: 10.10.10.[1:40]
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
newRootPassword: "3LEPnok84HxYc5"
```

> 运行

`--debug`输出`debug`日志，可选参数

```shell
easyctl set new-password -c config.yaml --debug
```

> 成功样例

```
[root@localhost ~]# ./easyctl set new-password -c config.yaml
I0926 15:14:56.431946  112411 log.go:184] 检测到配置文件中含有IP段，开始解析组装...
I0926 15:14:56.431999  112411 parse.go:113] 解析到IP子网网段为：10.10.1....
I0926 15:14:56.432021  112411 parse.go:117] 解析到IP区间为：1:2...
I0926 15:14:56.432026  112411 parse.go:121] 解析到起始IP为：10.10.1.1...
I0926 15:14:56.432031  112411 parse.go:125] 解析到末尾IP为：10.10.1.2...
I0926 15:14:56.432037  112411 exec.go:43] 开始并行执行命令...
I0926 15:14:56.432084  112411 exec.go:105] [10.10.1.2] 开始执行指令 ->
I0926 15:14:56.432114  112411 exec.go:105] [10.10.1.1] 开始执行指令 ->
I0926 15:14:56.634224  112411 log.go:184] <- 10.10.1.1执行命令成功...
I0926 15:14:56.634472  112411 log.go:184] <- 10.10.1.2执行命令成功...
| IP ADDRESS  |  CMD   | EXIT CODE | RESULT  |        OUTPUT        | EXCEPTION |
|-------------|--------|-----------|---------|----------------------|-----------|
| 10.10.1.1 | ****** |     0     | success | Changing password fo |           |
| 10.10.1.2 | ****** |     0     | success | Changing password fo |           |
```

> 测试

重新连接列表主机