package parametertable2

const (
	// DERIVED
	LOGIN_TOP_LEFT uint16 = iota
	LOGIN_TOP_RIGHT
	LOGIN_BOT_LEFT
	LOGIN_BOT_RIGHT
	// ROOT
	WINDOW_WIDTH
	WINDOW_HEIGHT
	_vallen
)

type Param struct {
	val      float32
	children []uint16
	calc     func() float32
}

var VAL = [_vallen]float32{}
