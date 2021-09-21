// Copyright © 2021 Mark Dumay. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be found in the LICENSE file.

package workspace

//======================================================================================================================
// region Import Statements
//======================================================================================================================

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Test Functions
//======================================================================================================================

func TestExists(t *testing.T) {
	arr := []string{"a", "b", "c"}

	assert.True(t, exists(arr, "a"))
	assert.False(t, exists(arr, "d"))
}

func TestNewDir(t *testing.T) {
	_, e := NewDir(Cache, "", []string{}, appName)
	require.Nil(t, e, "Unexpected result when initializing app directory")

	_, e = NewDir(Cache, "test", []string{}, appName)
	assert.EqualError(t, e, "cannot process relative path: test")

}

func TestAliases(t *testing.T) {
	d, e := NewDir(Cache, "", []string{}, appName)
	require.Nil(t, e, "Unexpected result when initializing app directory")

	arr := []string{"a", "b", "c"}
	d.AppendAliases(arr...)
	assert.Equal(t, arr, d.Aliases())

	d.RemoveAliases("a", "b", "c", "d")
	assert.Len(t, d.Aliases(), 0)
}

func TestString(t *testing.T) {
	type test struct {
		Type     DirType
		Expected string
	}

	tests := []test{
		{Type: Cache, Expected: "cache"},
		{Type: Config, Expected: "config"},
		{Type: Home, Expected: "home"},
		{Type: Workspace, Expected: "workspace"},
		{Type: Temp, Expected: "temp"},
		{Type: 0, Expected: ""},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, test.Type.String())
	}
}

func TestAbsPath(t *testing.T) {
	type test struct {
		BasePath string
		Path     string
		Expected string
	}

	home, e := os.UserHomeDir()
	require.Nil(t, e)

	tests := []test{
		{BasePath: "", Path: "~", Expected: home},
		{BasePath: home, Path: "test", Expected: filepath.Join(home, "test")},
		{BasePath: home, Path: "/test", Expected: "/test"},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, AbsPath(test.BasePath, test.Path))
	}
}

func TestRoot(t *testing.T) {
	type test struct {
		AppName  string
		Expected string
	}

	_, cmd := filepath.Split(os.Args[0])
	dir, e := os.Getwd()
	require.Nil(t, e)

	tests := []test{
		{AppName: cmd, Expected: dir},
		{AppName: "go-workspace", Expected: dir},
	}

	for _, test := range tests {
		got, e := Root(test.AppName)
		require.Nil(t, e)
		assert.Equal(t, test.Expected, got)
	}
}

//======================================================================================================================
// endregion
//======================================================================================================================
