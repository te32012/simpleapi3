package entityies

type Commands []Command

type Command struct {
	Id          int    `json:"Id,omitempty"`
	Description string `json:"Description"`
	Script      string `json:"Script,omitempty"`
}
type CommandCreated struct {
	Id int `json:"Id"`
}
