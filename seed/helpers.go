package seed

import (
	"strings"
	"sync"
)

type TypeVarMetadataList []*TypeVarMetadata

var pool = &sync.Pool{
	New: func() any {
		var b strings.Builder
		return &b
	},
}

func (tvs TypeVarMetadataList) Names() string {
	b := pool.Get().(*strings.Builder)
	defer pool.Put(b)
	b.Reset()
	for _, tv := range tvs {
		b.WriteString(tv.Name)
		b.WriteString(", ")
	}
	return b.String()
}
func (tvs TypeVarMetadataList) Types() string {
	b := pool.Get().(*strings.Builder)
	defer pool.Put(b)
	b.Reset()
	for _, tv := range tvs {
		b.WriteString(string(tv.Type))
		b.WriteString(", ")
	}
	return b.String()
}
func (tvs TypeVarMetadataList) NameAndTypes() string {
	b := pool.Get().(*strings.Builder)
	defer pool.Put(b)
	b.Reset()
	for _, tv := range tvs {
		b.WriteString(tv.Name)
		b.WriteRune(' ')
		b.WriteString(string(tv.Type))
		b.WriteString(", ")
	}
	return b.String()
}
