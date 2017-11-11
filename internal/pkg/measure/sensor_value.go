//  Copyright jean-françois PHILIPPE 2014-2016

package measure

import (
	"expvar"
	"github.com/jfphilippe/gocollecte/internal/pkg/config"
	"log"
	"sync"
)

const (
	// Taille initiale du tableau de valeurs.
	// see NewNodeValues
	initArraySize int = 5
)

var (
	valuesMetrics  *expvar.Map
	nbNodeValues   *expvar.Int
	nbSensorValues *expvar.Int
)

// SensorID identifiant de senseur
type SensorID uint16

// NodeID identifiant de Node.
type NodeID uint16

// SensorValue Represente un enregistrement de 'valeur'.
type SensorValue struct {
	Value  int64    `json:"value"`  // Valeur (entière)
	Sensor SensorID `json:"sensor"` // Type de valeur
}

// NodeValues Un ensemble de valeurs d un Node
// Valeur prises au meme moment.
type NodeValues struct {
	When   int64         `json:"when"`   // Date de la mesure
	Node   NodeID        `json:"node"`   // Id de Node
	Values []SensorValue `json:"values"` // Values
}

// NewNodeValues Reserve N slots de SensorValue
func NewNodeValues(when int64, nodeID NodeID) *NodeValues {
	return &NodeValues{When: when, Node: nodeID, Values: make([]SensorValue, 0, initArraySize)}
}

// AppendValue Ajout d une valeur.
func (v *NodeValues) AppendValue(Value int64, Sensor SensorID) {
	v.Values = append(v.Values, SensorValue{Value, Sensor})
}

// ValuesWriter interface de composant enregistrant des valeurs
type ValuesWriter interface {
	Write(values *NodeValues) error
	Close() error
}

// ValuesWriterCtor Constructeur de ValuesWriter
type ValuesWriterCtor func(conf *config.Config) (ValuesWriter, error)

// ValuesHandler Charge de traiter des Valeurs vers un ValuesWriter
type ValuesHandler struct {
	done   chan struct{}
	Chan   chan *NodeValues
	writer ValuesWriter
	nbInst int
}

// Start Launch the component.
func (h *ValuesHandler) Start() {
	var wg sync.WaitGroup
	routine := func() {
		defer wg.Done()
		log.Println("run.start")
		for {
			select {
			case val := <-h.Chan:
				// val == nil if Chan closed !
				if val != nil {
					nbNodeValues.Add(1)
					nbSensorValues.Add(int64(len(val.Values)))
					err := h.writer.Write(val)
					if err != nil {
						log.Println("ValuesHandler::Echec save", err)
					}
				}
			case <-h.done:
				log.Println("ValuesHandler::Start.end")
				return
			}
		}
	}
	wg.Add(h.nbInst)
	for i := 0; i < h.nbInst; i++ {
		go routine()
	}

	go func() {
		wg.Wait()
		err := h.writer.Close()
		if err != nil {
			log.Println("Error on  writer.Close", err)
		}
	}()
}

func init() {
	nbNodeValues = new(expvar.Int)
	nbSensorValues = new(expvar.Int)
	valuesMetrics = expvar.NewMap("values")
	valuesMetrics.Set("nodes", nbNodeValues)
	valuesMetrics.Set("sensors", nbSensorValues)
}

/* vi:set fileencodings=utf-8 tabstop=4 ai sw=2: */
