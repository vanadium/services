// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: kvstore

// Package kvstore implements a simple key-value store used for
// testing the groups-based authorization.
//
//nolint:revive
package kvstore

import (
	v23 "v.io/v23"
	"v.io/v23/context"
	"v.io/v23/rpc"
	"v.io/v23/security/access"
	"v.io/v23/vdl"
)

var initializeVDLCalled = false
var _ = initializeVDL() // Must be first; see initializeVDL comments for details.

// Interface definitions
// =====================

// StoreClientMethods is the client interface
// containing Store methods.
type StoreClientMethods interface {
	Get(_ *context.T, key string, _ ...rpc.CallOpt) (string, error)
	Set(_ *context.T, key string, value string, _ ...rpc.CallOpt) error
}

// StoreClientStub embeds StoreClientMethods and is a
// placeholder for additional management operations.
type StoreClientStub interface {
	StoreClientMethods
}

// StoreClient returns a client stub for Store.
func StoreClient(name string) StoreClientStub {
	return implStoreClientStub{name}
}

type implStoreClientStub struct {
	name string
}

func (c implStoreClientStub) Get(ctx *context.T, i0 string, opts ...rpc.CallOpt) (o0 string, err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Get", []interface{}{i0}, []interface{}{&o0}, opts...)
	return
}

func (c implStoreClientStub) Set(ctx *context.T, i0 string, i1 string, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Set", []interface{}{i0, i1}, nil, opts...)
	return
}

// StoreServerMethods is the interface a server writer
// implements for Store.
type StoreServerMethods interface {
	Get(_ *context.T, _ rpc.ServerCall, key string) (string, error)
	Set(_ *context.T, _ rpc.ServerCall, key string, value string) error
}

// StoreServerStubMethods is the server interface containing
// Store methods, as expected by rpc.Server.
// There is no difference between this interface and StoreServerMethods
// since there are no streaming methods.
type StoreServerStubMethods StoreServerMethods

// StoreServerStub adds universal methods to StoreServerStubMethods.
type StoreServerStub interface {
	StoreServerStubMethods
	// DescribeInterfaces the Store interfaces.
	Describe__() []rpc.InterfaceDesc
}

// StoreServer returns a server stub for Store.
// It converts an implementation of StoreServerMethods into
// an object that may be used by rpc.Server.
func StoreServer(impl StoreServerMethods) StoreServerStub {
	stub := implStoreServerStub{
		impl: impl,
	}
	// Initialize GlobState; always check the stub itself first, to handle the
	// case where the user has the Glob method defined in their VDL source.
	if gs := rpc.NewGlobState(stub); gs != nil {
		stub.gs = gs
	} else if gs := rpc.NewGlobState(impl); gs != nil {
		stub.gs = gs
	}
	return stub
}

type implStoreServerStub struct {
	impl StoreServerMethods
	gs   *rpc.GlobState
}

func (s implStoreServerStub) Get(ctx *context.T, call rpc.ServerCall, i0 string) (string, error) {
	return s.impl.Get(ctx, call, i0)
}

func (s implStoreServerStub) Set(ctx *context.T, call rpc.ServerCall, i0 string, i1 string) error {
	return s.impl.Set(ctx, call, i0, i1)
}

func (s implStoreServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implStoreServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{StoreDesc}
}

// StoreDesc describes the Store interface.
var StoreDesc rpc.InterfaceDesc = descStore

// descStore hides the desc to keep godoc clean.
var descStore = rpc.InterfaceDesc{
	Name:    "Store",
	PkgPath: "v.io/x/ref/services/groups/groupsd/testdata/kvstore",
	Methods: []rpc.MethodDesc{
		{
			Name: "Get",
			InArgs: []rpc.ArgDesc{
				{Name: "key", Doc: ``}, // string
			},
			OutArgs: []rpc.ArgDesc{
				{Name: "", Doc: ``}, // string
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Read"))},
		},
		{
			Name: "Set",
			InArgs: []rpc.ArgDesc{
				{Name: "key", Doc: ``},   // string
				{Name: "value", Doc: ``}, // string
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Write"))},
		},
	},
}

// initializeVDL performs vdl initialization.  It is safe to call multiple times.
// If you have an init ordering issue, just insert the following line verbatim
// into your source files in this package, right after the "package foo" clause:
//
//	var _ = initializeVDL()
//
// The purpose of this function is to ensure that vdl initialization occurs in
// the right order, and very early in the init sequence.  In particular, vdl
// registration and package variable initialization needs to occur before
// functions like vdl.TypeOf will work properly.
//
// This function returns a dummy value, so that it can be used to initialize the
// first var in the file, to take advantage of Go's defined init order.
func initializeVDL() struct{} {
	if initializeVDLCalled {
		return struct{}{}
	}
	initializeVDLCalled = true

	return struct{}{}
}
