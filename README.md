# Fort Wayne AWS User Group - Into to AWS Lambda 

The repo contains the slice and demo code for the Fort Wayne AWS User Group's presention Into to AWS Lambda

## Getting Started

These instructions will get you setup to run the demo both in AWS and locally.


### Prerequisites

The demo code requires that you have an AWS account and have installed the AWS CLI.

[How do I create and activate a new Amazon Web Services account?](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/)
[Installing the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)

If you would like to run it on your local machine you'll also need the Docker and the AWS SAM CLI. 

[Installing Docker](https://docs.docker.com/install/)
[Installing the AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

The demo is written in Go. If you want want to compile it yourself you'll need to install the Go compiler, but we've include pre-compiled binary as well if do not want to compile it yourself.

[Download Go](https://golang.org/dl/)

### Makefile

A Makefile has been included to make running the commands easier. The top of this files contains a few variables that you'll need to update for your own setup.

* PROFILE - This is the AWS profile you wish to use for your AWS CLI commands. It is initally set to *default* but you many wish to use a differnt profile that you have setup.
* REGION - This is the AWS region you wish to use. It is initiall set to *us-west-1*.
* S3_BUCKET - The deployment process requires an S3 bucket for staging the Lambda binary. This variable is initially set to *intro-to-aws-lambda*, but S3 bucket names must be globally unique, so this name may already be taken. You'll to create a new name for the bucket. **NOTE:** You do not need to create this bucket yourself. There are commands in the Makeflie and creating and deleting this bucket.
* STACK - SAM deployments use AWS CloudFormation and require an name fo the CloudFormation stack. This is initally set to *Fort-Wayne-AWS-User-Group-Intro-To-AWS-Lambda* but you can make it anything you want.

## Deploying the code to AWS

### Setup

This step create the S3 bucket that will be a stagging area for you lamdba code before it is deployed to the AWS Lamdba service

```
make setup
```
or
```
aws s3 mb s3://intro-to-aws-lambda --profile default --region us-west-1
```

### Compiling the Demo

This is an optional step that requires the Go compiler be installed on your machine. If you would like, you can skip this step and use the included pre-compiled binary. The binary that is created will be a Linux AMD64 executable regardless of the type of machine on which it is build. This is a requirement for AWS Lambda. You will not need to run this binary directly so don't worry if you are not using a Linux machine.

```
make mod
make build
```
or
```
go mod tidy
go mod verify
GOOS=linux GOARCH=amd64 go build -o ./cmd/guestbook/guestbook ./cmd/guestbook
```

### Packaging the Demo

This step will create a ZIP file of the Lambda binary and will upload it to your S3 bucket. It will also create an AWS CloudFormation template named *packaged.yaml* from the AWS SAM template named *template.yaml*

```
make package
```
or
```
aws cloudformation package --template-file template.yaml --s3-bucket intro-to-aws-lambda --output-template-file packaged.yaml --profile default --region us-west-1
```

### Deploying the Demo

This step will create all the required AWS resources needed for the demo and will deploy the AWS Lambda. You can monitor the creating of these resources in the AWS CloudFormation service.

```
make deploy
```
or
```
aws cloudformation deploy --template-file packaged.yaml --stack-name Fort-Wayne-AWS-User-Group-Intro-To-AWS-Lambda --capabilities CAPABILITY_IAM --profile default --region us-west-1
```


### Testing in AWS

Once deployed the Lamdba can be viewed in the AWS Console under the Lambda service. Clicking on the *API Gateway* trigger in the designer section will reveal the URL you can use to test your lambda. This URL can be loaded in any browers, in a REST client like [Postman](https://chrome.google.com/webstore/detail/postman/fhbjgbiflinjbdggehcddcbncdddomop?hl=en) or via the curl command a terminal. If successful you should receive a message stating your name was added to the guestbook. You can view the guestbook in the DynamoDB table named *Fort_Wayne_AWS_User_Group_Guestbook*. You can also view logs for your Lambda in CloudWatch and can trace the executing of your Lambda in AWS X-Ray service.

### Cleanup

When you are done, you'll want to delete the AWS resources created for this demo.

```
make teardown
```
or
```
aws cloudformation delete-stack --stack-name Fort-Wayne-AWS-User-Group-Intro-To-AWS-Lambda --profile default --region us-west-1
aws s3 rb --force s3://intro-to-aws-lambda --profile default --region us-west-1
```

## Running the Demo Locally

This demo can be run locally on your machine. This provides you with the opportiunity to test any changes you make to the lambda before you deploy them to AWS. If you do make any code changes, be sure to recompile the demo (see instructions above). The include pre-compiled binary can also be run locally.

You will need to be running docker to run this demo on your local machine.

### Start Local DynamoDB container

This step will be to be performed in a separate terminal as the command prompt is help with the container is running

```
make dynamodb
```
or
```
docker run -p 8000:8000 amazon/dynamodb-local
```

### Creating the local DynamoDB table

This step creates the table in your locally running DynamoDB

```
make create_table
```
or
```
aws dynamodb create-table --cli-input-json file://create_table.json --endpoint-url http://localhost:8000/
```

### Starting API Gateway and Lambda Locally

This step will spin up both AWS API Gateway and Lamba on your local machine. It will output the URL needed to test the locally running version of the demo. Like in AWS, this can be done via a web browser, a REST client of curl on the command line.

```
make local
```
or
```
sam local start-api
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
