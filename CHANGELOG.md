## [1.13.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.12.0...v1.13.0) (2025-11-05)

### Features

* **domain:** add non iso parse fn ([4eadbe9](https://github.com/talDoFlemis/triprotocol-benchmark/commit/4eadbe9dd269400e22bd0f5f5f2a616e75a2a755))
* **proto_requests:** parse header ([936c1b5](https://github.com/talDoFlemis/triprotocol-benchmark/commit/936c1b54f9200ef8b1f293a7eb305b4c0ee71192))

### Bug Fixes

* **protobuf_serde:** enter body interface value ([dff4907](https://github.com/talDoFlemis/triprotocol-benchmark/commit/dff4907a0b490e81cd57ab6ce6df1582441c8e68))
* **string_serde:** only add [ if string dont starts with it ([f0e2928](https://github.com/talDoFlemis/triprotocol-benchmark/commit/f0e2928d24ff8cba1c981ef374a7d1808033cbc1))

## [1.12.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.11.1...v1.12.0) (2025-11-04)

### Features

* add second ticker for showing logs with updated time ([05500ff](https://github.com/talDoFlemis/triprotocol-benchmark/commit/05500ff5942c3bb9f38950155679cba2e323f483))

## [1.11.1](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.11.0...v1.11.1) (2025-11-03)

### Bug Fixes

* **string_serde:** oob error on len data ([63aa3a9](https://github.com/talDoFlemis/triprotocol-benchmark/commit/63aa3a937d11b31031d36eabcebb933594b5b794))

## [1.11.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.10.3...v1.11.0) (2025-11-03)

### Features

* **string_serde:** add handling of missing end of token stream from shitty server ([5f63f8c](https://github.com/talDoFlemis/triprotocol-benchmark/commit/5f63f8c1a9e33c17ef132ab89058b8bd45321570))

## [1.10.3](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.10.2...v1.10.3) (2025-11-02)

### Code Refactoring

* remove old http error ([252aa0c](https://github.com/talDoFlemis/triprotocol-benchmark/commit/252aa0c60790bd6bacf7394df41d69a5f3d7bafb))

## [1.10.2](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.10.1...v1.10.2) (2025-11-02)

### Documentation

* add README ([93c00de](https://github.com/talDoFlemis/triprotocol-benchmark/commit/93c00de5264801905a8e1ac2900ce5201958b7ca))

### Code Refactoring

* remove handler ([596f26c](https://github.com/talDoFlemis/triprotocol-benchmark/commit/596f26c2742478af437f6aea31bced2c59f0311c))

## [1.10.1](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.10.0...v1.10.1) (2025-11-02)

### Bug Fixes

* **release:** missing service name ([7c33a3b](https://github.com/talDoFlemis/triprotocol-benchmark/commit/7c33a3b118024b3dbb952b330c68d7f2e968f72a))

## [1.10.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.9.1...v1.10.0) (2025-11-02)

### Features

* add Dockerfile ([c8d18e4](https://github.com/talDoFlemis/triprotocol-benchmark/commit/c8d18e4282db6c4a284c9454096aa6c75b61f018))
* add duration and size logs to app_layer ([9ceee64](https://github.com/talDoFlemis/triprotocol-benchmark/commit/9ceee642f4387d51eef72e383de5c6b2990397bf))
* add release ([cb46138](https://github.com/talDoFlemis/triprotocol-benchmark/commit/cb461388eec6d4f607eb9216909a1960b39979e0))
* **tui:** add slog debug logger ([9cf6ec2](https://github.com/talDoFlemis/triprotocol-benchmark/commit/9cf6ec229c3ce7f7abadf250224fb7db8844fe35))

### Bug Fixes

* **tui:** missing operation, params and protocol when failing in auth ([ec43c2c](https://github.com/talDoFlemis/triprotocol-benchmark/commit/ec43c2c33dce632758813dd97edc932c27153a34))

### Code Refactoring

* **string_serde:** make get str field representation a fn ([b29e9fc](https://github.com/talDoFlemis/triprotocol-benchmark/commit/b29e9fc84444ff0c7fc7b9c9ff67a241317e44f3))

## [1.9.1](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.9.0...v1.9.1) (2025-11-02)

### Code Refactoring

* **tui:** run in reverse order ([5030e61](https://github.com/talDoFlemis/triprotocol-benchmark/commit/5030e61ea526ceb3e708b58ec0d0448407f69b62))

## [1.9.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.8.0...v1.9.0) (2025-11-02)

### Features

* add last operation logs ([2e67d4f](https://github.com/talDoFlemis/triprotocol-benchmark/commit/2e67d4fa15b34473c2cd83fa37b1a9dbdb210a5c))
* add progress bar and glamoour render ([51b8588](https://github.com/talDoFlemis/triprotocol-benchmark/commit/51b85882ea7bb90d3ed712a1267282f269ae9293))
* add scrolling to viewport ([83b3c4a](https://github.com/talDoFlemis/triprotocol-benchmark/commit/83b3c4abdd15919ba0ac961f27482de5f44ab61f))
* add viewport rendering ([9b4d988](https://github.com/talDoFlemis/triprotocol-benchmark/commit/9b4d98824b0fa51bfa197b1302b534a6d329e9c8))

### Bug Fixes

* **auth:** dont use pointer ([e8babf4](https://github.com/talDoFlemis/triprotocol-benchmark/commit/e8babf46a27ca734990efaa3136488bacc3c57c9))
* **round_triper:** remove null terminator string ([e3f1d55](https://github.com/talDoFlemis/triprotocol-benchmark/commit/e3f1d55c5ff3d23aba9fb333c68a90c43990a831))
* **string_serde:** follow pointers ([5e25e07](https://github.com/talDoFlemis/triprotocol-benchmark/commit/5e25e072c77966f4e02f4a4ae7002d1271342b47))

## [1.8.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.7.0...v1.8.0) (2025-11-01)

### Features

* add json requests script ([76137be](https://github.com/talDoFlemis/triprotocol-benchmark/commit/76137be32c4849e6a4a935b1474e8752332bc4e5))
* add json serde with string serialization ([4ed9456](https://github.com/talDoFlemis/triprotocol-benchmark/commit/4ed94563bf6167d0a14c884edf7a6002ef6387a0))
* add unix timestamp handling ([6ad848b](https://github.com/talDoFlemis/triprotocol-benchmark/commit/6ad848bd8b251c1e149d1c127a47eef205270d14))
* add unmarshall to json serde ([5d7539b](https://github.com/talDoFlemis/triprotocol-benchmark/commit/5d7539b03470bfca990da00499e5aba3516e2f42))
* **domain:** add student data to auth response for json ([3369cc5](https://github.com/talDoFlemis/triprotocol-benchmark/commit/3369cc545f8a5ec2990db40ac8660f8d24730551))

### Bug Fixes

* **string_serde:** use getFieldTagValue on marshalling ([ef90038](https://github.com/talDoFlemis/triprotocol-benchmark/commit/ef900381b8a0049416552b021a44ed3bff792a6b))

### Code Refactoring

* **domain:** change types and tags ([b2cfe54](https://github.com/talDoFlemis/triprotocol-benchmark/commit/b2cfe54c67093a1e1236e03341441c54ee911ee5))
* **domain:** make number sum list different captalization ([4062792](https://github.com/talDoFlemis/triprotocol-benchmark/commit/4062792a68f9895727408ed374b88ee69b86f8eb))
* **string_request:** add slice operations ([c8c81b4](https://github.com/talDoFlemis/triprotocol-benchmark/commit/c8c81b44617d3db920c1555c5bd45934fcdd8a50))
* **sum_response:** use float for values ([755bce0](https://github.com/talDoFlemis/triprotocol-benchmark/commit/755bce0bd92499f8bc792c930d563ff3eee49bf0))
* use non std time.Time ([65a7ea7](https://github.com/talDoFlemis/triprotocol-benchmark/commit/65a7ea7682833f0679139c68a0476a2092b614cf))

## [1.7.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.6.0...v1.7.0) (2025-10-31)

### Features

* add better nil pointer handling ([fac7c81](https://github.com/talDoFlemis/triprotocol-benchmark/commit/fac7c81c8d7979e2ea0fa53f9c3f099373631f7d))
* add map parsing, optional fields and struct parsing ([2f68b2e](https://github.com/talDoFlemis/triprotocol-benchmark/commit/2f68b2eb326bc6cb539d1e55bbaebdd7d446bebf))
* add string serde testes ([d2db7ce](https://github.com/talDoFlemis/triprotocol-benchmark/commit/d2db7ce8dab5e68eb90ed844ae79dd1591d9faa0))
* add strings request ([f24b483](https://github.com/talDoFlemis/triprotocol-benchmark/commit/f24b483bb4611c1e9ba3ff65a20bf83d87825d6b))

### Bug Fixes

* bug on parsing long numbers as float ([bdefa3b](https://github.com/talDoFlemis/triprotocol-benchmark/commit/bdefa3ba6f905998dbed1819703af5bc2dffb15b))
* **string_serde:** skip omitempty ([5db490d](https://github.com/talDoFlemis/triprotocol-benchmark/commit/5db490d25367da21d464594f962f5980bd7ef395))

## [1.6.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.5.0...v1.6.0) (2025-10-30)

### Features

* add server per protocol on settings ([a5f4706](https://github.com/talDoFlemis/triprotocol-benchmark/commit/a5f47062045ad97590ed1612e006685d81b7a051))
* add server to main.go ([a64367d](https://github.com/talDoFlemis/triprotocol-benchmark/commit/a64367d9eeb9e2d0b786cfd24307ba61cb1ef233))
* **handler:** group api to /api/v1 ([fcbc953](https://github.com/talDoFlemis/triprotocol-benchmark/commit/fcbc9533999fb5a4e04a2caa1c58541b0df6d0e4))

### Code Refactoring

* **domain:** update definitions fomr actual api ([6108634](https://github.com/talDoFlemis/triprotocol-benchmark/commit/6108634f255dc2439be8e8e70ae58fab461b4d7b))
* **domain:** update timestamp response to match actual api ([603e0ab](https://github.com/talDoFlemis/triprotocol-benchmark/commit/603e0abb5d76af88e44efcac09336e82453cac4a))
* **error:** handle validation errors first ([28b5b0d](https://github.com/talDoFlemis/triprotocol-benchmark/commit/28b5b0daa09ff48038336c5c22c2f18720fb4125))
* **handler:** add getProtocolAddress ([ed189c8](https://github.com/talDoFlemis/triprotocol-benchmark/commit/ed189c8e76f0cfe69f3f2c7f1faf81d8c6894a3d))
* **string_serde:** add bool and int serializing handling ([d888f36](https://github.com/talDoFlemis/triprotocol-benchmark/commit/d888f3626ef02a0dec7e59aca5e0726922bb8597))
* **string_serde:** use new standard to unmarshall ([247ef2d](https://github.com/talDoFlemis/triprotocol-benchmark/commit/247ef2d8ce6e5af1c150beae71dcee251cd85c49))

## [1.5.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.4.0...v1.5.0) (2025-10-29)

### Features

* add roundtrip logs ([ceb6a23](https://github.com/talDoFlemis/triprotocol-benchmark/commit/ceb6a23c8c3616428bb0600e4af29048791e66d6))

## [1.4.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.3.0...v1.4.0) (2025-10-29)

### Features

* add app layer ([887449e](https://github.com/talDoFlemis/triprotocol-benchmark/commit/887449ee538a124469e79982a0313b75510a00ba))
* add custom validator ([5c92bb3](https://github.com/talDoFlemis/triprotocol-benchmark/commit/5c92bb3a7c08bd308ae4bb316662fb2c00f8dc59))
* add git lfs to pdf ([e0d0947](https://github.com/talDoFlemis/triprotocol-benchmark/commit/e0d0947108fbdd64afe22a0aca39db76ef0c40e1))
* add handler auth request ([f829711](https://github.com/talDoFlemis/triprotocol-benchmark/commit/f8297110b2c61e5bc087fc0fb842cb2546a54f0a))
* add history and logout domain entities ([e8ae711](https://github.com/talDoFlemis/triprotocol-benchmark/commit/e8ae711fc5e82d6601c7eb0ae2c22280de6c7f19))
* add protobuf definition ([18d1a38](https://github.com/talDoFlemis/triprotocol-benchmark/commit/18d1a38125baf112360324e518fadb71d1787989))
* add protogenerated files ([8760669](https://github.com/talDoFlemis/triprotocol-benchmark/commit/8760669d5dff5723a64805b45a9e0268b7631cc7))
* add roundtripper ([16f0476](https://github.com/talDoFlemis/triprotocol-benchmark/commit/16f047652fcd7494eeb99872d1dd7cbfd9acba4e))
* add string serde tests ([f7a434f](https://github.com/talDoFlemis/triprotocol-benchmark/commit/f7a434fb04ba74cbad062f86fa67471fc5ad5797))
* register custom validator and use app layer ([b972aba](https://github.com/talDoFlemis/triprotocol-benchmark/commit/b972aba80a1f2508c23ed2b7380017324b0a767e))
* spec of project ([853684f](https://github.com/talDoFlemis/triprotocol-benchmark/commit/853684fcca7a99d951d8e812a960a2be6cab5666))

## [1.3.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.2.1...v1.3.0) (2025-10-28)

### Features

* add domain ([b4678fe](https://github.com/talDoFlemis/triprotocol-benchmark/commit/b4678fecf4f72282326cd8c0ac7d6a93947e2390))

## [1.2.1](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.2.0...v1.2.1) (2025-10-23)

### Bug Fixes

* **handler:** wront tcp timeout ([2a89db3](https://github.com/talDoFlemis/triprotocol-benchmark/commit/2a89db397f091d5a9a45232dec3b353522e5adfc))

## [1.2.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.1.0...v1.2.0) (2025-10-23)

### Features

* add base serde ([1a5e63e](https://github.com/talDoFlemis/triprotocol-benchmark/commit/1a5e63eba41e1005bfe04dd9b3d34c86e8ee0df2))
* add new tcp and protocol settings ([d391251](https://github.com/talDoFlemis/triprotocol-benchmark/commit/d3912518fc1c4d47820011ff1cfae04f73d89c51))

## [1.1.0](https://github.com/talDoFlemis/triprotocol-benchmark/compare/v1.0.0...v1.1.0) (2025-10-23)

### Features

* add air config ([4d72101](https://github.com/talDoFlemis/triprotocol-benchmark/commit/4d72101ecb0a8501ac23e04b20dba89ec3ce8f60))
* add main handler ([7d97837](https://github.com/talDoFlemis/triprotocol-benchmark/commit/7d97837061e13a0d930439822a749509fc0bfac5))
* add settings ([da99a18](https://github.com/talDoFlemis/triprotocol-benchmark/commit/da99a18be181f5e2e746f6e44d9d8d037460973f))

## 1.0.0 (2025-10-23)

### Features

* add releaser ([239b82d](https://github.com/talDoFlemis/triprotocol-benchmark/commit/239b82d47c9d2cee847e3b087d60fa0e87e01ceb))
