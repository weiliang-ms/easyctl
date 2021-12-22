## 清理redis服务及文件

### 版本&兼容性

> 版本支持

- [v0.7.11-alpha以上](https://github.com/weiliang-ms/easyctl/releases/)

> 兼容性

- [x] `CentOS6`
- [x] `CentOS7`

### 使用方式

> 参考以下链接进行安装

- [安装说明文档](../-安装文档/README.md)

> 生成默认配置文件

```shell
$ easyctl clean redis
I1001 11:13:06.384839  126576 track.go:50] 检测到配置文件参数为空，生成配置文件样例 -> config.yaml
```

> 修改配置文件

`config.yaml`, 修改主机列表。`easyctl`根据主机列表`ssh`远程至目标主机执行清理`job`

```yaml
server:
  - host: 10.10.10.[1:40]
    username: root
    privateKeyPath: "" # ~/.ssh/id_rsa，为空默认走password登录；不为空默认走密钥登录
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
```

> 执行

```shell
$ easyctl clean redis -c config.yaml --debug
```