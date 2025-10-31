# Changelog

## [0.2.0](https://github.com/repsejnworb/keysej/compare/keysej-v0.1.0...keysej-v0.2.0) (2025-10-31)


### Features

* **cli:** add global pretty error handling with colorized usage hints and silence cobra spam ([f27e96c](https://github.com/repsejnworb/keysej/commit/f27e96c83b6ed90d90e4c0e482106b3861f55c01))
* **cli:** add initial Cobra + Bubble Tea scaffold ([7ce1cb2](https://github.com/repsejnworb/keysej/commit/7ce1cb23d7f5459d4ce46fa58f4ca6581c2d31da))
* **sshconf:** add tidy to normalize keysej files (sort blocks, trim, final newline) ([aa12eed](https://github.com/repsejnworb/keysej/commit/aa12eed2aab159d0bcf9b200a9dcbbb74a3d3eda))
* **sshconf:** add validate to lint keysej config files (CIDR/Host, IdentityFile, tag) ([a2866b6](https://github.com/repsejnworb/keysej/commit/a2866b6c4c62599142e5770163598c37469f4682))
* **sshconf:** guard new with private key existence and keysej tag; add --force to bypass ([867c27a](https://github.com/repsejnworb/keysej/commit/867c27a0394e8b7d23f7e70df8860aa3ceb090bf))
* **sshconf:** manage ~/.ssh/config.d/keysej.&lt;key&gt;.conf with list/new/delete (Host/CIDR) ([ad331c3](https://github.com/repsejnworb/keysej/commit/ad331c331b749b40eb9f1d0912e3ff5a15deef5a))


### Fixes

* **cli:** improve global Cobra arg error handling and flag parse UX ([56570f7](https://github.com/repsejnworb/keysej/commit/56570f7873746f61420d5ead98cd59c34b3f142f))
* sorted sshDir defaults ([feeec3e](https://github.com/repsejnworb/keysej/commit/feeec3ef21e6ac33a2796f6cb3555af3bd973c69))
* **sshconf:** friendly error messages + colors for list when files/hosts are missing ([cffbde8](https://github.com/repsejnworb/keysej/commit/cffbde803618fe9b1268292e25182d8628d32153))


### CI

* automerge release please prs ([4661bee](https://github.com/repsejnworb/keysej/commit/4661beedc44190f80cdb7ee026b6172e259d6c17))
* build and release ([56d8faa](https://github.com/repsejnworb/keysej/commit/56d8faacd671f1a9282d545fb05e77a3a12a4a52))
* fixed incorrect branch matching ([e9d9e45](https://github.com/repsejnworb/keysej/commit/e9d9e45d47fae755d2da28cf1ba9b228843b25aa))
* only automerge on main, also fixed check-reference ([ab5d21d](https://github.com/repsejnworb/keysej/commit/ab5d21d962e8ead063def2b01b39c53ffd1000e0))
* trigger automerge as before but wait for ci to finish ([28571ec](https://github.com/repsejnworb/keysej/commit/28571ece77d934477f39fb3f469c834ebce3a928))
* trigger automerge-workflow after ci finishes ([127153e](https://github.com/repsejnworb/keysej/commit/127153ecaede8b610aef5b68e54baffba90781a2))
* use the correct botname in matching for automerge ([61075cc](https://github.com/repsejnworb/keysej/commit/61075cce86101cc682ab5aef7fe5c452ab246445))


### Chores

* initial fixes and code ([c53f4f4](https://github.com/repsejnworb/keysej/commit/c53f4f44600bcc7d9731ca6272a63caae0b1ca35))
* made listing a little nicer ([4f140f6](https://github.com/repsejnworb/keysej/commit/4f140f658c21e0ea2aa8103eb1a5ee9d8fa3708c))
