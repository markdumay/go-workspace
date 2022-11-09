// Copyright © 2021 Mark Dumay. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be found in the LICENSE file.

// workspace is a Go package to simplify the access to the Cache, Config, Home, Workspace, and Temp folders for an
// application. It uses common settings for Unix, macOS, Plan 9, and Windows. In addition, it supports the substitution
// of configurable keywords, such as $CACHE, $HOME, $workspaceRoot, and $TEMP. Finally, go-workspace sets the workspace
// folder to the correct path when ran from source.
package workspace

//======================================================================================================================
// region Import Statements
//======================================================================================================================

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Constants
//======================================================================================================================

// Defines a pseudo enumeration of possible application directories.
const (
	// Cache is the OS's user-specific cache directory. On Unix, this is either '$XDG_CACHE_HOME' or '$HOME/.cache'. On
	// macOS, this is '$HOME/Library/Caches'. On Plan 9, the cache directory is '$home/lib/cache'. And lastly, on
	// Windows the cache directory is derived from '%LocalAppData%'.
	Cache DirType = iota + 1

	// Config is the directory containing the main application configuration file, if any. It is derived from
	// viper.ConfigFileUsed().
	Config

	// Home is the default, fully expanded user home directory.
	Home

	// Workspace is working directory of the repository or the running command. It is typically initialized by
	// WorkspaceRoot().
	Workspace

	// Temp is the OS-specific temp directory used by ShellDoc. The path is set to either '$TMPDIR' (on Unix or macOS)
	// or '/tmp' (on Unix, macOS or Plan 9). On Windows, the directory can be either '%TMP%', '%TEMP%',
	// '%USERPROFILE%', or the Windows directory.
	//
	// The path is not guaranteed to exist. Use RecreateTempDir() to recreate the directory prior to accessing it, and
	// use RemoveTempDir() once done.
	Temp
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Private Variables
//======================================================================================================================

var (
	defaultCache     = []string{"$CACHE", "${CACHE}"}
	defaultConfig    = []string{}
	defaultHome      = []string{"$HOME", "${HOME}"}
	defaultTemp      = []string{"$TEMP", "${TEMP}", "$TMP", "${TMP}", "$TMPDIR", "${TMPDIR}", "$TEMPDIR", "${TEMPDIR}"}
	defaultWorkspace = []string{"$workspaceRoot", "${workspaceRoot}", "$PWD", "${PWD}"}
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Private Types
//======================================================================================================================

// aliasesOption associates specific aliases for initialization of a new application directory.
type aliasesOption struct {
	Aliases []string
}

// pathOption associates a specific path for initialization of a new application directory.
type pathOption struct {
	Path string
}

// options defines the optional arguments when creating a new application directory.
type options struct {
	path    string
	aliases []string
}

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Types
//======================================================================================================================

// Dir holds a reference to a specific application directory and it's aliases (keywords).
type Dir struct {
	// dirType indicates the type of directory, either Cache, Config, Home, Workspace, or Temp.
	dirType DirType

	// path is the absolute path associated with the directory.
	path string

	// aliases holds a collection of the keywords associated with a directory.
	aliases []string
}

// DirType defines the type of directory to be configured.
type DirType int

// Option defines an optional argument for creating new application directory.
type Option interface {
	apply(*options)
}

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Private Functions
//======================================================================================================================

// apply associates an optional path for initialization of a new application directory.
func (o aliasesOption) apply(opts *options) {
	opts.aliases = o.Aliases
}

// apply associates an optional path for initialization of a new application directory.
func (o pathOption) apply(opts *options) {
	opts.path = o.Path
}

// exists validates if a specific item exists within an array.
func exists(arr []string, item string) bool {
	for _, a := range arr {
		if a == item {
			return true
		}
	}
	return false
}

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Functions
//======================================================================================================================

// NewDir creates a new Dir instance for the provided arguments. NewDir supports two optional parameters, set by
// WithAliases and WithPath respectively. WithAliases associates specific aliases with the application directory.
// WithPath initializes the application directory for a specific path. If omitted, both parameters revert to a default
// value pending the dir type.
func NewDir(dirType DirType, appName string, opts ...Option) (dir *Dir, err error) {
	// init the options
	options := options{}
	options.aliases = make([]string, 0)
	for _, o := range opts {
		o.apply(&options)
	}

	// init the path
	if options.path != "" {
		if !filepath.IsAbs(options.path) {
			return nil, fmt.Errorf("cannot process relative path: %s", options.path)
		}
	} else {
		switch dirType {
		case Cache:
			options.path, err = os.UserCacheDir()
			options.path = filepath.Join(options.path, appName)

		case Config, Workspace:
			options.path, err = Root(appName)

		case Home:
			options.path, err = os.UserHomeDir()

		case Temp:
			options.path = filepath.Join(os.TempDir(), appName)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("cannot initialize directory: %s", dirType.String())
	}

	// init the aliases
	if len(options.aliases) == 0 {
		switch dirType {
		case Cache:
			options.aliases = defaultCache

		case Config:
			options.aliases = defaultConfig

		case Workspace:
			options.aliases = defaultWorkspace

		case Home:
			options.aliases = defaultHome

		case Temp:
			options.aliases = defaultTemp
		}
	}

	// create a new Dir and return the value
	dir = &Dir{
		dirType: dirType,
		path:    filepath.Clean(options.path),
		aliases: options.aliases,
	}

	return
}

// Aliases retrieves a collection of the aliases (keywords) associated with a directory.
func (d *Dir) Aliases() (a []string) {
	a = make([]string, len(d.aliases))
	copy(a, d.aliases)
	return
}

// AppendAliases appends one or more aliases to the collection of aliases (keywords) associated with a directory.
func (d *Dir) AppendAliases(aliases ...string) {
	// append each alias if it does not exist already
	for _, a := range aliases {
		if !exists(d.aliases, a) {
			d.aliases = append(d.aliases, a)
		}
	}

	// sort the aliases alphabetically
	sort.Strings(d.aliases)
}

// DirType retrieves the type of configured directory, either Cache, Config, Home, Workspace, or Temp.
func (d *Dir) DirType() DirType {
	return d.dirType
}

// Path retrieves the absolute path associated with the directory.
func (d *Dir) Path() string {
	return d.path
}

// RemoveAliases removes one or more aliases from the collection of aliases. Unrecognized aliases are ignored.
func (d *Dir) RemoveAliases(aliases ...string) {
	for _, a := range aliases {
		for i, curr := range d.aliases {
			if a == curr {
				d.aliases = append(d.aliases[:i], d.aliases[i+1:]...)
				break
			}
		}
	}
}

// String converts a directory type to it's string representation.
func (d DirType) String() string {
	if d < Cache || d > Temp {
		return ""
	}
	return [...]string{"cache", "config", "home", "workspace", "temp"}[d-1]
}

// AbsPath returns the absolute path for a given base path and path. If path is relative it is joined with the base
// path, otherwise the path itself is returned. AbsPath calls filepath.Clean on the result. The special character "~"
// is expanded to the user's home directory (if set as prefix).
func AbsPath(base string, path string) string {
	if runtime.GOOS != "windows" && strings.HasPrefix(path, "~") {
		dir, e := os.UserHomeDir()
		if e != nil {
			dir = "~"
		}
		path = strings.Replace(path, "~", dir, 1)
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	return filepath.Clean(filepath.Join(base, path))
}

// Root returns the working directory of the repository or the running command. In debugging mode, the current working
// directory may actually be a sub directory, such as 'src' or 'cmd'. In these cases, the workspace root is set to the
// nearest parent directory containing a ".git" repository. When running a compiled binary, the function returns the
// current working directory.
func Root(appName string) (path string, err error) {
	_, cmd := filepath.Split(os.Args[0])
	dir, e := os.Getwd()
	if e != nil {
		return "", e
	}

	// return the current working directory when running a compiled binary
	if cmd == appName {
		return dir, nil
	}

	// traverse the current path for a workspace marker in reverse order
	isRoot := false
	for {
		// return the current path if it contains a ".git" directory
		s, err := os.Stat(filepath.Join(dir, ".git"))
		if err == nil && s.IsDir() {
			return dir, nil
		}

		// stop when at the root of the path
		if isRoot {
			return "", errors.New("cannot identify workspace root (no .git repository found)")
		}

		// TODO: test Windows compatibility
		// traverse one level up the path hierarchy
		dir = filepath.Dir(dir)
		if dir == filepath.VolumeName(dir)+string(os.PathSeparator) || dir == string(os.PathSeparator) {
			isRoot = true
		}
	}
}

// WithAliases associates optional aliases to be used by the application directory. A default value is used if omitted.
func WithAliases(aliases []string) Option {
	return aliasesOption{Aliases: aliases}
}

// WithPath associates an optional path to be used by the application directory. A default value is used if omitted.
func WithPath(path string) Option {
	return pathOption{Path: path}
}

//======================================================================================================================
// endregion
//======================================================================================================================
