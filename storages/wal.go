package storages

import (
	"go.uber.org/zap"
	"io"
	"os"
)

type WAL struct {
	File     *os.File
	Filename string
	OnDecode func(reader io.Reader)
	Logger   *zap.SugaredLogger
}

func NewWAL(lg *zap.SugaredLogger, filename string, onDecode func(io.Reader)) *WAL {
	w := new(WAL)
	w.Filename = filename
	w.OnDecode = onDecode
	w.Logger = lg
	return w
}

func (w *WAL) Write(serializedEvent []byte) {
	_, err := w.File.Write(serializedEvent)
	if err != nil {
		w.Logger.Fatalf("Unable to write to wal file: %s", err)
	}
}

func (w *WAL) Recover() {
	walFileName := w.Filename
	w.Logger.Infof("Trying to recover WAL %s", walFileName)
	if _, err := os.Stat(walFileName); err != nil {
		walFile, err := os.Open(walFileName)
		if err != nil && !os.IsExist(err) {
			w.Logger.Fatalf("Unable to open wal file: %s", err)
		}

		w.OnDecode(walFile)

		w.File, err = os.OpenFile(walFileName, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			w.Logger.Fatalf("Unable to open wal file with rights: %s", err)
		}
	} else {
		w.Logger.Infof("No WAL %s was found, creating...", walFileName)

		w.File, err = os.Create(walFileName)
		if err != nil {
			w.Logger.Fatalf("Unable to create wal file: %s", err)
		}
	}

}

func (w *WAL) Recreate() {
	err := os.Remove(w.Filename)
	if err != nil && err != os.ErrNotExist {
		w.Logger.Fatalf("Unable to remove wal file during recreation: %s", err)
	}

	w.File, err = os.Create(w.Filename)
	if err != nil {
		w.Logger.Fatalf("Unable to create wal file during recreation: %s", err)
	}
}
