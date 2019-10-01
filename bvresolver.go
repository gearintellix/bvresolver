package bvresolver

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

const (
	scheme = "bv"
)

func NewBivrostResolver(namespace string) (resolver.Builder, error) {
	cln, err := newK8sClient(namespace)
	if err != nil {
		return nil, err
	}

	stop := make(chan bool)
	err = cln.watchEndpoints("koinp2p-svc", stop)
	if err != nil {
		return nil, err
	}

	return bvResolverBuilder{
		namespace: namespace,
	}, nil
}

type bvResolverBuilder struct {
	namespace string
}

func (b bvResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (res resolver.Resolver, err error) {
	if target.Scheme != scheme {
		err = fmt.Errorf("Cannot solve schema %s", target.Scheme)
		return res, err
	}

	return nil, nil
}

func (ox bvResolverBuilder) Scheme() string {
	return scheme
}

type bvResolver struct {
	service   string
	cc        resolver.ClientConn
	endpoints []resolver.Address
	stop      chan bool
}

func (ox bvResolver) ResolveNow(opt resolver.ResolveNowOption) {
	//
}

func (ox *bvResolver) Close() {
	ox.stop <- true
}

func (ox *bvResolver) Start() {
}
