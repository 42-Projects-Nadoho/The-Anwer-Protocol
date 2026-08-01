package protocol


// The class Command represent an action parsing from the player entry
type Command struct {
	Action	string
	Args	[]string
}

// The class Reponse represent a response from an action
type Reponse struct {
	Type	string
	Payload	string
}
