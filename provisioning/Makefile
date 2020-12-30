# Usage: NHOST="IP_ADDR_HERE" make TARGET
HOST=${NHOST}
USER=medge
KUBE_MANAGER_IP=192.168.1.131
KUBE_METALLB_CIDR=192.168.1.144/28

# Provision a new node
.PHONY: provision
provision:
	  # Enable passwordless sudo for service account
		scp config/sudoers ${HOST}:/home/${USER}/sudoers
		@echo "SSH'ing into box. Run:"
		@echo "sudo rsync --remove-source-files --chown=root:root sudoers /etc/sudoers.d/medge"
		@echo "Then exit to continue installation"
		ssh ${HOST}
	  # Necessary installs & setup
		scp install.sh ${HOST}:/home/${USER}/install.sh
		ssh ${HOST} "chmod +x /home/${USER}/install.sh && sudo /home/${USER}/install.sh && rm /home/${USER}/install.sh"
		# SSH lockdown
		scp config/sshd_config ${HOST}:/tmp/sshd_config
		ssh ${HOST} 'sudo mv /tmp/sshd_config /etc/ssh/sshd_config && sudo systemctl restart sshd'

.PHONY: kube-manager
kube-manager: provision
	ssh ${HOST} 'export INSTALL_K3S_EXEC=" --no-deploy servicelb --no-deploy traefik"; curl -sfL https://get.k3s.io | sh -'
	ssh ${HOST} "sudo cp /etc/rancher/k3s/k3s.yaml /home/${USER}/k3s.yaml && sudo chown ${USER}:${USER} /home/${USER}/k3s.yaml"
	scp ${HOST}:/home/${USER}/k3s.yaml ./kubectl.yaml.tmp
	sed -e "s/127.0.0.1:6443/${KUBE_MANAGER_IP}:6443/g" ./kubectl.yaml.tmp > ./kubectl.yaml
	rm ./kubectl.yaml.tmp
	ssh ${HOST} "rm /home/${USER}/k3s.yaml"
	@echo 'Now, copy ./kubectl.yaml cluster config to ~/.kube/config'
	@echo 'If all good, kubectl get nodes should now work'

.PHONY: kube-worker
kube-worker: provision
	export TOKEN=$$(ssh ${KUBE_MANAGER_IP} 'sudo cat /var/lib/rancher/k3s/server/node-token') && \
		ssh ${HOST} "export K3S_URL=https://${KUBE_MANAGER_IP}:6443;export K3S_TOKEN='$${TOKEN}';curl -sfL https://get.k3s.io | sh -"

.PHONY: init-k8s
init-k8s:
	kubectl apply -f k8s/metallb-v0.9.5.yml
	kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$$(openssl rand -base64 128)" || exit 0
	kubectl apply -f k8s/metallb-config.yml
	# ingress-nginx
	kubectl apply -f k8s/ingress-nginx.yml
	kubectl wait --namespace ingress-nginx \
            --for=condition=ready pod \
            --selector=app.kubernetes.io/component=controller \
            --timeout=120s
	# kubectl apply --validate=false -f k8s/cert-manager-v1.0.3.yml
	# kubectl apply -f k8s/lets-encrypt-issuer.yml

.PHONY: test-k8s
test-k8s:
	kubectl apply -f test-cluster.yml
	kubectl wait --namespace default --for=condition=ready pod --selector=app=test-cluster --timeout=120s
	@echo "You should see an NGINX HTML output if all is working"
	export INGRESS_IP=$$(kubectl get svc -n ingress-nginx | grep ingress-nginx | awk '{ print $$4 }'); curl http://$${INGRESS_IP}/test
	kubectl delete -f test-cluster.yml

.PHONY: k3s-logs
k3s-logs:
	ssh ${HOST} 'sudo journalctl -feu k3s'
