## Schemas

- Issue: https://github.com/k14s/ytt/issues/103
- Status: **Being written** | Being implemented | Included in release | Rejected

### Examples

- Schemas for [cf-for-k8s's values](cf-for-k8s/values.yml)
  - [JSON schema (as json)](cf-for-k8s/json-schema.json)
  - [JSON schema (as yaml)](cf-for-k8s/json-schema.yml)
  - [ytt native schema](cf-for-k8s/ytt-schema.yml)

### Defining a schema document and what type of documents to apply it to

```yaml
#@schema attach="data/values"
---
#! Schema contents
```

Will create a new schema document. The attach keyword is used to specify what types of documents the schema will apply to, in this case, `data/values` documents.
The attach argument will default to `data/values` if none is provided.

### Schema annotations

#### Basic annotations

- `@schema/type`
  This annotation asserts on the type of a keys value. For example,
  ```yaml
  #@schema/type "array"
  app_domains: []
  ```
  will validate that any data value `app_domain` will be of type array. These strings are a predefined set: string, array, int, etc...
  ytt will also infer the type based on the key's value given in a schema document if the annotation is not provided. For example,
  ```yaml
  app_domains: []
  ```
  will also result in the `app_domains` key requiring value type array

- `@schema/validate`
  This annotation can be used in order to run a validation function on the value of the key it is applied to. For example,
  ```yaml
  #@schema/validate number_is_even, min=2
  replicas: 6
  ```
  will use the `min` predefined validator as well as the user provided `number_is_even` function to validate the value of `replicas`. The funtion signature should match what gets passed to overlay/assert annotations.

- `@schema/allow-empty`
  This annotation asserts that a key can have an empty value of its type i.e. "", [], 0, etc. This is useful for defining keys that are optional and could be avoided as data values.
  ```yaml
  #@schema/allow-empty
  system_domain: ""
  ```
  Because the schema values are extracted and used as defaults, by setting the value of `sytem_domain` to allow empty string, the user is allowed to provide a empty `system_domain` through their data values (or not provide it at all since it's an empty string by default). ytt uses starlark [type truth values](https://github.com/google/starlark-go/blob/master/doc/spec.md#data-types) to determine if a value is empty.

#### Describing annotations

- `@schema/title -> title for node (can maybe infer from the key name ie app_domain -> App domain)`
  This annotation provides a way to add a short title to the schema applied to a key.
  ```yaml
  #@schema/title "User Password"
  user_password: ""
  ```
  If the annotation is not present, the title will be infered from the key name by replacing special characters with spaces and capitalizing the first letter similar to rails humanize functionality.

- `@schema/doc`
  This annotation is a way to add a longer description of the schema applied to a key, similar to the JSON Schema description field.
  ```yaml
  #@schema/doc "The user password used to log in to the system"
  user_password: ""
  ```

- `@schema/example`
  Examples will take one arguments which consist of the example
  ```yaml
  #@schema/example "my_domain.example.com"
  system_domain: ""
  ```
  In this example, the example string "my_domain.example.com" will be attached to the `system_domain` key.

- `@schema/examples`
  Examples will take one or more tuple arguments which consist of {Title, Example value}
  ```yaml
  #@schema/examples ("Title 1", value1), ("Title 2", title2_example())
  foo: bar
  ```
  In this example, the Title 1 and Title 2 examples and their values are attached with the key `foo`.

#### Map key presence annotations

- `@schema/any-key`
  Applies to items of a map. Allows users to assert on structure while maintaining freedom to have any key
  Allows a schema to assert on the structure of map items without making any assertions on the value of the key
  ```yaml
  connection_options:
    #@schema/any-key
    _: 
    - ""
  ```
  This example will allow the map to contain any key names that is an array of strings.

- `@schema/key-may-be-present`
  This annotation asserts a key is allowed by the schema but is not guaranteed. This allows schemas to validate contents of a structure in cases where
  the contents are not referenced directly. For example,
  ```yaml
  connection_options:
    #@schema/key-may-be-present
    pooled: true
  ```
  will assert that the key `pooled` is allowed under the schema but not guaranteed. This would be useful when accessing something like #@ data.values.connection_options in a template
  instead of #@ data.values.connection_options.pooled. See more advanced examples below for more.

#### Complex schema annotations

- `@schema/any-of`
  Requires the key to satisy _at least one_ of the provided schemas
  ```yaml
  #@schema/any-of schema1(), schema2()
  foo: ""
  ```
  Note, any values passed to the `any-of` call will be interpreted as schema documents.

- `@schema/all-of`
  Requires the key to satisy _every_ provided schema
  ```yaml
  #@schema/all-of schema1(), schema2()
  foo: ""
  ```
  Note, any values passed to the `all-of` call will be interpreted as schema documents.

- `@schema/one-of`
  Requires the key to satisy _exactly one_ of provided schema
  ```yaml
  #@schema/one-of schema1(), schema2()
  foo: ""
  ```
  Note, any values passed to the `one-of` call will be interpreted as schema documents.

### Sequence of events
1. Extract defaults from the provided schemas
2. 

assert on higher level structure with optional lower key
```yaml
#@schema
---
foo: ""

#@data/values
---
foo: "val"

config.yml:
---
foo: #@ data.values.foo #! never errors with: no member 'foo' on data.values
```

```yaml
#@schema
---
#@schema/key-may-be-there
foo:

#@data/values
---
foo: "val"

config.yml:
---
foo: #@ yaml.encode(data.values) #! => {}, {"foo": "val"}
```

```yaml
#@schema
---
#@schema/key-may-be-there
foo:
  max_connections: 100
  username: ""

#@data/values
---
foo:
  username: val

config.yml:
---
foo: #@ yaml.encode(data.values) #! => {},  {"foo": {"username": "val", "max_connections":100}}
```
