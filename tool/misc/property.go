package misc

type Property interface {
	~int | ~int32 | ~int64 | ~uint32 | ~uint64 | ~uint | ~uint8 | ~uint16 | ~int8 | ~int16
}

func HasProperty[T Property](p T, properties ...T) bool {
	for _, v := range properties {
		if p&v == 0 {
			return false
		}
	}
	return true
}

func HasOneProperty[T Property](p T, properties ...T) bool {
	for _, v := range properties {
		if p&v != 0 {
			return true
		}
	}
	return false
}

func AddProperty[T Property](p T, properties ...T) T {
	for _, v := range properties {
		p = p | v
	}
	return p
}

func CreateProperty[T Property](properties ...T) T {
	return AddProperty(0, properties...)
}
