K8s搭建

1.1 安装Docker

    sudo yum install -y yum-utils \
    device-mapper-persistent-data \
    lvm2

    yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo

1.2 可用版本

    yum list docker-ce --showduplicates | sort -r


    yum install docker-ce-18.06.2.ce-3.el7 docker-ce-cli-18.06.2.ce-3.el7 containerd.io

1.3 启动docker

    systemctl start docker
    systemctl enable docker

2. 安装kubeadm, kubelet and kubectl

2.1 配置kubernetes的yum源

    cat <<EOF > /etc/yum.repos.d/kubernetes.repo
    [kubernetes]
    name=Kubernetes
    baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
    enabled=1
    gpgcheck=0
    repo_gpgcheck=0
    gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
    EOF

2.2 安装

    # Set SELinux in permissive mode (effectively disabling it)
    setenforce 0
    sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

    yum install -y kubeadm

    systemctl enable --now kubelet


    # very important
    cat <<EOF >  /etc/sysctl.d/k8s.conf
    net.bridge.bridge-nf-call-ip6tables = 1
    net.bridge.bridge-nf-call-iptables = 1
    EOF
    sysctl --system

    hostnamectl set-hostname k8s-node1

2.3 加入集群

    单节点集群：
        kubeadm init --image-repository=registry.aliyuncs.com/google_containers --pod-network-cidr=10.244.0.0/16
        kubectl apply -f <add-on.yaml>
        kubectl taint nodes --all node-role.kubernetes.io/master-

    kubeadm join 172.27.0.16:6443 --token hlb95v.8kflzbkzfv5w1xtu --discovery-token-ca-cert-hash sha256:cc4eb10ffbc912d2ddc438e907b5e1d3a5e19b237a69af64d4e8a79134e5871d


2.4 Dashboard

    kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta4/aio/deploy/recommended.yaml
    kubectl proxy
    curl http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/

2.4.1
    dashboard需要跳过登录
    配置admin角色权限