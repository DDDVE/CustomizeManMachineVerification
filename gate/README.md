# 业务网关介绍
此项目为整个系统的网关

## 代码结构

gate

    handler     处理层，路由匹配后进入

    schedule    定时任务或周期任务模块

    test        单元测试代码

    utils       工具类部分，调用层级中处于最底层

    main.go     项目入口，负责初始化各个模块和分配路由处理函数

## 项目功能
已实现的有：
1.  路由判断
2.  路由过滤
3.  jwt验证
4.  token定时灰度更新
5.  api网关动态注册和发现
6.  定时检查api网关存活情况
7.  各个既定路由的请求重定向或转发
8.  从本地读取黑名单并过滤
9.  桶令牌限流
10. 根据CPU使用率调整服务级别

待完成的有：

