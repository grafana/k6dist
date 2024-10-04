Build k6 distribution with extensions.

k6dist enables building k6 distributions with extensions. The extensions to be embedded can be specified with an extension registry. The location of the registry can be specified with http(s) URL or filesystem path.

The location of the registry can be specified as a command line argument. If omitted, the URL of the registry subset containing the official k6 extensions is the default.

Each flag has a meaningful default value.

k6dist builds k6 for the specified platforms with the extensions listed in the registry. After that, it also creates archive files corresponding to the platform (`.zip` for Windows, otherwise `.tar.gz` format).

In addition to the executable, the readme and the license file can be included in the archive. Their names are detected, but can also be specified using the `--readme` and `--license` flags.

The locations of executables and archives can be specified with the go template expression using the `--executable` and `--archive` flags.

A release notes file is also prepared, in which the extensions are listed with their current and previous versions and their description. The release notes file can be customized, using the `--notes-template` flag you can specify a go template file to generate the release notes.
