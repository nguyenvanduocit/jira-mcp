# Changelog

## [1.1.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.0.3...v1.1.0) (2026-03-13)


### Features

* add security leak scanning with gitleaks including git history ([f999dcc](https://github.com/nguyenvanduocit/jira-mcp/commit/f999dcc81debb6404990778efcf8a0423699b268))

## [1.0.3](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.0.2...v1.0.3) (2026-03-13)


### Bug Fixes

* add explicit archive IDs for homebrew brew formula ([81e502a](https://github.com/nguyenvanduocit/jira-mcp/commit/81e502a6cea0f51626f60da5f9cf4277267dfc37))

## [1.0.2](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.0.1...v1.0.2) (2026-03-13)


### Bug Fixes

* checkout tagged commit in goreleaser to fix tag mismatch ([d08ee11](https://github.com/nguyenvanduocit/jira-mcp/commit/d08ee11c485bfe255a59fd4dd4b8cc8d0a112a70))

## [1.0.1](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.0.0...v1.0.1) (2026-03-13)


### Bug Fixes

* use merge instead of squash to preserve commit SHAs for release-please ([adddad6](https://github.com/nguyenvanduocit/jira-mcp/commit/adddad621feda67eb1a0e8ef0ef7bbe1c6adc55e))

## 1.0.0 (2026-03-13)


### Features

* add ARM64 support for Docker images ([a3cdb5e](https://github.com/nguyenvanduocit/jira-mcp/commit/a3cdb5ef5dda5d519bd88266dce8a251bb6a8380))
* add attachment support with download tool and ADF media rendering ([2086709](https://github.com/nguyenvanduocit/jira-mcp/commit/2086709249b4222282366769ae90ae0b73b8be4b)), closes [#45](https://github.com/nguyenvanduocit/jira-mcp/issues/45)
* add build information support and prompts functionality ([3e624dc](https://github.com/nguyenvanduocit/jira-mcp/commit/3e624dcf485b08fcd230d4f8c31a28e9a2068df2))
* add CLI package and update README with CLI usage ([83420c1](https://github.com/nguyenvanduocit/jira-mcp/commit/83420c193d9b598a02d354346892511eb72006d0))
* add development information tool for Jira issues ([05e9734](https://github.com/nguyenvanduocit/jira-mcp/commit/05e97349065ee0ec6a5302e30a2d23fb565f84bd))
* add get_active_sprint tool ([6b4faba](https://github.com/nguyenvanduocit/jira-mcp/commit/6b4fabaf1ddeac5eaedfc13821975008d5a5208e))
* add gitleaks workflow for secret scanning ([9bb46a0](https://github.com/nguyenvanduocit/jira-mcp/commit/9bb46a0f63793934470a1701a42d2413f29898f8))
* add homebrew installation instructions ([bd44986](https://github.com/nguyenvanduocit/jira-mcp/commit/bd44986528b12cadccc1ca0901210f540077198e))
* add jira_ prefix to all tool names for better LLM discoverability ([0cb6ef7](https://github.com/nguyenvanduocit/jira-mcp/commit/0cb6ef7df7419c7e33532c566dd4c4d61934039a))
* bump version ([4066906](https://github.com/nguyenvanduocit/jira-mcp/commit/40669060ad88876cc2cb2c1af4b0bef1f88c0239))
* **docker:** add Docker support with build and run instructions ([dd3d2bf](https://github.com/nguyenvanduocit/jira-mcp/commit/dd3d2bf055a144ff6c165e9a31acb7e1c71b8b32))
* **docker:** add Docker support with build and run instructions ([a71f5ba](https://github.com/nguyenvanduocit/jira-mcp/commit/a71f5ba059f15e4c53b857e1d2ce1ddbd44a0723))
* **docker:** add GitHub Container Registry support ([1708e4e](https://github.com/nguyenvanduocit/jira-mcp/commit/1708e4eea933567a34de2079a45f97e49817f981))
* **docs:** rewrite README with clear USP and 2‑minute quick start ([4f4f085](https://github.com/nguyenvanduocit/jira-mcp/commit/4f4f085a7ac3e7c258081be631c3ef2c3b8e2dc8))
* init ([6e960c0](https://github.com/nguyenvanduocit/jira-mcp/commit/6e960c048f69fe61baee42c3061aef0a44602be3))
* **jira:** add issue history and relationship tools ([acb7e3a](https://github.com/nguyenvanduocit/jira-mcp/commit/acb7e3a4cf2015a07d5b02286d51437b8bf664f5))
* **jira:** add worklog functionality ([943e76b](https://github.com/nguyenvanduocit/jira-mcp/commit/943e76b204da601eec3d6ab00a4af62a62a7dfee))
* **jira:** enhance issue retrieval with changelog and story point estimate ([62d04fa](https://github.com/nguyenvanduocit/jira-mcp/commit/62d04fab2129076011030df3c2047219bf327044))
* **sprint:** add get_active_sprint tool ([62eb381](https://github.com/nguyenvanduocit/jira-mcp/commit/62eb381676f61951ae9076205259ce81abb931ad))
* **sprint:** add get_sprint tool to retrieve details of a specific sprint ([8f72256](https://github.com/nguyenvanduocit/jira-mcp/commit/8f722564fc3021a59359fda0c310f942f1f7c0a8))
* **sprint:** enhance list_sprints and get_active_sprint tools with project_key support ([5b595e0](https://github.com/nguyenvanduocit/jira-mcp/commit/5b595e06713fbaee2e713b391f1117fca556f061))
* support streamable http ([ee6c103](https://github.com/nguyenvanduocit/jira-mcp/commit/ee6c1032d94e67957dc44b92d746696b37d3c353))
* **tools:** enhance issue tools with flexible field selection ([0964458](https://github.com/nguyenvanduocit/jira-mcp/commit/096445875f177b6b09dd7d98263a377a1e0e6544))
* **tools:** simplify tool naming and add issue relationship features ([a3f4347](https://github.com/nguyenvanduocit/jira-mcp/commit/a3f43478e85d8ccd7dcb9c7c0aa36fdefbde631a))
* update foỏmat of sprint ([d364ded](https://github.com/nguyenvanduocit/jira-mcp/commit/d364ded79e6ec8b63385c61a92e6058bee82f960))
* upgrade to Jira API v3 ([dbd0678](https://github.com/nguyenvanduocit/jira-mcp/commit/dbd067899e4ae6051643ef311e49344e923f3751))


### Bug Fixes

* add --repo flag to gh pr merge to fix auto-merge ([706dce8](https://github.com/nguyenvanduocit/jira-mcp/commit/706dce8eb55a49c2ffe405744df8db390cc3aeab))
* add concurrency control to prevent release race conditions ([960f138](https://github.com/nguyenvanduocit/jira-mcp/commit/960f1383f1b23f02e0235486b9dd95dd421df7a8))
* **api:** update Jira comment API implementation ([e3aa7c1](https://github.com/nguyenvanduocit/jira-mcp/commit/e3aa7c1d295a579b9fcfb56767bd181160b7e035))
* correct ADF structure for issue descriptions and add delete tool ([449a6b2](https://github.com/nguyenvanduocit/jira-mcp/commit/449a6b250bf24bc83deaa3ad1e79380271dd3246))
* disable provenance and SBOM to prevent unknown/unknown platform entry ([618495b](https://github.com/nguyenvanduocit/jira-mcp/commit/618495b7d15db303b18464ac25ab4449ae0a778f))
* migrate JQL search to new API endpoint ([a444486](https://github.com/nguyenvanduocit/jira-mcp/commit/a44448623f4aa0c22ec85b2a998e6d2ec0f44e8e))
* render ADF comment body as readable text in jira_get_comments tool ([8f27f27](https://github.com/nguyenvanduocit/jira-mcp/commit/8f27f27726f7133e30727eda77414b4971a4565b))
* render ADF comment body as readable text in jira_get_comments tool ([ca373c1](https://github.com/nguyenvanduocit/jira-mcp/commit/ca373c100cb09086859f61e9a37ea87b050a3b8f))
* render ADF description as markdown instead of raw structs ([d97d343](https://github.com/nguyenvanduocit/jira-mcp/commit/d97d343a7df8393c7992debf6852657549c48030))
* retrieve all issue comments ([ebc958d](https://github.com/nguyenvanduocit/jira-mcp/commit/ebc958dac00f3924e8f3710edad94d5e80f03ec6))

## [1.17.4](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.17.3...v1.17.4) (2026-02-03)


### Bug Fixes

* render ADF comment body as readable text in jira_get_comments tool ([79fc6cd](https://github.com/nguyenvanduocit/jira-mcp/commit/79fc6cd05801759376f0b14ac99c12ad6bb19e2d))
* render ADF comment body as readable text in jira_get_comments tool ([5125236](https://github.com/nguyenvanduocit/jira-mcp/commit/512523641e8ea55ae9621d2d939fdf0ac9460bff))

## [1.17.3](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.17.2...v1.17.3) (2025-11-10)


### Bug Fixes

* correct ADF structure for issue descriptions and add delete tool ([f432b34](https://github.com/nguyenvanduocit/jira-mcp/commit/f432b346b5ee08c08a5fde3d089ef1df56c8e1f1))

## [1.17.2](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.17.1...v1.17.2) (2025-10-21)


### Bug Fixes

* render ADF description as markdown instead of raw structs ([3f915c6](https://github.com/nguyenvanduocit/jira-mcp/commit/3f915c6e8d5e35099d6d64bd8b3a44b28e4b298f))

## [1.17.1](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.17.0...v1.17.1) (2025-10-09)


### Bug Fixes

* disable provenance and SBOM to prevent unknown/unknown platform entry ([2acf42f](https://github.com/nguyenvanduocit/jira-mcp/commit/2acf42f1c1fbb43f078c2c4828bf0fd23493c026))

## [1.17.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.16.0...v1.17.0) (2025-10-09)


### Features

* add ARM64 support for Docker images ([d535e52](https://github.com/nguyenvanduocit/jira-mcp/commit/d535e52a3c966ca2af3f67d386282ed640e7b387))

## [1.16.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.15.0...v1.16.0) (2025-10-07)


### Features

* add build information support and prompts functionality ([d3df256](https://github.com/nguyenvanduocit/jira-mcp/commit/d3df256f80a181e80eca0510583e9771c47b9a58))

## [1.15.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.14.0...v1.15.0) (2025-10-07)


### Features

* add development information tool for Jira issues ([b0c2ec8](https://github.com/nguyenvanduocit/jira-mcp/commit/b0c2ec84de8d71be7ddb7cdaf976c31ea803ef8b))

## [1.14.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.13.0...v1.14.0) (2025-09-23)


### Features

* add jira_ prefix to all tool names for better LLM discoverability ([f44cd03](https://github.com/nguyenvanduocit/jira-mcp/commit/f44cd03a7c467255192f65404cfc15cc53ae0b76))

## [1.13.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.12.1...v1.13.0) (2025-09-23)


### Features

* **docs:** rewrite README with clear USP and 2‑minute quick start ([69b7b20](https://github.com/nguyenvanduocit/jira-mcp/commit/69b7b20166c424efe485fa2403193105012f2973))

## [1.12.1](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.12.0...v1.12.1) (2025-09-23)


### Bug Fixes

* migrate JQL search to new API endpoint ([bf6b6a2](https://github.com/nguyenvanduocit/jira-mcp/commit/bf6b6a2a105c5e34a0cbc6c3495f711d70cd47aa))

## [1.12.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.11.0...v1.12.0) (2025-09-08)


### Features

* upgrade to Jira API v3 ([1b65607](https://github.com/nguyenvanduocit/jira-mcp/commit/1b6560750ac0f4d37b5cc2fdf004d15932b4346b))

## [1.11.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.10.0...v1.11.0) (2025-06-11)


### Features

* support streamable http ([a67657a](https://github.com/nguyenvanduocit/jira-mcp/commit/a67657a9649b53ec02a5b2a4eb0789810bc7b372))

## [1.10.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.9.0...v1.10.0) (2025-06-11)


### Features

* update foỏmat of sprint ([d378f75](https://github.com/nguyenvanduocit/jira-mcp/commit/d378f7510fb21b8f4fd3ee677d3a1cf54c2ad025))

## [1.9.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.8.0...v1.9.0) (2025-06-06)


### Features

* add get_active_sprint tool ([3095319](https://github.com/nguyenvanduocit/jira-mcp/commit/30953192ce5e88d4940fea09f7a2331ae2516b54))
* **sprint:** add get_active_sprint tool ([89724b2](https://github.com/nguyenvanduocit/jira-mcp/commit/89724b26b7e9c19dfc16d1c8d81455801145fb27))
* **sprint:** enhance list_sprints and get_active_sprint tools with project_key support ([f31cac4](https://github.com/nguyenvanduocit/jira-mcp/commit/f31cac4701ef0e88e2c2b73cc950770c6b5cda02))

## [1.8.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.7.0...v1.8.0) (2025-06-05)


### Features

* **tools:** enhance issue tools with flexible field selection ([40fef1d](https://github.com/nguyenvanduocit/jira-mcp/commit/40fef1dcaa4ef0ed7833cd6db24e3e31c0a35f73))

## [1.7.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.6.0...v1.7.0) (2025-06-05)


### Features

* **jira:** enhance issue retrieval with changelog and story point estimate ([6e9902c](https://github.com/nguyenvanduocit/jira-mcp/commit/6e9902c464430ffc759124792fe0907697d80fab))


### Bug Fixes

* retrieve all issue comments ([a7cfae5](https://github.com/nguyenvanduocit/jira-mcp/commit/a7cfae5e459fd50b3a62f80223bddfc659a5453b))

## [1.6.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.5.0...v1.6.0) (2025-05-06)


### Features

* bump version ([cd2a85c](https://github.com/nguyenvanduocit/jira-mcp/commit/cd2a85c42c8594240e8718524ac0082acf1b7db7))

## [1.5.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.4.0...v1.5.0) (2025-04-18)


### Features

* **docker:** add GitHub Container Registry support ([f939b62](https://github.com/nguyenvanduocit/jira-mcp/commit/f939b629e764d4fe470f6954cc0d281eccde913f))
* **sprint:** add get_sprint tool to retrieve details of a specific sprint ([6756a1c](https://github.com/nguyenvanduocit/jira-mcp/commit/6756a1c79ed0692aeac9d12287fd92ef6bc5f1c2))

## [1.4.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.3.0...v1.4.0) (2025-04-18)


### Features

* **tools:** simplify tool naming and add issue relationship features ([02279ea](https://github.com/nguyenvanduocit/jira-mcp/commit/02279ead729b7a9bd6e78d6ed7903931d39c1580))

## [1.3.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.2.0...v1.3.0) (2025-04-17)


### Features

* **jira:** add issue history and relationship tools ([fc44bd4](https://github.com/nguyenvanduocit/jira-mcp/commit/fc44bd4384775260bf8ea7a0373c89d7053b6450))

## [1.2.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.1.0...v1.2.0) (2025-04-15)


### Features

* **docker:** add Docker support with build and run instructions ([979cd45](https://github.com/nguyenvanduocit/jira-mcp/commit/979cd459663c0004c566cda658efbf9fca50bf52))
* **docker:** add Docker support with build and run instructions ([078062e](https://github.com/nguyenvanduocit/jira-mcp/commit/078062ed7ba2686483a9df4c6000462d5b4fed3a))

## [1.1.0](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.0.1...v1.1.0) (2025-04-04)


### Features

* **jira:** add worklog functionality ([82ea857](https://github.com/nguyenvanduocit/jira-mcp/commit/82ea85767653ea5de4f20beb6585d9f694696a9a))

## [1.0.1](https://github.com/nguyenvanduocit/jira-mcp/compare/v1.0.0...v1.0.1) (2025-04-04)


### Bug Fixes

* **api:** update Jira comment API implementation ([12798ee](https://github.com/nguyenvanduocit/jira-mcp/commit/12798ee285f0b8d5c70689db87fa60b74e72376d))

## 1.0.0 (2025-03-25)


### Features

* add gitleaks workflow for secret scanning ([9bb46a0](https://github.com/nguyenvanduocit/jira-mcp/commit/9bb46a0f63793934470a1701a42d2413f29898f8))
* init ([6e960c0](https://github.com/nguyenvanduocit/jira-mcp/commit/6e960c048f69fe61baee42c3061aef0a44602be3))
