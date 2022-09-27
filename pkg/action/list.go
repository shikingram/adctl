package action

type List struct {
	cfg *Configuration

	ReleaseName string
}

func NewList(cfg *Configuration) *List {
	return &List{cfg: cfg}
}
