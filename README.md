## harborctl使用帮助

>根据环境要求 该工具支持harbor 2.2.1 版本 其他api变化可能导致报错 请谨慎使用!!

### 下载地址:

```shell
wget -cP /usr/local/bin https://github.com/bcsimple/harborctl/releases/download/v1.0/harborctl && chmod 755 /usr/local/bin/harborctl
```

### 1.配置文件

```shell
#NAME 为该harbor起一个名字
#USER 指定该harbor的链接使用的用户
#PASSWD 指定该harbor的链接使用的密码
#IP:PORT 指定该harbor服务器地址和端口必须带上端口

#注意仓库可以添加多个 但是名字必须唯一 不然会报错
harborctl config set-context NAME -s IP:PORT -u USER -p PASSWD 
harborctl config set-context harbor01 -s 10.0.0.1:1121 -u admin p Harbor12345
#查看
harborctl config view 
```

### 2.使用帮助

```
1.查询镜像
harborctl search image  hello,kube  // 支持逗号分割模糊匹配多关键字

2.查询chart
harborctl search chart  hello,kube  

3.查询复制规则
harborctl search replication  hello

4.查询执行器
harborctl search execution hello -i PEPLICATION_ID

5.查询项目
harborctl search project  hello,kube

6.查询任务
harborctl search task   hello  -i EXECUTION_ID

7.查询仓库列表
harborctl search registry   hello  



下载:
1.下载镜像:
harborctl download image -n hello  -p  /tmp      // -p 指定下载的目录(不指定默认为当前目录下)  -n 模糊匹配

2.下载chart:
harborctl download chart -n hello  -p  /tmp      // -p 指定下载的目录(不指定默认为当前目录下)  -n 模糊匹配



推送:
1.批量推送
harborctl start -i PEPLICATION_ID -p  xxx.json   // -p 制定读取配置文件路径  PEPLICATION_ID 为复制规则ID

xxx.json

格式如下:  //严格按照格式走 不然直接报错 或者产生一些不好的事件!
[
  {
    "name": "kube_system/**",
    "tag": "",
    "resource": "image",
    "dst_namespace": "public"
  },
  {
    "name": "public/**",
    "tag": "",
    "resource": "chart",
    "dst_namespace": "kube_system"
  }
]
字段解释:
	name: 过滤规则 支持* 或者 ? 通配符
	tag:  指定制品的版本号 支持* 或者 ? 通配符
	resource: 指定推送制品的类型 常用的就是image和chart
	dst_namespace: 二级仓库的项目名 即想把镜像推送到二级的那个项目下 不写默认为name相同的项目下


2.单个replication推送
harborctl update -i REPLICATION_ID   //修改单个规则  最后一个交互会出现y/n来操作是否直接推送 也可以使用下面的start进行手动复制

harborctl start -i REPLICATION_ID    //执行规则 开始复制操作


创建
1.创建仓库
harborctl create registry 
harborctl create registry -n NAME  使用配置文件中指定的harbor的名字或者别名

2.创建规则 默认是push策略
harborctl create replication -i REGISTRY_ID   //创建规则指定仓库id 即REGISTRY_ID
harborctl create replication -i REGISTRY_ID -p  //指定-p 则会创建pull类型的策略

3.创建项目
harborctl create project NAME 

删除:
1.删除项目
harborctl delete project NAME

2.删除仓库或者策略规则 这里的ID必须事先查出来 harborctl search registry/replication NAME 
harborctl delete replication  -i ID // harborctl delete registry -i ID

```

### 3.scan说明

```shell
1.指定扫描文件 并从文件中列出
harborctl scan -s xxx.yaml  -F   

2.指定扫描文件 并从harbor对比
harborctl scan -s xxx.yaml  -C -c harbor01 -f kube //可以指定-c用于对比的harbor

3.指定扫描文件 并从harbor差异
harborctl scan -s xxx.yaml  -d -c harbor01 -f kube //可以指定-c用于对比的harbor

随着版本的增加需要指定-r参数 默认当前为paas_v20230101

对比完成后修复 将harbor中的拉取下来修改成脚本所需要的镜像文件  下面用到的IP:PORT必须修改掉再执行
拉取
harborctl scan -s  ./config/image.yaml  -d  -f kube | tail -n +2 | awk '{printf ("docker pull IP:PORT/%s:%s\n",$4,$5)}' | sh
修改并推送:
harborctl scan -s  ./config/image.yaml  -d  -f kube | tail -n +2 |  tail -n +2 |awk '{printf ("docker tag IP:PORT/%s:%s  IP:PORT/%s:%s  && docker push IP:PORT/%s:%s \n",$4,$5,$2,$3,$2,$3)}' | sh
```

### 4.全局参数说明

```shell
-c context  操作指定的harbor的名字或者别名 如果不指定则默认使用配置文件中的current-context指定的harbor,需要指定名称或者别名
-f format   执行输出的格式 默认为表格 table 可以指定 kube 参数
```













