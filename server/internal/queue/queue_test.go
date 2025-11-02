package queue

import "testing"

func TestDequePushPopFrontBack(t *testing.T) {
	var d Deque[int]
	for i := 0; i < 5; i++ {
		d.PushBack(i)
	}
	if got := d.Len(); got != 5 {
		t.Fatalf("Len() = %d, want 5", got)
	}
	if front := d.Front(); front != 0 {
		t.Fatalf("Front() = %d, want 0", front)
	}
	if back := d.Back(); back != 4 {
		t.Fatalf("Back() = %d, want 4", back)
	}
	for i := 0; i < 5; i++ {
		if v := d.PopFront(); v != i {
			t.Fatalf("PopFront() = %d, want %d", v, i)
		}
	}
	if got := d.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}
}

func TestDequeWrapAndGrow(t *testing.T) {
	var d Deque[int]
	for i := 0; i < 3; i++ {
		d.PushBack(i)
	}
	for i := 0; i < 2; i++ {
		if v := d.PopFront(); v != i {
			t.Fatalf("PopFront() wrap = %d, want %d", v, i)
		}
	}
	for i := 3; i < 10; i++ {
		d.PushBack(i)
	}
	if front := d.Front(); front != 2 {
		t.Fatalf("Front() after wrap = %d, want 2", front)
	}
	if back := d.Back(); back != 9 {
		t.Fatalf("Back() after grow = %d, want 9", back)
	}
	for i := 2; i < 10; i++ {
		if v := d.PopFront(); v != i {
			t.Fatalf("PopFront() final = %d, want %d", v, i)
		}
	}
	if got := d.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}
}

func TestDequeClear(t *testing.T) {
	var d Deque[int]
	for i := 0; i < 4; i++ {
		d.PushBack(i)
	}
	d.Clear()
	if got := d.Len(); got != 0 {
		t.Fatalf("Len() after Clear = %d, want 0", got)
	}
	d.PushBack(10)
	if front := d.Front(); front != 10 {
		t.Fatalf("Front() after Clear = %d, want 10", front)
	}
}
