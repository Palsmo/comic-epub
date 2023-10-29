package types

// hook
type Option struct {
	option
}

// option data type
type option struct {
	Flag   string              // 'label
	Nargs  int                 // number of arguments (-1, set from flag number suffix)
	Info   []string            // 'description, each entry is a row
	Action func(...string) int // trigger as result of activation
}

// option factory
func NewOption(flag string, nargs int, info []string) Option {
	o := option{flag, nargs, info, nil}
	return Option{o}
}
