module meetup.com/fortwayne/guestbook

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.3.3 // indirect
	github.com/aws/aws-lambda-go v1.13.2
	github.com/aws/aws-sdk-go v1.23.21
	// github.com/aws/aws-xray-sdk-go v0.9.4
	github.com/aws/aws-xray-sdk-go v1.0.0-rc.14 // X-Ray support for Lambda wasn't added till v1.0.0-rc.1
	golang.org/x/net v0.0.0-20190912160710-24e19bdeb0f2 // indirect
)
