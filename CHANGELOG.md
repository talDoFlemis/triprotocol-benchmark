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
