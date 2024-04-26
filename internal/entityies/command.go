package entityies

type Commands []Command

type Command struct {
	Id          string
	Description string
	Script      string
}
type CommandCreated struct {
	Id int
}
