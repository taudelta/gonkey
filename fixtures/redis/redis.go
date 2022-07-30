package redis

import (
    "context"
    "time"

    "github.com/lamoda/gonkey/fixtures/redis/parser"

    redis "github.com/go-redis/redis/v9"
)

type loader struct {
    locations []string
    debug     bool
    client    *redis.Client
}

type LoaderOptions struct {
    Debug  bool
    Redis  *redis.Options
}

func New(location string, debug bool, opts LoaderOptions) *loader {
    client := redis.NewClient(opts.Redis)
    return &loader{
        locations: []string{location},
        client:    client,
        debug:     debug,
    }
}

func (l *loader) Load(names []string) error {
    ctx := parser.NewContext()
    fileParser := parser.New(l.locations)
    fixtureList, err := fileParser.ParseFiles(ctx, names)
    if err != nil {
        return err
    }
    return l.loadData(fixtureList)
}

func (l *loader) loadData(fixtures []*parser.Fixture) error {
    truncatedDatabases := make(map[int16]struct{})

    for _, redisFixture := range fixtures {
        for dbID, db := range redisFixture.Databases {
            pipeContext := context.Background()
            pipe := l.client.Pipeline()
            err := pipe.Select(pipeContext, int(dbID)).Err()
            if err != nil {
                return err
            }

            if _, ok := truncatedDatabases[dbID]; !ok {
                if err := pipe.FlushDB(pipeContext).Err(); err != nil {
                    return err
                }
                truncatedDatabases[dbID] = struct{}{}
            }

            if db.Keys != nil {
                for k, v := range db.Keys.Values {
                    if err := pipe.Set(pipeContext, k, v.Value, time.Duration(v.Expiration)*time.Millisecond).Err(); err != nil {
                        return err
                    }
                }
            }

            if db.Sets != nil {
                for setKey, v := range db.Sets.Values {
                    for v := range v.Values {
                        if err := pipe.SAdd(pipeContext, setKey, v).Err(); err != nil {
                            return err
                        }
                    }
                }
            }

            if db.Maps != nil {
                for mapKey, v := range db.Maps.Values {
                    for k, v := range v.Values {
                        if err := pipe.HSet(pipeContext, mapKey, k, v).Err(); err != nil {
                            return err
                        }
                    }
                }
            }

            if _, err := pipe.Exec(pipeContext); err != nil {
                return err
            }
        }
    }
    return nil
}
