package modules

import (
	"fmt"
	"github.com/ginx-contribs/ginx-server/internal/common/types"
	"github.com/ginx-contribs/ginx-server/internal/modules/system"
	"github.com/google/wire"
	"reflect"
	"sync"
)

// Module defines a minimum set of methods that custom module need to implement
type Module interface {
	// Name return module name
	Name() string
	// Init should do something initial.
	Init(injector types.Injector) error
	// Close should release and clean up the module resource
	Close() error
}

// Provider is provider for all modules
var Provider = wire.NewSet(
	system.Provider,
	wire.Struct(new(Modules), "*"),
)

// Modules holds all app modules
type Modules struct {
	System system.Module
}

func NewModuleManager(modules *Modules) *ModuleManager {
	return &ModuleManager{Modules: modules}
}

// ModuleManager manager app modules, it is not thread safe.
type ModuleManager struct {
	Modules *Modules
	Mods    []Module

	once sync.Once
}

func (m *ModuleManager) lazyInit() {
	m.once.Do(func() {
		modRef := reflect.ValueOf(m.Modules).Elem()
		numMods := modRef.NumField()
		TypeModule := reflect.TypeOf(new(Module)).Elem()
		// collect mods
		var mods []Module
		for i := range numMods {
			field := modRef.Field(i).Type()
			if !field.Implements(TypeModule) {
				panic(fmt.Sprintf("was not implement interface modules.Module: %s", field.Name()))
			}
			mods = append(mods, modRef.Field(i).Interface().(Module))
		}
		m.Mods = mods
	})
}

// AllMods return all app modules
func (m *ModuleManager) AllMods() ([]Module, error) {
	m.lazyInit()
	return m.Mods, nil
}

// Init init all modules
func (m *ModuleManager) Init(injector types.Injector) error {
	m.lazyInit()
	for _, module := range m.Mods {
		err := module.Init(injector)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close closes all modules
func (m *ModuleManager) Close() error {
	m.lazyInit()
	for _, module := range m.Mods {
		err := module.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
