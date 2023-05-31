package maplib

import "github.com/iancoleman/orderedmap"

func GetOrCreate(d *orderedmap.OrderedMap, k string) (*orderedmap.OrderedMap, bool) {
	if v, ok := d.Get(k); ok {
		return v.(*orderedmap.OrderedMap), ok
	}

	v := orderedmap.New()
	d.Set(k, v)
	return v, false
}
