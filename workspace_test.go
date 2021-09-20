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
	dirs := AppDirs{}
	path, e := Root(appName)
	require.Nil(t, e, "Unexpected result when initializing workspace directory")

	// test a default assign
	d, e := NewDir(Workspace, "", []string{}, appName)
	require.Nil(t, e, "Unexpected result when initializing default workspace directory")
	dirs.Assign(*d)
	for _, keyword := range defaultWorkspace {
		assert.Equal(t, path, dirs.keywords[keyword])
	}

	// test a custom assign
	d, e = NewDir(Workspace, path, []string{"$CUSTOM_DIR"}, appName)
	require.Nil(t, e, "Unexpected result when initializing custom workspace directory")
	dirs.Assign(*d)
	assert.Len(t, dirs.keywords, 1)
	assert.Equal(t, path, dirs.keywords["$CUSTOM_DIR"])
}

func TestCache(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	expectedCache, _ := os.UserCacheDir()
	expectedCache = filepath.Join(expectedCache, appName)
	assert.Equal(t, expectedCache, dirs.Cache())
}

func TestConfig(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	expectedConfig, _ := Root(appName)
	assert.Equal(t, expectedConfig, dirs.Config())
}

func TestHome(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	expectedHome, _ := os.UserHomeDir()
	assert.Equal(t, expectedHome, dirs.Home())
}

func TestTemp(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	expectedTemp := filepath.Join(os.TempDir(), appName)
	assert.Equal(t, expectedTemp, dirs.Temp())
}

func TestWorkspace(t *testing.T) {
	dirs, err := NewAppDirs(appName)
	require.Nil(t, err, "Unexpected result when initializing app directories")

	expectedWorkspace, _ := Root(appName)
	assert.Equal(t, expectedWorkspace, dirs.Workspace())
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

	var tests = [...]test{
		{input: "test", expected: filepath.Join(dirs.Workspace(), "test")},
		{input: "/test", expected: "/test"},
		{input: "./test", expected: filepath.Join(dirs.Workspace(), "test")},
		{input: "~/test", expected: filepath.Join(dirs.Home(), "test")}, //TODO: not on Windows
		{input: "test/../test/", expected: filepath.Join(dirs.Workspace(), "test")},
		{input: "$CACHE/test", expected: filepath.Join(dirs.Cache(), "test")},
		{input: "${CACHE}/test", expected: filepath.Join(dirs.Cache(), "test")},
		{input: "$HOME/test", expected: filepath.Join(dirs.Home(), "test")},
		{input: "${HOME}/test", expected: filepath.Join(dirs.Home(), "test")},
		{input: "$TEMP/test", expected: filepath.Join(dirs.Temp(), "test")},
		{input: "${TEMP}/test", expected: filepath.Join(dirs.Temp(), "test")},
		{input: "$TMP/test", expected: filepath.Join(dirs.Temp(), "test")},
		{input: "${TMP}/test", expected: filepath.Join(dirs.Temp(), "test")},
		{input: "$TMPDIR/test", expected: filepath.Join(dirs.Temp(), "test")},
		{input: "${TMPDIR}/test", expected: filepath.Join(dirs.Temp(), "test")},
		{input: "$TEMPDIR/test", expected: filepath.Join(dirs.Temp(), "test")},
		{input: "${TEMPDIR}/test", expected: filepath.Join(dirs.Temp(), "test")},
		{input: "$workspaceRoot/test", expected: filepath.Join(dirs.Workspace(), "test")},
		{input: "${workspaceRoot}/test", expected: filepath.Join(dirs.Workspace(), "test")},
		{input: "$PWD/test", expected: filepath.Join(dirs.Workspace(), "test")},
		{input: "${PWD}/test", expected: filepath.Join(dirs.Workspace(), "test")},
		{input: "$TEMPtest", expected: filepath.Join(dirs.Workspace(), "$TEMPtest")},
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

	var tests = [...]test{
		{input: filepath.Join(dirs.Workspace(), "test"), expected: "$workspaceRoot/test"},
		{input: "test", expected: "test"},
		{input: "/test", expected: "/test"},
		{input: "/test/", expected: "/test"},
		{input: filepath.Join(dirs.Workspace(), "test/../test/"), expected: "$workspaceRoot/test"},
		{input: "$TEMPtest", expected: "$TEMPtest"},
		{input: filepath.Join(dirs.Cache(), "test"), expected: "$CACHE/test"},
		{input: filepath.Join(dirs.Home(), "test"), expected: "$HOME/test"},
		{input: filepath.Join(dirs.Temp(), "test"), expected: "$TEMP/test"},
		{input: filepath.Join(dirs.Workspace(), "test"), expected: "$workspaceRoot/test"},
	}

	for _, curr := range tests {
		got := dirs.Parameterize(dirs.Workspace(), curr.input)
		assert.Equal(t, curr.expected, got)
	}
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

	var tests = [...]test{
		{input: filepath.Join(os.TempDir(), appName), expected: ""},
		{input: "", expected: ""},
		{input: "/temp", expected: "temp directory is considered unsafe"},
		{input: os.TempDir(), expected: "expected a subdirectory within the temp directory"},
		{input: filepath.Join(os.TempDir(), string(os.PathSeparator)), expected: "expected a subdirectory within the temp directory"},
	}

	for _, curr := range tests {
		d, e := NewDir(Temp, curr.input, []string{}, appName)
		require.Nil(t, e, "Unexpected result when initializing custom temp directory")
		dirs.Assign(*d)

		got := ""
		if e := dirs.RemoveTemp(""); e != nil {
			got = e.Error()
		}
		assert.Equal(t, curr.expected, got)
	}
}

//======================================================================================================================
// endregion
//======================================================================================================================
