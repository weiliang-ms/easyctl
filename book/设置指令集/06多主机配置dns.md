### 多主机配置dns

> 生成默认配置文件

```shell
easyctl set dns
```

> 修改配置文件

`config.yaml`

- 调整主机信息
- 调整`dns`地址列表

```yaml
server:
  - host: 10.10.10.[1:40]
    username: root
    privateKeyPath: "" # ~/.ssh/id_rsa，为空默认走password登录；不为空默认走密钥登录
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
dns:
  - 114.114.114.114
  - 8.8.8.8
```

> 运行

`--debug`输出`debug`日志，可选参数

```shell
easyctl set dns -c config.yaml --debug
```

> 测试

任意主机列表内的主机执行：

```shell
cat /etc/hosts
```