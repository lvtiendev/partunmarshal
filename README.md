## Partially unmarshal Golang structs

In Golang, structs are usually used to marshal / unmarshal from different encoding such as `JSON`, `XML`, etc. This package provides a way to specfiy updatable struct's fields and only unmarshal encoding values into these fields.

Imagine we are building a management service for fruits:

```go
type fruit struct {
	Name      string    `json:"name"`
	Color     string    `json:"color"`
}
```

On an update request, we usually unmarshal the client payload into the struct to update its fields:

```go
f := &fruit{
    Name: "watermelon"
}

payload := `{"color":"green","name":"cucumber"}`
_ = json.Unmarshal(byte(payload), f)
// f.Name == "cucumber"
// f.Color == "green"
```

However, what if we want to keep `name` intact regardless of client input ?

With the package, we can specfiy only `color` field to be updatable by annotation tag `u:"true"`.

```go
type fruit struct {
	Name      string    `json:"name"`
	Color     string    `json:"color" u:"true"`
}

f := &fruit{
    Name: "watermelon"
}

payload := `{"color":"green","name":"cucumber"}`
_ = partunmarshal.JSON(byte(payload), f)
// f.Name == "watermelon"
// f.Color == "green"
```

### TODO

- Nested structs support
