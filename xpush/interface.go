package xpush

type (
	IPushMod interface {
		Push(title string, content string) error
		PushMarkDown(title string, content string) error
		SetSetting(setting interface{}) error
	}
)
