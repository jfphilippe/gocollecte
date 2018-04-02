//  Copyright jean-françois PHILIPPE 2014-2018

package plugin

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jfphilippe/gocollecte/internal/pkg/collecte"
	"github.com/jfphilippe/gocollecte/internal/pkg/config"
)

// SignalPlugin plugin de gestion des signaux
type SignalPlugin struct {
	done    chan struct{}
	cmd     chan string
	signals chan os.Signal
}

// NewSignalPlugin Un PluginCtor
// Chargé de créer une instance de Plugin
func NewSignalPlugin(conf *config.Config, core *collecte.Collecte) (collecte.Plugin, error) {
	return &SignalPlugin{done: core.Done, cmd: core.Cmd, signals: make(chan os.Signal)}, nil
}

// Goroutine lancee par Start
func (p *SignalPlugin) run() {
	signal.Notify(p.signals, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGTSTP)
	for {
		select {
		case s := <-p.signals:
			switch s {
			case syscall.SIGHUP:
				p.cmd <- "log_rotate"
			case syscall.SIGUSR1:
				p.cmd <- "memdump"
			case syscall.SIGUSR2:
				p.cmd <- "nope"
			default:
				p.cmd <- "stop"
			}

		case <-p.done: // Detecte en fait la fermeture du chan dans core !
			return
		}
	}
}

// Stop Arret du composant
// Mais se base sur le chan done pour arreter la goroutine
func (p *SignalPlugin) Stop() error {
	log.Println("SignalPlugin::Stop")
	return nil
}

// Start Demarre le composant
func (p *SignalPlugin) Start() error {
	log.Println("SignalPlugin::Start")
	go p.run()
	return nil
}

// Enregistre le PluginCtor
func init() {
	collecte.Register("signal", NewSignalPlugin)
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
