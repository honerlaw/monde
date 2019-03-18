.PHONY: docker deploy-core deploy-lambda

TAG=latest
URL=106480132517.dkr.ecr.us-east-1.amazonaws.com/monde:$(TAG)

docker:
	eval $(aws ecr get-login --no-include-email --region=us-east-1)
	docker build -t $(URL) .
	docker push $(URL)

deploy-core:
	cd terraform/src/core/; terraform init; terraform apply -var-file=../../env/core/example.tfvars

destroy-core:
	cd terraform/src/core/; terraform init; terraform destroy -var-file=../../env/core/example.tfvars

deploy-lambda:
	$(MAKE) -C ./transcoder-lambda deploy

destroy-lambda:
	$(MAKE) -C ./transcoder-lambda destroy
