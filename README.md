# ✉️ GraphQL Postman

Converts a GraphQL schema to a Postman Collection v2.1 which can be used in GitLab CI for API Fuzzing Tests.

## 🖊 Postman Collection Format v2.1

https://schema.postman.com/json/collection/v2.1.0/docs/index.html

## 🚀 How to run

1. Have [go](https://golang.org/) installed
2. Have [make](https://www.gnu.org/software/make/) installed
3. Build with `make build`
4. Run with `ENDPOINT="http://GRAPHQL_ENDPOINT" make run`

Note: step three and four can be combined with the command: `ENDPOINT="http://GRAPHQL_ENDPOINT" make full`.

## 🚩 Flags

| Name                    | Description                         | Flag        | Default                                | Required |
|-------------------------|-------------------------------------|-------------|----------------------------------------|----------|
| GraphQL Endpoint        | GraphQL endpoint to connect to.     | `-endpoint` | -                                      | yes      |
| Output File             | The file to write the result to.    | `-output`   | `api.postman_collection.json`          | no       |
| Postman Collection ID   | The Postman Collection ID to use.   | `-id`       | `00000000-0000-0000-0000-000000000000` | no       |
| Postman Collection Name | The Postman Collection name to use. | `-name`     | `GraphQL Postman`                      | no       |

## 🐳 Docker
The image of this project is available on docker hub: <https://hub.docker.com/r/robincp/graphql-postman>

## ⚠️ Known issues

- Does not support lists with more than two dimensions
- Assumes a schema has queries *and* mutations, no subscriptions
- No support for interfaces, objects, and unions

## 🚨 Important note
This software was written to make automated GitLab API Fuzzing testing possible for our GraphQL API. The features this project contains, are limited to what our GraphQL API consists out of. Therefore, the known issues will not be fixed unless they become relevant for us (or if a very nice person comes around and opens a merge/pull request with the features 😉).

## ⚙️ How to use it for GitLab API fuzzing

1. Add a CI job that runs before the GitLab Fuzzing job 
2. Make the new job use the graphql-postman docker image
3. Add this line to the script: `/go/src/bin/graphql-postman -endpoint "${FUZZAPI_TARGET_URL}/gql"`
4. Expose the artifact, by default called `api.postman_collection.json`

### 📑 Example

```yaml
include:
  - template: Security/API-Fuzzing.gitlab-ci.yml

stages:
  - test
  - build
  - etc.
  - prepare-fuzz
  - fuzz

variables:
  FUZZAPI_PROFILE: Long-100
  FUZZAPI_POSTMAN_COLLECTION: ./api.postman_collection.json
  FUZZAPI_TARGET_URL: http://example.com

...

prepare-fuzz:
  stage: prepare-fuzz
  image:
    name: gitlab.example.org/pace/graphql-postman:master
  before_script: []
  script:
    - /go/src/bin/graphql-postman -endpoint "${FUZZAPI_TARGET_URL}/gql"
  after_script: []
  artifacts:
    expire_in: 1 week
    name: "$CI_COMMIT_REF_NAME_postman_collection"
    expose_as: "postman_collection"
    paths:
      - api.postman_collection.json

apifuzzer_fuzz:
  needs:
    - job: prepare-fuzz
      artifacts: true
```

### 🔨 Building the docker image yourself

You can also build the docker image yourself and push it to your GitLab Docker Image repo with CI, [how?](https://docs.gitlab.com/ee/ci/docker/using_kaniko.html#building-a-docker-image-with-kaniko) 

## 🧠 GraphQL Introspection Query

```graphql
query IntrospectionQuery {
  __schema {
    queryType {
      name
    }
    mutationType {
      name
    }
    types {
      name
      fields(includeDeprecated: false) {
        name
        args {
          name
          type {
            ...TypeRef
          }
        }
        type {
          ...TypeRef
        }
      }
      inputFields {
        name
        type {
          ...TypeRef
        }
      }
      enumValues(includeDeprecated: false) {
        name
      }
      possibleTypes {
        ...TypeRef
      }
    }
  }
}

fragment TypeRef on __Type {
  kind
  name
  ofType {
    kind
    name
    ofType {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
              }
            }
          }
        }
      }
    }
  }
}
```
