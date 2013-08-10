package vw

import (
	"bytes"
	"errors"
	"strconv"
)

type Data struct {
	Label      float64
	Importance float64
	Tag        string

	Features map[string]map[string]float64
}

func (d *Data) readHeader(src []byte) (err error) {
	// label / important / tag are the first split
	parts := bytes.Split(src, []byte(" "))

	// read the label
	d.Label, err = strconv.ParseFloat(string(parts[0]), 64)
	if err != nil {
		err = errors.New("Error parsing label: " + err.Error())
		return
	}

	if len(parts) == 1 || (len(parts) == 2 && len(parts[1]) == 0) {
		d.Importance = 1
		return nil
	}

	d.Importance, err = strconv.ParseFloat(string(parts[1]), 64)
	if err != nil {
		err = errors.New("Error parsing importance: " + err.Error())
		return
	}

	if len(parts) == 2 {
		return nil
	}

	d.Tag = string(parts[2])
	return nil
}

func (d *Data) readNamespace(src []byte) (err error) {
	parts := bytes.Split(src, []byte(" "))

	namespace := ""
	multiplier := 1.0

	for idx, part := range parts {
		key, weight, err := readFeature(part)
		if err != nil {
			return errors.New("Error reading feature from '" + string(part) + "': " + err.Error())
		}

		if idx == 0 {
			namespace = key
			multiplier = weight

			if _, ok := d.Features[namespace]; !ok {
				d.Features[namespace] = map[string]float64{}
			}
		} else {
			weight *= multiplier

			d.Features[namespace][key] = d.Features[namespace][key] + weight
		}
	}

	return nil
}

func readFeature(src []byte) (key string, weight float64, err error) {
	subParts := bytes.Split(src, []byte(":"))

	key = string(subParts[0])
	weight = 1.0
	if len(subParts) == 2 {
		weight, err = strconv.ParseFloat(string(subParts[1]), 64)
		if err != nil {
			return key, weight, errors.New("Error parsing weight from '" + string(src) + "': " + err.Error())
		}
	} else if len(subParts) > 2 {
		return key, weight, errors.New("Bad feature: " + string(src))
	}

	return
}

// Parse an instance description into a label and features. Feature namespaces are collapsed into the feature key.
// descriptions follow the scheme defined at https://github.com/JohnLangford/vowpal_wabbit/wiki/Input-format
func Parse(description []byte) (d Data, err error) {
	d = Data{Features: map[string]map[string]float64{}}
	split := bytes.Split(description, []byte("|"))

	if err = d.readHeader(split[0]); err != nil {
		return d, err
	}

	for _, part := range split[1:] {
		if err = d.readNamespace(part); err != nil {
			return d, err
		}
	}

	return d, nil
}
