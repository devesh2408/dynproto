// symbol.go -- symbol types from protoc-gen-go

package proto

import "strings"

// symbol is an interface representing an exported Go symbol.
type symbol interface {
	// GenerateAlias should generate an appropriate alias
	// for the symbol from the named package.
	GenerateAlias(g *Generator, pkg string)
}

type messageSymbol struct {
	sym                         string
	hasExtensions, isMessageSet bool
	hasOneof                    bool
	getters                     []getterSymbol
}

type getterSymbol struct {
	name     string
	typ      string
	typeName string // canonical name in proto world; empty for proto.Message and similar
	genType  bool   // whether typ contains a generated type (message/group/enum)
}

func (ms *messageSymbol) GenerateAlias(g *Generator, pkg string) {
	remoteSym := pkg + "." + ms.sym

	g.P("type ", ms.sym, " ", remoteSym)
	g.P("func (m *", ms.sym, ") Reset() { (*", remoteSym, ")(m).Reset() }")
	g.P("func (m *", ms.sym, ") String() string { return (*", remoteSym, ")(m).String() }")
	g.P("func (*", ms.sym, ") ProtoMessage() {}")
	if ms.hasExtensions {
		g.P("func (*", ms.sym, ") ExtensionRangeArray() []", g.Pkg["proto"], ".ExtensionRange ",
			"{ return (*", remoteSym, ")(nil).ExtensionRangeArray() }")
		if ms.isMessageSet {
			g.P("func (m *", ms.sym, ") Marshal() ([]byte, error) ",
				"{ return (*", remoteSym, ")(m).Marshal() }")
			g.P("func (m *", ms.sym, ") Unmarshal(buf []byte) error ",
				"{ return (*", remoteSym, ")(m).Unmarshal(buf) }")
		}
	}
	if ms.hasOneof {
		// Oneofs and public imports do not mix well.
		// We can make them work okay for the binary format,
		// but they're going to break weirdly for text/JSON.
		enc := "_" + ms.sym + "_OneofMarshaler"
		dec := "_" + ms.sym + "_OneofUnmarshaler"
		size := "_" + ms.sym + "_OneofSizer"
		encSig := "(msg " + g.Pkg["proto"] + ".Message, b *" + g.Pkg["proto"] + ".Buffer) error"
		decSig := "(msg " + g.Pkg["proto"] + ".Message, tag, wire int, b *" + g.Pkg["proto"] + ".Buffer) (bool, error)"
		sizeSig := "(msg " + g.Pkg["proto"] + ".Message) int"
		g.P("func (m *", ms.sym, ") XXX_OneofFuncs() (func", encSig, ", func", decSig, ", func", sizeSig, ", []interface{}) {")
		g.P("return ", enc, ", ", dec, ", ", size, ", nil")
		g.P("}")

		g.P("func ", enc, encSig, " {")
		g.P("m := msg.(*", ms.sym, ")")
		g.P("m0 := (*", remoteSym, ")(m)")
		g.P("enc, _, _, _ := m0.XXX_OneofFuncs()")
		g.P("return enc(m0, b)")
		g.P("}")

		g.P("func ", dec, decSig, " {")
		g.P("m := msg.(*", ms.sym, ")")
		g.P("m0 := (*", remoteSym, ")(m)")
		g.P("_, dec, _, _ := m0.XXX_OneofFuncs()")
		g.P("return dec(m0, tag, wire, b)")
		g.P("}")

		g.P("func ", size, sizeSig, " {")
		g.P("m := msg.(*", ms.sym, ")")
		g.P("m0 := (*", remoteSym, ")(m)")
		g.P("_, _, size, _ := m0.XXX_OneofFuncs()")
		g.P("return size(m0)")
		g.P("}")
	}
	for _, get := range ms.getters {

		if get.typeName != "" {
			g.RecordTypeUse(get.typeName)
		}
		typ := get.typ
		val := "(*" + remoteSym + ")(m)." + get.name + "()"
		if get.genType {
			// typ will be "*pkg.T" (message/group) or "pkg.T" (enum)
			// or "map[t]*pkg.T" (map to message/enum).
			// The first two of those might have a "[]" prefix if it is repeated.
			// Drop any package qualifier since we have hoisted the type into this package.
			rep := strings.HasPrefix(typ, "[]")
			if rep {
				typ = typ[2:]
			}
			isMap := strings.HasPrefix(typ, "map[")
			star := typ[0] == '*'
			if !isMap { // map types handled lower down
				typ = typ[strings.Index(typ, ".")+1:]
			}
			if star {
				typ = "*" + typ
			}
			if rep {
				// Go does not permit conversion between slice types where both
				// element types are named. That means we need to generate a bit
				// of code in this situation.
				// typ is the element type.
				// val is the expression to get the slice from the imported type.

				ctyp := typ // conversion type expression; "Foo" or "(*Foo)"
				if star {
					ctyp = "(" + typ + ")"
				}

				g.P("func (m *", ms.sym, ") ", get.name, "() []", typ, " {")
				g.In()
				g.P("o := ", val)
				g.P("if o == nil {")
				g.In()
				g.P("return nil")
				g.Out()
				g.P("}")
				g.P("s := make([]", typ, ", len(o))")
				g.P("for i, x := range o {")
				g.In()
				g.P("s[i] = ", ctyp, "(x)")
				g.Out()
				g.P("}")
				g.P("return s")
				g.Out()
				g.P("}")
				continue
			}
			if isMap {
				// Split map[keyTyp]valTyp.
				bra, ket := strings.Index(typ, "["), strings.Index(typ, "]")
				keyTyp, valTyp := typ[bra+1:ket], typ[ket+1:]
				// Drop any package qualifier.
				// Only the value type may be foreign.
				star := valTyp[0] == '*'
				valTyp = valTyp[strings.Index(valTyp, ".")+1:]
				if star {
					valTyp = "*" + valTyp
				}

				typ := "map[" + keyTyp + "]" + valTyp
				g.P("func (m *", ms.sym, ") ", get.name, "() ", typ, " {")
				g.P("o := ", val)
				g.P("if o == nil { return nil }")
				g.P("s := make(", typ, ", len(o))")
				g.P("for k, v := range o {")
				g.P("s[k] = (", valTyp, ")(v)")
				g.P("}")
				g.P("return s")
				g.P("}")
				continue
			}
			// Convert imported type into the forwarding type.
			val = "(" + typ + ")(" + val + ")"
		}

		g.P("func (m *", ms.sym, ") ", get.name, "() ", typ, " { return ", val, " }")
	}

}

type enumSymbol struct {
	name   string
	proto3 bool // Whether this came from a proto3 file.
}

func (es enumSymbol) GenerateAlias(g *Generator, pkg string) {
	s := es.name
	g.P("type ", s, " ", pkg, ".", s)
	g.P("var ", s, "_name = ", pkg, ".", s, "_name")
	g.P("var ", s, "_value = ", pkg, ".", s, "_value")
	g.P("func (x ", s, ") String() string { return (", pkg, ".", s, ")(x).String() }")
	if !es.proto3 {
		g.P("func (x ", s, ") Enum() *", s, "{ return (*", s, ")((", pkg, ".", s, ")(x).Enum()) }")
		g.P("func (x *", s, ") UnmarshalJSON(data []byte) error { return (*", pkg, ".", s, ")(x).UnmarshalJSON(data) }")
	}
}

type constOrVarSymbol struct {
	sym  string
	typ  string // either "const" or "var"
	cast string // if non-empty, a type cast is required (used for enums)
}

func (cs constOrVarSymbol) GenerateAlias(g *Generator, pkg string) {
	v := pkg + "." + cs.sym
	if cs.cast != "" {
		v = cs.cast + "(" + v + ")"
	}
	g.P(cs.typ, " ", cs.sym, " = ", v)
}
