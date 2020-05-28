package requestpatterncondition

import (
	"github.com/herb-go/herb/service/httpservice/requestmatching"

	"github.com/herb-go/herb/middleware/middlewarefactory"
)

//NewConditionFactory create new requestrule condition
func NewConditionFactory() middlewarefactory.ConditionFactory {
	return func(loader func(v interface{}) error) (middlewarefactory.Condition, error) {
		c := &requestmatching.PatternAllConfig{}
		err := loader(c)
		if err != nil {
			return nil, err
		}
		return c.CreatePattern()
	}
}
