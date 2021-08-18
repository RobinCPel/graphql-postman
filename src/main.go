package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/RobinCPel/graphql-postman/src/internal/graphql/introspection"
	"github.com/RobinCPel/graphql-postman/src/internal/graphql/introspection/reformatted"
	"github.com/RobinCPel/graphql-postman/src/internal/graphql/kind"
	"github.com/RobinCPel/graphql-postman/src/internal/graphql/scalar"
	"github.com/RobinCPel/graphql-postman/src/internal/postman"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"strings"
)

var types map[string]reformatted.Type

func init() {
	log.SetLevel(log.DebugLevel)
}

func getDummyValueOfScalar(scalarName string) string {
	switch strings.ToLower(scalarName) {
	case scalar.Integer:
		return `4200`
	case scalar.Float:
		return `96.031`
	case scalar.String:
		return `"This is a test!1!"`
	case scalar.Boolean:
		return `true`
	case scalar.ID:
		return `"V2llRGl0TGVlc3RJc0dlaw=="`
	}

	log.Warning(`scalar "` + scalarName + `" does not exist, using NULL as dummy value`)
	return `NULL`
}

func getDummyValueOfType(typeName, typeKind string) (string, error) {
	t, ok := types[typeName]
	if !ok {
		return "", errors.New(`could not find the type "` + typeName + `" in the types map`)
	}

	emptyResponse := func() (string, error) {
		return "NULL", nil
	}

	switch strings.ToLower(typeKind) {

	case kind.Scalar: // Has only a name
		return getDummyValueOfScalar(t.Name), nil

	case kind.InputObject: // Has only a name and input fields
		dummyValue := `{`
		var count int
		for key, val := range t.InputFields {
			count++

			typeDummyVal, err := getDummyValueOfType(val.Name, val.Kind)
			if err != nil {
				log.Fatal(`Could not get the dummy value of type with name "` + val.Name + `"`)
			}

			if val.List {
				if val.TwoDList {
					dummyValue += `"` + key + `":[[` + typeDummyVal + `]]`
				} else {
					dummyValue += `"` + key + `":[` + typeDummyVal + `]`
				}
			} else {
				dummyValue += `"` + key + `":` + typeDummyVal
			}

			if count != len(t.InputFields) {
				dummyValue += `,`
			}
		}
		dummyValue += `}`
		return dummyValue, nil

	case kind.Enum: // Has only a name and enum values
		// Return a random element from the enum values slice
		if len(t.EnumValues) > 0 {
			return `"` + t.EnumValues[rand.Intn(len(t.EnumValues))] + `"`, nil
		} else {
			log.WithField("name", t.Name).Warning("Found an enum without values")
			return emptyResponse()
		}

	case kind.Interface: // Has only a name, fields, and possible types
		log.Warning("Interfaces are only present in mutation and query responses, so a dummy value for an interface is not needed")
		return emptyResponse()

	case kind.Object: // Has only a name and fields
		log.Warning("Objects are only present in mutation and query responses, so a dummy value for an object is not needed")
		return emptyResponse()

	case kind.Union: // Has only a name and possible types
		log.Warning("Unions are only present in mutation and query responses, so a dummy value for a union is not needed")
		return emptyResponse()

	default:
		log.Fatal(`type of kind "` + typeKind + `" does not exist`)
	}

	return "", errors.New(`type of kind "` + typeKind + `" does not exist`)
}

func gqlInputFromOperation(o reformatted.Operation, operationName string) (*postman.GqlInput, error) {
	input := postman.GqlInput{Name: o.Name}

	// Assemble the query
	var count int
	var argLine1, argLine2 string
	for k, v := range o.Arguments {
		count++

		argLine1 += `$` + k + `: ` + v.Name
		if v.NonNull {
			argLine1 += `!`
		}

		argLine2 += k + `: $` + k

		// Not the last one? Add the delimiter
		if count != len(o.Arguments) {
			argLine1 += `, `
			argLine2 += `, `
		}
	}
	input.Query = operationName + ` ` + o.Name + `(` + argLine1 + `) { ` + o.Name + `(` + argLine2 + `) { __typename } }`

	// Assemble the dummy variables
	count = 0
	for k, v := range o.Arguments {
		count++

		if v.List {
			log.Warn("list graphql operation arguments are not supported")
		}

		input.Variables += `{"` + k + `":`

		dummyVal, err := getDummyValueOfType(v.Name, v.Kind)
		if err != nil {
			return nil, err
		}
		input.Variables += dummyVal

		if count == len(o.Arguments) {
			// Last one, close off the json
			input.Variables += `}`
		} else {
			// Not the last one? Add the delimiter
			input.Variables += `, `
		}
	}

	return &input, nil
}

func main() {
	fmt.Print(`
                                       (     
                                       (        (                                       
 (                            )   (    )\ )     )\ )             )                      
 )\ )    (       )         ( /( ( )\  (()/(    (()/(          ( /(    )       )         
(()/(    )(   ( /(  '  )   )\()))((_)  /(_))    /(_)) (   (   )\())  (     ( /(   (     
 /(_))_ (()\  )(_)) /(/(  ((_)\((_)_  (_))     (_))   )\  )\ (_))/   )\  ' )(_))  )\ )  
(_)) __| ((_)((_)_ ((_)_\ | |(_)/ _ \ | |      | _ \ ((_)((_)| |_  _((_)) ((_)_  _(_/(  
  | (_ || '_|/ _' || '_ \)| ' \| (_) || |__    |  _// _ \(_-<|  _|| '  \()/ _' || ' \)) 
   \___||_|  \__,_|| .__/ |_||_|\__\_\|____|   |_|  \___//__/ \__||_|_|_| \__,_||_||_|  
                   |_|

`)

	// Define flags
	var url, outputFileName, postmanCollectionID, postmanCollectionName string
	flag.StringVar(&url, "endpoint", "", "graphql endpoint to connect to")
	flag.StringVar(&outputFileName, "output", "api.postman_collection.json", "the file to write the result to")
	flag.StringVar(&postmanCollectionID, "id", "00000000-0000-0000-0000-000000000000", "the Postman Collection ID to use")
	flag.StringVar(&postmanCollectionName, "name", "GraphQL Postman", "the Postman Collection name to use")
	flag.Parse()

	// Check if the endpoint is defined
	if url == "" {
		log.Fatal(`an endpoint needs to be specified with the flag "-endpoint"`)
	}

	// Introspect
	log.Info("Running the GraphQL Introspection...")
	raw, err := introspection.Introspect(url)
	if err != nil {
		log.WithError(err).Fatal("could not introspect the graphql endpoint")
	}

	log.Info("Reformatting the GraphQL Introspected models...")
	model := reformatted.Reformat(raw)
	if model == nil {
		log.Fatal("reformatted graphql introspection model equals nil")
		return
	}

	// Copy the types so they can be accessed globally
	types = model.Types

	gqlInputs := make([]postman.GqlInput, 0, len(model.Mutations)+len(model.Queries))

	// Convert the mutations
	log.Info("Converting the mutations...")
	for _, m := range model.Mutations {

		gqlInput, err := gqlInputFromOperation(m, "mutation")
		if err != nil {
			log.WithField("name", m.Name).WithError(err).
				Warning("failed to convert a mutation to a GQL Input, skipping")
			continue
		}

		gqlInputs = append(gqlInputs, *gqlInput)
	}

	// Convert Queries
	log.Info("Converting the queries...")
	for _, q := range model.Queries {
		gqlInput, err := gqlInputFromOperation(q, "query")
		if err != nil {
			log.WithField("name", q.Name).WithError(err).
				Warning("failed to convert a query to a GQL Input, skipping")
			continue
		}

		gqlInputs = append(gqlInputs, *gqlInput)
	}

	// Convert GQL Inputs to postman collection
	log.Info("Storing the mutations and queries in a postman collection...")
	col := postman.CreateCollection(gqlInputs, postmanCollectionID, postmanCollectionName)
	data, err := json.MarshalIndent(col, "", "    ")
	if err != nil {
		log.WithError(err).Fatal("failed to encode the postman collection as json")
	}

	if err = ioutil.WriteFile(outputFileName, data, 0644); err != nil {
		log.WithError(err).Fatal("failed to write the postman collection to a file")
	}
	log.Info(`Writing the result to "` + outputFileName + `"...`)
	log.Info("All done!")
	fmt.Println()
}
