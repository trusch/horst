package links

import "encoding/json"

type Link struct {
	From       string `json:"from"`
	FromOutput string `json:"fromOut"`
	To         string `json:"to"`
	ToInput    string `json:"toIn"`
}

type LinkMap map[string]map[string]*Link

func (m LinkMap) Add(from, fromOut, to, toIn string) {
	m.add(&Link{from, fromOut, to, toIn})
}

func (m LinkMap) add(link *Link) {
	outputMap, ok := m[link.From]
	if !ok {
		outputMap = make(map[string]*Link)
		m[link.From] = outputMap
	}
	outputMap[link.FromOutput] = link
}

func (m LinkMap) Get(from, fromOutput string) (*Link, bool) {
	outputMap, ok := m[from]
	if !ok {
		return nil, ok
	}
	link, ok := outputMap[fromOutput]
	return link, ok
}

func (m LinkMap) ToJSON() (string, error) {
	bs, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (m LinkMap) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), &m)
}

func (m LinkMap) String() string {
	str, err := m.ToJSON()
	if err != nil {
		str = "[error]"
	}
	return str
}
