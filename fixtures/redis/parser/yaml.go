package parser

import (
    "errors"
    "fmt"
    "io/ioutil"

    "gopkg.in/yaml.v3"
)

type redisYamlParser struct {
    rootParser *fileParser
}

func (p *redisYamlParser) Copy(rootParser *fileParser) FixtureFileParser {
    cp := &(*p)
    cp.rootParser = rootParser
    return cp
}

func (p *redisYamlParser) buildTemplate(ctx *Context, f Fixture) error {
    for refName, tplData := range f.Templates.Keys {
        if _, ok := ctx.KeyRefs[refName]; ok {
            return fmt.Errorf("unable to load template %s: duplicating ref name", refName)
        }
        if tplData.Extend != "" {
            baseRecord, err := p.resolveKeyReference(ctx.KeyRefs, tplData.Extend)
            if err != nil {
                return err
            }
            for k, v := range tplData.Values {
                baseRecord.Values[k] = v
            }
            tplData.Values = baseRecord.Values
        }

        keyRef := Keys{
            Values: make(map[string]*KeyValue, len(tplData.Values)),
        }
        for k, v := range tplData.Values {
            keyRef.Values[k] = v
        }
        ctx.KeyRefs[refName] = keyRef
    }

    for refName, tplData := range f.Templates.Sets {
        if _, ok := ctx.SetRefs[refName]; ok {
            return fmt.Errorf("unable to load template %s: duplicating ref name", refName)
        }
        if tplData.Extend != "" {
            baseRecord, err := p.resolveSetReference(ctx.SetRefs, tplData.Extend)
            if err != nil {
                return err
            }
            for k, v := range tplData.Values {
                baseRecord.Values[k] = v
            }
            tplData.Values = baseRecord.Values
        }

        setRef := SetRecordValue{
            Values: make(map[string]*SetValue),
        }
        for k, v := range tplData.Values {
            var valueCopy *SetValue
            if v != nil {
                valueCopy = &(*v)
            }
            setRef.Values[k] = valueCopy
        }
        ctx.SetRefs[refName] = setRef
    }

    for refName, tplData := range f.Templates.Maps {
        if _, ok := ctx.MapRefs[refName]; ok {
            return fmt.Errorf("unable to load template %s: duplicating ref name", refName)
        }
        if tplData.Extend != "" {
            baseRecord, err := p.resolveMapReference(ctx.MapRefs, tplData.Extend)
            if err != nil {
                return err
            }
            for k, v := range tplData.Values {
                baseRecord.Values[k] = v
            }
            tplData.Values = baseRecord.Values
        }
        mapRef := MapRecordValue{
            Values: make(map[string]string, len(tplData.Values)),
        }
        for k, v := range tplData.Values {
            mapRef.Values[k] = v
        }
        ctx.MapRefs[refName] = mapRef
    }

    return nil
}

func (p *redisYamlParser) resolveKeyReference(refs map[string]Keys, refName string) (*Keys, error) {
    refTemplate, ok := refs[refName]
    if !ok {
        return nil, fmt.Errorf("ref not found: %s", refName)
    }
    keysCopy := &Keys{
        Values: make(map[string]*KeyValue),
    }
    for k, v := range refTemplate.Values {
        keysCopy.Values[k] = v
    }
    return keysCopy, nil
}

func (p *redisYamlParser) resolveMapReference(refs map[string]MapRecordValue, refName string) (*MapRecordValue, error) {
    refTemplate, ok := refs[refName]
    if !ok {
        return nil, fmt.Errorf("ref not found: %s", refName)
    }
    copy_ := &MapRecordValue{
        Values: make(map[string]string),
    }
    for k, v := range refTemplate.Values {
        copy_.Values[k] = v
    }
    return copy_, nil
}

func (p *redisYamlParser) resolveSetReference(refs map[string]SetRecordValue, templateName string) (*SetRecordValue, error) {
    refTemplate, ok := refs[templateName]
    if !ok {
        return nil, errors.New("ref not found")
    }
    copy_ := &SetRecordValue{
        Values: make(map[string]*SetValue),
    }
    for k, v := range refTemplate.Values {
        var setValue *SetValue
        if v != nil {
            setValue = &(*v)
        }
        copy_.Values[k] = setValue
    }
    return copy_, nil
}

func (p *redisYamlParser) buildKeys(ctx *Context, data *Keys) error {
    if data == nil {
        return nil
    }
    if data.Extend != "" {
        baseRecord, err := p.resolveKeyReference(ctx.KeyRefs, data.Extend)
        if err != nil {
            return err
        }
        for k, v := range data.Values {
            var keyValue *KeyValue
            if v != nil {
                keyValue = &(*v)
            }
            baseRecord.Values[k] = keyValue
        }
        data.Values = baseRecord.Values
    }
    return nil
}

func (p *redisYamlParser) buildMaps(ctx *Context, data *Maps) error {
    if data == nil {
        return nil
    }
    for _, v := range data.Values {
        if v.Extend != "" {
            baseRecord, err := p.resolveMapReference(ctx.MapRefs, v.Extend)
            if err != nil {
                return err
            }
            for k, v := range v.Values {
                baseRecord.Values[k] = v
            }
            v.Values = baseRecord.Values
        }
        if v.Name != "" {
            mapRef := MapRecordValue{
                Values: make(map[string]string, len(v.Values)),
            }
            for k, v := range v.Values {
                mapRef.Values[k] = v
            }
            ctx.MapRefs[v.Name] = mapRef
        }
    }
    return nil
}

func (p *redisYamlParser) buildSets(ctx *Context, data *Sets) error {
    if data == nil {
        return nil
    }
    for _, v := range data.Values {
        if v.Extend != "" {
            baseRecord, err := p.resolveSetReference(ctx.SetRefs, v.Extend)
            if err != nil {
                return err
            }
            for k, v := range v.Values {
                baseRecord.Values[k] = v
            }
            v.Values = baseRecord.Values
        }
        if v.Name != "" {
            setRef := SetRecordValue{
                Values: make(map[string]*SetValue),
            }
            for k, v  := range v.Values {
                var setValue *SetValue
                if v != nil {
                    setValue = &(*v)
                }
                setRef.Values[k] = setValue
            }
            ctx.SetRefs[v.Name] = setRef
        }
    }
    return nil
}

func (p *redisYamlParser) Parse(ctx *Context, filename string) (*Fixture, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var fixture Fixture
    if err := yaml.Unmarshal(data, &fixture); err != nil {
        return nil, err
    }

    for _, parentFixture := range fixture.Inherits {
        _, err := p.rootParser.ParseFiles(ctx, []string{parentFixture})
        if err != nil {
            return nil, err
        }
    }

    if err = p.buildTemplate(ctx, fixture); err != nil {
        return nil, err
    }

    for _, databaseData := range fixture.Databases {
        if err := p.buildKeys(ctx, databaseData.Keys); err != nil {
            return nil, err
        }
        if err := p.buildMaps(ctx, databaseData.Maps); err != nil {
            return nil, err
        }
        if err := p.buildSets(ctx, databaseData.Sets); err != nil {
            return nil, err
        }
    }

    return &fixture, nil
}
