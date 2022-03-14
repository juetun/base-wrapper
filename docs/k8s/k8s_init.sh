#!/bin/sh

# 参考资料
# 1、https://www.cnblogs.com/w1sh/p/15509347.html
# 2、https://www.cnblogs.com/tangy1/p/14925216.html
# 关闭防火墙
systemctl stop firewalld
systemctl disable firewalld

# 关闭selinux
sed -i 's/enforcing/disabled/' /etc/selinux/config  #永久

setenforce 0  #临时

# 关闭swap（k8s禁止虚拟内存以提高性能）
sed -ri 's/.*swap.*/#&/' /etc/fstab #永久

swapoff -a #临时

# 在master添加hosts

cat >> /etc/hosts << EOF
10.10.2.xxx k8smaster
10.10.2.xxx k8snode1
10.10.2.xxx k8snode2
10.10.2.xxx k8snode3
EOF

# 设置网桥参数
cat > /etc/sysctl.d/k8s.conf << EOF

net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward=1
vm.max_map_count=262144
EOF

sysctl --system  #生效
# 时间同步

yum install ntpdate -y
ntpdate time.windows.com


yum install wget -y

wget https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo -O /etc/yum.repos.d/docker-ce.repo
yum install docker-ce.x86_64 3:20.10.13-3.el7 -y


cat > /etc/sysctl.d/k8s.conf << EOF

{
"registry-mirrors": ["https://registry.docker-cn.com"],
"exec-opts": ["native.cgroupdriver=systemd"]
}
EOF


# 启动docker
systemctl enable docker && systemctl start docker

# 安装1.16版本
yum install -y kubelet kubeadm kubectl

# 添加k8s的阿里云YUM源
cat > /etc/yum.repos.d/kubernetes.repo << EOF

[kubernetes]

name=Kubernetes

baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64

enabled=1

gpgcheck=0

repo_gpgcheck=0

gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg

EOF

yum install kubelet kubeadm kubectl -y

systemctl enable kubelet.service

#6.初始化配置文件，只在master01上执行
mkdir ~/k8s-install && cd ~/k8s-install

#生成配置文件
kubeadm config print init-defaults > ~/k8s-install/kubeadm.yaml


# variable=`ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:"`

# sed -i "s/enforcing/${variable}/" ~/k8s-install/kubeadm.yaml  #永久

#提前下载镜像
kubeadm config images list --config ~/k8s-install/kubeadm.yaml

docker images|grep aliyun

kubeadm init --config ~/k8s-install/kubeadm.yaml

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config