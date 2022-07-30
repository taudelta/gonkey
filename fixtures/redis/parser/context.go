package parser

type Context struct {
    Fixtures []LoadedFixture
    KeyRefs  map[string]Keys
    MapRefs  map[string]MapRecordValue
    SetRefs  map[string]SetRecordValue
}

func NewContext() *Context {
    return &Context{
        KeyRefs: make(map[string]Keys),
        MapRefs: make(map[string]MapRecordValue),
        SetRefs: make(map[string]SetRecordValue),
    }
}
