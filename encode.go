package wavwrite

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

// Encode writes all audio streamed from s to w in WAVE format.
//
// Format precision must be 1 or 2 bytes.
func Encode(w io.WriteSeeker, s Streamer, format Format) (err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "wav")
		}
	}()

	if format.NumChannels <= 0 {
		return errors.New("wav: invalid number of channels (less than 1)")
	}
	if format.Precision != 1 && format.Precision != 2 {
		return errors.New("wav: unsupported precision, 1 or 2 is supported")
	}

	h := header{
		RiffMark:      [4]byte{'R', 'I', 'F', 'F'},
		FileSize:      -1, // finalization
		WaveMark:      [4]byte{'W', 'A', 'V', 'E'},
		FmtMark:       [4]byte{'f', 'm', 't', ' '},
		FormatSize:    16,
		FormatType:    1,
		NumChans:      int16(format.NumChannels),
		SampleRate:    int32(format.SampleRate),
		ByteRate:      int32(int(format.SampleRate) * format.NumChannels * format.Precision),
		BytesPerFrame: int16(format.NumChannels * format.Precision),
		BitsPerSample: int16(format.Precision) * 8,
		DataMark:      [4]byte{'d', 'a', 't', 'a'},
		DataSize:      -1, // finalization
	}
	if err := binary.Write(w, binary.LittleEndian, &h); err != nil {
		return err
	}

	var (
		bw      = bufio.NewWriter(w)
		buffer  = make([]byte, 512*format.Width())
		written int
	)
	for {
		n, ok := s.Stream(buffer)
		if !ok {
			break
		}
		nn, err := bw.Write(buffer[:n])
		if err != nil {
			return err
		}
		written += nn
	}
	if err := bw.Flush(); err != nil {
		return err
	}

	// finalize header
	h.FileSize = int32(44 + written) // 44 is the size of the header
	h.DataSize = int32(written)
	if _, err := w.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h); err != nil {
		return err
	}
	if _, err := w.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	return nil
}

type header struct {
	RiffMark      [4]byte
	FileSize      int32
	WaveMark      [4]byte
	FmtMark       [4]byte
	FormatSize    int32
	FormatType    int16
	NumChans      int16
	SampleRate    int32
	ByteRate      int32
	BytesPerFrame int16
	BitsPerSample int16
	DataMark      [4]byte
	DataSize      int32
}
