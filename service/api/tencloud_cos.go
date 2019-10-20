package api

type ITencloudCos interface {
	PushImageObject(name string, data []byte) error
}
