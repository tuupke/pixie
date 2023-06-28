package loader

type (
	Loadable interface {
		Load()
	}

	value[T any] struct {
		key    string
		toLoad T
	}

	propSet[T Loadable] struct {
		mp map[string]value[T]
	}

	Loader interface {
	}
)

func New[T Loadable](loadFuncs ...func()) Loader {
	return nil
}
