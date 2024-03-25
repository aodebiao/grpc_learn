package main

import (
	"fmt"
	"os"
)

const defaultValueC = 1

type ServerConfig struct {
	A string
	B string
	C int
	X struct{}
	Y Info
}

type Info struct {
	addr string
}

func NewServerConfig1(a, b string, c int) *ServerConfig {
	return &ServerConfig{
		A: a,
		B: b,
		C: c,
	}
}

// NewServerConfig2 a,b参数必传,c可选

func NewServerConfig2(a, b string, c ...int) *ServerConfig {
	valueC := defaultValueC
	if len(c) > 0 {
		valueC = c[0]
	}
	return &ServerConfig{
		A: a,
		B: b,
		C: valueC,
	}
}

type FuncServiceConfigOption func(config *ServerConfig)

func NewConfigServer3(a, b string, opts ...FuncServiceConfigOption) *ServerConfig {
	sc := &ServerConfig{
		A: a,
		B: b,
		C: defaultValueC,
	}
	for _, opt := range opts {
		opt(sc)
	}
	return sc
}

// 针对可选配置实现专用方法
// WithC x
func WithC(c int) FuncServiceConfigOption {
	return func(sc *ServerConfig) {
		sc.C = c
	}
}

func WithInfo(info Info) FuncServiceConfigOption {
	return func(sc *ServerConfig) {
		sc.Y = info
	}
}
func newConfig(age int, opts ...ConfigOption) config {
	cfg := config{age: age}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	return cfg
}
func main() {
	//sc := NewConfigServer3("hello", "world")
	//fmt.Printf("sc:%#v\n", sc)

	// 普通 option模式
	//sc := NewConfigServer3("hello", "world", WithC(100))
	//fmt.Printf("sc:%#v\n", sc)

	//sc := NewConfigServer3("hello", "world", WithC(100), WithInfo(Info{
	//	addr: "华软",
	//}))
	//fmt.Printf("sc:%#v\n", sc)

	// 进阶option模式
	c := newConfig(1, WithConfigName("rest test"))
	fmt.Printf("%#v\n", c)
}

type config struct {
	name string
	age  int
}

type ConfigOption interface {
	apply(*config)
}

type funcOption struct {
	f func(*config)
}

func (f funcOption) apply(cfg *config) {
	f.f(cfg)
}
func newFuncOption(f func(c *config)) funcOption {
	return funcOption{f: f}
}
func WithConfigName(name string) ConfigOption {
	return newFuncOption(func(c *config) {
		c.name = name
	})
}

func test() {
	ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	go func() {

	}()
	<-ch
}
