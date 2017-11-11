//  Copyright jean-françois PHILIPPE 2014-2016

package plugin

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/jfphilippe/gocollecte/internal/pkg/collecte"
	"github.com/jfphilippe/gocollecte/internal/pkg/config"
	"github.com/jfphilippe/gocollecte/internal/pkg/pipe"
)

// PipePlugin Plugin d'ouverture d'un pipe de communication (gestion de comandes)
type PipePlugin struct {
	done     chan struct{}
	cmd      chan string
	conn     *net.UnixConn
	filename string
}

// NewPipePlugin Un PluginCtor
// Chargé de créer une instance de Plugin
func NewPipePlugin(conf *config.Config, core *collecte.Collecte) (collecte.Plugin, error) {
	filename, err := conf.String("filename", pipe.SocketName)
	if err != nil {
		return nil, err
	}

	// Remove file if exists (prevent some already ijn use error ?)
	os.Remove(filename)

	conn, err := net.ListenUnixgram("unixgram", &net.UnixAddr{Name: filename, Net: "unixgram"})
	if err != nil {
		return nil, err
	}
	return &PipePlugin{done: core.Done, cmd: core.Cmd, conn: conn, filename: filename}, nil
}

// Lancé par Start.
// Lit les commandes envoyes, les valide et les transmet au chan de commandes.
func (p *PipePlugin) run() {
	for {
		var buf [1024]byte
		n, err := p.conn.Read(buf[:])
		if err != nil {
			if ne, ok := err.(net.Error); ok && !ne.Temporary() {
				log.Println("PopePlugin.run closed socket detected")
				return
			}
			// Error will be logged below
		} else {
			// Parse json !
			msg := pipe.Msg{}
			err = json.Unmarshal(buf[0:n], &msg)
			if err == nil {
				err = msg.Validate()
				if err == nil && p.conn != nil {
					// Send msg to cmde pipe
					p.cmd <- msg.Cmde
				}
			}
		}

		// may be non nil upon read, unmarshaling or validate.
		if err != nil {
			// Log Error, but don't stop
			log.Println("PipePlugin.run error ", err)
		}

	}
}

// Stop Arret du composant
// Mais se base sur le chan done pour arreter la goroutine
func (p *PipePlugin) Stop() error {
	log.Println("PipePlugin::Stop")
	p.conn.Close() // Ignore Error
	p.conn = nil
	os.Remove(p.filename)
	return nil
}

// Start Demarre le composant
func (p *PipePlugin) Start() error {
	log.Println("PipePlugin::Start")
	go p.run()
	return nil
}

// Enregistre le PluginCtor
func init() {
	collecte.Register("pipe", NewPipePlugin)
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
