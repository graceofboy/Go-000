# 工程管理

## 工程项目结构
> [project-layout/README_zh.md at master · golang-standards/project-layout · GitHub](https://github.com/golang-standards/project-layout/blob/master/README_zh.md)  
> https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html  
> https://blog.golang.org/wire : explicitly providing components with all of   

the dependencies that need to work.
> [GitHub - facebook/ent: An entity framework for Go](https://github.com/facebook/ent)  
> [kratos/app.go at v2 · go-kratos/kratos · GitHub](https://github.com/go-kratos/kratos/blob/v2/app.go)  
> 目标：请基于实际理解去随机应变*灵活*，主要是多人协同的规范管理，使项目更便于devops。  
> 核心思想是围绕着DDD来走。  
<a href='%E5%B7%A5%E7%A8%8B%E7%AE%A1%E7%90%86/4.%20Go%20%E5%B7%A5%E7%A8%8B%E5%8C%96%E5%AE%9E%E8%B7%B5.pptx'>4. Go 工程化实践.pptx</a>

```
.
├── LICENSE
├── Makefile
├── README.md
├── api
│   ├── README.md
│   └── helloworld
│       ├── helloworld.pb.go
│       ├── helloworld.proto
│       ├── helloworld_grpc.pb.go
│       └── helloworld_http.pb.go
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── biz
│   │   └── README.md
│   ├── data
│   │   └── README.md
│   └── service
│       ├── README.md
│       └── greeter.go
└── pkg
    ├── cache
    │   ├── memcache
    │   └── redis
    └── conf
        ├── dsn
        ├── env
        ├── flagvar
        └── paladin
```

- cmd 项目主干，负责启动、关闭、配置初始化等，不应该放过多代码
/cmd/myapp/main.go
- pkg可以导入其他项目中使用的代码，功能库「kit,util」
- internal/pkg 一般用于项目内的扩多个应用的公共共享代码，但其作用域仅在单个项目工程内。
- internal 自己项目使用的代码,Go的编译器本身会限制第三方项目去导入internal里面的项目。
### Kit 项目必须具备的特点：
1. 统一
2. 标准库方式布局
3. 高度抽象
4. 支持插件
### Service Application Project Layout
- /api 协议定义目录
- /configs 配置文档
	- toml、yaml等
- /test 单元测试集合
- /app 每个微服务一个app
	- interface 对外的BFF服务，接受来自用户的请求，比如最终暴露了HTTP/gRPC接口
		- /dao 依赖model
		- /service 依赖dao ::「service也经常依赖model怎么办？」::
		- /server 依赖service
	- service 对内的微服务，仅接受来自内部其他服务或者网关的请求，比如暴露了只对内的gRPC服务
	- admin 区别于service，更多是面向运营侧的服务，通常数据权限更高，隔离带来更好的代码级别安全。
	- job 流式任务处理的服务，上游一般依赖message broker
	- task 定时任务，类似crontab

## API设计
> 依赖倒置： 上层模块不应该依赖于下层模块，它们共同依赖于一个抽象；抽象不能依赖于具象，具象依赖于抽象。  
> 依赖注入： 方便做单次初始化和依赖注入。  
> 充血模型 | 贫血模型  
> /Facebook/ent : 类似GORM  
> 「git LGTM」  

### API 管理
为了统一检索和规范API，内部简历一个统一的bapis仓库，整合所有对内对外API。
- API仓库，方便跨部门协作
- 版本管理，基于git控制
- 规范化检查， API lint
- API design review, 变更diff
- 权限管理，目录 OWNERS

#### API Compatibility「也适用于数据库」
- 向后兼容的修改
	1. 新增API接口
	2. Request添加字段
	3. Response添加字段
- 向后不兼容的修改
	1. 删除或重命名服务、字段、方法或枚举值
	2. 修改字段类型
	3. 修改现有请求的可见行为
	4. 给资源消息添加 读取/写入 字段

#### API Naming Conventions
> [API 设计指南  |  Google Cloud](https://cloud.google.com/apis/design/?hl=zh-cn)  
> 为输入(requestBody)和输出(responseBody)都要定义一个message对象  
```
//RequestURL:/<package_name>.<version>.<service_name>/method
```
- Create…
- Get…
- List…
- Delete…
- Update…

#### API Primitive Fields
> https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/wrappers.proto  

#### API Errors

> 标准的HTTP状态码，方便运维组件监控  
> 一个404可以定义各种找不到  
> 「全局错误码」是不合适的  
> 上下游服务间的错误码应该拦截不要透传，遇到未知的应该报错  

- 使用一小组标准错误码配和大量资源
- 错误传播
	- 隐藏实现详细信息和机密信息
	- 调整负责该错误的一方
- 设计错误码分两层
	- 第一层 http原始错误码
	- 第二层 业务封装的错误码


### DTO （data transfer object）「毛剑见解」
- model ： 放对应“存储层”结构体。 「贫血模型：避免service逻辑一大坨，列举model.user的逻辑判定」
- dao ： 数据库和缓存处理，包括cache miss处理。
- service:  组和各种数据访问来构建业务逻辑
- server：依赖proto定义的服务作为入参，提供快捷的启动服务全局方法。
- api:  定义了API proto文件，和生成stud代码，生成interface,其实现逻辑在service中。

### DO 领域对象
> 就是从现实世界中抽象出来的有形或无形的业务实体，缺乏DTO->DO的对象转换。  
> DDD 的思想是偏向业务的  
- /biz 业务逻辑的组装层，类似DDD的domain层。 
- /data 处理业务数据，类似dao层处理数据库数据。
- /service 实现api定义的服务层。

### PO 持久化对象
> 持久化对象，跟持久层的数据结构形成一一对应的映射关系，如果持久层是关系型数据库，那么数据表中的每一个字段就对应PO的一个属性。  

### Lifecycle 程序生命周期
> [GitHub - go-kratos/service-layout: Kratos Service Layout](https://github.com/go-kratos/service-layout)  


## 配置管理
> field mask  
> 《google api》文档  
1. 不建议热加载，热重启配置文件
2. 基础库不好做reload
3. 作为一个公共函数尽量不传nil 
4. 位置文件区分可选必选
5. 区分环境变量和配置变量，业务程序配置文件尽可能少

# 毕业总结

	毛大这13周的课程里，带着我们去深入像工程管理、微服务架构、内存以及并发模型、或者是当前比较主流的架构比如Caffco等，给我搭建了一套完善的知识框架体系，以至于今后几年的学习有了更明确的目标。

	在这其中我觉得对我触发比较大的毛大所坚持的一些思想，比如看英文文档，编程素养，架构思想等这些需要靠去悟的东西，上完课我觉得自己确实有进入到另外一个层面，体现出来的一个点就是，以前遇到自己不懂的知识点，回到处去找各种各样的文档但是都没有回溯到源头以至于浪费了大量时间还学不到知识，现在我首先第一步就是回到代码中去找到问题的源头，再针对性的去翻阅一些比较权威的知识点。
