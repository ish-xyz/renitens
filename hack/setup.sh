#!/bin/bash
# Setting up external node for minikube

set -e

echo "Creating directory structure & certs.."
mkdir -p /etc/kubernetes/
mkdir -p /var/lib/kubelet/pki
mkdir -p /var/lib/kubelet/certs
openssl genrsa -out /var/lib/kubelet/pki/custom-node.key 2048
openssl req -new -key /var/lib/kubelet/pki/custom-node.key -out  /var/lib/kubelet/pki/custom-node.csr -subj "/O=system:nodes/CN=system:node:custom-node"
openssl x509 -req -in /var/lib/kubelet/pki/custom-node.csr -CA /home/ish/.minikube/ca.crt -CAkey /home/ish/.minikube/ca.key -CAcreateserial -out /var/lib/kubelet/pki/custom-node.crt -days 365 -extensions v3_req

echo "Exporting configs..."
docker cp minikube:/etc/kubernetes/kubelet.conf /etc/kubernetes/kubelet.conf
docker cp -L minikube:/var/lib/minikube/certs/ca.crt /var/lib/kubelet/certs/ca.crt
docker cp minikube:/var/lib/kubelet/config.yaml /var/lib/kubelet/config.yaml

echo "Customizing exported config..."
sed -i 's#unix:///var/run/cri-dockerd.sock#unix:///run/crio/crio.sock#g' /var/lib/kubelet/config.yaml
sed -i 's#/var/lib/minikube/certs/ca.crt#/var/lib/kubelet/certs/ca.crt#g' /var/lib/kubelet/config.yaml
sed -i 's#client-certificate: /var/lib/kubelet/pki/kubelet-client-current.pem#client-certificate: /var/lib/kubelet/pki/custom-node.crt#g' /etc/kubernetes/kubelet.conf
sed -i 's#client-key: /var/lib/kubelet/pki/kubelet-client-current.pem#client-key: /var/lib/kubelet/pki/custom-node.key#g' /etc/kubernetes/kubelet.conf

echo "Setting up kubelet systemd unit..."
cp kubelet.service /usr/lib/systemd/system/kubelet.service
systemctl daemon-reload

echo "Restarting crio..."
systemctl restart crio

echo "Restarting kubelet..."
systemctl restart kubelet
