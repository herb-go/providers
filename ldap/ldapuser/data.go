package ldapuser

import (
	"github.com/herb-go/herb/cache/cachestore"
)

func profileCreateor() interface{} {
	return map[string][]string{}
}
func profileLoader(c *Config, Fields ...string) func(keys ...string) (map[string]interface{}, error) {
	return func(keys ...string) (map[string]interface{}, error) {
		var result = map[string]interface{}{}
		l, err := c.DialBound()
		if err != nil {
			return result, err
		}
		defer l.Close()
		for _, v := range keys {
			data, err := c.search(l, v, Fields...)
			if err != nil {
				return result, err
			}
			result[v] = data

		}
		return result, nil

	}
}

func newProfileProvider(c *Config, Field ...string) *cachestore.DataSource {
	s := cachestore.NewDataSource()
	s.Creator = profileCreateor
	s.SourceLoader = profileLoader(c)
	return s
}
