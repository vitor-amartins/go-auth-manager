package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type AuthManager struct {
	client       *cognitoidentityprovider.CognitoIdentityProvider
	userPoolId   string
	clientId     string
	clientSecret string
}

func getAttributeValue(attributes []*cognitoidentityprovider.AttributeType, attribute string) (string, error) {
	var err error
	for i := 0; i < len(attributes); i++ {
		if *attributes[i].Name == attribute {
			return *attributes[i].Value, err
		}
	}
	return "", &MappedError{
		StatusCode: 500,
		Message:    fmt.Sprintf("Couldn't find %s on user's attributes", attribute),
		ErrorCode:  "ERR.5.0002",
	}
}

func (authManager *AuthManager) GetUserId(email string) (string, error) {
	res, err := authManager.client.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &authManager.userPoolId,
		Username:   &email,
	})
	if err != nil {
		return "", err
	}
	return getAttributeValue(res.UserAttributes, "sub")
}

func (authManager *AuthManager) GetUserGroups(email string) ([]string, error) {
	res, err := authManager.client.AdminListGroupsForUser(&cognitoidentityprovider.AdminListGroupsForUserInput{
		UserPoolId: &authManager.userPoolId,
		Username:   &email,
	})
	g := []string{}
	if err != nil {
		return g, err
	}
	for i := 0; i < len(res.Groups); i++ {
		g = append(g, *res.Groups[i].GroupName)
	}
	return g, err
}

func (authManager *AuthManager) GetUserAuthInfo(email string) (UserAuthInfo, error) {
	var uai UserAuthInfo
	userOutput, err := authManager.client.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &authManager.userPoolId,
		Username:   &email,
	})
	g, err := authManager.GetUserGroups(email)
	if err != nil {
		return uai, err
	}
	if err != nil {
		return uai, err
	}
	id, err := getAttributeValue(userOutput.UserAttributes, "sub")
	if err != nil {
		return uai, err
	}
	f, err := getAttributeValue(userOutput.UserAttributes, "given_name")
	if err != nil {
		return uai, err
	}
	l, err := getAttributeValue(userOutput.UserAttributes, "family_name")
	if err != nil {
		return uai, err
	}
	return UserAuthInfo{
		Id:        id,
		Email:     email,
		FirstName: f,
		LastName:  l,
		Groups:    g,
	}, err
}
