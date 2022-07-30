package parser

type Fixture struct {
	Name      string
	Inherits  []string           `yaml:"inherits" toml:"inherits" json:"inherits"`
	Templates DatabaseTemplate   `yaml:"templates" toml:"templates" json:"templates"`
	Databases map[int16]Database `yaml:"databases" toml:"databases" json:"databases"`
	parent    map[string]*Fixture
}

type DatabaseTemplate struct {
	Maps map[string]*MapRecordValue `yaml:"maps" toml:"maps" json:"maps"`
	Sets map[string]*SetRecordValue `yaml:"sets" toml:"sets" json:"sets"`
	Keys map[string]*Keys           `yaml:"keys" toml:"keys" json:"keys"`
}

type Database struct {
	Maps *Maps `yaml:"maps" toml:"maps" json:"maps"`
	Sets *Sets `yaml:"sets" toml:"sets" json:"sets"`
	Keys *Keys `yaml:"keys" toml:"keys" json:"keys"`
}

type MapRecordValue struct {
	Name   string            `yaml:"$name" toml:"$name" json:"$name"`
	Extend string            `yaml:"$extend" json:"$extend"`
	Values map[string]string `yaml:"values" toml:"values" json:"values"`
}

type Maps struct {
	Values map[string]*MapRecordValue `yaml:"values" json:"values"`
}

type SetValue struct {
	Expiration int `yaml:"expiration" toml:"expiration" json:"expiration"`
}

type SetRecordValue struct {
	Name   string               `yaml:"$name" toml:"$name" json:"$name"`
	Extend string               `yaml:"$extend" json:"$extend"`
	Values map[string]*SetValue `yaml:"values" toml:"values" json:"values"`
}

type Sets struct {
	Values map[string]*SetRecordValue `yaml:"values" json:"values"`
}

type KeyValue struct {
	Value      string `yaml:"value" toml:"value" json:"value"`
	Expiration int    `yaml:"expiration" toml:"expiration" json:"expiration"`
}

type Keys struct {
	Extend string               `yaml:"$extend" json:"$extend"`
	Values map[string]*KeyValue `yaml:"values" json:"values"`
}
