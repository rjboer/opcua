package queue

// Deque is a generic double-ended queue backed by a ring buffer.
type Deque[T any] struct {
	buf   []T
	head  int
	tail  int
	count int
}

// Len returns the number of elements stored in the deque.
func (d *Deque[T]) Len() int {
	return d.count
}

// PushBack appends v to the tail of the deque.
func (d *Deque[T]) PushBack(v T) {
	if d.count == len(d.buf) {
		d.grow()
	}
	d.buf[d.tail] = v
	d.tail++
	if d.tail == len(d.buf) {
		d.tail = 0
	}
	d.count++
}

// PopFront removes and returns the element at the head of the deque.
func (d *Deque[T]) PopFront() T {
	if d.count == 0 {
		panic("queue: PopFront called on empty deque")
	}
	v := d.buf[d.head]
	var zero T
	d.buf[d.head] = zero
	d.head++
	if d.head == len(d.buf) {
		d.head = 0
	}
	d.count--
	return v
}

// PopBack removes and returns the element at the tail of the deque.
func (d *Deque[T]) PopBack() T {
	if d.count == 0 {
		panic("queue: PopBack called on empty deque")
	}
	d.tail--
	if d.tail < 0 {
		d.tail = len(d.buf) - 1
	}
	v := d.buf[d.tail]
	var zero T
	d.buf[d.tail] = zero
	d.count--
	return v
}

// Front returns the element at the head of the deque without removing it.
func (d *Deque[T]) Front() T {
	if d.count == 0 {
		panic("queue: Front called on empty deque")
	}
	return d.buf[d.head]
}

// Back returns the element at the tail of the deque without removing it.
func (d *Deque[T]) Back() T {
	if d.count == 0 {
		panic("queue: Back called on empty deque")
	}
	idx := d.tail - 1
	if idx < 0 {
		idx = len(d.buf) - 1
	}
	return d.buf[idx]
}

// Clear removes all elements from the deque.
func (d *Deque[T]) Clear() {
	if d.count == 0 {
		return
	}
	var zero T
	for i := 0; i < d.count; i++ {
		idx := (d.head + i) % len(d.buf)
		d.buf[idx] = zero
	}
	d.head = 0
	d.tail = 0
	d.count = 0
}

func (d *Deque[T]) grow() {
	newCap := len(d.buf)
	if newCap == 0 {
		newCap = 1
	} else {
		newCap *= 2
	}
	newBuf := make([]T, newCap)
	if d.count > 0 {
		if d.head < d.tail {
			copy(newBuf, d.buf[d.head:d.tail])
		} else {
			n := copy(newBuf, d.buf[d.head:])
			copy(newBuf[n:], d.buf[:d.tail])
		}
	}
	d.buf = newBuf
	d.head = 0
	d.tail = d.count
}
