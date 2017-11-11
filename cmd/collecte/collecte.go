//Copyright jean-françois PHILIPPE 2014-2016
//
//Paquet principal de l application de collecte.

package main

import (
	"flag"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/jfphilippe/gocollecte/internal/pkg/collecte"
	// Force le chargement du paquet pour executer l initialisation de ce dernier.
	_ "github.com/jfphilippe/gocollecte/internal/pkg/collecte/plugin"
	"github.com/jfphilippe/gocollecte/internal/pkg/config"
)

// Variables de conf.
var (
	configFile = flag.String("config", "./collecte.conf", "Fichier de conf a utiliser")
	logDir     = flag.String("log_dir", "./logs", "Répertoire de logs")
	logLvl     = flag.String("log_lvl", "W", "Niveau de logs (D|I|W|E)")
	cpuprofile = flag.String("cpuprofile", "", "Ecriture de cpu profile dans le fichier indiqué")
	memprofile = flag.String("memprofile", "", "Ecriture de mem profile dans le fichier indiqué")
	logger     *lumberjack.Logger
	Collector  *collecte.Collecte
)

// Initialize logs.
// Use lumberjack to rotate logs
func initLogs() {
	// Va cre un nom de fichiers.
	fname := filepath.Join(*logDir, filepath.Base(os.Args[0])+".log")
	logger = &lumberjack.Logger{
		Filename:   fname,
		MaxSize:    50, // megabytes
		MaxBackups: 7,
		MaxAge:     28, //days
	}
	log.SetOutput(logger)

	log.SetFlags(log.LUTC | log.LstdFlags)
}

func configureRuntime(conf *config.Config) {
	numProcs, _ := conf.Int64("maxProcs", 0)
	if numProcs < 0 {
		numProcs = int64(runtime.NumCPU())
	}
	log.Println("Set NumProcs ", numProcs)
	_ = runtime.GOMAXPROCS(int(numProcs))
}

func main() {
	// Desciption de l usage du programme.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Collecte\n")
		fmt.Fprintf(os.Stderr, "    version: %s\n", "0.1")
		fmt.Fprintf(os.Stderr, "    copyright: %s\n", "jeff")
		fmt.Fprintf(os.Stderr, "Usage de %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Signaux interprétés \n")
		fmt.Fprintf(os.Stderr, "    TERM   arret du processus\n")
		fmt.Fprintf(os.Stderr, "    HUP    rotation du fichier de logs\n")
		fmt.Fprintf(os.Stderr, "    USR1   Donnéee de mémoire dans les logs\n")
		fmt.Fprintf(os.Stderr, "    USR2   Cree un fichier profile de la mémoire(cf option memprofile)\n")
	}

	flag.Parse() // Parse parameters

	initLogs() // Create log file

	// If requested do some CPUProfile Things.
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("Echec creation du fichier", *cpuprofile, err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	conf := config.New(nil)
	err := conf.LoadFile(*configFile)
	if err != nil {
		log.Fatal("Echec chargement configuration", *configFile, err)
	}

	// Check for Runtime config
	configureRuntime(conf.Section("runtime"))

	Collector, err = collecte.New(conf)

	if err == nil {
		Collector.Start()

		// Boucle d attente...
	run:
		for {
			select {
			case s := <-Collector.Cmd:
				switch s {
				case "log_rotate":
					log.Println("Rotate!")
					logger.Rotate()
				case "stop":
					log.Println("Stop !")
					break run
				default:
					log.Println("Reception '", s, "' ignore")
				}
			}
		}

		// Notifie les composants de la "fermeture"
		err = Collector.Stop()
	}

	log.Println("Done")
}
