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


### Refactors

* drop internal/version; embed version via ldflags and expose --version ([1ec85f3](https://github.com/repsejnworb/keysej/commit/1ec85f36d7fde3727b79d8e33afe4b1e80aabc5a))


### CI

* added temporary manual trigger for automerge while debugging ([03fa2d2](https://github.com/repsejnworb/keysej/commit/03fa2d27e61a1a11c7b61d5a15b370b02c1a80cf))
* automerge release please prs ([4661bee](https://github.com/repsejnworb/keysej/commit/4661beedc44190f80cdb7ee026b6172e259d6c17))
* build and release ([56d8faa](https://github.com/repsejnworb/keysej/commit/56d8faacd671f1a9282d545fb05e77a3a12a4a52))
* debug automerge not executing ([9be582c](https://github.com/repsejnworb/keysej/commit/9be582c7c701de3bf070f34d513caef94549148c))
* fixed incorrect branch matching ([e9d9e45](https://github.com/repsejnworb/keysej/commit/e9d9e45d47fae755d2da28cf1ba9b228843b25aa))
* made automerge the most complex fucking beast ever.. this is stupid ([bb44bf4](https://github.com/repsejnworb/keysej/commit/bb44bf4b9200fe1e0fe2559b00f8329e8b6d4cd0))
* merge release-please PR immediately via gh (no auto-merge / no branch protection required) ([fe6988f](https://github.com/repsejnworb/keysej/commit/fe6988fdee7a50a9d7d4a6e7afa7361c04a1744d))
* only automerge on main, also fixed check-reference ([ab5d21d](https://github.com/repsejnworb/keysej/commit/ab5d21d962e8ead063def2b01b39c53ffd1000e0))
* removed automerging, ill just click the button ([09d1ab5](https://github.com/repsejnworb/keysej/commit/09d1ab5a51808cadbe800466134466fa71851adf))
* sorted release-please syntax and use gh-cli ([d75e1f3](https://github.com/repsejnworb/keysej/commit/d75e1f33df8551f45bb2d01a8886af182a44571f))
* sorted workflow dependency ([0ae73dd](https://github.com/repsejnworb/keysej/commit/0ae73ddd687d219ee24e773461bbb54ac48e4371))
* tell gh-cli what repo we are on about ([ba7cd01](https://github.com/repsejnworb/keysej/commit/ba7cd0146db6943a62722e53e7ff44df5ec42a5c))
* trigger automerge as before but wait for ci to finish ([28571ec](https://github.com/repsejnworb/keysej/commit/28571ece77d934477f39fb3f469c834ebce3a928))
* trigger automerge-workflow after ci finishes ([127153e](https://github.com/repsejnworb/keysej/commit/127153ecaede8b610aef5b68e54baffba90781a2))
* use correct release-please outputs ([88ca9af](https://github.com/repsejnworb/keysej/commit/88ca9af5738f206dff390489ed029a58fdf9e7b3))
* use the correct action version.. ([56d1f5a](https://github.com/repsejnworb/keysej/commit/56d1f5a61718a0b8451be0580015e0d3ba1daeca))
* use the correct botname in matching for automerge ([61075cc](https://github.com/repsejnworb/keysej/commit/61075cce86101cc682ab5aef7fe5c452ab246445))


### Chores

* added CODEOWNERS ([d401bbb](https://github.com/repsejnworb/keysej/commit/d401bbbdd79a0ec150fe9a755e93e54eec7bb94e))
* initial fixes and code ([c53f4f4](https://github.com/repsejnworb/keysej/commit/c53f4f44600bcc7d9731ca6272a63caae0b1ca35))
* made listing a little nicer ([4f140f6](https://github.com/repsejnworb/keysej/commit/4f140f658c21e0ea2aa8103eb1a5ee9d8fa3708c))
* trigger release-please ([7ebd016](https://github.com/repsejnworb/keysej/commit/7ebd01658f22c0783fa1d0d51ff629de7ab195bd))
