.PHONY: deploy destroy deploy-core deploy-destroy deploy-lambda destroy-lambda deploy-server destroy-server

deploy: deploy-core deploy-lambda deploy-server

destroy: destroy-server destroy-lambda destroy-core

deploy-core:
	cd terraform/src/core/; terraform init; terraform apply -auto-approve -var-file=../../env/core/example.tfvars

destroy-core:
	cd terraform/src/core/; terraform init; terraform destroy -auto-approve -var-file=../../env/core/example.tfvars

deploy-lambda:
	$(MAKE) -C ./services/transcoder-lambda deploy

destroy-lambda:
	$(MAKE) -C ./services/transcoder-lambda destroy

deploy-server:
	$(MAKE) -C ./services/server deploy

destroy-server:
	$(MAKE) -C ./services/server destroy
