// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: server

//nolint:revive
package server

import (
	"v.io/v23/security/access"
	"v.io/v23/services/groups"
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
	vdlTypeMap2    *vdl.Type = nil
	vdlTypeSet3    *vdl.Type = nil
	vdlTypeString4 *vdl.Type = nil
)

// Type definitions
// ================
// groupData represents the persistent state of a group. (The group name is
// persisted as the store entry key.)
type groupData struct {
	Perms   access.Permissions
	Entries map[groups.BlessingPatternChunk]struct{}
}

func (groupData) VDLReflect(struct {
	Name string `vdl:"v.io/x/ref/services/groups/internal/server.groupData"`
}) {
}

func (x groupData) VDLIsZero() bool { //nolint:gocyclo
	if len(x.Perms) != 0 {
		return false
	}
	if len(x.Entries) != 0 {
		return false
	}
	return true
}

func (x groupData) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct1); err != nil {
		return err
	}
	if len(x.Perms) != 0 {
		if err := enc.NextField(0); err != nil {
			return err
		}
		if err := x.Perms.VDLWrite(enc); err != nil {
			return err
		}
	}
	if len(x.Entries) != 0 {
		if err := enc.NextField(1); err != nil {
			return err
		}
		if err := vdlWriteAnonSet1(enc, x.Entries); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func vdlWriteAnonSet1(enc vdl.Encoder, x map[groups.BlessingPatternChunk]struct{}) error {
	if err := enc.StartValue(vdlTypeSet3); err != nil {
		return err
	}
	if err := enc.SetLenHint(len(x)); err != nil {
		return err
	}
	for key := range x {
		if err := enc.NextEntryValueString(vdlTypeString4, string(key)); err != nil {
			return err
		}
	}
	if err := enc.NextEntry(true); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *groupData) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = groupData{}
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
			if err := x.Perms.VDLRead(dec); err != nil {
				return err
			}
		case 1:
			if err := vdlReadAnonSet1(dec, &x.Entries); err != nil {
				return err
			}
		}
	}
}

func vdlReadAnonSet1(dec vdl.Decoder, x *map[groups.BlessingPatternChunk]struct{}) error {
	if err := dec.StartValue(vdlTypeSet3); err != nil {
		return err
	}
	var tmpMap map[groups.BlessingPatternChunk]struct{}
	if len := dec.LenHint(); len > 0 {
		tmpMap = make(map[groups.BlessingPatternChunk]struct{}, len)
	}
	for {
		switch done, key, err := dec.NextEntryValueString(); {
		case err != nil:
			return err
		case done:
			*x = tmpMap
			return dec.FinishValue()
		default:
			if tmpMap == nil {
				tmpMap = make(map[groups.BlessingPatternChunk]struct{})
			}
			tmpMap[groups.BlessingPatternChunk(key)] = struct{}{}
		}
	}
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
	vdl.Register((*groupData)(nil))

	// Initialize type definitions.
	vdlTypeStruct1 = vdl.TypeOf((*groupData)(nil)).Elem()
	vdlTypeMap2 = vdl.TypeOf((*access.Permissions)(nil))
	vdlTypeSet3 = vdl.TypeOf((*map[groups.BlessingPatternChunk]struct{})(nil))
	vdlTypeString4 = vdl.TypeOf((*groups.BlessingPatternChunk)(nil))

	return struct{}{}
}
