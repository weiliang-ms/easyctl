### 配置主机禁用selinux

> 生成默认配置文件

```shell
easyctl deny selinux
```

> 修改配置文件

`config.yaml`

```yaml
server:
  - host: 10.10.1.[1:3]
    username: root
    password: 111111
    port: 22
excludes:
  - 192.168.235.132
```

> 配置

`--debug`输出`debug`日志，可选参数

```shell
easyctl deny selinux -c config.yaml --debug
```