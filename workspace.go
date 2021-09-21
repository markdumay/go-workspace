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
	"sort"
	"strings"
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Types
//======================================================================================================================

// AppDirs holds a reference to the initialized directories for the application cache, configuration directory, home
// directory, workspace directory, and the application's temp directory.
type AppDirs struct {
	cache     *Dir
	config    *Dir
	home      *Dir
	temp      *Dir
	workspace *Dir

	keywords        map[string]string //TODO: add make to init?
	keywordsReverse map[string]string
}

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Private Functions
//======================================================================================================================

func init() {
	if runtime.GOOS != "windows" {
		defaultHome = append(defaultHome, "~")
	}
}

func (a *AppDirs) initKeywords() {
	var dirs []*Dir
	a.keywords = make(map[string]string)        // clear the current keywords
	a.keywordsReverse = make(map[string]string) // clear the current reverse keyword map

	if a.cache != nil {
		dirs = append(dirs, a.cache)
	}
	if a.config != nil {
		dirs = append(dirs, a.config)
	}
	if a.home != nil {
		dirs = append(dirs, a.home)
	}
	if a.temp != nil {
		dirs = append(dirs, a.temp)
	}
	if a.workspace != nil {
		dirs = append(dirs, a.workspace)
	}

	for _, d := range dirs {
		for i, alias := range d.Aliases() {
			a.keywords[alias] = d.Path()
			if i == 0 {
				a.keywordsReverse[d.Path()] = alias
			}
		}
	}
}

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Functions
//======================================================================================================================

// NewAppDirs initializes a AppDirs type with default values for the application-specific cache, config, home, temp,
// and workspace directories. Default aliases are added to enable keyword expansion. The keywords follow POSIX string
// expansion rules, using "$" as sigil and optional braces. The following keywords are supported: $HOME, $CACHE, $PWD,
// $TEMP, $TMP, $TMPDIR, $TEMPDIR, and $workspaceRoot. The special character '~' is expanded to the home directory
// (unless the OS is Windows).
func NewAppDirs(appName string) (dirs *AppDirs, err error) {
	var d AppDirs

	cache, e := NewDir(Cache, appName)
	if e != nil {
		return nil, e
	}
	d.cache = cache

	config, e := NewDir(Config, appName)
	if e != nil {
		return nil, e
	}
	d.config = config

	home, e := NewDir(Home, appName)
	if e != nil {
		return nil, e
	}
	d.home = home

	temp, e := NewDir(Temp, appName)
	if e != nil {
		return nil, e
	}
	d.temp = temp

	workspace, e := NewDir(Workspace, appName)
	if e != nil {
		return nil, e
	}
	d.workspace = workspace

	d.initKeywords()

	return &d, nil
}

// Assign initializes a new application-specific directory and updates the internal keyword map to enable
// parameterization of paths. Default aliases are added when no aliases are provided. The full keyword map is updated
// when an existing entry is updated, otherwise the new keywords are appended. Assign does not check for potential
// duplicate keywords.
func (a *AppDirs) Assign(d Dir) {
	var updated bool
	switch d.DirType() {
	case Cache:
		updated = a.cache != nil
		if len(d.Aliases()) == 0 {
			d.AppendAliases(defaultCache...)
		}
		a.cache = &d

	case Config:
		updated = a.config != nil
		if len(d.Aliases()) == 0 {
			d.AppendAliases(defaultConfig...)
		}
		a.config = &d

	case Home:
		updated = a.home != nil
		if len(d.Aliases()) == 0 {
			d.AppendAliases(defaultHome...)
		}
		a.home = &d

	case Temp:
		updated = a.temp != nil
		if len(d.Aliases()) == 0 {
			d.AppendAliases(defaultTemp...)
		}
		a.temp = &d

	case Workspace:
		updated = a.workspace != nil
		if len(d.Aliases()) == 0 {
			d.AppendAliases(defaultWorkspace...)
		}
		a.workspace = &d
	}

	// update the keywords maps
	if updated {
		a.initKeywords()
	} else {
		// initialize keyword maps if needed
		if a.keywords == nil {
			a.keywords = make(map[string]string)
		}
		if a.keywordsReverse == nil {
			a.keywordsReverse = make(map[string]string)
		}

		for i, alias := range d.Aliases() {
			a.keywords[alias] = d.Path()
			if i == 0 {
				a.keywordsReverse[d.Path()] = alias // use the first alias for a reverse substitution
			}
		}
	}
}

// Cache retrieves the current cache directory. It returns an empty string if the directory is not set. Use Assign() to
// initialize a new Cache directory.
func (a *AppDirs) Cache() string {
	if a.cache != nil {
		return a.cache.Path()
	}
	return ""
}

// Config retrieves the current config directory. It returns an empty string if the directory is not set. Use Assign()
// to initialize a new Config directory.
func (a *AppDirs) Config() string {
	if a.cache != nil {
		return a.config.Path()
	}
	return ""
}

// CreateTemp creates the application's temp directory, with mode set to 0755. Nothing happens if the directory
// already exists.
func (a *AppDirs) CreateTemp() (err error) {
	// identify the temp dir path
	path := a.Temp()
	if path == "" {
		// return an error when no temp dir is defined, probably a was not initialized using NewAppDirs
		return fmt.Errorf("cannot create temp directory, invalid state")
	}

	// check if the path already exists, return an error if it's a file or invalid path
	info, e := os.Stat(path)
	if e == nil {
		if info.IsDir() {
			return nil
		}
		return fmt.Errorf("cannot create temp directory: '%s'", path)
	}

	// create the temp directory
	if e := os.Mkdir(path, 0755); e != nil {
		return fmt.Errorf("cannot create temp directory: %s", path)
	}

	return err
}

// Home retrieves the current home directory. It returns an empty string if the directory is not set. Use Assign() to
// initialize a new Home directory.
func (a *AppDirs) Home() string {
	if a.home != nil {
		return a.home.Path()
	}
	return ""
}

// MakeAbsolute returns the absolute path for a given input. It replaces supported keywords with their replacement
// values and converts a relative path to an absolute path. MakeAbsolute calls filepath.Clean on the result.
func (a *AppDirs) MakeAbsolute(basePath string, input string) (path string) {
	segments := strings.Split(input, string(os.PathSeparator))
	var result string

	for _, segment := range segments {
		s := a.keywords[segment]
		if s != "" {
			result = filepath.Join(result, s)
		} else {
			result = filepath.Join(result, segment)
		}
	}

	// prepend the leading `/` if needed
	if filepath.IsAbs(input) && runtime.GOOS != "windows" && !filepath.IsAbs(result) {
		result = string(os.PathSeparator) + result
	}

	return AbsPath(basePath, result)
}

// MakeRelative returns the path for a given input relative to a base path. It replaces supported keywords with their
// replacement values. If input cannot be made relative to the base path, the input itself is returned as result.
// MakeRelative calls filepath.Clean on the result.
func (a *AppDirs) MakeRelative(basePath string, input string) (path string) {
	abs := a.MakeAbsolute(basePath, input)

	rel, e := filepath.Rel(basePath, abs)
	if e == nil {
		return rel
	}
	return filepath.Clean(input)
}

// Parameterize returns the path for a given input relative to the provided base directory, if applicable. Matched path
// segments are replaced with their parameter alias. A non-deterministic match is returned in case of duplicate
// keywords. The first alias is returned when multiple aliases are defined for a directory. Parameterize calls
// filepath.Clean on the result.
func (a *AppDirs) Parameterize(basePath string, input string) (path string) {
	// create an list of all key/value pairs, sorted by key length in descending order
	type item struct {
		key   string
		value string
	}
	ordered := make([]item, len(a.keywordsReverse))
	for k, v := range a.keywordsReverse {
		ordered = append(ordered, item{key: k, value: v})
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		return len(ordered[i].key) > len(ordered[j].key)
	})

	// substitute the paths with their keyword
	for _, o := range ordered {
		input = strings.ReplaceAll(input, o.key, o.value)
	}

	// remove any trailing '/'
	input = strings.TrimSuffix(input, string(os.PathSeparator))

	if !filepath.IsAbs(input) {
		result, err := filepath.Rel(basePath, input)
		if err != nil {
			return filepath.Clean(input)
		}
		return result
	}
	return input
}

// RecreateTemp recreates a subdirectory of the application's temp directory, deleting all existing files. Leave
// subdir empty to recreate the entire application's temp directory. It uses RemoveTempDir to safely remove the
// directory. The mode is set to 0755.
func (a *AppDirs) RecreateTemp(subdir string) (err error) {
	if e := a.RemoveTemp(subdir); e != nil {
		return e
	}

	// create the temp dir
	path := filepath.Join(a.temp.Path(), subdir)
	if e := os.Mkdir(path, 0755); e != nil {
		return fmt.Errorf("cannot create temp directory: %s", path)
	}

	return err
}

// RemoveTemp removes the configured temp dir, deleting all existing files. It uses a failsafe to ensure the
// configured temp dir is valid and within the scope of the system's default temp directory. The expected base paths
// are '$TMPDIR' (on Unix or macOS) or '/tmp' (on Unix, macOS or Plan 9). On Windows, the directories can be either
// '%TMP%' or '%TEMP%'.
func (a *AppDirs) RemoveTemp(subdir string) (err error) {

	// validate the configured temp directory is valid and safe
	if a.temp.Path() == "" {
		return fmt.Errorf("temp directory is not configured correctly")
	}
	tmp := filepath.Clean(os.TempDir())
	current := filepath.Join(a.temp.Path(), subdir)

	if !strings.HasPrefix(current, tmp) {
		return fmt.Errorf("temp directory is considered unsafe")
	}

	if current == tmp {
		return fmt.Errorf("expected a subdirectory within the temp directory")
	}

	// remove the temp dir if it exists
	if e := os.RemoveAll(current); e != nil {
		return e
	}

	return err
}

// Temp retrieves the current temp directory. It returns an empty string if the directory is not set. Use Assign() to
// initialize a new Temp directory.
func (a *AppDirs) Temp() string {
	if a.temp != nil {
		return a.temp.Path()
	}
	return ""
}

// Workspace retrieves the current workspace directory. It returns an empty string if the
// directory is not set. Use Assign() to initialize a new Workspace directory.
func (a *AppDirs) Workspace() string {
	if a.workspace != nil {
		return a.workspace.Path()
	}
	return ""
}

//======================================================================================================================
// endregion
//======================================================================================================================
