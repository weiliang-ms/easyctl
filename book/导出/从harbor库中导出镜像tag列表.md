## harbor镜像tag列表导出

从`harbor`中批量导出镜像`tag`列表

> 背景说明

一些场景需要获取镜像`tag`列表（比如：批量导出镜像时）

> 生成默认配置文件

```shell
[root@localhost ~]# easyctl export harbor-image-list
I0928 21:19:46.803428   10628 export.go:41] 检测到配置文件参数为空，生成配置文件样例 -> config.yaml
```

> 修改配置文件

`vi config.yaml`

```yaml
harbor-repo:
  schema: http                      # 不可修改（暂不支持https harbor）
  address: 192.168.1.1:80           # harbor连接地址
  domain: harbor.wl.io              # harbor域
  user: admin                       # harbor用户
  password: Harbor-12345            # harbor用户密码
  preserve-dir: harbor-image-list   # 不建议修改，持久化tag
  projects:                         # 导出哪些项目下的镜像tag（如果为空表示全库导出）
    - apache                        # project名称
    - weaveworks
  excludes:                         # 配置'projects'空值使用，过滤某些project
    - ddd
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

### 单项目导出

导出单项目镜像列表

> 配置信息（部分内容已脱敏）

```yaml
harbor:
  schema: http
  address: *.*.*.*
  domain: harbor.wl.io
  user: admin
  password: ******
  preserve-dir: harbor-image-list
  projects:
    - apache
```

> 导出语句，`--debug`输出`debug`日志，可选参数

```shell
easyctl export harbor-image-list -c config.yaml --debug
```

> 目录文件结构如下：

```shell
/work                               # 执行命令的目录
├── config.yaml                     # 配置文件
└── harbor-image-list               # 存放镜像列表文件的目录（内部按项目建立文件夹进行隔离）
    ├── apache                      # apache项目目录，存放apache下镜像列表文件
    │   └── image-list.txt    # apache下镜像列表文件
    └── images-list.txt             # 导出项目下的所有镜像列表（x/image-list.txt xx/image-list.txt等内容合集）
```

- `images/images-list.txt`内容
```shell
[root@localhost work]# cat images/images-list.txt
harbor.wl.io/apache/skywalking-java-agent:8.6.0-alpine
harbor.wl.io/apache/skywalking-oap-server:8.6.0-es7
harbor.wl.io/apache/skywalking-ui:8.6.0
```

- `images/apache/image-list.txt`内容
```shell
[root@localhost work]# cat images/apache/image-list.txt
harbor.wl.io/apache/skywalking-java-agent:8.6.0-alpine
harbor.wl.io/apache/skywalking-oap-server:8.6.0-es7
harbor.wl.io/apache/skywalking-ui:8.6.0
```

### 多项目导出

导出一个以上项目下镜像列表

> 配置信息（部分内容已脱敏）

```yaml
harbor:
  schema: http
  address: *.*.*.*
  domain: harbor.wl.io
  user: admin
  password: ******
  preserve-dir: harbor-image-list
  projects:
    - apache
    - b2i
```

> 导出语句，`--debug`输出`debug`日志，可选参数

```shell
easyctl export harbor-image-list -c config.yaml --debug
```

目录文件结构如下：

```shell
/work/                              # 执行命令的目录
├── config.yaml                     # 配置文件
└── harbor-image-list               # 存放镜像列表文件的目录（内部按项目建立文件夹进行隔离）
    ├── apache                      # apache项目目录，存放apache下镜像列表文件
    │   └── image-list.txt    # apache下镜像列表文件
    ├── b2i                         # b2i项目目录，存放b2i下镜像列表文件
    │   └── image-list.txt    # b2下镜像列表文件
    └── images-list.txt             # 导出项目下的所有镜像列表（x/image-list.txt xx/image-list.txt等内容合集）
```

- `images/images-list.txt`内容
```shell
[root@localhost work]# cat images/images-list.txt
harbor.wl.io/apache/skywalking-java-agent:8.6.0-alpine
harbor.wl.io/apache/skywalking-oap-server:8.6.0-es7
harbor.wl.io/apache/skywalking-ui:8.6.0
harbor.wl.io/b2i/binary-nginx-builder:latest
harbor.wl.io/b2i/nginx-centos7-s2ibuilder:latest
harbor.wl.io/b2i/java-8-runtime:base-alpha
harbor.wl.io/b2i/java-8-runtime:base
harbor.wl.io/b2i/java-8-runtime:advance
harbor.wl.io/b2i/java-8-centos7:base
harbor.wl.io/b2i/java-8-centos7:advance
harbor.wl.io/b2i/tomcat9-java8-runtime:latest
harbor.wl.io/b2i/tomcat8-java8-runtime:latest
harbor.wl.io/b2i/tomcat8-java8-centos7:latest
harbor.wl.io/b2i/tomcat9-java8-centos7:latest
```

- `images/apache/image-list.txt`内容
```shell
[root@localhost work]# cat images/apache/image-list.txt
harbor.wl.io/apache/skywalking-java-agent:8.6.0-alpine
harbor.wl.io/apache/skywalking-oap-server:8.6.0-es7
harbor.wl.io/apache/skywalking-ui:8.6.0
```

- `images/b2i/image-list.txt`内容
```shell
[root@localhost work]# cat images/apache/image-list.txt
harbor.wl.io/b2i/binary-nginx-builder:latest
harbor.wl.io/b2i/nginx-centos7-s2ibuilder:latest
harbor.wl.io/b2i/java-8-runtime:base-alpha
harbor.wl.io/b2i/java-8-runtime:base
harbor.wl.io/b2i/java-8-runtime:advance
harbor.wl.io/b2i/java-8-centos7:base
harbor.wl.io/b2i/java-8-centos7:advance
harbor.wl.io/b2i/tomcat9-java8-runtime:latest
harbor.wl.io/b2i/tomcat8-java8-runtime:latest
harbor.wl.io/b2i/tomcat8-java8-centos7:latest
harbor.wl.io/b2i/tomcat9-java8-centos7:latest
```

### 全项目导出

导出全部项目下镜像列表

> 配置信息（部分内容已脱敏）

```yaml
harbor:
  schema: http
  address: *.*.*.*
  domain: harbor.wl.io
  user: admin
  password: ******
  export-all: true
  preserve-dir: harbor-image-list
  projects:
```

导出语句

```shell
easyctl export harbor-image-list -c config.yaml --debug
```

目录文件结构如下：

```shell
/work/                              # 执行命令的目录
├── config.yaml                     # 配置文件
└── harbor-image-list               # 存放镜像列表文件的目录（内部按项目建立文件夹进行隔离）
    ├── apache                      # apache项目目录，存放apache下镜像列表文件
    │   └── image-list.txt    # apache下镜像列表文件
    ├── b2i                         # b2i项目目录，存放b2i下镜像列表文件
    │   └── image-list.txt    # b2下镜像列表文件
    └── images-list.txt             # 导出项目下的所有镜像列表（x/image-list.txt xx/image-list.txt等内容合集）
    ├── ceph-csi
    │   └── image-list.txt
    ├── champ
    │   └── image-list.txt
    ├── charts
    │   └── image-list.txt
    ├── csiplugin
    │   └── image-list.txt
    ├── elastic
    │   └── image-list.txt
    ├── elasticsearch
    │   └── image-list.txt
    ├── grafana
    │   └── image-list.txt
    ├── hsa-cep
    │   └── image-list.txt
    ├── hsa-k8s-public
    │   └── image-list.txt
    ├── images-list.txt
    ├── istio
    │   └── image-list.txt
    ├── jaegertracing
    │   └── image-list.txt
    ├── jenkins
    │   └── image-list.txt
    ├── jimmidyson
    │   └── image-list.txt
    ├── kubernetes
    │   └── image-list.txt
    ├── kubesphere
    │   └── image-list.txt
    ├── library
    │   └── image-list.txt
    ├── minio
    │   └── image-list.txt
    ├── openebs
    │   └── image-list.txt
    ├── openpitrix
    │   └── image-list.txt
    ├── osixia
    │   └── image-list.txt
    ├── paas
    │   └── image-list.txt
    ├── prom
    │   └── image-list.txt
    └── weaveworks
        └── image-list.txt
```