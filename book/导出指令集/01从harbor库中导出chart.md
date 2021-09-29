### chart导出

从`harbor`中批量下载`chart`文件

> 背景说明

通常`chart`应用存放于`harbor`进行管理，但`harbor`（v2.1.4）只支持单个`chart`文件下载。

对于离线批量分发场景，比较难受，特开发该功能。

> 生成默认配置文件

```shell
[root@localhost ~]# easyctl export chart
I0927 09:30:47.685246   14817 export.go:44] 检测到配置文件参数为空，生成配置文件样例 -> config.yaml
```

> 修改配置文件

`vi config.yaml`

```yaml
helm-repo:
  endpoint: 10.10.1.3:80   # harbor访问地址
  domain: harbor.wl.io      # harbor域
  username: admin           # harbor用户
  password: 123456          # harbor密码
  preserveDir: /root/charts # chart包持久化目录
  package: true             # 是否打成tar包
  repo-name: charts         # chart repo harbor内的名称
```

> 配置

`--debug`输出`debug`日志，可选参数

```shell
[root@node1 ~]# easyctl export chart -c config.yaml
INFO[0000] 解析chart仓库配置...
INFO[0000] 待导出chart数量为: 135
INFO[0000] 导出chart...
INFO[0000] 创建目录: /root/charts
INFO[0000] 逐一导出chart中...
INFO[0002] 导出完毕，chart总数为:135
```