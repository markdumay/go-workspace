# go-workspace

<!-- Tagline -->
<p align="center">
    <b>Simplify the platform-aware access to the Cache, Config, Home, Workspace, and Temp folders</b>
    <br />
</p>


<!-- Badges -->
<p align="center">
    <a href="https://pkg.go.dev/go.markdumay.org/workspace" alt="Go Package">
        <img src="https://pkg.go.dev/badge/go.markdumay.org/workspace.svg" alt="Go Reference" />
    </a>
    <a href="https://www.codefactor.io/repository/github/markdumay/go-workspace" alt="CodeFactor">
        <img src="https://img.shields.io/codefactor/grade/github/markdumay/go-workspace" />
    </a>
    <a href="https://github.com/markdumay/go-workspace/commits/main" alt="Last commit">
        <img src="https://img.shields.io/github/last-commit/markdumay/go-workspace.svg" />
    </a>
    <a href="https://github.com/markdumay/go-workspace/issues" alt="Issues">
        <img src="https://img.shields.io/github/issues/markdumay/go-workspace.svg" />
    </a>
    <a href="https://github.com/markdumay/go-workspace/pulls" alt="Pulls">
        <img src="https://img.shields.io/github/issues-pr-raw/markdumay/go-workspace.svg" />
    </a>
    <a href="https://github.com/markdumay/go-workspace/blob/main/LICENSE" alt="License">
        <img src="https://img.shields.io/github/license/markdumay/go-workspace" />
    </a>
</p>


<!-- Table of Contents -->
<p align="center">
  <a href="#about">About</a> •
  <a href="#built-with">Built With</a> •
  <a href="#prerequisites">Prerequisites</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#donate">Donate</a> •
  <a href="#license">License</a>
</p>


## About
go-workspace is a Go package to simplify the access to the Cache, Config, Home, Workspace, and Temp folders for an application. It uses common settings for Unix, macOS, Plan 9, and Windows. In addition, it supports the substitution of configurable keywords, such as `$CACHE`, `$HOME`, `$workspaceRoot`, and `$TEMP`. Finally, go-workspace sets the workspace folder to the correct path when ran from source.


## Built With
The project uses the following core software components:
* [Testify][testify_url] - Go unit-testing toolkit with common assertions and mocks.


## Prerequisites
go-workspace requires Go version 1.16 or later to be installed on your system.


## Installation
```console
go get -u go.markdumay.org/workspace
```


## Usage
Import go-workspace into your application to start using the package. The following code snippet illustrates the basic usage of go-workspace. Please refer to the [package documentation][package] for more details.

<!-- TODO: add example -->
```go
	const appName = "my_app"

	// initialize the application directories
	dirs, e := NewAppDirs(appName)
	if e != nil {
		fmt.Println("ERROR: cannot initialize application directories")
		os.Exit(1)
	}

	// show the cache directory for 'user', expected output (macOS):
	// /Users/user/Library/Caches/my_app
	fmt.Println(dirs.Cache())

	// show the absolute path of 'test' relative to the home directory, expected output (macOS):
	// /Users/user/test
	fmt.Println(dirs.MakeAbsolute(dirs.Workspace(), "$HOME/test"))

	// show the path of a custom keyword, expected output (macOS):
	// $MYDIR
	w, e := NewDir(Workspace, appName, WithPath("/mydir"), WithAliases([]string{"$MYDIR"}))
	if e != nil {
		fmt.Println("ERROR: cannot initialize workspace directory")
		os.Exit(1)
	}
	dirs.Assign(*w)
	fmt.Println(dirs.Parameterize(dirs.Workspace(), "/mydir"))
```

### Supported Folders
`go-workspace` supports the following five types of folders.

| Type      | Description |
|-----------|-------------|
| Cache     | User-specific cache directory |
| Config    | Current directory (when running from console) or project root (when running from source) |
| Home      | User home directory |
| Workspace | Current directory (when running from console) or project root (when running from source) |
| Temp      | Temp directory |

Unfold one of the below operating systems to see the mapping of the folders to their physical location. The locations are prioritized from left to right in case multiple locations are specified.

<details>
<summary>Unix</summary>

| Type      | Default location                                        |
|-----------|---------------------------------------------------------|
| Cache     | `$XDG_CACHE_HOME/$APP_NAME` or `$HOME/.cache/$APP_NAME` |
| Config    | `$PWD`                                                  |
| Home      | `$HOME/.$APP_NAME`                                      |
| Workspace | `$PWD`                                                  |
| Temp      | `$TMPDIR` or `/tmp`                                     |
</details>

<details>
<summary>macOS</summary>

| Type      | Default location                                        |
|-----------|---------------------------------------------------------|
| Cache     | `$HOME/Library/Caches/$APP_NAME` |
| Config    | `$PWD`                                                  |
| Home      | `$HOME/.$APP_NAME`                                      |
| Workspace | `$PWD`                                                  |
| Temp      | `$TMPDIR` or `/tmp`                                     |
</details>

<details>
<summary>Plan 9</summary>

| Type      | Default location                                        |
|-----------|---------------------------------------------------------|
| Cache     | `$home/lib/cache/$APP_NAME`                             |
| Config    | `$pwd`                                                  |
| Home      | `$home/.$APP_NAME`                                      |
| Workspace | `$pwd`                                                  |
| Temp      | `/tmp`                                                  |
</details>

<details>
<summary>Windows</summary>

| Type      | Default location                                                                                  |
|-----------|---------------------------------------------------------------------------------------------------|
| Cache     | `%LocalAppData%\$APP_NAME`                                                                        |
| Config    | `%cd%`                                                                                            |
| Home      | `%HOME%\$APP_NAME`, `%HOMEDRIVE%\$APP_NAME`, `%HOMEPATH%\$APP_NAME`, or `%USERPROFILE%\$APP_NAME` |
| Workspace | `%cd%`                                                                                            |
| Temp      | `%TMP%`, `%TEMP%`, `%USERPROFILE%`, or the Windows directory                                      |
</details>


## Contributing
go-workspace welcomes contributions of any kind. It is recommended to create an issue to discuss your intended contribution before submitting a larger pull request though. Please consider the following guidelines when contributing:
- Address all linting recommendations from `golangci-lint run` (using `.golangci.yml` from the repository).
- Ensure the code is covered by one or more unit tests (using Testify when applicable).
- Follow the recommendations from [Effective Go][effective_go] and the [Uber Go Style Guide][uber_go_guide].

The following steps decribe how to submit a Pull Request:
1. Clone the repository and create a new branch 
    ```console
    $ git checkout https://github.com/markdumay/go-workspace.git -b name_for_new_branch
    ```
2. Make and test the changes
3. Submit a Pull Request with a comprehensive description of the changes


## Donate
<a href="https://www.buymeacoffee.com/markdumay" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/lato-orange.png" alt="Buy Me A Coffee" style="height: 51px !important;width: 217px !important;"></a>


## License
The go-workspace codebase is released under the [MIT license][license]. The documentation (including the "README") is licensed under the Creative Commons ([CC BY-NC 4.0)][cc-by-nc-4.0] license.

<!-- MARKDOWN PUBLIC LINKS -->
[cc-by-nc-4.0]: https://creativecommons.org/licenses/by-nc/4.0/
[effective_go]: https://golang.org/doc/effective_go
[testify_url]: https://github.com/stretchr/testify
[uber_go_guide]: https://github.com/uber-go/guide/

<!-- MARKDOWN MAINTAINED LINKS -->
<!-- TODO: add blog link
[blog]: https://markdumay.com
-->
[blog]: https://github.com/markdumay
[license]: https://github.com/markdumay/go-workspace/blob/main/LICENSE
[package]: https://pkg.go.dev/go.markdumay.org/workspace
[repository]: https://github.com/markdumay/go-workspace.git