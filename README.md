# 一个广播数据库变更的小工具

go-mysql-elasticsearch是一个golang包，可以以一个MysqlServer的salve的身份，监听mysql的数据变更，并将变更导入到elasticsearch。
我看到这个包的想法是，是否可以稍做修改，当mysql数据库有变更时，将变更广播到消息队列中去,增强及时性？

昨天睡前喝了点小酒，拿着elasticsearch的代码就开始动刀了，我把访问elasticsearch的那一部分删掉，改成提交数据到ws.z.12zan.net 到凌晨1点的时候已经基本完成了。

源代码放在:[https://github.com/yuanfenxi/yuanlicast](https://github.com/yuanfenxi/yuanlicast)上了。

用法和go-mysql-elasticsearch基本一致，启动是: ``` ./bin/yuancast -config=./etc/river.toml ```

### 配置文件
配置文件也在./etc/river.toml里。
配置里几个关键项如下:
* ``_bulk_size = 1`` 
 
这个请修改成1，这样每次有数据变更就提交。不然就失去项目存在的意义了。

* ``` yfx_gateway = "https://ws.z.12zan.net/dbcast/some_channel_as_you_wish" ``` 

每个不同的URL地址就是一个不同的频道。这里可以随便填，请确保你填写的和别人的是不一样的频道地址。

* ``` yfx_secret = hellloafsdfefss9843ru93f ```

广播的时候，对数据进行一下加密，用的是这个密钥。

### 客户端

客户端解密的代码需要根据自己的应用场景来写了。
我提供了一个nodejs的示例，在 ```jsClientExample ``` 目录下:
```
npm install
node ./test.js
```




