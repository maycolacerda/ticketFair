package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/models"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"userID":   &graphql.Field{Type: graphql.String},
		"email":    &graphql.Field{Type: graphql.String},
		"username": &graphql.Field{Type: graphql.String},
	},
})
var userQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "userQuery",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type:        userType,
			Description: "Get user by ID",
			Args: graphql.FieldConfigArgument{
				"userID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				UserID, ok := p.Args["userID"].(string)
				if !ok {
					return nil, nil
				}
				var user models.User
				if err := database.DB.First(&user, "user_id = ?", UserID).Error; err != nil {
					return nil, err
				}
				return user, nil

			},
		},
	},
})

var userMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "userMutation",
	Fields: graphql.Fields{
		"createUser": &graphql.Field{
			Type:        userType,
			Description: "Create new user",
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			username, _ := p.Args["username"].(string)
			email, _ := p.Args["email"].(string)
			password, _ := p.Args["password"].(string)
			user := models.User{
				Username: username,
				Email:    email,
				Password: password,
			}
			if err := database.DB.Create(&user).Error; err != nil {
				return nil, err
			}
			return user, nil
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    userQuery,
	Mutation: userMutation,
})
