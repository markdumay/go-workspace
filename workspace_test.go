// Copyright © 2021 Mark Dumay. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be found in the LICENSE file.

package workspace

//======================================================================================================================
// region Import Statements
//======================================================================================================================

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Private Constants
//======================================================================================================================

const appName = "Test"

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Test Functions
//======================================================================================================================

func TestAssign(t *testing.T) {
	type test struct {
		DirType  DirType
		Path     string
		Aliases  []string
		AppName  string
		Expected []string
	}

	dirs := AppDirs{}
	path, e := Root(appName)
	require.Nil(t, e)

	tests := []test{
		{
			DirType:  Cache,
			Path:     path,
			Aliases:  []string{},
			AppName:  appName,
			Expected: defaultCache,
		},
		{
			DirType:  Config,
			Path:     path,
			Aliases:  []string{},
			AppName:  appName,
			Expected: defaultConfig,
		},
		{
			DirType:  Home,
			Path:     path,
			Aliases:  []string{},
			AppName:  appName,
			Expected: defaultHome,
		},
		{
			DirType:  Workspace,
			Path:     path,
			Aliases:  []string{},
			AppName:  appName,
			Expected: defaultWorkspace,
		},
		{
			DirType:  Temp,
			Path:     path,
			Aliases:  []string{},
			AppName:  appName,
			Expected: defaultTemp,
		},
		{
			DirType:  Workspace,
			Path:     path,
			Aliases:  []string{"$CUSTOM_DIR"},
			AppName:  appName,
			Expected: []string{"$CUSTOM_DIR"},
		},
	}

	for _, test := range tests {
		d, e := NewDir(test.DirType, test.AppName, WithPath(test.Path), WithAliases(test.Aliases))
		require.Nil(t, e)
		dirs.Assign(*d)

		for _, keyword := range test.Expected {
			assert.Equal(t, path, dirs.keywords[keyword])
		}
	}
}

func TestCache(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	expectedCache, _ := os.UserCacheDir()
	expectedCache = filepath.Join(expectedCache, appName)
	assert.Equal(t, expectedCache, dirs.Cache())

	dirs = &AppDirs{}
	assert.Equal(t, "", dirs.Cache())
}

func TestConfig(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	expectedConfig, _ := Root(appName)
	assert.Equal(t, expectedConfig, dirs.Config())

	dirs = &AppDirs{}
	assert.Equal(t, "", dirs.Config())
}

func TestHome(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	expectedHome, _ := os.UserHomeDir()
	assert.Equal(t, expectedHome, dirs.Home())

	dirs = &AppDirs{}
	assert.Equal(t, "", dirs.Home())
}

func TestTemp(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	expectedTemp := filepath.Join(os.TempDir(), appName)
	assert.Equal(t, expectedTemp, dirs.Temp())

	dirs = &AppDirs{}
	assert.Equal(t, "", dirs.Temp())
}

func TestWorkspace(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	expectedWorkspace, _ := Root(appName)
	assert.Equal(t, expectedWorkspace, dirs.Workspace())

	dirs = &AppDirs{}
	assert.Equal(t, "", dirs.Workspace())
}

func TestNewAppDirs(t *testing.T) {
	_, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")
}

func TestMakeAbsolute(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	type test struct {
		input    string
		expected string
	}

	var tests = []test{
		{input: "test", expected: filepath.Join(dirs.Workspace(), "test")},
		{input: filepath.Join("", "test"), expected: filepath.Join(dirs.Workspace(), "test")},
		{input: filepath.Join("test", "..", "test"), expected: filepath.Join(dirs.Workspace(), "test")},
		{input: filepath.Join("$CACHE", "test"), expected: filepath.Join(dirs.Cache(), "test")},
		{input: filepath.Join("${CACHE}", "test"), expected: filepath.Join(dirs.Cache(), "test")},
		{input: filepath.Join("$HOME", "test"), expected: filepath.Join(dirs.Home(), "test")},
		{input: filepath.Join("${HOME}", "test"), expected: filepath.Join(dirs.Home(), "test")},
		{input: filepath.Join("$TEMP", "test"), expected: filepath.Join(dirs.Temp(), "test")},
		{input: filepath.Join("${TEMP}", "test"), expected: filepath.Join(dirs.Temp(), "test")},
		{input: filepath.Join("$TMP", "test"), expected: filepath.Join(dirs.Temp(), "test")},
		{input: filepath.Join("${TMP}", "test"), expected: filepath.Join(dirs.Temp(), "test")},
		{input: filepath.Join("$TMPDIR", "test"), expected: filepath.Join(dirs.Temp(), "test")},
		{input: filepath.Join("${TMPDIR}", "test"), expected: filepath.Join(dirs.Temp(), "test")},
		{input: filepath.Join("$TEMPDIR", "test"), expected: filepath.Join(dirs.Temp(), "test")},
		{input: filepath.Join("${TEMPDIR}", "test"), expected: filepath.Join(dirs.Temp(), "test")},
		{input: filepath.Join("$workspaceRoot", "test"), expected: filepath.Join(dirs.Workspace(), "test")},
		{input: filepath.Join("${workspaceRoot}", "test"), expected: filepath.Join(dirs.Workspace(), "test")},
		{input: filepath.Join("$PWD", "test"), expected: filepath.Join(dirs.Workspace(), "test")},
		{input: filepath.Join("${PWD}", "test"), expected: filepath.Join(dirs.Workspace(), "test")},
		{input: filepath.Join("$TEMPtest"), expected: filepath.Join(dirs.Workspace(), "$TEMPtest")},
	}

	if runtime.GOOS != "windows" {
		tests = append(tests, test{input: "/test", expected: "/test"})
		tests = append(tests, test{input: "~/test", expected: filepath.Join(dirs.Home(), "test")})
	} else {
		tests = append(tests, test{input: fmt.Sprintf("c:%c%s", filepath.Separator, "test"), expected: fmt.Sprintf("c:%c%s", filepath.Separator, "test")})
	}

	for _, curr := range tests {
		got := dirs.MakeAbsolute(dirs.Workspace(), curr.input)
		assert.Equal(t, curr.expected, got)
	}
}

func TestParameterize(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	type test struct {
		input    string
		expected string
	}

	var tests = []test{
		{input: filepath.Join(dirs.Workspace(), "test"), expected: filepath.Join("$workspaceRoot", "test")},
		{input: "test", expected: "test"},
		{input: filepath.Join(dirs.Workspace(), "test", "..", "test"), expected: filepath.Join("$workspaceRoot", "test")},
		{input: "$TEMPtest", expected: "$TEMPtest"},
		{input: filepath.Join(dirs.Cache(), "test"), expected: filepath.Join("$CACHE", "test")},
		{input: filepath.Join(dirs.Home(), "test"), expected: filepath.Join("$HOME", "test")},
		{input: filepath.Join(dirs.Temp(), "test"), expected: filepath.Join("$TEMP", "test")},
		{input: filepath.Join(dirs.Workspace(), "test"), expected: filepath.Join("$workspaceRoot", "test")},
	}

	if runtime.GOOS != "windows" {
		tests = append(tests, test{input: "/test", expected: "/test"})
		tests = append(tests, test{input: "/test/", expected: "/test"})
	} else {
		tests = append(tests, test{input: fmt.Sprintf("c:%c%s", filepath.Separator, "test"), expected: fmt.Sprintf("c:%c%s", filepath.Separator, "test")})
	}

	for _, curr := range tests {
		got := dirs.Parameterize(dirs.Workspace(), curr.input)
		assert.Equal(t, curr.expected, got)
	}
}

func TestCreateTemp(t *testing.T) {
	dirs := &AppDirs{}

	// test invalid state
	e := dirs.CreateTemp()
	assert.EqualError(t, e, "cannot create temp directory, invalid state")

	// test default behavior
	dirs, e = NewAppDirs(appName)
	require.Nil(t, e)
	e = dirs.CreateTemp()
	require.Nil(t, e)
}

func TestRecreateTemp(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	err = dirs.RecreateTemp("")
	require.Nil(t, err)
}

func TestRemoveTemp(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	type test struct {
		input    string
		expected string
	}

	var tests = []test{
		{input: filepath.Join(os.TempDir(), appName), expected: ""},
		{input: "", expected: ""},
		{input: os.TempDir(), expected: "expected a subdirectory within the temp directory"},
		{input: filepath.Join(os.TempDir(), string(os.PathSeparator)), expected: "expected a subdirectory within the temp directory"},
	}

	if runtime.GOOS != "windows" {
		tests = append(tests, test{input: "/unsupported", expected: "temp directory is considered unsafe"})
	} else {
		tests = append(tests, test{input: fmt.Sprintf("c:%c%s", filepath.Separator, "unsupported"), expected: "temp directory is considered unsafe"})
	}

	for _, curr := range tests {
		d, e := NewDir(Temp, appName, WithPath(curr.input))
		require.Nil(t, e, "Unexpected result when initializing custom temp directory")
		dirs.Assign(*d)

		got := ""
		if e := dirs.RemoveTemp(""); e != nil {
			got = e.Error()
		}
		assert.Equal(t, curr.expected, got)
	}
}

func TestMakeRelative(t *testing.T) {
	dirs, e := NewAppDirs(appName)
	require.Nil(t, e)

	type test struct {
		BasePath string
		Input    string
		Expected string
	}

	var tests []test
	if runtime.GOOS != "windows" {
		tests = []test{
			{
				BasePath: "", Input: "", Expected: ".",
			},
			{
				BasePath: "/", Input: "/test", Expected: "test",
			},
			{
				BasePath: "/x", Input: "/test", Expected: "../test",
			},
		}
	} else {
		tests = []test{
			{
				BasePath: "", Input: "", Expected: ".",
			},
			{
				BasePath: `c:\`, Input: `c:\\test`, Expected: "test",
			},
			{
				BasePath: `c:\\x`, Input: `c:\\test`, Expected: "..\\test",
			},
		}
	}

	for _, test := range tests {
		got := dirs.MakeRelative(test.BasePath, test.Input)
		assert.Equal(t, test.Expected, got)
	}

}

//======================================================================================================================
// endregion
//======================================================================================================================
