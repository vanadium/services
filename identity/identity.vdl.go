// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: identity

// Package identity defines interfaces for Vanadium identity providers.
//
//nolint:revive
package identity

import (
	v23 "v.io/v23"
	"v.io/v23/context"
	"v.io/v23/rpc"
	"v.io/v23/security"
	"v.io/v23/vdl"
)

var initializeVDLCalled = false
var _ = initializeVDL() // Must be first; see initializeVDL comments for details.

// Hold type definitions in package-level variables, for better performance.
// Declare and initialize with default values here so that the initializeVDL
// method will be considered ready to initialize before any of the type
// definitions that appear below.
//
//nolint:unused
var (
	vdlTypeStruct1 *vdl.Type = nil
	vdlTypeList2   *vdl.Type = nil
)

// Type definitions
// ================
// BlessingRootResponse is the struct representing the JSON response provided
// by the "blessing-root" route of the identity service.
type BlessingRootResponse struct {
	// Names of the blessings.
	Names []string
	// Base64 der-encoded public key.
	PublicKey string
}

func (BlessingRootResponse) VDLReflect(struct {
	Name string `vdl:"v.io/x/ref/services/identity.BlessingRootResponse"`
}) {
}

func (x BlessingRootResponse) VDLIsZero() bool { //nolint:gocyclo
	if len(x.Names) != 0 {
		return false
	}
	if x.PublicKey != "" {
		return false
	}
	return true
}

func (x BlessingRootResponse) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct1); err != nil {
		return err
	}
	if len(x.Names) != 0 {
		if err := enc.NextField(0); err != nil {
			return err
		}
		if err := vdlWriteAnonList1(enc, x.Names); err != nil {
			return err
		}
	}
	if x.PublicKey != "" {
		if err := enc.NextFieldValueString(1, vdl.StringType, x.PublicKey); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func vdlWriteAnonList1(enc vdl.Encoder, x []string) error {
	if err := enc.StartValue(vdlTypeList2); err != nil {
		return err
	}
	if err := enc.SetLenHint(len(x)); err != nil {
		return err
	}
	for _, elem := range x {
		if err := enc.NextEntryValueString(vdl.StringType, elem); err != nil {
			return err
		}
	}
	if err := enc.NextEntry(true); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *BlessingRootResponse) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = BlessingRootResponse{}
	if err := dec.StartValue(vdlTypeStruct1); err != nil {
		return err
	}
	decType := dec.Type()
	for {
		index, err := dec.NextField()
		switch {
		case err != nil:
			return err
		case index == -1:
			return dec.FinishValue()
		}
		if decType != vdlTypeStruct1 {
			index = vdlTypeStruct1.FieldIndexByName(decType.Field(index).Name)
			if index == -1 {
				if err := dec.SkipValue(); err != nil {
					return err
				}
				continue
			}
		}
		switch index {
		case 0:
			if err := vdlReadAnonList1(dec, &x.Names); err != nil {
				return err
			}
		case 1:
			switch value, err := dec.ReadValueString(); {
			case err != nil:
				return err
			default:
				x.PublicKey = value
			}
		}
	}
}

func vdlReadAnonList1(dec vdl.Decoder, x *[]string) error {
	if err := dec.StartValue(vdlTypeList2); err != nil {
		return err
	}
	if len := dec.LenHint(); len > 0 {
		*x = make([]string, 0, len)
	} else {
		*x = nil
	}
	for {
		switch done, elem, err := dec.NextEntryValueString(); {
		case err != nil:
			return err
		case done:
			return dec.FinishValue()
		default:
			*x = append(*x, elem)
		}
	}
}

// Interface definitions
// =====================

// MacaroonBlesserClientMethods is the client interface
// containing MacaroonBlesser methods.
//
// MacaroonBlesser returns a blessing given the provided macaroon string.
type MacaroonBlesserClientMethods interface {
	// Bless uses the provided macaroon (which contains email and caveats)
	// to return a blessing for the client.
	Bless(_ *context.T, macaroon string, _ ...rpc.CallOpt) (blessing security.Blessings, _ error)
}

// MacaroonBlesserClientStub embeds MacaroonBlesserClientMethods and is a
// placeholder for additional management operations.
type MacaroonBlesserClientStub interface {
	MacaroonBlesserClientMethods
}

// MacaroonBlesserClient returns a client stub for MacaroonBlesser.
func MacaroonBlesserClient(name string) MacaroonBlesserClientStub {
	return implMacaroonBlesserClientStub{name}
}

type implMacaroonBlesserClientStub struct {
	name string
}

func (c implMacaroonBlesserClientStub) Bless(ctx *context.T, i0 string, opts ...rpc.CallOpt) (o0 security.Blessings, err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Bless", []interface{}{i0}, []interface{}{&o0}, opts...)
	return
}

// MacaroonBlesserServerMethods is the interface a server writer
// implements for MacaroonBlesser.
//
// MacaroonBlesser returns a blessing given the provided macaroon string.
type MacaroonBlesserServerMethods interface {
	// Bless uses the provided macaroon (which contains email and caveats)
	// to return a blessing for the client.
	Bless(_ *context.T, _ rpc.ServerCall, macaroon string) (blessing security.Blessings, _ error)
}

// MacaroonBlesserServerStubMethods is the server interface containing
// MacaroonBlesser methods, as expected by rpc.Server.
// There is no difference between this interface and MacaroonBlesserServerMethods
// since there are no streaming methods.
type MacaroonBlesserServerStubMethods MacaroonBlesserServerMethods

// MacaroonBlesserServerStub adds universal methods to MacaroonBlesserServerStubMethods.
type MacaroonBlesserServerStub interface {
	MacaroonBlesserServerStubMethods
	// DescribeInterfaces the MacaroonBlesser interfaces.
	Describe__() []rpc.InterfaceDesc
}

// MacaroonBlesserServer returns a server stub for MacaroonBlesser.
// It converts an implementation of MacaroonBlesserServerMethods into
// an object that may be used by rpc.Server.
func MacaroonBlesserServer(impl MacaroonBlesserServerMethods) MacaroonBlesserServerStub {
	stub := implMacaroonBlesserServerStub{
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

type implMacaroonBlesserServerStub struct {
	impl MacaroonBlesserServerMethods
	gs   *rpc.GlobState
}

func (s implMacaroonBlesserServerStub) Bless(ctx *context.T, call rpc.ServerCall, i0 string) (security.Blessings, error) {
	return s.impl.Bless(ctx, call, i0)
}

func (s implMacaroonBlesserServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implMacaroonBlesserServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{MacaroonBlesserDesc}
}

// MacaroonBlesserDesc describes the MacaroonBlesser interface.
var MacaroonBlesserDesc rpc.InterfaceDesc = descMacaroonBlesser

// descMacaroonBlesser hides the desc to keep godoc clean.
var descMacaroonBlesser = rpc.InterfaceDesc{
	Name:    "MacaroonBlesser",
	PkgPath: "v.io/x/ref/services/identity",
	Doc:     "// MacaroonBlesser returns a blessing given the provided macaroon string.",
	Methods: []rpc.MethodDesc{
		{
			Name: "Bless",
			Doc:  "// Bless uses the provided macaroon (which contains email and caveats)\n// to return a blessing for the client.",
			InArgs: []rpc.ArgDesc{
				{Name: "macaroon", Doc: ``}, // string
			},
			OutArgs: []rpc.ArgDesc{
				{Name: "blessing", Doc: ``}, // security.Blessings
			},
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

	// Register types.
	vdl.Register((*BlessingRootResponse)(nil))

	// Initialize type definitions.
	vdlTypeStruct1 = vdl.TypeOf((*BlessingRootResponse)(nil)).Elem()
	vdlTypeList2 = vdl.TypeOf((*[]string)(nil))

	return struct{}{}
}
