## 持久层，负责检索和修改域模型的数据。 它只依赖（depend）于模型（model）层。 数据服务可以通过RPC或RESTFul调用从数据库或其他微服务获取数据。