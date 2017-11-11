//  Copyright jean-françois PHILIPPE 2014-2016

package measure

import (
	"database/sql"
	"errors"
	"github.com/jfphilippe/gocollecte/internal/pkg/config"
)

var (
	ctors = make(map[string]ValuesWriterCtor, 5)
)

// NewValuesHandler cree une nouvelle instance.
func NewValuesHandler(conf *config.Config, done chan struct{}) (*ValuesHandler, error) {
	nbInst, _ := conf.Int64("nbInstances", 1)
	chanSize, _ := conf.Int64("chanSize", 1)
	writer, err := conf.String("writer", "-dontexists-")
	if err != nil {
		return nil, errors.New("No Writer configured 'writer' param missing")
	}

	writerConf := conf.Section(writer)
	valType, err := writerConf.String("type")

	if err != nil {
		return nil, errors.New("Writer '" + writer + "' type not defined")
	}

	ctor := ctors[valType]
	if ctor == nil {
		return nil, errors.New("Writer '" + writer + "' type '" + valType + "' unknown !")
	}

	valuesWriter, err := ctor(writerConf)
	if err == nil {
		return &ValuesHandler{done, make(chan *NodeValues, chanSize), valuesWriter, int(nbInst)}, nil
	}
	return nil, err
}

// Register Enregistre un construtor...
func Register(name string, ctor ValuesWriterCtor) {
	ctors[name] = ctor
}

func newDbValuesWriter(conf *config.Config) (ValuesWriter, error) {
	strcon, err := conf.String("database")
	if err == nil {
		db, err := sql.Open("postgres", strcon)
		if err == nil {
			return &DbValuesWriter{PGInsertStmt, db}, nil
		}
	}
	return nil, err
}

func newFileValuesWriter(conf *config.Config) (ValuesWriter, error) {
	daily, _ := conf.Bool("daily", true)
	maxrecords, _ := conf.Int64("maxrecords", 0)
	basename, err := conf.String("basename")
	if err != nil {
		return nil, err
	}
	return &JSONValuesWriter{Basename: basename, Daily: daily, Maxrecords: uint(maxrecords)}, nil
}

func init() {
	Register("database", newDbValuesWriter)
	Register("file", newFileValuesWriter)
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
