package templates

import (
	"errors"
	"github.com/jucardi/infuse/config"
)

// ErrTypeNotFound is returned when the template type does not match a defined template implementation
var (
	ErrTypeNotFound = errors.New("type not found")

	instance IFactory
)

type factory struct {
	ctors map[string]func(...string) ITemplate
}

// Factory returns the templates factory
func Factory() IFactory {
	if instance == nil {
		instance = &factory{ctors: map[string]func(...string) ITemplate{}}
	}
	return instance
}

func (f *factory) New(name ...string) ITemplate {
	if t, err := f.Create(config.Get().DefaultType); err == nil {
		return t
	}
	for _, ctor := range f.ctors {
		return ctor(name...)
	}
	return nil
}

func (f *factory) Create(typeStr string, name ...string) (ITemplate, error) {
	ctor, ok := f.ctors[typeStr]
	if !ok {
		return nil, ErrTypeNotFound
	}
	return ctor(name...), nil
}

func (f *factory) Register(name string, constructor func(name ...string) ITemplate) {
	f.ctors[name] = func(name ...string) ITemplate {
		ret := constructor(name...)
		return &baseTemplate{
			ITemplate: ret,
		}
	}
}

func (f *factory) GetAvaliableTypes() []string {
	keys := make([]string, 0, len(f.ctors))
	for k := range f.ctors {
		keys = append(keys, k)
	}
	return keys
}

func (f *factory) Contains(typeStr string) bool {
	_, ok := f.ctors[typeStr]
	return ok
}
