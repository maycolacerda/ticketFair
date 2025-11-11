package graphql

import (
	"log"

	"github.com/graphql-go/graphql"
	"github.com/maycolacerda/ticketfair/database"
	"github.com/maycolacerda/ticketfair/models"
	"gorm.io/gorm"
)

// userType define a estrutura do objeto User no Schema GraphQL.
var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"Userid": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "User ID",
			},
			"UserEmail": &graphql.Field{
				Type:        graphql.String,
				Description: "User Email",
			},
			"Username": &graphql.Field{
				Type:        graphql.String,
				Description: "User Name",
			},
			"Password": &graphql.Field{
				Type:        graphql.String,
				Description: "User Password",
			},
		},
	},
)

// rootUser define as Queries (leituras) para o Schema.
var rootUser = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootUser",
	Fields: graphql.Fields{
		// QUERY para buscar um único usuário por ID
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"Userid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				idQuery, ok := p.Args["Userid"].(string)
				if !ok {
					return nil, nil
				}
				var user models.User
				if err := database.DB.First(&user, "Userid = ?", idQuery).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return nil, nil
					}
					log.Println("Error fetching user:", err)
					return nil, err
				}
				return user, nil
			},
		},
		// QUERY para buscar todos os usuários
		"users": &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(p graphql.ResolveParams) (any, error) {
				var users []models.User
				if err := database.DB.Find(&users).Error; err != nil {
					log.Println("Error fetching users:", err)
					return nil, err
				}
				return users, nil
			},
		},
	},
})

// userMutation define as Mutations (escritas) para o Schema.
var userMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "userMutation",
		Fields: graphql.Fields{
			// MUTATION: Cria um novo usuário
			"createUser": &graphql.Field{
				Type:        userType,
				Description: "Cria um novo usuário no banco de dados.",
				Args: graphql.FieldConfigArgument{
					"userID":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"username": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					newUser := models.User{
						UserID:   p.Args["userID"].(string),
						Email:    p.Args["email"].(string),
						Password: p.Args["password"].(string),
						Username: p.Args["username"].(string),
					}

					if err := database.DB.Create(&newUser).Error; err != nil {
						return nil, err
					}

					return newUser, nil
				},
			},
			// MUTATION: Atualiza um usuário existente (apenas campos fornecidos)
			"updateUser": &graphql.Field{
				Type:        userType,
				Description: "Atualiza um usuário existente no banco de dados. Campos email, password e username são opcionais.",
				Args: graphql.FieldConfigArgument{
					"userID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					// Campos opcionais:
					"email":    &graphql.ArgumentConfig{Type: graphql.String},
					"password": &graphql.ArgumentConfig{Type: graphql.String},
					"username": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					userID := p.Args["userID"].(string)
					var existingUser models.User

					// 1. Encontra o usuário existente
					if err := database.DB.First(&existingUser, "Userid = ?", userID).Error; err != nil {
						return nil, err // Retorna erro se não encontrar
					}

					// 2. Constrói um mapa apenas com os campos que foram realmente fornecidos
					updates := make(map[string]any)
					if email, ok := p.Args["email"].(string); ok && email != "" {
						updates["UserEmail"] = email
					}
					if password, ok := p.Args["password"].(string); ok && password != "" {
						updates["Password"] = password
					}
					if username, ok := p.Args["username"].(string); ok && username != "" {
						updates["Username"] = username
					}

					// Se não houver nada para atualizar, retorna o usuário existente
					if len(updates) == 0 {
						return existingUser, nil
					}

					// 3. Executa a atualização parcial usando GORM Updates
					if err := database.DB.Model(&existingUser).Where("Userid = ?", userID).Updates(updates).Error; err != nil {
						return nil, err
					}

					// 4. Recarrega (ou garante) que o usuário retornado tenha os novos dados
					// Note: O GORM Updates geralmente atualiza a struct `existingUser`
					return existingUser, nil

				},
			},
			// MUTATION: Deleta um usuário existente
			"deleteUser": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Deleta um usuário existente pelo ID.",
				Args: graphql.FieldConfigArgument{
					"userID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					userID := p.Args["userID"].(string)
					// Tenta deletar o usuário
					result := database.DB.Delete(&models.User{}, "Userid = ?", userID)
					if result.Error != nil {
						return false, result.Error
					}
					// Retorna true se pelo menos uma linha foi afetada
					return result.RowsAffected > 0, nil
				},
			},
		},
	},
)

// Define o Schema completo
var schemaUser, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    rootUser,
		Mutation: userMutation,
	},
)

// QueryUser executa a query GraphQL usando o Schema definido.
func QueryUser(params graphql.Params) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schemaUser,
		RequestString: params.RequestString,
	})
	if len(result.Errors) > 0 {
		log.Println("Error executing query:", result.Errors)
	}
	return result
}
