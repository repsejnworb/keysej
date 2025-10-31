# 1.0.0 (2025-10-31)


### Bug Fixes

* **cli:** improve global Cobra arg error handling and flag parse UX ([56570f7](https://github.com/repsejnworb/keysej/commit/56570f7873746f61420d5ead98cd59c34b3f142f))
* **goreleaser:** use v2 schema key "checksum"  (singular) to unblock releases ([3a2388a](https://github.com/repsejnworb/keysej/commit/3a2388a74625e8b70fac68f5a1d77f5094a5a0b3))
* sorted sshDir defaults ([feeec3e](https://github.com/repsejnworb/keysej/commit/feeec3ef21e6ac33a2796f6cb3555af3bd973c69))
* **sshconf:** friendly error messages + colors for list when files/hosts are missing ([cffbde8](https://github.com/repsejnworb/keysej/commit/cffbde803618fe9b1268292e25182d8628d32153))


### Features

* **cli:** add global pretty error handling with colorized usage hints and silence cobra spam ([f27e96c](https://github.com/repsejnworb/keysej/commit/f27e96c83b6ed90d90e4c0e482106b3861f55c01))
* **cli:** add initial Cobra + Bubble Tea scaffold ([7ce1cb2](https://github.com/repsejnworb/keysej/commit/7ce1cb23d7f5459d4ce46fa58f4ca6581c2d31da))
* **sshconf:** add tidy to normalize keysej files (sort blocks, trim, final newline) ([aa12eed](https://github.com/repsejnworb/keysej/commit/aa12eed2aab159d0bcf9b200a9dcbbb74a3d3eda))
* **sshconf:** add validate to lint keysej config files (CIDR/Host, IdentityFile, tag) ([a2866b6](https://github.com/repsejnworb/keysej/commit/a2866b6c4c62599142e5770163598c37469f4682))
* **sshconf:** guard new with private key existence and keysej tag; add --force to bypass ([867c27a](https://github.com/repsejnworb/keysej/commit/867c27a0394e8b7d23f7e70df8860aa3ceb090bf))
* **sshconf:** manage ~/.ssh/config.d/keysej.<key>.conf with list/new/delete (Host/CIDR) ([ad331c3](https://github.com/repsejnworb/keysej/commit/ad331c331b749b40eb9f1d0912e3ff5a15deef5a))
