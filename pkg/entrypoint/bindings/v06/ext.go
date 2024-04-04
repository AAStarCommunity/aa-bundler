package entrypoint

const VERSION = "0.6"

func (e *IStakeManagerDepositInfo) GetVersion() string {
	return VERSION
}

func (e *Entrypoint) GetVersion() string {
	return VERSION
}
