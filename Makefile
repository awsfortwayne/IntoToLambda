PROFILE=default
REGION=us-west-1
S3_BUCKET=intro-to-aws-lambda
STACK=Fort-Wayne-AWS-User-Group-Intro-To-AWS-Lambda

TEMPLATE=template.yaml
PACKAGED=packaged.yaml

.PHONY: all setup teardown mod build clean package deploy local dynamodb create_table

all: build package deploy
	
mod:
	go mod tidy
	go mod verify

clean: 
	rm -f ./cmd/guestbook/guestbook
	rm -f ./$(PACKAGED)
	
build:
	GOOS=linux GOARCH=amd64 go build -o ./cmd/guestbook/guestbook ./cmd/guestbook

package:
	aws cloudformation package --template-file $(TEMPLATE) --s3-bucket $(S3_BUCKET) --output-template-file $(PACKAGED) --profile $(PROFILE) --region $(REGION)

deploy:
	aws cloudformation deploy --template-file $(PACKAGED) --stack-name $(STACK) --capabilities CAPABILITY_IAM --profile $(PROFILE) --region $(REGION)

local:
	sam local start-api

dynamodb:
	docker run -p 8000:8000 amazon/dynamodb-local

create_table:
	aws dynamodb create-table --cli-input-json file://create_table.json --endpoint-url http://localhost:8000/

setup:
	aws s3 mb s3://$(S3_BUCKET) --profile $(PROFILE) --region $(REGION)

teardown:
	aws cloudformation delete-stack --stack-name $(STACK) --profile $(PROFILE) --region $(REGION)
	aws s3 rb --force s3://$(S3_BUCKET) --profile $(PROFILE) --region $(REGION)