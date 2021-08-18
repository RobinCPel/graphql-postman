package introspection

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const request = `{"operationName":"IntrospectionQuery","variables":{},"query":"query IntrospectionQuery {\n  __schema {\n    queryType {\n        name\n    }\n    mutationType {\n        name\n    }\n    types {\n      name\n      fields(includeDeprecated: false) {\n        name\n        args {\n          name\n          type {\n            ...TypeRef\n          }\n        }\n        type {\n          ...TypeRef\n        }\n      }\n      inputFields {\n        name\n        type {\n          ...TypeRef\n        }\n      }\n      enumValues(includeDeprecated: false) {\n        name\n      }\n      possibleTypes {\n        ...TypeRef\n      }\n    }\n  }\n}\n\nfragment TypeRef on __Type {\n  kind\n  name\n  ofType {\n    kind\n    name\n    ofType {\n      kind\n      name\n      ofType {\n        kind\n        name\n        ofType {\n          kind\n          name\n          ofType {\n            kind\n            name\n            ofType {\n              kind\n              name\n              ofType {\n                kind\n                name\n              }\n            }\n          }\n        }\n      }\n    }\n  }\n}\n"}`

// Introspect introspects a graphql endpoint and returns the result in structs.
func Introspect(url string) (*Model, error) {
	r, err := http.Post(url, "application/json", strings.NewReader(request))
	if err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.New("response status is not 200 (OK)")
	}
	if r.Body == nil {
		return nil, errors.New("response body is empty")
	}

	// Convert the data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var model Model
	if err = json.Unmarshal(body, &model); err != nil {
		return nil, err
	}

	return &model, nil
}
