package prototype_test

import (
	"testing"

	"github.com/podhmo/gos/prototype"
)

func TestToString(t *testing.T) {
	b := prototype.NewBuilder()

	tests := []struct {
		name string
		typ  prototype.TypeBuilder
		want string
	}{
		{"primitive-0", b.String(), "string"},
		{"primitive-1", b.Integer(), "integer"},
		{"new-type-primitive", prototype.Define("Name", b.String()), "Name"},
		{"array-string", b.Array(b.String()), "array[string]"},
		{"array-array-string", b.Array(b.Array(b.String())), "array[array[string]]"},
		{"object", b.Object(
			b.Field("name", b.String()),
			b.Field("age", b.String()).Required(false),
		), "object{name, age?}"},
		{"new-type-object", prototype.Define("Person", b.Object(
			b.Field("name", b.String()),
			b.Field("age", b.String()).Required(false),
		)), "Person{name, age?}"},
		// action
		{"func()", b.Action(), "action ()"},
		{"func(string)int", b.Action(
			b.Input(b.Param("name", b.String())),
			b.Output(b.Integer()),
		), "action (string)integer"},
		{"func(strin,int)string", b.Action(
			b.Input(
				b.Param("name", b.String()),
				b.Param("age", b.Integer()),
			),
			b.Output(b.String()),
		), "action (string, integer)string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prototype.ToString(tt.typ)
			if got != tt.want {
				t.Errorf("ToString() = %v, but want is %v", got, tt.want)
			}
		})
	}
}
