package patterns

// TODO: Convert to separate saveState files and use via load

type PatternType string

const (
	Stable     PatternType = "stable"
	Oscillator PatternType = "oscillator"
	Generator  PatternType = "generator"
	Periodic   PatternType = "periodic"
	Chaotic    PatternType = "chaotic"
)

func GetAvailablePatterns() map[PatternType]map[string]string {
	return map[PatternType]map[string]string{
		Stable: {
			"Block":           "block",
			"Beehive":         "beehive",
			"Loaf":            "loaf",
			"Boat":            "boat",
			"Tub":             "tub",
			"AircraftCarrier": "aircraft_carrier",
		},
		Oscillator: {
			"Blinker": "blinker",
			"Toad":    "toad",
			"Beacon":  "beacon",
		},
		Periodic: {
			"Pulsar":         "pulsar",
			"PentaDecathlon": "penta_decathlon",
		},
		Chaotic: {
			"R-Pentomino": "r_pentomino",
		},
	}
}
