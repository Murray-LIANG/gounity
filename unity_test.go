package gounity

import (
	"testing"
)

type WrapperOfStorage struct {
	unity Storage
}

func (w *WrapperOfStorage) SetStorage(s Storage) {
	w.unity = s
}

func TestStorageInterface(t *testing.T) {
	ctx, _ := newTestContext()
	w := &WrapperOfStorage{}
	//Use below call to make all methods in Storage interface are implemented by Unity strcut.
	w.SetStorage(ctx.unity)
}
