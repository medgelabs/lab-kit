# Usage: NHOST="IP_ADDR_HERE" make TARGET
NODE_HOST=${HOST}
USER=medge
KUBE_MANAGER_IP=192.168.1.131
METALLB_RANGE=192.168.1.200-192.168.1.240

# Provision a new node
.PHONY: provision
provision:
	# Enable passwordless sudo for service account
	scp config/sudoers ${NODE_HOST}:/home/${USER}/sudoers
	@echo "SSH'ing into box. Run:"
	@echo "sudo rsync --remove-source-files --chown=root:root sudoers /etc/sudoers.d/medge"
	@echo "Then exit to continue installation"
	ssh ${NODE_HOST}
	# Necessary installs & setup
	scp install.sh ${NODE_HOST}:/home/${USER}/install.sh
	ssh ${NODE_HOST} "chmod +x /home/${USER}/install.sh && sudo /home/${USER}/install.sh && rm /home/${USER}/install.sh"
	# SSH lockdown
	scp config/sshd_config ${NODE_HOST}:/tmp/sshd_config
	ssh ${NODE_HOST} 'sudo mv /tmp/sshd_config /etc/ssh/sshd_config && sudo systemctl restart sshd'
	ssh ${NODE_HOST} 'rm /home/${USER}/sudoers'

.PHONY: kube-manager
kube-manager: provision
	ssh ${NODE_HOST} 'export INSTALL_K3S_EXEC=" --write-kubeconfig-mode 644 --no-deploy servicelb --no-deploy traefik"; curl -sfL https://get.k3s.io | sh -'
	scp ${NODE_HOST}:/etc/rancher/k3s/k3s.yaml ./kubectl.yaml.tmp
	sed -e "s/127.0.0.1:6443/${KUBE_MANAGER_IP}:6443/g" kubectl.yaml.tmp > ./kubectl.yaml
	rm ./kubectl.yaml.tmp
	@echo 'Now, copy ./kubectl.yaml cluster config to ~/.kube/config'
	@echo 'If all good, kubectl get nodes should now work'

.PHONY: kube-worker
kube-worker: provision
	export TOKEN=$$(ssh ${KUBE_MANAGER_IP} 'sudo cat /var/lib/rancher/k3s/server/node-token') && \
		ssh ${NODE_HOST} "export K3S_URL=https://${KUBE_MANAGER_IP}:6443;export K3S_TOKEN='$${TOKEN}';curl -sfL https://get.k3s.io | sh -"

.PHONY: init-k8s
init-k8s:
	# metallb Load Balancer
	kubectl apply -f k8s/metallb-v0.9.5.yml
	kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$$(openssl rand -base64 128)" || exit 0
	sed -e "s/METALLB_RANGE/${METALLB_RANGE}/g" k8s/metallb-config.yml | kubectl apply -f -
	kubectl wait --namespace metallb-system \
            --for=condition=ready pod \
            --selector=app=metallb \
            --timeout=120s
	# ingress-nginx
	kubectl apply -f k8s/ingress-nginx.yml
	kubectl wait --namespace ingress-nginx \
            --for=condition=ready pod \
            --selector=app.kubernetes.io/component=controller \
            --timeout=120s
	# TLS certs with LetsEncrypt
	kubectl apply --validate=false -f k8s/cert-manager.yml
	kubectl wait --namespace cert-manager \
            --for=condition=ready pod \
            --selector=app.kubernetes.io/instance=cert-manager \
            --timeout=120s
	kubectl apply -f k8s/lets-encrypt-issuer.yml

.PHONY: test-k8s
test-k8s:
	kubectl apply -f test-cluster.yml
	kubectl wait --namespace default --for=condition=ready pod --selector=name=test-cluster --timeout=120s
	export INGRESS_IP=$$(kubectl get svc -n ingress-nginx | grep ingress-nginx | awk '{ print $$4 }'); curl https://$${INGRESS_IP}/test-cluster
	@echo -n "Press ENTER when done testing" && read
	kubectl delete -f test-cluster.yml

.PHONY: k3s-logs
k3s-logs:
	ssh ${NODE_HOST} 'sudo journalctl -feu k3s'
