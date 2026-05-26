# Release Process

PR Maven CLI releases are driven by Git tags.

## Create a Release

1. Ensure `main` is green.
2. Choose a semantic version such as `v0.1.0`.
3. Create and push the tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow will:

- build Linux, macOS, and Windows binaries;
- package archives;
- generate SHA-256 checksums;
- create a GitHub release;
- generate release notes from GitHub metadata.

## Validate a Local Build

```bash
sh scripts/build.sh dist dev
./dist/prmaven version
```

On Windows PowerShell:

```powershell
.\scripts\build.ps1 -Version dev
.\dist\prmaven.exe version
```

## Version Contract

Release builds embed the tag in the CLI:

```bash
prmaven version
```

Development builds report `dev` unless a script or workflow passes a specific version.
