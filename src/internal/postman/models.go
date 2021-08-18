package postman

type Info struct {
	PostManID string `json:"_postman_id"` // Configurable, by default: "00000000-0000-0000-0000-000000000000"
	Name      string `json:"name"`        // Configurable, by default: "GraphQL Postman"
	Schema    string `json:"schema"`      // Always "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
}

type Graphql struct {
	Query     string `json:"query"`     // The actual query
	Variables string `json:"variables"` // The variables
}

type Body struct {
	Mode    string  `json:"mode"` // Always "graphql"
	GraphQL Graphql `json:"graphql"`
}

type Url struct {
	Raw      string   `json:"raw"`      // Always "http://localhost/gql"
	Protocol string   `json:"protocol"` // Always "http"
	Host     []string `json:"host"`     // Always "localhost"
	Path     []string `json:"path"`     // Always ["gql"]
}

type Request struct {
	Method string        `json:"method"` // Always "POST"
	Header []interface{} `json:"header"` // Always empty
	Body   Body          `json:"body"`
	URL    Url           `json:"url"`
}

type Item struct {
	Name     string        `json:"name"` // Name of the GQL query
	Request  Request       `json:"request"`
	Response []interface{} `json:"response"` // Always empty
}

// Collection represents a Postman Collection v2.1.
type Collection struct {
	Info Info   `json:"info"`
	Item []Item `json:"item"`
}

// GqlInput examples:
//
// Name:      "CreateShip" => doesn't really matter tbh
// Query:	  "mutation CreateShip($input: CreateShipInput!) { createShip(input: $input) { __typename }"
// Variables: `{"input":{"name":"anything","speed":3}}`
//
// Name:      "Node"
// Query:     "query Node($id: ID!) { node(id: $id) { __typename } }"
// Variables: `{"id": "anything"}`
type GqlInput struct {
	Name      string
	Query     string
	Variables string
}

// createItems takes GqlInput data and stuffs it into a Postman Collection item.
func createItems(gql []GqlInput) []Item {
	items := make([]Item, len(gql), len(gql))

	for i, entry := range gql {
		items[i] = Item{
			Name: entry.Name,
			Request: Request{
				Method: "POST",
				Header: []interface{}{},
				Body: Body{
					Mode: "graphql",
					GraphQL: Graphql{
						Query:     entry.Query,
						Variables: entry.Variables,
					},
				},
				URL: Url{
					Raw:      "http://localhost/gql",
					Protocol: "http",
					Host:     []string{"localhost"},
					Path:     []string{"gql"},
				},
			},
			Response: []interface{}{},
		}
	}

	return items
}

// CreateCollection returns a collection with all of the default values already set.
func CreateCollection(gql []GqlInput, postmanID, name string) Collection {
	return Collection{
		Info: Info{
			PostManID: postmanID,
			Name:      name,
			Schema:    "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item: createItems(gql),
	}
}
