package prototype_test

import (
	"encoding/json"
	"os"

	"github.com/podhmo/gos/prototype"
)

func ExampleToSchema() {
	b := prototype.NewBuilder()

	Name := prototype.Define("Name", b.String().MinLength(1))

	prototype.Define("Person", b.Object(
		b.Field("name", b.String()).Doc("name of person"),
		b.Field("age", b.Integer().Format("int32")),
		b.Field("nickname", b.Reference(Name)).Required(false),
		b.Field("father", b.ReferenceByName("Person")).Required(false),
		b.Field("friends", b.Array(b.ReferenceByName("Person"))).Required(false),
	)).Doc("person object")

	d, _ := prototype.ToSchema(b)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "@@")
	enc.Encode(d)

	// Output:
	// {
	// @@"components": {
	// @@@@"schemas": {
	// @@@@@@"Name": {
	// @@@@@@@@"type": "string",
	// @@@@@@@@"minlength": 1
	// @@@@@@},
	// @@@@@@"Person": {
	// @@@@@@@@"type": "object",
	// @@@@@@@@"description": "person object",
	// @@@@@@@@"properties": {
	// @@@@@@@@@@"name": {
	// @@@@@@@@@@@@"type": "string",
	// @@@@@@@@@@@@"description": "name of person"
	// @@@@@@@@@@},
	// @@@@@@@@@@"age": {
	// @@@@@@@@@@@@"type": "integer",
	// @@@@@@@@@@@@"format": "int32"
	// @@@@@@@@@@},
	// @@@@@@@@@@"nickname": {
	// @@@@@@@@@@@@"$ref": "#/components/schemas/Name"
	// @@@@@@@@@@},
	// @@@@@@@@@@"father": {
	// @@@@@@@@@@@@"$ref": "#/components/schemas/Person"
	// @@@@@@@@@@},
	// @@@@@@@@@@"friends": {
	// @@@@@@@@@@@@"type": "array",
	// @@@@@@@@@@@@"items": {
	// @@@@@@@@@@@@@@"$ref": "#/components/schemas/Person"
	// @@@@@@@@@@@@}
	// @@@@@@@@@@}
	// @@@@@@@@},
	// @@@@@@@@"required": [
	// @@@@@@@@@@"name",
	// @@@@@@@@@@"age"
	// @@@@@@@@],
	// @@@@@@@@"additionalProperties": false
	// @@@@@@}
	// @@@@}
	// @@}
	// }
}
