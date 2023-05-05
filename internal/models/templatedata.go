package models

// template Data holds data set from Handler send to template
type TemplateData struct {
	StringMap map[string]string
	IntMag    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // interface = dont know what it will be
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}
