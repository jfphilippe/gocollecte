// Copyright jean-françois PHILIPPE 2014-2016
// Paquet principal de l application de collecte.

package collecte

import (
	"github.com/jfphilippe/gocollecte/internal/pkg/config"
	"github.com/jfphilippe/gocollecte/internal/pkg/measure"
	"log"
	"time"
)

var (
	// Tableau blanc des contructeurs de plugin.
	ctors = make(map[string]PluginCtor, 5)
)

// Plugin interface de plugin
type Plugin interface {
	Start() error
	Stop() error
}

// PluginCtor Constructeur de Plugin
type PluginCtor func(conf *config.Config, core *Collecte) (Plugin, error)

// Register Enregistre un Constructeur sous un nom
func Register(name string, ctor PluginCtor) {
	ctors[name] = ctor
}

// Collecte Composant principal de l'application.
type Collecte struct {
	valuesHandler *measure.ValuesHandler
	plugins       []Plugin
	Done          chan struct{}
	Values        chan *measure.NodeValues
	Cmd           chan string
}

// New creation de collecte
func New(conf *config.Config) (*Collecte, error) {
	result := &Collecte{nil, make([]Plugin, 0, 5), make(chan struct{}), nil, make(chan string)}

	err := result.loadValues(conf)

	if err != nil {
		return nil, err
	}

	err = result.loadPlugins(conf.Section("plugins"))

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Start Demarre le composant de collecte
func (c *Collecte) Start() error {
	log.Println("Collecte::Start")

	// Demarre le Value_Handler
	c.valuesHandler.Start()

	// Demarre les plugin
	for _, plugin := range c.plugins {
		err := plugin.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

// Stop Arrete le composant.
func (c *Collecte) Stop() error {
	log.Println("Collecte::Stop")

	// Stoppe les plugins dan l'ordre "inverse"
	for i := len(c.plugins) - 1; i >= 0; i-- {
		c.plugins[i].Stop()
	}

	close(c.Done)
	// Laisse le temps de prendre en compte la fermeture
	time.Sleep(100 * time.Millisecond)

	// Clot le Chan des valeurs
	if c.valuesHandler != nil {
		close(c.valuesHandler.Chan)
	}

	close(c.Cmd)

	return nil
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
