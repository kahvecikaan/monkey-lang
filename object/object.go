package object

import (
	"bytes"
	"fmt"
	"github.com/kahvecikaan/monkey-lang/ast"
	"hash/fnv"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

var (
	TRUE  = &Boolean{Value: true, hashKey: nil}
	FALSE = &Boolean{Value: false, hashKey: nil}
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type Integer struct {
	Value   int64
	hashKey *HashKey //Private field to store the cached hash key
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	// If we have a cached value, return it
	if i.hashKey != nil {
		return *i.hashKey
	}

	// Otherwise, compute the hash
	hash := HashKey{Type: i.Type(), Value: uint64(i.Value)}
	// Cache if for future use
	i.hashKey = &hash
	return hash
}
func NewInteger(value int64) *Integer {
	return &Integer{Value: value, hashKey: nil}
}

type Boolean struct {
	Value   bool
	hashKey *HashKey //Private field to store the cached hash key
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	// If we have a cached value, return it
	if b.hashKey != nil {
		return *b.hashKey
	}

	// Otherwise, compute the hash
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	hash := HashKey{Type: b.Type(), Value: value}
	// Cache it for future use
	b.hashKey = &hash
	return hash
}
func GetBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewClosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment // to allow for closures
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value   string
	hashKey *HashKey //Private field to store the cached hash key
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	// If we have a cached value, return it
	if s.hashKey != nil {
		return *s.hashKey
	}

	// Otherwise compute the hash key using fnv hash
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	hash := HashKey{Type: s.Type(), Value: h.Sum64()}
	// Cache it for future use
	s.hashKey = &hash
	return hash
}
func NewString(value string) *String {
	return &String{Value: value, hashKey: nil}
}

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

func compareObjects(a, b Object) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a := a.(type) {
	case *String:
		return a.Value == b.(*String).Value
	case *Integer:
		return a.Value == b.(*Integer).Value
	case *Boolean:
		return a.Value == b.(*Boolean).Value
	default:
		return false
	}
}

type HashChain []HashPair

func (chain HashChain) FindPair(key Object) (HashPair, bool) {
	for _, pair := range chain {
		if compareObjects(pair.Key, key) {
			return pair, true
		}
	}

	return HashPair{}, false
}

// Hash uses HashKey as the map key rather that just using the hash (uint64) directly because it prevents
// collisions between different types.
type Hash struct {
	Pairs map[HashKey]HashChain
}

func NewHash() *Hash {
	return &Hash{Pairs: make(map[HashKey]HashChain)}
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, chain := range h.Pairs {
		for _, pair := range chain {
			pairs = append(pairs, fmt.Sprintf("%s: %s",
				pair.Key.Inspect(), pair.Value.Inspect()))
		}
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// Add adds or updates a key-value pair in the hash table.
// If the key already exists, its value is updated.
// If the key hashes to an existing value but is different, it's added to the chain.
func (h *Hash) Add(key, value Object) error {
	hashKey, ok := key.(Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", key.Type())
	}

	hashed := hashKey.HashKey()
	chain := h.Pairs[hashed]
	newPair := HashPair{Key: key, Value: value}

	// Check if we're updating an existing key in the chain
	for i, pair := range chain {
		if compareObjects(pair.Key, key) {
			chain[i] = newPair
			h.Pairs[hashed] = chain
			return nil
		}
	}

	// If key wasn't found, append to chain
	chain = append(chain, newPair)
	h.Pairs[hashed] = chain
	return nil
}
