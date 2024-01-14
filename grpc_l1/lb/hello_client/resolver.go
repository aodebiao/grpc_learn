package main

import (
	"fmt"
	"google.golang.org/grpc/resolver"
)

// 自定义name resolver

const (
	myScheme   = "aodeibiao"
	myEndpoint = "resolver.aodeibiao.com"
)

var addrs = []string{"127.0.0.1:8991", "127.0.0.1:8080"}

// aodebiaoResolver 自定义name resolver，实现Resolver接口

type aodebiaoResolver struct {
	target    resolver.Target
	cc        resolver.ClientConn
	addrStore map[string][]string
}

func (a *aodebiaoResolver) ResolveNow(o resolver.ResolveNowOptions) {
	addrStrs := a.addrStore[a.target.Endpoint()]
	addrList := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrList[i] = resolver.Address{
			Addr: s,
		}
	}
	a.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (a *aodebiaoResolver) Close() {

}

// aodebiaoResolverBuilder 需实现 Builder 接口

type aodebiaoResolverBuilder struct {
}

func (a *aodebiaoResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &aodebiaoResolver{
		target: target,
		cc:     cc,
		addrStore: map[string][]string{
			myEndpoint: addrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (a *aodebiaoResolverBuilder) Scheme() string {
	return myScheme
}

func init() {
	// 注册 q1miResolverBuilder
	fmt.Printf("11111111")
	resolver.Register(&aodebiaoResolverBuilder{})
}
