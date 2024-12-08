package menu

type menuTraining struct {
	result Result
}

func trainingNew() (*menuTraining, error) {
	res := &menuTraining{
		result: ResultContinue,
	}
	return res, nil
}

func (r *menuTraining) End() {
}

func (r *menuTraining) Update() bool {
	// debug: とりあえずgo_battleを返す
	r.result = ResultGoBattle
	return true

	// if inputs.CheckKey(inputs.KeyCancel) == 1 {
	// 	sound.On(resources.SECancel)
	// 	return true
	// }
	// return false
}

func (r *menuTraining) Draw() {
}

func (r *menuTraining) GetResult() Result {
	return r.result
}
