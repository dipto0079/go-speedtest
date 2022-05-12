
package speedtest

import "io"
import "fmt"

// A JunkReader produces junk-ish data
type JunkReader struct {
	Data []byte
	Size int
	Pos  int
}


func NewJunkReader(size int) JunkReader {
	return JunkReader{
		Size: size,
		Pos:  0,
	}
}


func (r *JunkReader) Read(p []byte) (n int, err error) {
	for {
		if r.Size >= 0 && r.Pos >= r.Size {
			// Outta data
			return n, io.EOF
		} else if n < len(p) {
			p[n] = byte(n)
			n++
			r.Pos++
		} else {
			// Outta buffer space
			return n, err
		}
	}
}


type CallbackWriter struct {
	Callback func(n int) error
}

func NewCallbackWriter(callback func(n int) error) CallbackWriter {
	return CallbackWriter{callback}
}


func (b CallbackWriter) Write(p []byte) (n int, err error) {
	return len(p), b.Callback(len(p))
}


func NiceRate(rate int) string {
	bps := float64(8 * rate)
	kbps := bps / 1024
	mbps := kbps / 1024

	if mbps > 0.1 {
		return fmt.Sprintf("%.2fmbps", mbps)
	} else if kbps > 0.1 {
		return fmt.Sprintf("%.2fkbps", kbps)
	} else {
		return fmt.Sprintf("%.2fbps", bps)
	}
}
