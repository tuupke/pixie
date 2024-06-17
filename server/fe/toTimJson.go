package main

type (
	SequenceType      string
	SequenceAxis      string
	SequenceDirection string

	SequenceInterface struct {
		Typ              SequenceType      `json:"type"`
		Num              float64           `json:"num"`
		Axis             SequenceAxis      `json:"axis"`
		Dir              SequenceDirection `json:"dir"`
		Radius           float64           `json:"radius"`
		Separation       float64           `json:"separation"`
		EquivalentSpaced bool              `json:"equivalentSpaced"`
	}
	ElementInterface struct {
		Base    RotationCoordinateInterface `json:"base"`
		Repeats []SequenceInterface         `json:"repeats"`
	}
	PathInterface struct {
		Start CoordinateInterface `json:"start"`
		End   CoordinateInterface `json:"end"`
	}

	RoomInterface struct {
		Name     string
		Outline  []CoordinateInterface
		Elements []ElementInterface
		Paths    []PathInterface
	}
	CoordinateInterface struct {
		X float64
		Y float64
	}

	RotationCoordinateInterface struct {
		CoordinateInterface `json:"embed"`
		Rotation            float64 `json:"rotation"`
	}
)

const (
	SequenceLine   SequenceType = "line"
	SequenceCircle SequenceType = "circle"

	SequenceVertical   SequenceAxis = "vertical"
	SequenceHorizontal SequenceAxis = "horizontal"

	SequencePositive SequenceDirection = "positive"
	SequenceNegative SequenceDirection = "negative"
)
