# GitVersion
[![Build Status][wercker-image]][wercker-url] [![Open Issues][issues-image]][issues-url]

> A helper for bumping versions via git tags.

## Usage

```bash
gitversion --prefix v bump patch
```

Only [semver](http://semver.org/)-style versions with optional prefix are
supported at this time (major.minor.patch).

`gitversion` will filter all tags of the format
`<prefix><major>.<minor>.<patch>`, sort them, and increment the requested
field (patch in this example) on the largest version. It then tags the
current revision with the result.

### e.g.

```bash
> git tag
v1.2.3
v1.2.4

> gitversion --prefix v bump patch
Bumping patch for version 1.2.4
1.2.5

> git tag
v1.2.3
v1.2.4
v1.2.5

> gitversion --prefix v show
v1.2.5
```

### Auto

Auto is a special field that will determine the proper field to bump
based on the contents of the last commit message.  It looks for anything matching:

> [major] or [major bump]
>
> [minor] or [minor bump]

Example:
```bash
> git tag
1.2.3
1.2.4

> git log -1
[minor] Added show, major, minor, and auto features

> gitversion bump auto
1.3.0
```

And will default to patch if none found or if the commit is already tagged.

## Testing

```bash
go test ./...
```

## License

Code licensed under the BSD 3-Clause license. See LICENSE file for terms.

[issues-image]: https://img.shields.io/github/issues/screwdriver-cd/gitversion.svg
[issues-url]: https://github.com/screwdriver-cd/gitversion/issues
[wercker-image]: https://app.wercker.com/status/28e7d21d5c6bfe687a26689ea48e53a7
[wercker-url]: https://app.wercker.com/project/bykey/28e7d21d5c6bfe687a26689ea48e53a7
