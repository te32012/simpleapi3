package entityies

type Commands []Command

type Command struct {
	Id          int
	Description string
	Script      string
}
type CommandCreated struct {
	Id int
}
