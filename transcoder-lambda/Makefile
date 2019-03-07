.PHONY: clean build deploy

TERRAFORM_SRC_PATH=../deploy/terraform/src/transcoder/
LAMBDA_PATH=$(shell pwd)
LAMBDA_ZIP_FILE=transcoder-lambda.zip
LAMBDA_FILE_PATH=$(LAMBDA_PATH)/$(LAMBDA_ZIP_FILE)
TERRAFORM_VAR_FILE=../deploy/terraform/env/transcoder/prod.tfvars

clean:
	rm -rf ./bin $(LAMBDA_ZIP_FILE)

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/transcoder ./src/lambda/main.go
	zip $(LAMBDA_ZIP_FILE) bin/transcoder .env

deploy: clean build
	terraform init $(TERRAFORM_SRC_PATH)
	terraform apply -var-file=$(TERRAFORM_VAR_FILE) -var="lambda_file_path=$(LAMBDA_FILE_PATH)" $(TERRAFORM_SRC_PATH)

destroy: clean build
	terraform destroy -var-file=$(TERRAFORM_VAR_FILE) -var="lambda_file_path=$(LAMBDA_FILE_PATH)" $(TERRAFORM_SRC_PATH)
	$(MAKE) clean
