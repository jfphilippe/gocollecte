#!/bin/bash
# scripts to compile
# Copyright jean-francois PHILIPPE  2011-2018

# Ok Force Format and vet first
go fmt github.com/jfphilippe/gocollecte/...
go vet github.com/jfphilippe/gocollecte/...

# build each exe you can find !
go install github.com/jfphilippe/gocollecte/cmd/...
