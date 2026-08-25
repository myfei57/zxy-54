package audit

type Saver interface {
	SaveJSON(name string, value interface{}) error
	LoadJSON(name string, value interface{}) error
}
