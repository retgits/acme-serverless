# -----------------------------------------------------------------------------
# Description: Makefile
# Author(s): retgits <https://github.com/retgits/>
# Last updated: 2020-03-02
# 
# This software may be modified and distributed under the terms of the
# MIT license. See the LICENSE file for details.
# -----------------------------------------------------------------------------

## The stage to deploy to
stage         	:= dev

## The name of the project
project_name  	:= serverless-pos

## The author of the project
author          := retgits

## The name of the team
team			:= vcs

## The version of the project, either uses the current commit hash, or will default to "dev"
version       	:= $(strip $(if $(shell git describe --tags --always --dirty="-dev"),$(shell git describe --tags --always --dirty="-dev"),dev))

## The Amazon S3 bucket to upload files to
aws_bucket    	?= $$AWS_S3_BUCKET

# Suppress checking files and all Make output
.PHONY: help get clean deploy destroy
.SILENT: help get clean deploy destroy

# Targets
help: ## Displays the help for each target (this message)
	echo
	echo Usage: make [TARGET]
	echo
	echo Makefile targets
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	echo

get: ## Download the source code of the POS service
	echo
	git clone https://github.com/retgits/acme-serverless-pos
	echo

clean: ## Remove all generated files
	echo
	-rm s3-packaged.yaml
	echo

deploy: ## Deploy the table to Amazon S3
	echo
	aws cloudformation package --template-file s3-template.yaml --output-template-file s3-packaged.yaml --s3-bucket $(aws_bucket)
	aws cloudformation deploy \
		--template-file s3-packaged.yaml \
		--stack-name $(project_name)-$(stage) \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides Version=$(version) \
		Author=$(author) \
		Team=$(team) \
		Stage=$(stage)
	aws cloudformation describe-stacks --stack-name $(project_name)-$(stage) --query 'Stacks[].Outputs'
	cd ../../src && aws s3 sync `pwd` s3://acmeserverless-pos-$(author)-$(stage)
	echo

destroy: ## Deletes the CloudFormation stack and all created resources
	echo
	aws cloudformation delete-stack --stack-name $(project_name)-$(stage)
	echo