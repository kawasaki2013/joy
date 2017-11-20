package parser

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"

	"github.com/matthewmueller/golly/internal/dom/def"
	"github.com/matthewmueller/golly/internal/dom/defs"
	"github.com/matthewmueller/golly/internal/dom/index"
	"github.com/pkg/errors"
)

// apis struct
type apis struct {
	CallbackFunctions  []*defs.Callback          `xml:"callback-functions>callback-function"`
	CallbackInterfaces []*defs.CallbackInterface `xml:"callback-interfaces>interface"`
	Dictionaries       []*defs.Dictionary        `xml:"dictionaries>dictionary"`
	Enums              []*defs.Enum              `xml:"enums>enum"`
	Interfaces         []*defs.Interface         `xml:"interfaces>interface"`
	MixinInterfaces    []*defs.Interface         `xml:"mixin-interfaces>interface"`
	TypeDefs           []*defs.TypeDef           `xml:"typedefs>typedef"`
}

// docs
type docs struct {
	Interfaces []*struct {
		Name    string `json:"name"`
		Members struct {
			Properties []struct {
				Name    string `json:"name"`
				Comment string `json:"comment"`
			} `json:"property"`
			Methods []struct {
				Name    string `json:"name"`
				Comment string `json:"comment"`
			} `json:"method"`
		} `json:"members"`
	} `json:"interfaces"`
}

// Parse the data
func Parse(api string, doc string) ([]def.Definition, error) {
	var a apis
	if e := xml.Unmarshal([]byte(api), &a); e != nil {
		return nil, errors.Wrap(e, "error parsing xml")
	}

	var d docs
	if e := json.Unmarshal([]byte(doc), &d); e != nil {
		return nil, errors.Wrap(e, "error parsing json")
	}

	docmap := map[string]int{}
	methodmap := map[string]int{}
	propmap := map[string]int{}
	for i, iface := range d.Interfaces {
		if _, isset := docmap[iface.Name]; isset {
			return nil, fmt.Errorf("duplicate doc %s", iface.Name)
		}
		docmap[iface.Name] = i

		for i, method := range iface.Members.Methods {
			methodmap[iface.Name+" "+method.Name] = i
		}
		for i, prop := range iface.Members.Properties {
			propmap[iface.Name+" "+prop.Name] = i
		}
	}

	// add in the documentation if it's present
	for _, iface := range a.Interfaces {
		name := iface.InterfaceName

		docidx, isset := docmap[name]
		if !isset {
			continue
		}
		doc := d.Interfaces[docidx]
		if doc == nil {
			continue
		}

		for _, m := range iface.Methods {
			if i, isset := methodmap[name+" "+m.Name]; isset {
				m.Comment = doc.Members.Methods[i].Comment
			}
		}

		for _, m := range iface.Properties {
			if i, isset := propmap[name+" "+m.Name]; isset {
				m.Comment = doc.Members.Properties[i].Comment
			}
		}
	}

	idx := index.Index{}

	// add to the index
	for _, d := range a.CallbackFunctions {
		d.Index = idx
		idx[d.ID()] = d
	}
	for _, d := range a.CallbackInterfaces {
		d.Index = idx
		idx[d.ID()] = d
	}
	for _, d := range a.Dictionaries {
		d.Index = idx
		idx[d.ID()] = d
	}
	for _, d := range a.Enums {
		d.Index = idx
		idx[d.ID()] = d
	}
	for _, d := range a.Interfaces {
		d.Index = idx
		idx[d.ID()] = d
	}
	for _, d := range a.MixinInterfaces {
		d.Index = idx
		idx[d.ID()] = d
	}
	for _, d := range a.TypeDefs {
		d.Index = idx
		idx[d.ID()] = d
	}

	return sortDefs(idx), nil
}

func sortDefs(m map[string]def.Definition) (sorted []def.Definition) {
	var order []string
	for k := range m {
		order = append(order, k)
	}

	sort.Strings(order)

	for _, o := range order {
		sorted = append(sorted, m[o])
	}

	return sorted
}