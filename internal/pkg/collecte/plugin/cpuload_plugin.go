//  Copyright jean-françois PHILIPPE 2014-2018

package plugin

import (
	"errors"
	"github.com/jfphilippe/gocollecte/internal/pkg/collecte"
	"github.com/jfphilippe/gocollecte/internal/pkg/config"
	"github.com/jfphilippe/gocollecte/internal/pkg/measure"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

// CPUPlugin Plugin de mesure de la charge CPU
type CPUPlugin struct {
	filename string
	freq     time.Duration
	done     chan struct{}
	values   chan *measure.NodeValues
	tick     *time.Ticker
}

// NewCPULoad Un PluginCtor
// Chargé de créer une instance de Plugin
func NewCPULoad(conf *config.Config, core *collecte.Collecte) (collecte.Plugin, error) {
	freq, err := conf.Duration("frequency", time.Minute*15)
	log.Println("frequency ", freq)
	if err != nil {
		return nil, err
	}
	filename, err := conf.String("filename", "/proc/loadavg")
	if err != nil {
		return nil, err
	}
	return &CPUPlugin{filename: filename, freq: freq, done: core.Done, values: core.Values}, nil
}

// Stop Arret du composant
// Stoppe le Ticker.
// Mais se base sur le chan done pour arreter la goroutine
func (p *CPUPlugin) Stop() error {
	log.Println("CpuPlugin::Stop")
	if p.tick != nil {
		p.tick.Stop()
	}
	return nil
}

// Goroutine lancee par Start
func (p *CPUPlugin) run() {
	for {
		select {
		case <-p.tick.C:
			str, err := p.readFile()
			if err == nil {
				val := p.newValues(str)
				if val != nil {
					p.values <- val
				}
			} else {
				log.Println("Can not readFile ", p.filename, " : ", err)
			}
		case <-p.done: // Detecte en fait la fermeture du chan dans core !
			return
		}
	}
}

func (p *CPUPlugin) readFile() (string, error) {
	dat, err := ioutil.ReadFile(p.filename)
	if err == nil {
		return string(dat), nil
	}
	return "", err
}

func (p *CPUPlugin) newValues(val string) *measure.NodeValues {
	vals := strings.Fields(val)
	if len(vals) > 3 {
		when := time.Now().Unix()
		v := measure.NewNodeValues(when, 0)
		var avg float64
		for i := 0; i < 3; i++ {
			avg, _ = strconv.ParseFloat(vals[i], 64)
			avg = avg*100.0 + 0.5
			avg = math.Floor(avg)
			v.AppendValue(int64(avg), measure.SensorID(10+i))
		}

		return v
	}

	return nil
}

// Start Demarre le composant
func (p *CPUPlugin) Start() error {
	log.Println("CpuPlugin::Start")
	p.tick = time.NewTicker(p.freq)
	if p.tick != nil {
		go p.run()
		return nil
	}
	return errors.New("Can not create Tick")
}

// Enregistre le PluginCtor
func init() {
	collecte.Register("loadavg", NewCPULoad)
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
