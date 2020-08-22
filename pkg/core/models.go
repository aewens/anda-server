package core

type SQLSelect struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value []byte `json:"value"`
	Flag  int    `json:"flag"`
}

type SQLType struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SQLEntity struct {
	Id     int    `json:"id"`
	TypeId int    `json:"type_id"`
	UUID   string `json:"uuid"`
}

type SQLAttribute struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ValueTypeId int    `json:"value_type_id"`
}

type SQLTypeAttribute struct {
	Id          int `json:"id"`
	TypeId      int `json:"type_id"`
	AttributeId int `json:"attribute_id"`
}

type SQLValueType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SQLValue struct {
	Id          int    `json:"id"`
	EntityId    int    `json:"entity_id"`
	AttributeId int    `json:"attribute_id"`
	Value       []byte `json:"value"`
}
