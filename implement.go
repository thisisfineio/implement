package implement


type Options struct {
	InputPath string
	OutputPath string
	OverwriteExisting bool
	Interfaces []string
	Implementors []string
}


func (o *Options) Validate() error {
	if len(o.Interfaces) != len(o.Implementors) {
		return ErrInvalidLengths
	}
	return nil
}
