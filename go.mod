module swan-provider

go 1.16

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/Khan/genqlient v0.5.0
	github.com/docker/docker v20.10.21+incompatible
	github.com/fatih/color v1.13.0
	github.com/filswan/go-swan-lib v0.2.141
	github.com/gin-gonic/gin v1.7.7
	github.com/google/uuid v1.3.0
	github.com/ipfs/go-log/v2 v2.5.1
	github.com/itsjamie/gin-cors v0.0.0-20160420130702-97b4a9da7933
	github.com/joho/godotenv v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/tendermint/tendermint v0.35.9
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi

replace github.com/filecoin-project/boostd-data => github.com/FogMeta/boostd-data v1.6.3

replace github.com/filswan/go-swan-lib => ./extern/go-swan-lib
