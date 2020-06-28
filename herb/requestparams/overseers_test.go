package requestparams_test

import (
	"testing"

	"github.com/herb-go/requestparams"
	"github.com/herb-go/worker"
)

func TestOverseer(t *testing.T) {
	var err error
	o := &worker.PlainOverseer{}
	err = requestparams.NewReaderFactoryOverseerConfig().ApplyTo(o)
	if err != nil {
		t.Fatal(err)
	}
	err = requestparams.NewReaderOverseerConfig().ApplyTo(o)
	if err != nil {
		t.Fatal(err)
	}
	err = requestparams.NewFormaterFactoryOverseerConfig().ApplyTo(o)
	if err != nil {
		t.Fatal(err)
	}
	err = requestparams.NewFormaterOverseerConfig().ApplyTo(o)
	if err != nil {
		t.Fatal(err)
	}
}
