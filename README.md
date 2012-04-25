Winhello
========

An example Windows GUI hello world application.

Linker options
--------------

The 8l linker has two options for linking on Windows:
    -Hwindows Writes console binaries (the default)
    -Hwindowsgui Writes GUI binaries

You can link GUI programs with -Hwindows. This will give a console window as well as a GUI window. To suppress the console window use:
go install -ldflags -Hwindowsgui [packages]

Strings
-------

Strings in Go are conventionally encoded in UTF-8 and are not null terminated.
Strings in Windows are encoded in UTF-16 and are null terminated.