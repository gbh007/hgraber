package modelV2

type RawAttribute struct {
	Parsed bool     `json:"parsed"`
	Values []string `json:"values,omitempty"`
}

func (attr RawAttribute) Copy() RawAttribute {
	copyAttr := RawAttribute{
		Parsed: attr.Parsed,
		Values: make([]string, len(attr.Values)),
	}

	copy(copyAttr.Values, attr.Values)

	return copyAttr
}
