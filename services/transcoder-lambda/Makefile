.PHONY: clean build deploy mediainfo

MEDIA_INFO_TAG=media-info-image
TERRAFORM_SRC_PATH=../../terraform/src/transcoder/
TERRAFORM_VAR_FILE=../../terraform/env/transcoder/example.tfvars
LAMBDA_PATH=$(shell pwd)
LAMBDA_ZIP_FILE=transcoder-lambda.zip
LAMBDA_FILE_PATH=$(LAMBDA_PATH)/$(LAMBDA_ZIP_FILE)

clean:
	rm -rf ./bin $(LAMBDA_ZIP_FILE)

mediainfo:
	docker build ./mediainfo -t $(MEDIA_INFO_TAG)
	mkdir -p ./bin/
	docker run --rm $(MEDIA_INFO_TAG):latest cat mediainfo > ./bin/mediainfo
	chmod +x ./bin/mediainfo

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/transcoder ./src/lambda/main.go
	$(MAKE) mediainfo
	zip $(LAMBDA_ZIP_FILE) bin/transcoder bin/mediainfo .env

deploy: clean build
	terraform init $(TERRAFORM_SRC_PATH)
	terraform apply -var-file=$(TERRAFORM_VAR_FILE) -var="lambda_file_path=$(LAMBDA_FILE_PATH)" $(TERRAFORM_SRC_PATH)

destroy: clean build
	terraform destroy -var-file=$(TERRAFORM_VAR_FILE) -var="lambda_file_path=$(LAMBDA_FILE_PATH)" $(TERRAFORM_SRC_PATH)
	$(MAKE) clean
