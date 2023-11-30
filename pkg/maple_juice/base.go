package maple_juice

type KV struct {
	Key   string
	Value interface{}
}

/*
Struct for maple stage
*/
type MapleFunc = func(*KV) (*KV, error)

/*
	Struct for juice stage
*/

type JuiceFunc = func([]*KV) (*KV, error)
