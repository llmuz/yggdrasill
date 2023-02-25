package uzapimpl

import (
	"testing"
)

func TestParserLevel(t *testing.T) {

	level := ParserLevel("not found")
	if level != Info {
		t.Error("should be info, but get ", level)
		return
	}
	level = ParserLevel("Silent")
	if level != Silent {
		t.Error("should be silent, but get ", level)
	}

	level = ParserLevel("InFo")
	if level != Info {
		t.Error("should be info, but get ", level)
	}

	level = ParserLevel("warn")
	if level != Warn {
		t.Error("should be warn, but get ", level)
	}

	level = ParserLevel("errOr")
	if level != Error {
		t.Error("should be error, but get ", level)
	}
}
