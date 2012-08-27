Kelvin
------

Kelvin is a [Amazon Glacier](http://aws.amazon.com/glacier/) client implemented in Go. I chose Go 
for 3 reasons:

1. I want to learn Go better.
2. I have really enjoyed programming in Go so far.
3. Go compiles to a static binary for Mac, Windows, and Linux

Glacier in Brief
================
Glacier is a service offering from Amazon designed for reliably storing large quantities of 
data in a long time frame.  Glacier has two basic abstractions, a "vault" and an "archive."
An archive is a dicrete storage item - an arbitrary blob of data. A vault collects archives
in a namespace.  Amazon automatically encrypts data at rest using AES256 and it manages the 
keys for you. (You may optionally tell kelvin to encrypt your files before transmission using
your own keys that you control.)

Setup
=====
The first time you run `kelvin` from the command line, it will build a configuration file for you,
prompting you to enter your Amazon AWS credentials.

It will also use a set of sensible default values.

If you want to reconfigure Kelvin, you can run `kelvin --configure` explicitly.

Basic Use
=========
Kelvin accepts a list of pathnames or files much like GNU tar. You may exclude files or paths
using the `-x` (or `--exclude`) flag.  By default, kelvin will treat different path specifications
as vault definitions.

    /home/foo would become vault "hostname-home-foo"
    /home/bar would become vault "hostname-home-bar"

    and so on...

The default operation is to recurse through a path specification and upload files found therein as
discrete archives.  You may optionally tell kelvin to aggregate files into a ZIP style 
archive before transmission.

Vault Operations
================

    list-vaults
    create-vault <name>
    delete-vault <name>

Retriving an archive
====================
XXX TODO XXX
