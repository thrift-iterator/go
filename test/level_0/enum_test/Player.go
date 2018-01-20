package enum_test

type Player int64
const (
	Player_JAVA Player = 0
	Player_FLASH Player = 1
)

func (p Player) String() string {
	switch p {
	case Player_JAVA: return "JAVA"
	case Player_FLASH: return "FLASH"
	}
	return "<UNSET>"
}