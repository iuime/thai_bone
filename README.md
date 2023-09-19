#更新

1、将查询到的内容保存成文件。

2、遍历所有数据库。

3、排除mysql默认的几个库
```
mysql
information_schema
sys
performance_schema
innodb_sys_data
innodb_sys_undo
```




# 泰式摸骨

## 项目目的

快速发现数据库内涉及隐私数据。为数据脱敏做准备工作。

## 安装

下载[release](https://github.com/SecurityPaper/thai_bone/releases)对应系统的二进制包

解压缩

```tar zxfv xxx.tar.gz```

```wget -O config.yaml  https://raw.githubusercontent.com/SecurityPaper/thai_bone/v1.1/config-example.yaml```

```nano config.yaml``` 修改数据库链接地址和账号密码

```./thai_bone```直接运行即可
