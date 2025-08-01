package utils

import (
	"cmp"
	"errors"
	"fmt"
	"io"

	"github.com/gabe-lee/OurSweeper/data_buffer"
)

type (
	WriteBuffer = data_buffer.WriteBuffer
)

type ErrorChecker struct {
	Err error
}

func (e *ErrorChecker) IsErr(err error) bool {
	e.Err = err
	return err != nil
}

type FirstError struct {
	Err error
}

func (e *FirstError) Add(err error) {
	if e.Err == nil {
		e.Err = err
	}
}

func (e *FirstError) Read(r io.Reader, dst []byte) (n int) {
	n, ee := r.Read(dst)
	e.Add(ee)
	return n
}

type ErrorCollector struct {
	Err error
}

func (e *ErrorCollector) Do(err error) {
	e.Err = errors.Join(e.Err, err)
}

func (e *ErrorCollector) Add(err error) {
	e.Err = errors.Join(e.Err, err)
}

func (e *ErrorCollector) AddFmt(format string, args ...any) {
	e.Err = errors.Join(e.Err, fmt.Errorf(format, args...))
}

type ErrorBuffer struct {
	buf    *WriteBuffer
	writer io.Writer
}

func NewErrorBuffer(writer io.Writer, buf *WriteBuffer) ErrorBuffer {
	return ErrorBuffer{
		buf:    buf,
		writer: writer,
	}
}

func (eb *ErrorBuffer) addLine() {
	eb.buf.WriteBytes('\n', ' ', 0xE2, 0x86, 0xB3, ' ') // ↳
}

func (eb *ErrorBuffer) IfErrAddErr(err error) {
	if err != nil {
		eb.addLine()
		eb.buf.WriteString(err.Error())
	}
}

func (eb *ErrorBuffer) IfErrAddErrWithStr(err error, format string, args ...any) {
	if err != nil {
		eb.addLine()
		eb.buf.WriteString(err.Error())
		eb.addLine()
		fmt.Fprintf(eb.buf, format, args...)
	}
}

func (eb *ErrorBuffer) AddStr(format string, args ...any) {
	eb.addLine()
	fmt.Fprintf(eb.buf, format, args...)
}

func (eb *ErrorBuffer) BytesRef() []byte {
	return eb.buf.BytesRef()
}

func (eb *ErrorBuffer) BytesCopy() []byte {
	return eb.buf.BytesCopy()
}

func (eb *ErrorBuffer) StringCopy() string {
	return eb.buf.StringCopy()
}

func (eb *ErrorBuffer) Error() string {
	return eb.buf.StringCopy()
}

func (eb *ErrorBuffer) Clear() {
	eb.buf.Reset()
}

func (eb *ErrorBuffer) Flush(prefix string) {
	if eb.Len() > 0 {
		eb.writer.Write([]byte(prefix))
		eb.writer.Write(eb.buf.BytesRef())
		eb.buf.Reset()
	}
}

func (eb *ErrorBuffer) Len() int {
	return eb.buf.Len()
}

func (eb *ErrorBuffer) Cap() int {
	return eb.buf.Cap()
}

func (eb *ErrorBuffer) EnsureSpace(pool *data_buffer.WriteBufferPool, space int) {
	eb.buf.EnsureSpace(space)
}

func (eb *ErrorBuffer) Close() error {
	return eb.buf.Close()
}

var QuickItoA = [100]string{
	"00", "01", "02", "03", "04", "05", "06", "07", "08", "09",
	"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
	"20", "21", "22", "23", "24", "25", "26", "27", "28", "29",
	"30", "31", "32", "33", "34", "35", "36", "37", "38", "39",
	"40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
	"50", "51", "52", "53", "54", "55", "56", "57", "58", "59",
	"60", "61", "62", "63", "64", "65", "66", "67", "68", "69",
	"70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99",
}

var quickItoX = [16]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'A', 'B', 'C', 'D', 'E', 'F',
}

var quickXtoI = [...]byte{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'A': 10,
	'B': 11,
	'C': 12,
	'D': 13,
	'E': 14,
	'F': 15,
}

func QuickIntToHexString(val uint64) (hexStr [16]byte, firstNonzeroChar int) {
	i := 15
	firstNonzeroChar = 15
	for i >= 0 {
		char := quickItoX[val&0b1111]
		hexStr[i] = char
		if char != '0' {
			firstNonzeroChar = i
		}
		val >>= 4
		i -= 1
	}
	return
}

func QuickHexStringToInt(hexStr string) uint64 {
	var val uint64
	i := 0
	for i < len(hexStr) {
		val <<= 4
		char := hexStr[i]
		if char >= '0' && char <= '9' {
			val |= uint64(char - '0')
		} else if char >= 'A' && char <= 'F' {
			val |= uint64(char - 'A')
		} else if char >= 'a' && char <= 'a' {
			val |= uint64(char - 'a')
		} else {
			return val
		}
		i += 1
	}
	return val
}

// Returns the index in the ordered data set where element `target`
// should exist, and bool `found` if the actual element at that index matches the target
func BinarySearch[T cmp.Ordered](data []T, target T) (idx int, found bool) {
	low := 0
	high := len(data) - 1
	for low <= high {
		idx := low + ((high - low) / 2)
		if data[idx] == target {
			return idx, true
		} else if data[idx] < target {
			low = idx + 1
		} else {
			high = idx - 1
		}
	}
	return low, data[low] == target
}
