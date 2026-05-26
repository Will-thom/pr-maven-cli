# Installation

PR Maven CLI can be installed from source today and from release artifacts after tagged releases are published.

The analyzer is local-first. It does not need GitHub tokens, CI API access, network services, or an AI provider. It reads Maven report artifacts that already exist in a workspace.

## Requirements

- Go 1.22 or newer to build from source.
- A Maven project with `pom.xml` files and test report artifacts.
- Maven 3.9.x as the production baseline for generated reports.

The project is documented against Maven 3.9.16. Apache Maven lists Maven 3.9.16 as the latest recommended release for all users, while Maven 4.x is still a preview line that is not safe for production use: <https://maven.apache.org/download.cgi>.

Maven itself is not required to run PR Maven CLI if the report files already exist. Maven is only required when you need to generate or regenerate reports.

## Install From Source

```bash
git clone https://github.com/Will-thom/pr-maven-cli.git
cd pr-maven-cli
go test ./...
go install ./cmd/prmaven
```

Confirm the binary is available:

```bash
prmaven version
```

If `prmaven` is not found, ensure the Go binary directory is on `PATH`.

On Unix-like systems:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

On Windows PowerShell:

```powershell
$env:Path = "$(go env GOPATH)\bin;$env:Path"
```

## Run Without Installing

From the repository root:

```bash
go run ./cmd/prmaven fails -project demo/multi-module-failure
go run ./cmd/prmaven why -project demo/multi-module-failure -format json
```

This is the fastest path for contributors and for trying the demo fixture.

## Build Local Binaries

Use the project scripts when you want local binaries under `dist/`.

Unix-like systems:

```bash
sh scripts/build.sh dist dev
./dist/prmaven version
```

Windows PowerShell:

```powershell
.\scripts\build.ps1 -Version dev
.\dist\prmaven.exe version
```

## Install From Releases

After a tagged release exists, download the archive for your operating system from:

```text
https://github.com/Will-thom/pr-maven-cli/releases
```

Then:

1. Extract the archive.
2. Move `prmaven` or `prmaven.exe` into a directory on `PATH`.
3. Run `prmaven version`.

Release artifacts are documented in [release.md](release.md).

