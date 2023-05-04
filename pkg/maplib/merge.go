package maplib

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/orderedmap"
)

func NewMap(vs ...any) *orderedmap.OrderedMap {
	m := orderedmap.New()
	for i := 0; i < len(vs); i += 2 {
		m.Set(vs[i].(string), vs[i+1])
	}
	return m
}

func Merge[T any](dst *orderedmap.OrderedMap, src *T) (*orderedmap.OrderedMap, error) {
	return MergeWith(dst, src, func(f reflect.StructField) string {
		if v, ok := f.Tag.Lookup("json"); ok {
			return v
		}
		return f.Name
	})
}

func MergeWith[T any](dst *orderedmap.OrderedMap, src *T, namer func(reflect.StructField) string) (*orderedmap.OrderedMap, error) {
	if dst == nil {
		dst = orderedmap.New()
	}
	rv := reflect.ValueOf(src)
	rt := reflect.TypeOf(src)
	if err := merge(dst, "", rt.Elem(), rv.Elem(), false, namer); err != nil {
		return dst, err
	}
	return dst, nil
}

func merge(dst *orderedmap.OrderedMap, k string, rt reflect.Type, rv reflect.Value, omitempty bool, namer func(reflect.StructField) string) error {
	// log.SetFlags(0)
	// log.Printf("@ kind:%s\tdst:%v\tk:%s\trt:%v\trv:%s", rt.Kind(), dst, k, rt, rv)

	switch rt.Kind() {
	case reflect.Bool:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, rv.Bool())
	case reflect.Int:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, int(rv.Int()))
	case reflect.Int8:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, int8(rv.Int()))
	case reflect.Int16:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, int16(rv.Int()))
	case reflect.Int32:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, int32(rv.Int()))
	case reflect.Int64:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, int64(rv.Int()))
	case reflect.Uint:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, uint(rv.Uint()))
	case reflect.Uint8:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, uint8(rv.Uint()))
	case reflect.Uint16:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, uint16(rv.Uint()))
	case reflect.Uint32:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, uint32(rv.Uint()))
	case reflect.Uint64:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, uint64(rv.Uint()))
	case reflect.Uintptr:
		dst.Set(k, uint64(rv.Uint()))
	case reflect.Float32:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, float32(rv.Float()))
	case reflect.Float64:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, float64(rv.Float()))
	case reflect.Complex64:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, complex64(rv.Complex()))
	case reflect.Complex128:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, complex128(rv.Complex()))
	case reflect.String:
		if omitempty && rv.IsZero() {
			return nil
		}
		dst.Set(k, rv.String())
	case reflect.Array, reflect.Slice:
		if rv.IsNil() {
			return nil
		}
		if rv.Len() == 0 {
			dst.Set(k, reflect.MakeSlice(rt, 0, 0).Interface()) // nil?
		}

		st := rt.Elem()
		for st.Kind() == reflect.Pointer {
			st = st.Elem()
		}

		switch st.Kind() {
		case reflect.Func, reflect.Chan:
			return nil // skip
		case reflect.Interface:
			isHetero := true
			r := make([]any, rv.Len())
			for i, n := 0, rv.Len(); i < n; i++ {
				sv := rv.Index(i)
				st := sv.Type()
				if err := add(r, i, st, sv, isHetero, namer); err != nil {
					return fmt.Errorf("[%d]%w", i, err)
				}
			}
			dst.Set(k, r)
			return nil
		case reflect.Array, reflect.Slice:
			isHetero := false
			r := make([]any, rv.Len())
			st := rt.Elem()
			for i, n := 0, rv.Len(); i < n; i++ {
				sv := rv.Index(i)
				if err := add(r, i, st, sv, isHetero, namer); err != nil {
					return fmt.Errorf("[%d]%w", i, err)
				}
			}
			dst.Set(k, r)
			return nil
		case reflect.Map, reflect.Struct:
			r := make([]*orderedmap.OrderedMap, rv.Len())
			st := rt.Elem()
			for i, n := 0, rv.Len(); i < n; i++ {
				m := orderedmap.New()
				r[i] = m
				if err := merge(m, "", st, rv.Index(i), false, namer); err != nil {
					return fmt.Errorf("[%d]%w", i, err)
				}
			}
			dst.Set(k, r)
			return nil
		default:
			dst.Set(k, rv.Interface()) // should be slices of primitive type
			return nil
		}
	case reflect.Map:
		m := dst
		if k != "" {
			if v, ok := m.Get(k); ok {
				if omap, ok := v.(*orderedmap.OrderedMap); ok {
					m = omap
				} else {
					m = orderedmap.New()
					if err := merge(m, "", reflect.TypeOf(v), reflect.ValueOf(v), false, namer); err != nil {
						return fmt.Errorf(".%s%w", k, err)
					}
					dst.Set(k, m)
				}
			} else {
				m = orderedmap.New()
				dst.Set(k, m)
			}
		}
		iter := rv.MapRange()
		for iter.Next() {
			sk := iter.Key().String() // map[string] only?
			sv := iter.Value()
			if err := merge(m, sk, sv.Type(), sv, false, namer); err != nil {
				return fmt.Errorf("[%s]%w", sk, err)
			}
		}
		return nil
	case reflect.Func, reflect.Chan:
		return nil // skip interface
	case reflect.Struct:
		m := dst
		if k != "" {
			if v, ok := m.Get(k); ok {
				if omap, ok := v.(*orderedmap.OrderedMap); ok {
					m = omap
				} else {
					m = orderedmap.New()
					if err := merge(m, "", reflect.TypeOf(v), reflect.ValueOf(v), false, namer); err != nil {
						return fmt.Errorf(".%s%w", k, err)
					}
					dst.Set(k, m)
				}
			} else {
				m = orderedmap.New()
				dst.Set(k, m)
			}
		}

		if rt == rOMapType {
			src := rv.Addr().Interface().(*orderedmap.OrderedMap)
			for _, sk := range src.Keys() {
				sv, ok := src.Get(sk)
				if !ok {
					continue
				}
				rsv := reflect.ValueOf(sv)
				if err := merge(m, sk, rsv.Type(), rsv, false, namer); err != nil {
					return fmt.Errorf("[%s]%w", sk, err)
				}
			}
			return nil
		}

		for i, n := 0, rt.NumField(); i < n; i++ {
			f := rt.Field(i)
			if !f.IsExported() {
				continue
			}

			// handling `json:"<name>,omitempty"`
			name := namer(f)
			omitempty := false
			if v, suffix, found := strings.Cut(name, ","); found {
				name = strings.TrimSpace(v)
				if strings.Contains(suffix, "omitempty") {
					omitempty = true
				}
			}
			if name == "-" {
				continue
			}

			if err := merge(m, name, f.Type, rv.Field(i), omitempty, namer); err != nil {
				return fmt.Errorf(".%s%w", name, err)
			}
		}
	case reflect.Pointer:
		if rv.IsNil() {
			if !omitempty {
				dst.Set(k, nil)
			}
			return nil
		}
		return merge(dst, k, rt.Elem(), rv.Elem(), omitempty, namer)
	case reflect.Interface:
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
		return merge(dst, k, rv.Type(), rv, omitempty, namer)
	default: // reflect.Invalid, reflect.UnsafePointer
		return fmt.Errorf("invalid type: %s: %v", rt, k)
	}
	return nil
}

func add(dst []any, i int, rt reflect.Type, rv reflect.Value, isHetero bool, namer func(reflect.StructField) string) error {
	switch rt.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String: // primitive types
		dst[i] = rv.Interface()
	case reflect.Array, reflect.Slice:
		if rv.IsNil() {
			dst[i] = nil
			return nil
		}
		if rv.Len() == 0 {
			dst[i] = []any{} // fill nil?
			return nil
		}

		st := rt.Elem()
		for st.Kind() == reflect.Pointer {
			st = st.Elem()
		}

		switch st.Kind() {
		case reflect.Func, reflect.Chan:
			return nil // skip
		case reflect.Interface:
			isHetero := true
			r := make([]any, rv.Len())
			for i, n := 0, rv.Len(); i < n; i++ {
				sv := rv.Index(i)
				st := rv.Type()
				if err := add(r, i, st, sv, isHetero, namer); err != nil {
					return fmt.Errorf("[%d]%w", i, err)
				}
			}
			dst[i] = r
			return nil
		case reflect.Array, reflect.Slice:
			isHetero := false
			r := make([]any, rv.Len())
			st := rt.Elem()
			for i, n := 0, rv.Len(); i < n; i++ {
				sv := rv.Index(i)
				if err := add(r, i, st, sv, isHetero, namer); err != nil {
					return fmt.Errorf("[%d]%w", i, err)
				}
			}
			dst[i] = r
			return nil
		case reflect.Map, reflect.Struct:
			r := make([]*orderedmap.OrderedMap, rv.Len())
			st := rt.Elem()
			for i, n := 0, rv.Len(); i < n; i++ {
				m := orderedmap.New()
				r[i] = m
				if err := merge(m, "", st, rv.Index(i), false, namer); err != nil {
					return fmt.Errorf("[%d]%w", i, err)
				}
			}
			dst[i] = r
			return nil
		default: // primitive types
			dst[i] = rv.Interface()
		}
	case reflect.Map:
		m := orderedmap.New()
		dst[i] = m
		iter := rv.MapRange()
		for iter.Next() {
			sk := iter.Key().String() // map[string] only?
			sv := iter.Value()
			if err := merge(m, sk, sv.Type(), sv, false, namer); err != nil {
				return fmt.Errorf("[%s]%w", sk, err)
			}
		}
		return nil
	case reflect.Func, reflect.Chan:
		return nil // skip interface
	case reflect.Struct:
		m := orderedmap.New()
		dst[i] = m

		if rt == rOMapType {
			src := rv.Addr().Interface().(*orderedmap.OrderedMap)
			for _, sk := range src.Keys() {
				sv, ok := src.Get(sk)
				if !ok {
					continue
				}
				rsv := reflect.ValueOf(sv)
				if err := merge(m, sk, rsv.Type(), rsv, false, namer); err != nil {
					return fmt.Errorf("[%s]%w", sk, err)
				}
			}
			return nil
		}

		for i, n := 0, rt.NumField(); i < n; i++ {
			f := rt.Field(i)
			if !f.IsExported() {
				continue
			}

			// handling `json:"<name>,isHetero"`
			name := namer(f)
			omitempty := false
			if v, suffix, found := strings.Cut(name, ","); found {
				name = strings.TrimSpace(v)
				if strings.Contains(suffix, "omitempty") {
					omitempty = true
				}
			}
			if name == "-" {
				continue
			}

			if err := merge(m, name, f.Type, rv.Field(i), omitempty, namer); err != nil {
				return fmt.Errorf(".%s%w", name, err)
			}
		}
	case reflect.Pointer:
		if rv.IsNil() {
			return nil
		}
		return add(dst, i, rt.Elem(), rv.Elem(), isHetero, namer)
	case reflect.Interface:
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
		return add(dst, i, rv.Type(), rv, isHetero, namer)
	default: // reflect.Invalid, reflect.UnsafePointer
		return fmt.Errorf("invalid type: %s: %d", rt, i)
	}
	return nil
}

var rOMapType = reflect.TypeOf((*orderedmap.OrderedMap)(nil)).Elem()
