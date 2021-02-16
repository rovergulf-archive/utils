package storages

import (
	"encoding/gob"
	"go.uber.org/zap"
	"os"
	"time"
)

const (
	backupSuffix = ".backup"
	tempSuffix   = ".temp"
)

type Dump struct {
	Filename        string
	FlushFunc       func()
	Deserialize     func([]byte)
	OnFlushComplete func()
	FlushDelay      time.Duration
	IsFlushing      bool
	Logger          *zap.SugaredLogger
}

func NewDump(lg *zap.SugaredLogger, filename string, flushDelay time.Duration, flushFunc func(), onFlushComplete func()) *Dump {
	d := new(Dump)
	d.Filename = filename
	d.FlushDelay = flushDelay
	d.FlushFunc = flushFunc
	d.OnFlushComplete = onFlushComplete
	d.Logger = lg.Named("dump")
	return d
}

func (d *Dump) Flush(version int64, data interface{}) {
	d.Logger.Infof("Flushing dump file: %s", d.Filename)

	d.IsFlushing = true

	startTime := time.Now()
	dumpFileName := d.Filename

	tempDumpFileName := dumpFileName + tempSuffix

	file, err := os.Create(tempDumpFileName)
	if err != nil {
		d.Logger.Fatalf("Unable to create temp dump '%s' : %s", tempDumpFileName, err)
	}

	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(version); err != nil {
		d.Logger.Fatalf("Unable to encode temp dump data : %s", err)
	}
	if err := encoder.Encode(data); err != nil {
		d.Logger.Fatalf("Unable to encode temp dump data : %s", err)
	}

	if err := file.Close(); err != nil {
		d.Logger.Fatalf("Unable to close temp dump '%s' : %s", tempDumpFileName, err)
	}

	_, err = os.Stat(dumpFileName)
	// if == nil <- watch out !!
	if err == nil {
		err := os.Rename(dumpFileName, dumpFileName+backupSuffix)
		if err != nil {
			d.Logger.Fatalf("Unable to flush the dump '%s' : %s", d.Filename, err)
		}
	}

	if err := os.Rename(tempDumpFileName, dumpFileName); err != nil {
		d.Logger.Fatalf("Unable to flush the dump '%s' : %s", d.Filename, err)
	}

	if d.OnFlushComplete != nil {
		d.OnFlushComplete()
	}

	d.IsFlushing = false

	durationTime := time.Since(startTime) / 1000000
	d.Logger.Infof("Duration of flush of %s: %d%s", d.Filename, durationTime, "ms")
}

func (d *Dump) Recover(version *int64, dataDst interface{}) error {
	dumpFileName := d.Filename
	d.Logger.Infof("Trying to recover dump %s", dumpFileName)

	_, err := os.Stat(dumpFileName)
	if err != nil {
		d.Logger.Errorf("Unable to check out file info: %s", err)
		return err
	}
	d.Logger.Infof("Dump exists. Restoring %s", dumpFileName)

	file, err := os.Open(dumpFileName)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(version); err != nil {
		d.Logger.Errorf("Failed to decode file version: %s", err)
		return err
	}
	if err := decoder.Decode(dataDst); err != nil {
		d.Logger.Errorf("Failed to decode file data: %s", err)
		return err
	}

	return file.Close()
}

func (d *Dump) StartFlushThread() {
	d.Logger.Infof("Starting flushing thread for %s", d.Filename)
	ticker := time.NewTicker(d.FlushDelay)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				d.Logger.Infof("Flushing a dump %s", d.Filename)
				d.FlushFunc()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (d *Dump) ForceFlush() {
	d.Logger.Info("Forced flush")
	if !d.IsFlushing {
		d.FlushFunc()
	}
}
