package compiler

type SymbolScope string

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer       *SymbolTable
	FreeSymbols []Symbol

	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Outer:          nil,
		store:          make(map[string]Symbol),
		numDefinitions: 0,
		FreeSymbols:    []Symbol{},
	}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Outer:          outer,
		store:          make(map[string]Symbol),
		numDefinitions: 0,
		FreeSymbols:    []Symbol{},
	}
}

func (s *SymbolTable) Define(name string) *Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	s.store[name] = symbol
	s.numDefinitions++

	return &symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if ok {
		return obj, ok
	}

	obj, ok = s.Outer.Resolve(name)
	if ok {
		return obj, ok
	}

	return Symbol{}, false
}
