package maplib_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/builder/maplib"
)

func TestMerge(t *testing.T) {
	newMap := maplib.NewMap
	type person struct {
		Name     string  `json:"name"`
		Age      int     `json:"age,omitempty"`
		NickName *string `json:"nickname"`
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
			dst: nil, src: newMap("name", "foo"),
			want: newMap("name", "foo"),
		},
		{msg: "nil-map",
			dst: nil, src: map[string]string{"name": "foo"},
			want: newMap("name", "foo"),
		},
		{msg: "nil-struct",
			dst: nil, src: person{Name: "foo"},
			want: newMap("name", "foo"),
		},
		{msg: "nil-map-map",
			dst: nil, src: map[string]any{"father": map[string]string{"name": "foo"}},
			want: newMap("father", newMap("name", "foo")),
		},
		{msg: "nil-int-slices",
			dst: nil, src: map[string]any{"values": []int{1, 2, 3}},
			want: newMap("values", []int{1, 2, 3}),
		},
		{msg: "nil-struct-slices",
			dst: nil, src: map[string]any{"people": []person{{Name: "foo"}, {Name: "bar"}}},
			want: newMap("people", []*orderedmap.OrderedMap{newMap("name", "foo"), newMap("name", "bar")}),
		},
		{msg: "nil-interface-slices",
			dst: nil, src: map[string]any{"people": []any{person{Name: "foo"}, &person{Name: "bar"}}},
			want: newMap("people", []any{newMap("name", "foo"), newMap("name", "bar")}),
		},
		{msg: "nil-int-slices-slices",
			dst: nil, src: map[string]any{"values-list": [][]int{{1, 2, 3}, {4, 5, 6}}},
			want: newMap("values-list", []any{[]int{1, 2, 3}, []int{4, 5, 6}}),
		},
		{msg: "nil-struct-slices-slices",
			dst: nil, src: map[string]any{"people": [][]person{{{Name: "foo"}}, {{Name: "bar"}}}},
			want: newMap("people", []any{[]*orderedmap.OrderedMap{newMap("name", "foo")}, []*orderedmap.OrderedMap{newMap("name", "bar")}}),
		},
		{msg: "empty-map-map",
			dst: newMap(), src: map[string]any{"father": map[string]string{"name": "foo"}},
			want: newMap("father", newMap("name", "foo")),
		},
		{msg: "map-override-map",
			dst: newMap("name", "foo"), src: map[string]string{"name": "bar"},
			want: newMap("name", "bar"),
		},
		{msg: "map-override-omap",
			dst: newMap("name", "foo"), src: newMap("name", "bar"),
			want: newMap("name", "bar"),
		},
		{msg: "map-append-omap",
			dst: newMap("name", "foo"), src: newMap("age", 20),
			want: newMap("name", "foo", "age", 20),
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.msg, func(t *testing.T) {
			got, err := maplib.Merge(c.dst, &c.src)
			if err != nil {
				t.Errorf("unexpected error")
			}
			// json.NewEncoder(os.Stdout).Encode(newMap("want", c.want, "got", got))
			if diff := cmp.Diff(c.want, got, opt); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
