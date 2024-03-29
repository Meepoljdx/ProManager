# ProManager

The Manager system for Prometheus

## 介绍

### 发生了什么？

这其实是一个搁浅了很多很多次的项目，一直以来都没有想好我到底要做什么，然后边写边改，最后什么也没做出来；

一开始我只是需要一个能够批量做Prometheus启停和部署的工具，所以我写了一个简易的后端用来做部署；

后来我觉得，也许我能写一个管理Prometheus的工具？所以我又写了一个版本，新的版本能够对Prometheus进行启停、安装、操作等，但是不能去进行配置的更改；

再后来，随着使用Prometheus越来越多，感觉需要的能力也越来越多，比如匹配AlertManager的告警规则之类的，于是，我又迷茫了。

我有太多想做的东西了，但能力又不足以支撑，所以到头来什么都没做出来。

### 到底做些什么？

其实到现在，我还是不知道到底做什么，所以先列下来，一条条来。

- Server端：Server用来对Agent进行管理以及发送一些操作；
- Agent端：Agent用来进行实际服务器上的操作控制，包含一些运维操作的能力，这个就等于是把我之前写的一些小工具集成起来；
- 服务管理：就是把之前Prometheus的各种操作集成起来，抽象成组件，暂时不去考虑采集器的部署和操作，以后可以的话一起集成进来；
- 文件感知：这是对服务器上的文件进行感知的能力，在Server进行编辑策略后，在Agent端统一进行管理，比如某个路径的关键文件，定期查看是否被修改，如果被修改了就进行记录；
- 脚本运行：通过Server端编辑计划任务，传递到agent端进行实际调度和操作；
- 用户管理：对用户的角色进行管理，允许给予系统管理员、业务管理员、普通用户三种权限；
- 标签管理：这个来源于Prometheus的启发，一直以来我都很困扰，如何对Prometheus进行分组，基于最基本的集群之下，不再提供别的组别类似的概念，就是用标签的方式，允许用户通过标签筛选来查找相关的Prometheus节点；
- 集群管理：允许创建集群

## 数据表结构

### User

|字段名|值类型|含义|约束|
|--|--|--|--|
|id|数值|用户id|主键|
|user_name|字符串|用户名||
|phone|字符串|电话|唯一|
|description|字符串|用户描述||
|create_at|时间类型|创建时间||
|update_at|时间类型|更新时间||

### Cluster

集群信息表

|字段名|值类型|含义|约束|
|--|--|--|--|
|id|数值|集群id|主键|
|cluster_name|字符串|集群名|唯一|
|description|字符串|集群描述||
|create_user|数值|创建集群的用户id|
|create_at|时间类型|创建时间|
|update_at|时间类型|更新时间|

### Role

角色信息表

|字段名|值类型|含义|约束|
|--|--|--|--|
|id|数值|角色id|主键|
|role_name|字符串|角色名||
|description|字符串|角色描述||
|create_at|时间类型|创建时间||
|update_at|时间类型|更新时间||

### Role&User&Cluster

定义用户在具体的集群中所属的角色，以此能够得到用户的权限

|字段名|值类型|含义|约束|
|--|--|--|--|
|id|数值|角色id|主键|
|role_id|数值|角色id|联合唯一|
|user_id|数值|用户id|联合唯一|
|cluster_id|数值|集群id|联合唯一|
|create_at|时间类型|创建时间||
|update_at|时间类型|更新时间||

### Label

|字段名|值类型|含义|约束|
|--|--|--|--|
|id|数值|角色id|主键|
|label_name|字符串|标签名称|联合唯一|
|cluster_id|数值|集群id|联合唯一|
|create_user|数值|创建用户id||
|create_at|时间类型|创建时间||
|update_at|时间类型|更新时间||

### Component

|字段名|值类型|含义|约束|
|--|--|--|--|
|id|数值|组件id|主键|
|component_name|字符串|组件名称|联合唯一|
|description|字符串|组件描述||
|version|字符串|组件版本|联合唯一|
|create_at|时间类型|创建时间||
|update_at|时间类型|更新时间||

### Host

|字段名|值类型|含义|约束|
|--|--|--|--|
|id|数值|服务器id|主键|
|cluster_id|数值|集群id|
|ip|字符串|ip|唯一|
|hostname|字符串|主机名|
|cpu|数值|cpu核数|
|mem|字符串|内存使用率|
|root|字符串|根目录利用率|
|os|字符串|操作系统|
|status|字符串|主机状态|
|create_at|时间类型|创建时间|
|update_at|时间类型|更新时间|

### Service

Service是component的实例，描述了一个组件具体存在于哪个节点，service关联host，通过host关联cluster

|字段名|值类型|含义|约束|
|--|--|--|--|
|id|数值|Serviceid|主键|
|host_id|数值|主机id||
|status|字符串|服务状态||
|create_at|时间类型|创建时间||
|update_at|时间类型|更新时间||

