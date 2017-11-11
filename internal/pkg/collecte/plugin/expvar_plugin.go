//  Copyright jean-françois PHILIPPE 2014-2016

package plugin

import (
	"github.com/jfphilippe/gocollecte/internal/pkg/collecte"
	"github.com/jfphilippe/gocollecte/internal/pkg/config"
	// expvar will declare HttpHandler
	_ "expvar"
	"log"
	"net"
	"net/http"
)

const (
	// ListenAddr default listen Addr.
	ListenAddr string = "localhost:8123"
)

// ExpVarPlugin plugin to open expvar http handler
type ExpVarPlugin struct {
	done     chan struct{}
	listener net.Listener
}

// NewExpVarPlugin Un PluginCtor
// Chargé de créer une instance de Plugin
func NewExpVarPlugin(conf *config.Config, core *collecte.Collecte) (collecte.Plugin, error) {
	laddr, err := conf.String("listen", ListenAddr)
	if err != nil {
		return nil, err
	}

	// Open socket
	sock, err := net.Listen("tcp", laddr)
	if err != nil {
		return nil, err
	}
	return &ExpVarPlugin{done: core.Done, listener: sock}, nil
}

// Lancé par Start.
// Lit les commandes envoyes, les valide et les transmet au chan de commandes.
func (p *ExpVarPlugin) run() {
	http.Serve(p.listener, nil)
}

// Stop Arret du composant
// Mais se base sur le chan done pour arreter la goroutine
func (p *ExpVarPlugin) Stop() error {
	log.Println("ExpVarPlugin::Stop")
	p.listener.Close() // Ignore Error
	return nil
}

// Start Demarre le composant
func (p *ExpVarPlugin) Start() error {
	log.Println("ExpVarPlugin::Start")
	go p.run()
	return nil
}

// Enregistre le PluginCtor
func init() {
	collecte.Register("expvar", NewExpVarPlugin)
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
