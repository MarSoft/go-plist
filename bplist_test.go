package plist

import (
	"bytes"
	"testing"
)

func BenchmarkBplistGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := newBplistGenerator(nilWriter(0))
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
	pval := d.parseDocument()
	if pval.kind != Integer || pval.value.(uint64) != expected {
		t.Error("Expected", expected, "received", pval.value)
	}
}

func TestVariousIllegalBplists(t *testing.T) {
	bplists := [][]byte{
		[]byte{0x62, 0x70, 0x6c, 0x69, 0x73, 0x74, 0x30, 0x30, 0x15, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x19},
		[]byte{0x62, 0x70, 0x6c, 0x69, 0x73, 0x74, 0x30, 0x30, 0x24, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x19},
		[]byte{0x62, 0x70, 0x6c, 0x69, 0x73, 0x74, 0x30, 0x30, 0xFF, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x19},
		[]byte{0x62, 0x70, 0x6c, 0x69, 0x73, 0x74, 0x40, 0x41},
		[]byte{0x62, 0x71, 0x6c, 0x69, 0x73, 0x74, 0x30, 0x30},
	}

	testDecode := func(bplist []byte) (e error) {
		defer func() {
			if err := recover(); err != nil {
				e = err.(error)
			}
		}()
		buf := bytes.NewReader(bplist)
		d := newBplistParser(buf)
		d.parseDocument()
		return nil
	}

	for _, bplist := range bplists {
		err := testDecode(bplist)
		t.Logf("Error: %v", err)
		if err == nil {
			t.Error("Expected error, received nothing.")
		}
	}
}