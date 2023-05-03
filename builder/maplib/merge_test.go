package maplib_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/builder/maplib"
)

func TestMerge(t *testing.T) {
	oMap := maplib.NewMap
	type person struct {
		Name     string  `json:"name"`
		Age      int     `json:"age,omitempty"`
		NickName *string `json:"nickname,omitempty"`
		Father   *person `json:"father,omitempty"`
	}

	type withIgnore struct {
		Name       string `json:"name"`
		Ignored    string `json:"-"`
		unexported string
		NOJSONTag  string
	}

	opt := cmp.Transformer("", func(src *orderedmap.OrderedMap) map[string]any {
		keys := src.Keys()
		dst := make(map[string]any, len(keys))
		for _, k := range keys {
			v, _ := src.Get(k)
			dst[k] = v
		}
		return dst
	})

	cases := []struct {
		msg  string
		dst  *orderedmap.OrderedMap
		src  any
		want *orderedmap.OrderedMap
	}{
		{msg: "nil-omap",
			dst: nil, src: oMap("name", "foo"),
			want: oMap("name", "foo"),
		},
		{msg: "nil-map",
			dst: nil, src: map[string]string{"name": "foo"},
			want: oMap("name", "foo"),
		},
		{msg: "nil-struct",
			dst: nil, src: person{Name: "foo"},
			want: oMap("name", "foo"),
		},
		{msg: "nil-map-map",
			dst: nil, src: map[string]any{"father": map[string]string{"name": "foo"}},
			want: oMap("father", oMap("name", "foo")),
		},
		{msg: "nil-int-slices",
			dst: nil, src: map[string]any{"values": []int{1, 2, 3}},
			want: oMap("values", []int{1, 2, 3}),
		},
		{msg: "nil-struct-slices",
			dst: nil, src: map[string]any{"people": []person{{Name: "foo"}, {Name: "bar"}}},
			want: oMap("people", []*orderedmap.OrderedMap{oMap("name", "foo"), oMap("name", "bar")}),
		},
		{msg: "nil-interface-slices",
			dst: nil, src: map[string]any{"people": []any{person{Name: "foo"}, &person{Name: "bar"}}},
			want: oMap("people", []any{oMap("name", "foo"), oMap("name", "bar")}),
		},
		{msg: "nil-int-slices-slices",
			dst: nil, src: map[string]any{"values-list": [][]int{{1, 2, 3}, {4, 5, 6}}},
			want: oMap("values-list", []any{[]int{1, 2, 3}, []int{4, 5, 6}}),
		},
		{msg: "nil-struct-slices-slices",
			dst: nil, src: map[string]any{"people": [][]person{{{Name: "foo"}}, {{Name: "bar"}}}},
			want: oMap("people", []any{[]*orderedmap.OrderedMap{oMap("name", "foo")}, []*orderedmap.OrderedMap{oMap("name", "bar")}}),
		},
		{msg: "empty-map-map",
			dst: oMap(), src: map[string]any{"father": map[string]string{"name": "foo"}},
			want: oMap("father", oMap("name", "foo")),
		},
		{msg: "map-override-map",
			dst: oMap("name", "foo"), src: map[string]string{"name": "bar"},
			want: oMap("name", "bar"),
		},
		{msg: "map-override-omap",
			dst: oMap("name", "foo"), src: oMap("name", "bar"),
			want: oMap("name", "bar"),
		},
		{msg: "map-override-person",
			dst: oMap("name", "foo"), src: person{Name: "bar"},
			want: oMap("name", "bar"),
		},
		{msg: "map-override-person",
			dst: oMap("name", "foo"), src: &person{Name: "bar"},
			want: oMap("name", "bar"),
		},
		{msg: "map-override-slices",
			dst: oMap("values", []int{1, 2, 3}), src: oMap("values", []int{10}),
			want: oMap("values", []int{10}),
		},
		{msg: "map-override-slices-struct",
			dst: oMap("people", []person{{Name: "foo"}}), src: oMap("people", []person{{Name: "foo"}, {Name: "bar"}}),
			want: oMap("people", []*orderedmap.OrderedMap{oMap("name", "foo"), oMap("name", "bar")}),
		},
		{msg: "map-append-omap",
			dst: oMap("name", "foo"), src: oMap("age", 20),
			want: oMap("name", "foo", "age", 20),
		},
		// nested
		{msg: "nested-map-append-omap",
			dst: oMap("name", "foo", "father", oMap("name", "boo")), src: oMap("age", 20, "father", oMap("age", 40)),
			want: oMap("name", "foo", "age", 20, "father", oMap("name", "boo", "age", 40)),
		},
		{msg: "nested-map-append-person",
			dst: oMap("age", 20, "father", oMap("age", 40)), src: person{Name: "foo", Father: &person{Name: "boo"}},
			want: oMap("name", "foo", "age", 20, "father", oMap("name", "boo", "age", 40)),
		},
		// tags
		{msg: "arrange by tags",
			dst: nil, src: withIgnore{},
			want: oMap("name", "", "NOJSONTag", ""), // ignore the value of json tag is "-" or unexported field
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.msg, func(t *testing.T) {
			got, err := maplib.Merge(c.dst, &c.src)
			if err != nil {
				t.Errorf("unexpected error")
			}
			// json.NewEncoder(os.Stdout).Encode(oMap("want", c.want, "got", got))
			if diff := cmp.Diff(c.want, got, opt); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
