package OrderedMap

type nocopy struct{}

func (n *nocopy) Lock()   {}
func (n *nocopy) Unlock() {}

type OrderMapInterface interface {
	Add(key, value string)
	Del(key string)
	Range() []string
}
