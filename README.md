## GThree

* *描述*
```
  GThree是一个多区服游戏运维工具,可以进行游戏区服的更新,启停,状态检查等。
  技术栈: golang + gin + grpc
  server & client架构:
  gtmaster接受web请求, 并通过grpc将请求转发到gtservant进行业务处理(多节点区服处理)
```
* *进度*
```
  开发中...
```
* *部署*
```
  运行build.sh后在release目录下生对应文件
  分别将release目录下的文件部署到运维服和游戏服
  将gtmaster二进制文件部署到运维服(目录结构):
    gtmaster/
    ├── gtmaster
    └── gtmaster.yaml
  将gtservant二进制文件部署到游戏服(目录结构):
    gtservant/
    ├── gtservant
    └── gtservant.yaml
```

<!-- 
* *生成grpc相关代码*
> protoc --go_out=../service --go_opt=paths=source_relative --go-grpc_out=../service/ --go-grpc_opt=paths=source_relative *.proto 


* *生成证书双向认证*
> 拷贝openssl.cnf文件到项目目录下的cert文件夹:cp /etc/ssl/openssl.cnf ./cert
>> 修改./cert/openssl.cnf文件
>> [req段落]取消注释:
>> req_extensions = v3_req # The extensions to add to a certificate request
>> 添加如下配置:
>> [ v3_req ]
>> # Extensions to add to a certificate request

>> basicConstraints = CA:FALSE
>> keyUsage = nonRepudiation, digitalSignature, keyEncipherment
>> subjectAltName = @alt_names

>> [ alt_names ]
>> DNS.1 = www.test.gthree.com

> 生成ca文件：
>> openssl genrsa -out ca.key 2048
>> openssl req -x509 -new -nodes -key ca.key -subj "/CN=test.gthree.com" -days 5000 -out ca.pem

> 生成服务端证书:
>> openssl req -new -nodes -subj "/C=CN/ST=Chengdu/L=Chengdu/O=grpcdev/OU=grpcdev/CN=www.test.gthree.com" -config <(cat openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:www.test.gthree.com")) -keyout server.key -out server.csr
>> openssl x509 -req -days 365000 -in server.csr -CA ca.pem -CAkey ca.key -CAcreateserial -extfile <(printf "subjectAltName=DNS:www.test.gthree.com") -out server.pem

> 生成客户端证书:
>> openssl req -new -nodes -subj "/C=CN/ST=Chengdu/L=Chengdu/O=grpcdev/OU=grpcdev/CN=www.test.gthree.com" -config <(cat openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:www.test.gthree.com")) -keyout client.key -out client.csr
>> openssl x509 -req -days 365000 -in client.csr -CA ca.pem -CAkey ca.key -CAcreateserial -extfile <(printf "subjectAltName=DNS:www.test.gthree.com") -out client.pem
-->