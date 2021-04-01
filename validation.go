package utils

import (
	"encoding/json"
	"regexp"
)

var (
	emailRegex         = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]$`)
	phoneRegex         = regexp.MustCompile(`^\+\d{11,15}$`)
	fullShortlinkRegex = regexp.MustCompile(`([a-zA-Z0-9]+[_-]?){1,256}$`)
	shortlinkRegex     = regexp.MustCompile(`([a-zA-Z0-9]+[_-]?){5,256}$`)
	resourceNameRegex  = regexp.MustCompile(`([a-z]+[_-]?){2,255}$`)
)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePhoneNumber(phone string) bool {
	return phoneRegex.MatchString(phone)
}

func ValidateShortlink(r string) bool {
	return shortlinkRegex.MatchString(r)
}

func ValidateFullShortlink(r string) bool {
	return fullShortlinkRegex.MatchString(r)
}

func ValidateResourceName(r string) bool {
	return resourceNameRegex.MatchString(r)
}

func CheckInt(old, new int) int {
	if new != old {
		if new != 0 {
			return new
		}
	}
	return old
}

func CheckInt8(old, new int8) int8 {
	if new != old {
		if new != 0 {
			return new
		}
	}
	return old
}

func CheckInt16(old, new int16) int16 {
	if new != old {
		if new != 0 {
			return new
		}
	}
	return old
}

func CheckInt32(old, new int32) int32 {
	if new != old {
		if new != 0 {
			return new
		}
	}
	return old
}

func CheckInt64(old, new int64) int64 {
	if new != old {
		if new != 0 {
			return new
		}
	}
	return old
}

func CheckUint8(old, new uint8) uint8 {
	if new != old {
		if new != 0 {
			return new
		}
	}
	return old
}

func CheckUint16(old, new uint16) uint16 {
	if new != old {
		if new != 0 {
			return new
		}
	}
	return old
}

func CheckUint32(old, new uint32) uint32 {
	if new != old {
		if new != 0 {
			return new
		}
	}
	return old
}

func CheckUint64(old, new uint64) uint64 {
	if new != old {
		if new != 0 {
			return new
		}
	}
	return old
}

func CheckString(old, new string) string {
	if new != old {
		if new != "" {
			return new
		}
	}
	return old
}

//jsonb
type PropertyMap map[string]interface{}

func (p *PropertyMap) Scan(source json.RawMessage) error {

	jsonStr := string(source)

	//p = make(map[string]interface{})

	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		return err
	}

	return nil
}

func CheckJson(old, new json.RawMessage) (json.RawMessage, error) {
	var o, n map[string]interface{}

	if old == nil || new == nil {
		return nil, nil
	}

	if err := json.Unmarshal(old, &o); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(new, &n); err != nil {
		return nil, err
	}

	for i := range n {
		newProp := n[i]
		for j := range o {
			oldProp := o[j]
			if newProp != oldProp {
				o[j] = newProp
			} else {
				o[i] = oldProp
			}
		}
	}

	res, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return res, nil
}
