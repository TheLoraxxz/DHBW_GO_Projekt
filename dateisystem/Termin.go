package dateisystem

type repeat int

const (
	taeglich repeat = iota
	woechentlich
	monatlich
	jaehrlich
	niemals
)

type Termin struct {
	title       string
	description string
	recurring   repeat
	time        string
	date        string
}
