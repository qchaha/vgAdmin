前端admin模板引用了"panjiacheng"的vue-admin-template，项目地址：https://github.com/PanJiaChen/vue-admin-template

后端golang的gin框架，项目地址：https://github.com/gin-gonic/gin

开发环境：
go version go1.11.2 linux/amd64
npm 3.5.2
node v8.10.0

前端：
```
# 项目部署：
git clone https://github.com/qchaha/vgAdmin.git

# 安装依赖
npm install

# 开发热部署，访问地址：localhost:9528
npm run dev

# 生产部署
npm run build
```

后端：
```
#启动后端服务
cd go && go run .
```

后端数据结构，执行src/sql下的sql语句