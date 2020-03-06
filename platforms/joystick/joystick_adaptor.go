package joystick

import (
	"errors"
	"strings"

	"gobot.io/x/gobot"

	"github.com/veandco/go-sdl2/sdl"
)

type joystick interface {
	Close()
	InstanceID() sdl.JoystickID
}

// Adaptor represents a connection to a joystick
type Adaptor struct {
	name     string
	joystick joystick
	connect  func(*Adaptor) (err error)
}

// NewAdaptor returns a new Joystick Adaptor.
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name: gobot.DefaultName("Joystick"),
		connect: func(j *Adaptor) (err error) {
			sdl.Init(sdl.INIT_JOYSTICK)
			if sdl.NumJoysticks() > 0 {
				j.joystick = sdl.JoystickOpen(0)
				return
			}
			return errors.New("No joystick available")
		},
	}
}

// Returns a new Joystick Adaptor for an implementation-specific named device
func NewNamedAdaptor(targetName string) *Adaptor {
	return &Adaptor{
		name: gobot.DefaultName("Joystick"),
		connect: func(adaptor *Adaptor) (err error) {
			initErr := sdl.Init(sdl.INIT_JOYSTICK)
			if initErr != nil {
				return initErr
			}
			foundIndex := -1
			for i := 0; i < sdl.NumJoysticks(); i++ {
				if strings.TrimSpace(targetName) == strings.TrimSpace(sdl.JoystickNameForIndex(i)) {
					adaptor.joystick = sdl.JoystickOpen(foundIndex)
					return
				}
			}
			return errors.New("cannot find joystick by specified name")
		},
	}
}

// Name returns the Adaptors name
func (j *Adaptor) Name() string { return j.name }

// SetName sets the Adaptors name
func (j *Adaptor) SetName(n string) { j.name = n }

// Connect connects to the joystick
func (j *Adaptor) Connect() (err error) {
	err = j.connect(j)
	return
}

// Finalize closes connection to joystick
func (j *Adaptor) Finalize() (err error) {
	j.joystick.Close()
	return
}
