BGPPeeringTool
============

Compiling from source
------------
go build

Usage
------------
-md5
  : Set a session password
-maxprefix
  : Override maximum prefixes obtained from PeeringDB
-template
  : Select a template

Example
        ./BGPPeeringTool -local=61049 -remote=16509
