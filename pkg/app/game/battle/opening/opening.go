package opening

type Opening interface {
	End()
	Update() bool
	Draw()
}
