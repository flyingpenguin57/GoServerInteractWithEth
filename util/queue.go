package queue

// 定义泛型队列结构体
type Queue[T any] struct {
	items []T
	maxSize int
}

//返回头部元素
func (q *Queue[T]) GetHeader() (T) {
	return q.items[0]
}

//返回items
func (q *Queue[T]) GetItems() ([]T) {
	return q.items
}

// NewQueue 创建一个新的泛型队列并指定最大长度
func NewQueue[T any](maxSize int) *Queue[T] {
	return &Queue[T]{maxSize: maxSize}
}

// Enqueue 向队列头部添加一个元素，如果长度超出则删除尾部元素
func (q *Queue[T]) Enqueue(item T) {
	// 将新元素添加到队列头部
	q.items = append([]T{item}, q.items...)
	// 检查队列长度是否超出最大长度
	if len(q.items) > q.maxSize {
		// 删除尾部元素
		q.items = q.items[:q.maxSize]
	}
}

// Dequeue 从队列移除并返回尾部元素
func (q *Queue[T]) Dequeue() (T, bool) {
	var zeroValue T
	if len(q.items) == 0 {
		return zeroValue, false
	}
	// 获取并移除尾部元素
	item := q.items[len(q.items)-1]
	q.items = q.items[:len(q.items)-1]
	return item, true
}

// IsEmpty 检查队列是否为空
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size 返回队列的大小
func (q *Queue[T]) Size() int {
	return len(q.items)
}