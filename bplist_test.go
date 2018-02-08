package plist

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"testing"

	"howett.net/plist/cf"
)

func BenchmarkBplistGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := newBplistGenerator(ioutil.Discard)
		d.generateDocument(plistValueTree)
	}
}

func BenchmarkBplistParse(b *testing.B) {
	buf := bytes.NewReader(plistValueTreeAsBplist)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		d := newBplistParser(buf)
		d.parseDocument()
		b.StopTimer()
		buf.Seek(0, 0)
	}
}

func TestBplistInt128(t *testing.T) {
	bplist := []byte{0x62, 0x70, 0x6c, 0x69, 0x73, 0x74, 0x30, 0x30, 0x14, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x19}
	expected := uint64(0x090a0b0c0d0e0f10)
	buf := bytes.NewReader(bplist)
	d := newBplistParser(buf)
	pval, _ := d.parseDocument()
	if pinteger, ok := pval.(*cf.Number); !ok || pinteger.Value != expected {
		t.Error("Expected", expected, "received", pval)
	}
}

func TestBplistLatin1ToUTF16(t *testing.T) {
	expectedPrefix := []byte{0x62, 0x70, 0x6c, 0x69, 0x73, 0x74, 0x30, 0x30, 0xd1, 0x01, 0x02, 0x51, 0x5f, 0x6f, 0x10, 0x80}
	expectedPostfix := []byte{0x00, 0x08, 0x00, 0x0b, 0x00, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x10}
	expectedBuf := bytes.NewBuffer(expectedPrefix)

	sBuf := &bytes.Buffer{}
	for i := uint16(0xc280); i <= 0xc2bf; i++ {
		binary.Write(sBuf, binary.BigEndian, i)
		binary.Write(expectedBuf, binary.BigEndian, i-0xc200)
	}

	for i := uint16(0xc380); i <= 0xc3bf; i++ {
		binary.Write(sBuf, binary.BigEndian, i)
		binary.Write(expectedBuf, binary.BigEndian, i-0xc300+0x0040)
	}

	expectedBuf.Write(expectedPostfix)

	var buf bytes.Buffer
	encoder := NewBinaryEncoder(&buf)

	data := map[string]string{
		"_": string(sBuf.Bytes()),
	}
	if err := encoder.Encode(data); err != nil {
		t.Error(err.Error())
	}

	if !bytes.Equal(buf.Bytes(), expectedBuf.Bytes()) {
		t.Error("Expected", expectedBuf.Bytes(), "received", buf.Bytes())
		return
	}
}
