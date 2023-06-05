package M



// name of something
type Name string



// Friends of something
type Friends []Person



type Tag struct {
    // name of tag
    Name Name `json:"name"`

    Doc string `json:"doc"`
}



// person object
type Person struct {
    // name of person
    Name Name `json:"name"`

    Age int64 `json:"age"`

    Father *Person `json:"father"`

    Children []Person `json:"children"`

    Friends Friends `json:"friends"`

    Tags map[string]Tag `json:"tags"`
}
