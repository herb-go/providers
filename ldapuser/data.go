package ldapuser

import "github.com/herb-go/herb/cache/cacheablemap"

type ProfileProvider struct {
	Config *Config
	Fields []string
}

type Profiles struct {
	Provider *ProfileProvider
	profile  map[string]map[string][]string
}

func (p *Profiles) NewMapElement(s string) error {
	p.profile[s] = map[string][]string{}
	return nil
}

//LoadMapElements method load element to map by give key list.
//Return any error if raised.
func (p *Profiles) LoadMapElements(keys ...string) error {
	l, err := p.Provider.Config.DialBound()
	if err != nil {
		return err
	}
	defer l.Close()
	for _, v := range keys {
		data, err := p.Provider.Config.search(l, v, p.Provider.Fields...)
		if err != nil {
			return err
		}
		p.profile[v] = data

	}
	return nil
}

//Map return cachable map
func (p *Profiles) Map() interface{} {
	return p.profile
}

func (p *ProfileProvider) Create() cacheablemap.Map {
	return &Profiles{
		Provider: p,
		profile:  map[string]map[string][]string{},
	}
}
