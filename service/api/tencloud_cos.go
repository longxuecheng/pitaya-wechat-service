package api

type ITencloudCos interface {
	PushData(name string, data []byte) error
}
