package data

type Task struct {
	/// @TODO: Is the data type ID needed as we have one when stored in the DB?
	Id         string `json:"id" bson:"id"`
	Command    string `json:"command" bson:"command"`
	StartedAt  string `json:"started_at" bson:"started_at"`
	FinishedAt string `json:"finished_at" bson:"finished_at"`
	Status     string `json:"status" bson:"status"`
	StdOut     string `json:"stdout" bson:"stdout"`
	StdErr     string `json:"stderr" bson:"stderr"`
	ExitCode   string `json:"exit_code" bson:"exit_code"`
}

func (t *Task) Empty() bool {
	isEmpty := true

	isEmpty = (isEmpty) && (t.Id == "")
	isEmpty = (isEmpty) && (t.Command == "")
	isEmpty = (isEmpty) && (t.StartedAt == "")
	isEmpty = (isEmpty) && (t.FinishedAt == "")
	isEmpty = (isEmpty) && (t.Status == "")
	isEmpty = (isEmpty) && (t.StdOut == "")
	isEmpty = (isEmpty) && (t.StdErr == "")
	isEmpty = (isEmpty) && (t.ExitCode == "")

	return isEmpty
}
