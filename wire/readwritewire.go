package wire

// type ReadWriteWire struct {
// 	rws ReadWriteSeeker
// 	bin Order
// 	err error
// }

// func NewReadWriteWire(rws io.ReadWriteSeeker, order Order) ReadWriteWire {
// 	return ReadWriteWire{
// 		rws: rws,
// 		bin: order,
// 	}
// }

// func (w *ReadWriteWire) StartReading() Incoming {
// 	return Incoming{
// 		rdr: w.rws,
// 		bin: w.bin,
// 		err: w.err,
// 	}
// }

// func (w *ReadWriteWire) StartWriting() Outgoing {
// 	return Outgoing{
// 		wtr: w.rws,
// 		bin: w.bin,
// 		err: w.err,
// 	}
// }

// func (w *ReadWriteWire) SeekFromStart(off int64) {
// 	w.err = nil
// 	_, w.err = w.rws.Seek(off, io.SeekStart)
// }

// func (w *ReadWriteWire) SeekFromEnd(off int64) {
// 	w.err = nil
// 	_, w.err = w.rws.Seek(off, io.SeekEnd)
// }

// func (w *ReadWriteWire) SeekRelative(off int64) {
// 	w.err = nil
// 	_, w.err = w.rws.Seek(off, io.SeekCurrent)
// }

// func (w *ReadWriteWire) Err() error {
// 	return w.err
// }

// func (w *ReadWriteWire) HasErr() bool {
// 	return w.err != nil
// }

// func (w *ReadWriteWire) AddErr(err error) {
// 	if err != nil {
// 		if w.err == nil {
// 			w.err = err
// 		} else {
// 			w.err = errors.Join(w.err, err)
// 		}
// 	}
// }

// func (w *ReadWriteWire) ClearErrs() {
// 	w.err = nil
// }

// func (w *ReadWriteWire) GetOrder() Order {
// 	return w.bin
// }

// func (w *ReadWriteWire) GetReader() io.Reader {
// 	return w.rws
// }

// func (w *ReadWriteWire) GetWriter() io.Writer {
// 	return w.rws
// }

// func (w *ReadWriteWire) GetSeeker() io.Seeker {
// 	return w.rws
// }

// func (w *ReadWriteWire) GetReadWriteSeeker() io.ReadWriteSeeker {
// 	return w.rws
// }
