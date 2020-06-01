package container

// モジュールの生成の方法を表す関数
type Builder interface{}

// モジュール名と生成の方法の構造体
type Definition struct {
	Name    string
	Builder Builder
}

// DI Container
type Container struct {
	store map[string]Builder
}

func NewContainer() *Container {
	return &Container{
		store: map[string]Builder{},
	}
}

// DI Containerにモジュールを登録する
func (c *Container) Register(d *Definition) {
	c.store[d.Name] = d.Builder
}

// DI Containerからモジュールを取り出す
func (c *Container) Get(key string) interface{} {
	builder, _ := c.store[key]
	return builder
}
