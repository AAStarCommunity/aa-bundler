package entrypoint

const VERSION = "0.7"

func (e *IStakeManagerDepositInfo) GetVersion() string {
	return VERSION
}

func (e *Entrypoint) GetVersion() string {
	return VERSION
}
