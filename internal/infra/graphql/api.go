package graphql

import (
	"net/url"
)

type introspection_response struct {
	schema introspection_schema `json:"__schema"`
}

type introspection_schema struct {
	query        introspection_root_type  `json:"queryType"`
	mutation     *introspection_root_type `json:"mutationType"`
	subscription *introspection_root_type `json:"subscriptionType"`
	types        []introspection_output_type
	directives   []introspection_directive
}

type introspection_root_type struct {
	name string
}

type introspection_directive struct {
	name        string
	description *string
	locations   []string
	args        []introspection_input_value
}

type introspection_input_value struct {
	name         string
	description  *string
	defaultValue interface{}
}

type introspection_output_type struct {
	name        string
	description *string
	fields      []introspection_field
	interfaces  []introspection_interface
}

type introspection_field struct {
}

type introspection_interface struct {
}

type graphql_api struct {
	endpoint url.URL
}

const introspectionQuery = `
query IntrospectionQuery {
  __schema {
    queryType {
      name
    }
    mutationType {
      name
    }
    subscriptionType {
      name
    }
    types {
      ...FullType
    }
    directives {
      name
      description
      locations
      args {
        ...InputValue
      }
    }
  }
}

fragment FullType on __Type {
  kind
  name
  description
  fields(includeDeprecated: true) {
    name
    description
    args {
      ...InputValue
    }
    type {
      ...TypeRef
    }
    isDeprecated
    deprecationReason
  }
  inputFields {
    ...InputValue
  }
  interfaces {
    ...TypeRef
  }
  enumValues(includeDeprecated: true) {
    name
    description
    isDeprecated
    deprecationReason
  }
  possibleTypes {
    ...TypeRef
  }
}

fragment InputValue on __InputValue {
  name
  description
  type {
    ...TypeRef
  }
  defaultValue
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
`
