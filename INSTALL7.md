### 系统更新及第三方安装
```
yum install iptables-services net-tools

systemctl mask firewalld
systemctl stop firewalld

systemctl enable iptables
systemctl enable ip6tables

systemctl start iptables
systemctl start ip6tables

yum update
yum groupinstall "Development Tools"

yum install epel-release
yum install python36
yum install python3-devel openssl-devel pcre-devel

yum install nc iptraf nethogs sharutils

yum install nginx
```

### 安装go
```
tar xvfz go1.15.6.linux-amd64.tar.gz
mv go /opt/
ln -s /opt/go/bin/go /usr/local/bin/
go env -w GO111MODULE=auto
go env -w GOPROXY=https://goproxy.cn
```

### 需要复制的配置文件（最好从同类服务器上复制）
```
/etc/sysconfig/iptables
/usr/local/nginx/conf/nginx.conf
/usr/local/nginx/html/log_cut.sh
```

### crontab -e 添加任务
```
0 5 * * * /usr/local/nginx/html/log_cut.sh > /tmp/log_cut.log
```

### 启动前需调整的配置文件
/usr/local/nginx/conf/nginx.conf


### Nginx关闭自动压缩日志
```
rm /etc/logrotate.d/nginx
```
