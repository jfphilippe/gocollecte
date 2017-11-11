//Copyright jean-françois PHILIPPE 2014-2016

package collecte

import (
	"github.com/jfphilippe/gocollecte/internal/pkg/config"
	"github.com/jfphilippe/gocollecte/internal/pkg/measure"
	"log"
	"strings"
)

func (c *Collecte) loadValues(conf *config.Config) error {
	// Check for Runtime config
	var err error
	c.valuesHandler, err = measure.NewValuesHandler(conf.Section("values"), c.Done)
	if err != nil {
		return err
	}
	c.Values = c.valuesHandler.Chan
	return nil
}

// Initialise les plugin a partir des sections de la config.
// Ignore les sections dont le nom debute par #
// Chaque section doit contenir un parametre 'type'.
func (c *Collecte) loadPlugins(conf *config.Config) error {
	log.Println("Collecte::Load.Begin")
	sections := conf.Sections()
	for name, conf := range sections {
		if !strings.HasPrefix(name, "#") {
			kind, err := conf.String("type", name)
			if err == nil {
				ctor := ctors[kind]
				if ctor != nil {
					plugin, err := ctor(conf, c)
					if err == nil {
						log.Println("Collecte::Creation plugin '" + name + "' succeded")
						c.plugins = append(c.plugins, plugin)
					} else {
						log.Println("Collecte::Plugin '"+name+"' error Ctor : %v", err)
						return err
					}
				} else {
					log.Println("Collecte::Plugin '" + name + "' type '" + kind + "' unknown")
				}
			} else {
				log.Println("Collecte::Section '" + name + "' parameter type not found")
			}
		}
	}
	return nil
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
