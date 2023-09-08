package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const (
	userPoolID = "ap-northeast-1_XXXXXXX"
	roleARN    = "arn:aws:iam::9999999999:role/cognito-assume-role"
	externalID = "sample"
	awsRegion  = "ap-northeast-1"
)

func main() {
	lambda.Start(handler)
}

// newCognitoClient creates a client for Cognito.
func newCognitoClient(cfg *aws.Config) *cognitoidentityprovider.Client {
	// Create a client for Cognito
	return cognitoidentityprovider.NewFromConfig(*cfg)
}

// newCognitoClientWithAssumeRole creates a client for Cognito with AssumeRole.
func newCognitoClientWithAssumeRole(cfg *aws.Config) *cognitoidentityprovider.Client {
	// Create a client for STS
	stsClient := sts.NewFromConfig(*cfg)

	// Create an AssumeRoleProvider
	assumeRoleProvider := stscreds.NewAssumeRoleProvider(stsClient, roleARN, func(o *stscreds.AssumeRoleOptions) {
		o.ExternalID = aws.String(externalID)
	})

	cfg.Credentials = aws.NewCredentialsCache(assumeRoleProvider)

	// Create a client for Cognito
	return cognitoidentityprovider.NewFromConfig(*cfg)
}

// listUsers outputs a list of users.
func listUsers(ctx context.Context, client *cognitoidentityprovider.Client) error {
	// Get a list of users
	users, err := client.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(userPoolID),
	})
	if err != nil {
		return err
	}

	// Output a list of users
	for _, user := range users.Users {
		fmt.Printf("user: %+v\n", user)
	}

	return nil

}

// handler lambda handler.
func handler(ctx context.Context) {
	// Create a config by specifying the region
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v\n", err)
		return
	}

	// Create a client for Cognito with AssumeRole
	client := newCognitoClientWithAssumeRole(&cfg)

	// Get a list of users from Cognito
	err = listUsers(ctx, client)
	if err != nil {
		fmt.Printf("Failed to list users: %v\n", err)
		return
	}

	return
}
