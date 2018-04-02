//  Copyright jean-françois PHILIPPE 2014-2018

package measure

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// JSONValuesWriter Encode Values to a Json File.
// Handle File rotation
// Given a 'base name', a new file is creted by appending date to basename.
type JSONValuesWriter struct {
	sync.Mutex
	Basename     string // mst point to an existing directory !
	writer       *json.Encoder
	file         *os.File // Current file or nil
	Maxrecords   uint     // Max records allowed in file, 0 to disable.
	nbrecords    uint     // Current nb of record in file
	Daily        bool     // true if file must "change" each day
	dailyOpenday int      // day in month of current file or 0
}

// SetRotateDaily configure la rotation quotidienne du fichier.
func (w *JSONValuesWriter) SetRotateDaily(daily bool) *JSONValuesWriter {
	w.Daily = daily
	return w
}

// SetMaxRecords configure le nbre max d enregistrements par fichier
func (w *JSONValuesWriter) SetMaxRecords(maxrecords uint) *JSONValuesWriter {
	w.Maxrecords = maxrecords
	return w
}

// SetBasename configure le nom du fichier
func (w *JSONValuesWriter) SetBasename(basename string) *JSONValuesWriter {
	w.Basename = basename
	return w
}

func (w *JSONValuesWriter) rotate() error {
	now := time.Now()
	name := fmt.Sprintf("%04d%02d%02d-%02d%02d%02d.json",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())
	filename := w.Basename + name

	log.Println("JsonValuesWriter::rotate Try to create ", filename)
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		w.file = nil
		w.writer = nil
		return err
	}
	w.file = fd
	w.writer = json.NewEncoder(fd)
	w.dailyOpenday = now.Day()
	w.nbrecords = 0
	return nil
}

// Write enregistrement d une valeur
func (w *JSONValuesWriter) Write(values *NodeValues) error {
	//
	w.Lock()
	defer w.Unlock()

	// Check for file rotation.
	now := time.Now()
	var err error
	if w.file == nil ||
		(w.Daily && w.dailyOpenday != now.Day()) ||
		(w.Maxrecords > 0 && w.nbrecords >= w.Maxrecords) {
		err = w.rotate()
		if err != nil {
			return err
		}
	}

	// Write record
	err = w.writer.Encode(*values)
	w.nbrecords++
	return err
}

// Close ferme le flux
func (w *JSONValuesWriter) Close() error {
	log.Println("JsonValuesWriter::Close")

	w.Lock()
	defer w.Unlock()

	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		w.writer = nil
		return err
	}
	return nil
}
