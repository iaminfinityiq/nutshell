package objects

func Mark(scope *Scope, zero_values *ZeroValues) {
	scope.Heap.Unmark()
	scope.Mark()
	zero_values.Hashmap.Mark()
}

func Sweep(heap *Heap) {
	heap.Sweep()
}
