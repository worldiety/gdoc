package api

import "sort"

type Import string

type Imports []Import

func (p Imports) Len() int           { return len(p) }
func (p Imports) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Imports) Less(i, j int) bool { return p[i] < p[j] }

func (p Imports) Sort() {
	sort.Sort(p)
}
