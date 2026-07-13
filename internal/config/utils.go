package config

func applyDefaults[T comparable](val *T, def T) {
	if *val == *new(T) {
		*val = def
	}
}

func stripDefaults[T comparable](val *T, def T) {
	if *val == def {
		*val = *new(T)
	}
}
