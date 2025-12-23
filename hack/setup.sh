#!/bin/bash
# Setting up external node for minikube
set -euo pipefail
shopt -s failglob

CONTAINER_NAME="minikube"
HOME=/home/ish

# NOTE: before running this script install crio, kubelet, docker and run minikube start --driver=docker --cni=calico
echo "NOTE: before running this script install crio, kubelet, docker and run minikube start --driver=docker --cni=calico"
echo "---"
echo "press <enter> to continue"
read pause_exec

echo ">> Creating directory structure & certs.."
mkdir -p /etc/kubernetes/
mkdir -p /var/lib/kubelet/pki
mkdir -p /var/lib/kubelet/certs
mkdir -p /custom-node/etc/cni/
mkdir -p /custom-node/opt/cni/

openssl genrsa -out /var/lib/kubelet/pki/custom-node.key 2048
openssl req -new -key /var/lib/kubelet/pki/custom-node.key -out  /var/lib/kubelet/pki/custom-node.csr -subj "/O=system:nodes/CN=system:node:custom-node"
openssl x509 -req -in /var/lib/kubelet/pki/custom-node.csr -CA $HOME/.minikube/ca.crt -CAkey $HOME/.minikube/ca.key -CAcreateserial -out /var/lib/kubelet/pki/custom-node.crt -days 365 -extensions v3_req

echo ">> Exporting configs..."
docker cp $CONTAINER_NAME:/etc/kubernetes/kubelet.conf /etc/kubernetes/kubelet.conf
docker cp $CONTAINER_NAME:/var/lib/minikube/certs/ca.crt /var/lib/kubelet/certs/ca.crt
docker cp $CONTAINER_NAME:/var/lib/kubelet/config.yaml /var/lib/kubelet/config.yaml
docker cp $CONTAINER_NAME:opt/cni/bin /custom-node/opt/cni/
docker cp $CONTAINER_NAME:/etc/cni/net.d /custom-node/etc/cni/

echo ">> Customizing exported config..."
sed -i 's#unix:///var/run/cri-dockerd.sock#unix:///run/crio/crio.sock#g' /var/lib/kubelet/config.yaml
sed -i 's#/var/lib/minikube/certs/ca.crt#/var/lib/kubelet/certs/ca.crt#g' /var/lib/kubelet/config.yaml
sed -i 's#client-certificate: /var/lib/kubelet/pki/kubelet-client-current.pem#client-certificate: /var/lib/kubelet/pki/custom-node.crt#g' /etc/kubernetes/kubelet.conf
sed -i 's#client-key: /var/lib/kubelet/pki/kubelet-client-current.pem#client-key: /var/lib/kubelet/pki/custom-node.key#g' /etc/kubernetes/kubelet.conf

echo ">> Setting up kubelet systemd unit..."
cp kubelet.service /usr/lib/systemd/system/kubelet.service
cp crio.conf /etc/crio/crio.conf.d/10-crio.conf
systemctl daemon-reload

echo ">> Restarting crio..."
systemctl restart crio

echo ">> Restarting kubelet..."
systemctl restart kubelet

echo ">> Labeling node.."
kubectl --kubeconfig=$HOME/.kube/config label no custom-node kubernetes.io/role=worker
