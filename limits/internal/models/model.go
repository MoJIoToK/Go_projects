package models

// ClientLimits represents normalized limits for single client.
// Contains all its instrument positions after processing.
type ClientLimits struct {
	ClientCode string
	Positions  []InstrumentPosition
}

// InstrumentPosition represents a single trading instrument position
// for single client, including its balance in this position and
// associated limits with this instrument.
type InstrumentPosition struct {
	SecId       string
	OpenBalance float64
	Limits      map[int]DepoLimit
}

// DepoLimit represents a deposit limit record for a trading instrument
type DepoLimit struct {
	// ClientCode links to clients limits for one instrument
	ClientCode string
	// SecCode is the ticker for trading instrument (e.g. "SBER")
	SecCode string
	// LimitKind defines the time horizon: 0(intraday), 1(T+1), 2(T+2), 365(long-term)
	LimitKind int
	// OpenLimit is the available limit amount
	OpenLimit float64
	// OpenBalance is the current open balance in this instrument
	OpenBalance float64
}

// MoneyLimit represents a monetary limit record (currently not used in normalization)
type MoneyLimit struct {
	ClientCode  string
	Leverage    float64
	LimitKind   int
	OpenLimit   float64
	OpenBalance float64
}
