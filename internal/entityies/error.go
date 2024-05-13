package entityies

type Error struct {
	E   error
	Err []byte `json:"Err"`
}
