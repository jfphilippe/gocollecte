// Copyright jean-françois PHILIPPE 2014-2016
//
// executable pour envoyer des commandes au programme de collecte

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/jfphilippe/gocollecte/internal/pkg/pipe"
)

var (
	filename = flag.String("socket", pipe.SocketName, "Name of socket file")
)

func main() {
	// Desciption de l usage du programme.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "send\n")
		fmt.Fprintf(os.Stderr, "    version: %s\n", "0.1")
		fmt.Fprintf(os.Stderr, "    copyright: %s\n", "jeff")
		fmt.Fprintf(os.Stderr, "Usage de %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse() // Parse parameters

	// Open socket
	conn, err := net.DialUnix("unixgram", nil, &net.UnixAddr{Name: *filename, Net: "unixgram"})

	if err == nil {
		// Create and send cmd
		cmdes := flag.Args()
		for _, cmd := range cmdes {
			msg, _ := pipe.New(cmd)
			b, _ := json.Marshal(msg)
			_, err = conn.Write(b)
			if err != nil {
				panic(err)
			}

		}
	} else {
		panic(err)
	}

}
