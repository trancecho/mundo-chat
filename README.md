# mundo-chat
mundo聊天室开发仓库

**配置文件格式可以参考如下**

```yaml
app:
  logFile: log/gin.log
  Port: "12387"
  URL:  "127.0.0.1:12387"
  Jwt: "1111111"


redis:
  addr: "localhost:6379"
  password: ""
  DB: 0
  poolSize: 30
  minIdleConns: 30
```
注意Jwt需要和生成的代码的部分保持一致，自行配置
然后放入./config/app.yaml