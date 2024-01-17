package data

type Sensor struct {
	Name string `json:"name"`
	Measures []Measure `json:"measure"` 
}

type Measure struct {
	Value uint8 `json:"value"`
}