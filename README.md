## 使用方法
安装方法
go build 
迁移工具支持etcdv2和etcdv3两个版本的数据的迁移
-dest-etcd-address 目标集群的服务器地址
-dest-cert  目标集群的证书
-dest-key   目标集群的私钥
-dest-cacert 签发目标集群的证书的跟证书
-version 需要迁移的集群的版本，默认是3
-src-etcd-address 源集群的服务器地址
-src-cert   源集群的客户端证书
-src-key    源集群的客户端的证书的私钥
-src-cacert  源证书的客户端证书的根证书
-overwrite   存在一样内容是否覆盖
首次使用需要将parttern文件里面的blacklist更新一次,默认是拒绝所有


例子:
```
#!/bin/bash
SRC_CERT_PATH=/home/larry/Work/k8s/cert
DEST_CERT_PATH=/home/larry/Work/k8s/cluster/cluster1/cert
./etcd_migrate_tools -dest-etcd-address https://xx.xx.xx.xx:2379 -dest-cert $DEST_CERT_PATH/admin/admin.pem -dest-key $DEST_CERT_PATH/admin/admin-key.pem -dest-cacert $DEST_CERT_PATH/ca/ca.pem -version 3 -src-etcd-address https://xx.xx.xx.xx:2379 -src-cert $CERT_PATH/admin/admin.pem -src-key $CERT_PATH/admin/admin-key.pem -src-cacert $CERT_PATH/ca/ca.pem
``

如果没有使用https可以不指定证书
