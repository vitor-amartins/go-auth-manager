package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/joho/godotenv"
)

type UserAuthInfo struct {
	Id        string
	Email     string
	FirstName string
	LastName  string
	Groups    []string
}

func main() {
	godotenv.Load()

	awsRegion := os.Getenv("REGION")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	if err != nil {
		log.Fatal(err)
	}

	authManager := AuthManager{
		client:       cognitoidentityprovider.New(sess, aws.NewConfig().WithRegion(awsRegion)),
		userPoolId:   os.Getenv("USER_POOL_ID"),
		clientId:     os.Getenv("CLIENT_ID"),
		clientSecret: os.Getenv("CLIENT_SECRET"),
	}

	email := "vitormartins@gravidadezero.space"

	var uid string
	uid, err = authManager.GetUserId(email)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User id: %s\n", uid)

	var g []string
	g, err = authManager.GetUserGroups(email)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User groups: %v\n", g)

	var uai UserAuthInfo
	uai, err = authManager.GetUserAuthInfo(email)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User auth info: %v\n", uai)
}
