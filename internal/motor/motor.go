package motor

import (
	"github.com/google/uuid"

	"minebelt/internal/belt"
	"minebelt/internal/roller"
)

type Motor struct {
	ID       string
	BeltID   string
	RatingKW float64
	Baseline float64
	Powered  bool
}

type FeederPort interface {
	Start(beltID string) error
	Stop(beltID string) error
	RecordOrder(beltID, stage string)
}

type BrakePort interface {
	Release(beltID string) error
	Engage(beltID string) error
	Engaged(beltID string) bool
}

type RatingStore interface {
	LoadRatings() (map[string]float64, error)
	SaveRating(beltID string, baseline float64) error
}

type Controller struct {
	line         *belt.Line
	monitor      *roller.Monitor
	ratings      RatingStore
	feeders      FeederPort
	brakes       BrakePort
	motors       map[string]*Motor
	baseline     map[string]float64
	powered      map[string]bool
	readySeq     map[string]int
	clearances   map[string]string
	clearanceLog map[string][]string
}

func NewController(line *belt.Line, monitor *roller.Monitor, ratings RatingStore) *Controller {
	c := &Controller{
		line:         line,
		monitor:      monitor,
		ratings:      ratings,
		motors:       make(map[string]*Motor),
		baseline:     make(map[string]float64),
		powered:      make(map[string]bool),
		readySeq:     make(map[string]int),
		clearances:   make(map[string]string),
		clearanceLog: make(map[string][]string),
	}
	if ratings != nil {
		rows, err := ratings.LoadRatings()
		if err == nil {
			for id, v := range rows {
				c.baseline[id] = v
			}
		}
	}
	return c
}

func (c *Controller) AddMotor(beltID string, ratingKW, baseline float64) *Motor {
	m := &Motor{ID: uuid.NewString(), BeltID: beltID, RatingKW: ratingKW, Baseline: baseline}
	c.motors[m.ID] = m
	c.baseline[beltID] = baseline
	return m
}

func (c *Controller) MotorForBelt(beltID string) (*Motor, bool) {
	for _, m := range c.motors {
		if m.BeltID == beltID {
			return m, true
		}
	}
	return nil, false
}

func (c *Controller) SetFeederPort(p FeederPort) {
	c.feeders = p
}

func (c *Controller) SetBrakePort(p BrakePort) {
	c.brakes = p
}
