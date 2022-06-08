package opening

type Opening interface {
	End()
	Process() bool
	Draw()
}
