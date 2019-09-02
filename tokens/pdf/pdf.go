package pdf

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/gobuffalo/packr/v2"
)

const StaticUrl = "http://abcdefghijklmnopqrstuvwxyz.zyxwvutsrqponmlkjihgfedcba.aceegikmoqsuwy.bdfhjlnprtvxz"
const PdfFile = "files/template.pdf"
const StreamOffset = 793

func assertSize(name string, want, is int) error {
	if want != is {
		e := fmt.Sprintf("assertSize() %s offset wants:%d != is:%d", name, want, is)
		return errors.New(e)
	}
	return nil
}

type PdfContext struct {
	Key    string
	Url    string
	Buffer bytes.Buffer
}

func NewPdfContext(url string, box *packr.Box) (pdfc *PdfContext, err error) {
	var bf bytes.Buffer
	var size int
	var bytesZlib []byte
	var retKey, retUrl string
	pdfc = &PdfContext{}
	size = 1

	pdfb, err := box.Find(PdfFile)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(".*\\/Length ([0-9]+)\\/.*")
	match := re.FindAll(pdfb[StreamOffset:], -1)

	lineLength := len(match[0]) - len("stream\n")
	if err := assertSize("lineLength", 67, lineLength); err != nil {
		return nil, err
	}

	stream_beg := StreamOffset + lineLength
	if err := assertSize("stream_beg", 860, stream_beg); err != nil {
		return nil, err
	}

	stream_start := stream_beg + 8
	if err := assertSize("stream_start", 868, stream_start); err != nil {
		return nil, err
	}

	stream_header := pdfb[StreamOffset:stream_start]

	submatch := re.FindAllSubmatch(pdfb[StreamOffset:], -1)
	s := fmt.Sprintf("%s", submatch[0][1])
	stream_size, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}

	if err := assertSize("stream_size", 110, stream_size); err != nil {
		return nil, err
	}

	stream := pdfb[stream_start : stream_start+stream_size]
	if err := assertSize("stream", 110, len(stream)); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	newreader := bytes.NewReader(stream)
	zr, err := zlib.NewReader(newreader)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(buf, zr)
	zr.Close()

	for len(bytesZlib) != len(buf.Bytes()) {
		retKey = RandomString(size)
		retUrl = url + "/t/" + retKey + "/i.png"
		rp := bytes.Replace([]byte(buf.String()), []byte(StaticUrl), []byte(retUrl), -1)

		var bf bytes.Buffer
		zw, err := zlib.NewWriterLevel(&bf, 6)
		if err != nil {
			log.Fatal(err)
		}
		_, err = zw.Write([]byte(rp))
		if err != nil {
			log.Fatal(err)
		}
		zw.Close()

		bytesZlib = bf.Bytes()
		if len(bytesZlib) > len(buf.Bytes()) {
			e := `
			This is a hacky method. We will try and fix it later and you shouldn't see this error. 
			The zlib compressed URL we are replacing is too large and won't fit in the pdf. 
			That shouldn't be the cast, BUT maybe you are using a HUGE domain name.
			If so, I'd like to understand more.
			`
			return nil, errors.New(e)
		}
		pdfc.Key = retKey
		pdfc.Url = retUrl
		size = size + 1
	}

	bw := bufio.NewWriter(&bf)

	_, err = bw.Write(pdfb[0:StreamOffset])
	if err != nil {
		return nil, err
	}

	err = bw.Flush()
	if err != nil {
		return nil, err
	}

	_, err = bw.Write(stream_header)
	if err != nil {
		return nil, err
	}

	err = bw.Flush()
	if err != nil {
		return nil, err
	}

	_, err = bw.Write(bytesZlib)
	if err != nil {
		return nil, err
	}

	err = bw.Flush()
	if err != nil {
		return nil, err
	}

	_, err = bw.Write(pdfb[stream_start+stream_size:])
	if err != nil {
		return nil, err
	}

	err = bw.Flush()
	if err != nil {
		return nil, err
	}

	pdfc.Buffer = bf
	return pdfc, nil

}

// RandomString isn't crypto secure or anything
// we just need a psuedo random key to be generated for our canary token
func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
