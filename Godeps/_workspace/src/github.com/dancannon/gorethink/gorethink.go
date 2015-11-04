package gorethink

import (
	"reflect"

	"github.com/kylemclaren/mongo-transporter/Godeps/_workspace/src/github.com/dancannon/gorethink/encoding"
)

func init() {
	// Set encoding package
	encoding.IgnoreType(reflect.TypeOf(Term{}))
}
