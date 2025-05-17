package objects

func Mark(scope *Scope) {
	scope.Heap.Unmark()
	scope.Mark()
}

func Sweep(heap *Heap) {
	heap.Sweep()
}
