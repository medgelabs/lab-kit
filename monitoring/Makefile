.PHONY: install
install:
	kubectl apply -f prometheus.yml

.PHONY: nuke
nuke:
	kubectl delete -f prometheus.yml

.PHONY: redeploy
redeploy: nuke install

.PHONY: expose
expose:
	kubectl get pods -n monitoring | grep prometheus | awk '{ print $$1 }' | xargs -I{} kubectl port-forward -n monitoring {} 9090

