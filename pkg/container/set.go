// 非线程安全set集合
package container

import "fmt"

type Set[T comparable] map[T]void
type void struct{}

func NewSet[E comparable](items []E) *Set[E] {
	// 获取Set的地址
	s := Set[E]{}
	// 声明map类型的数据结构
	for _, item := range items {
		// unsafe.Sizeof(void{}) // 结果为0
		s[item] = void{}
	}
	return &s
}

func (s *Set[T]) Add(i T) bool {
	if _, ok := (*s)[i]; ok {
		return false
	}

	// unsafe.Sizeof(void{}) // 结果为0
	(*s)[i] = void{}
	return true
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := (*s)[item]
	return ok
}

func (s *Set[T]) Size() int {
	return len(*s)
}

func (s *Set[T]) Elem() []T {
	keys := make([]T, 0, len(*s))
	for k := range *s {
		keys = append(keys, k)
	}
	return keys
}

func (s *Set[T]) Remove(i T) {
	delete(*s, i)
}

func (s *Set[T]) Clear() {
	*s = make(map[T]void)
}

// 相等
func (s *Set[T]) Equal(other *Set[T]) bool {
	// _ = other.(*Set)
	// 如果两者Size不相等，就不用比较了
	if s.Size() != other.Size() {
		return false
	}
	for elem := range *s {
		if !other.Contains(elem) {
			return false
		}
	}
	return true

}

// 子集
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	// s的size长于other，不用说了
	if s.Size() > other.Size() {
		return false
	}
	// 迭代遍历
	for key := range *s {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

// 并集
func (set *Set[T]) Union(other *Set[T]) *Set[T] {
	unionedSet := Set[T]{}
	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *other {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

// 交集
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {

	intersection := Set[T]{}
	// loop over smaller set
	if s.Size() < other.Size() {
		for elem := range *s {
			if other.Contains(elem) {
				intersection.Add(elem)
			}
		}
	} else {
		for elem := range *other {
			if s.Contains(elem) {
				intersection.Add(elem)
			}
		}
	}
	return &intersection
}

// 差集
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {

	difference := Set[T]{}
	for elem := range *s {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

// 函数遍历
func (s *Set[T]) Each(cb func(interface{}) bool) {
	for elem := range *s {
		if cb(elem) {
			break
		}
	}
}

// 返回string数组
func (s *Set[T]) StringElem() []string {
	items := make([]string, 0, len(*s))

	for elem := range *s {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return items
}
