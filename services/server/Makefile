.PHONY: build docker clean deploy destroy

TERRAFORM_SRC_PATH=../../terraform/src/server/
TERRAFORM_VAR_FILE=../../terraform/env/server/example.tfvars
TAG=latest
URL=106480132517.dkr.ecr.us-east-1.amazonaws.com/vueon:$(TAG)

clean:
	rm -rf bin

build:
	cd ./src/server/ && go get -d -v
	go build -o bin/server ./src/server/main.go

docker:
	eval $(shell aws ecr get-login --no-include-email --region=us-east-1)
	docker build -t $(URL) .
	docker push $(URL)

deploy:
	terraform init $(TERRAFORM_SRC_PATH)
	terraform apply -var-file=$(TERRAFORM_VAR_FILE) -var="image=$(URL)" $(TERRAFORM_SRC_PATH)

destroy:
	terraform init $(TERRAFORM_SRC_PATH)
	terraform destroy -var-file=$(TERRAFORM_VAR_FILE) -var="image=$(URL)" $(TERRAFORM_SRC_PATH)