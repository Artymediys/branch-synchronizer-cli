package cli

var (
	ErrorSpinner = "возникла ошибка при отображении интерфейса загрузки -> "
	ErrorForm    = "возникла ошибка при отображении пользовательского интерфейса -> "
)

type StageStatus int8

func (s StageStatus) Status() string {
	switch s {
	case 1:
		return "DONE"
	case 0:
		return "NOT STARTED"
	case -1:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

type Stage struct {
	MergeRequest StageStatus
	Notification StageStatus
}
