module isesol.com/iport

go 1.12

require (
	github.com/cheggaaa/pb/v3 v3.0.4 // indirect
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/json-iterator/go v1.1.9
	github.com/woshihot/go-lib v0.0.0-20200117095909-ec532273b11b
)

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20200115085410-6d4e4cb37c7d
	golang.org/x/mod => github.com/golang/mod v0.2.0
	golang.org/x/net => github.com/golang/net v0.0.0-20200114155413-6afb5195e5aa
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20200116001909-b77594299b42
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20200117012304-6edc0a871e69
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191204190536-9bdfabe68543
)
