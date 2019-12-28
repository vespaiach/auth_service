package share

//Direction type will be used in sorting
type Direction int

//Define sorting type
const (
	BiDirection Direction = iota
	Ascendant
	Descendant
)

//Boolean type allows nil value
type Boolean struct {
	IsSet bool
	Bool  bool
}

//DefaultLimit is default value of record per page
const DefaultLimit = 10
