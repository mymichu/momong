package server_todo

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"text": &graphql.Field{
			Type: graphql.String,
		},
		"done": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

// root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createTodo": &graphql.Field{
			Type: todoType, // the return type for this field
			Args: graphql.FieldConfigArgument{
				"text": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},

			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				// marshall and cast the argument value
				text, _ := params.Args["text"].(string)

				// perform mutation operation here
				// for e.g. create a Todo and save to DB.
				newTodo := &Todo{
					ID:   "id0001",
					Text: text,
					Done: false,
				}

				// return the new Todo object that we supposedly save to DB
				// Note here that
				// - we are returning a `Todo` struct instance here
				// - we previously specified the return Type to be `todoType`
				// - `Todo` struct maps to `todoType`, as defined in `todoType` ObjectConfig`
				return newTodo, nil
			},
		},
	},
})

// root query
// we just define a trivial example here, since root query is required.
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"lastTodo": &graphql.Field{
			Type: todoType,
		},
	},
})

// define schema
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func Start() {

	// define custom GraphQL ObjectType `todoType` for our Golang struct `Todo`
	// Note that
	// - the fields in our todoType maps with the json tags for the fields in our struct
	// - the field type matches the field type in our struct

	/*if err != nil {
		panic(err)
	}*/

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	// serve HTTP
	http.Handle("/graphql1", h)
	http.ListenAndServe(":12345", nil)

	// How to make a HTTP request using cUrl
	// -------------------------------------
	// In `graphql-go-handler`, based on the GET/POST and the Content-Type header, it expects the input params differently.
	// This behaviour was ported from `express-graphql`.
	//
	//
	// 1) using GET
	// $ curl -g -GET 'http://localhost:12345/graphql?query=mutation+M{newTodo:createTodo(text:"This+is+a+todo+mutation+example"){text+done}}'
	//
	// 2) using POST + Content-Type: application/graphql
	// $ curl -XPOST http://localhost:12345/graphql -H 'Content-Type: application/graphql' -d 'mutation M { newTodo: createTodo(text: "This is a todo mutation example") { text done } }'
	//
	// 3) using POST + Content-Type: application/json
	// $ curl -XPOST http://localhost:12345/graphql -H 'Content-Type: application/json' -d '{"query": "mutation M { newTodo: createTodo(text: \"This is a todo mutation example\") { text done } }"}'
	//
	// Any of the above would return the following output:
	// {
	//   "data": {
	// 	   "newTodo": {
	// 	     "done": false,
	// 	     "text": "This is a todo mutation example"
	// 	   }
	//   }
	// }
}
