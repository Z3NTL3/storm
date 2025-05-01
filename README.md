# storm
Layer 7 (HTTP/HTTPS) stress testing and server performance assessment utility

> [!NOTE]  
> This is a powerful layer 7 stress tester, please do not use it for illegal purposes, see [license](#license) and the [disclaimer](#disclaimer)

**Do not forget to star this repo, otherwise don't use my tool**

##### Features
- Supports tunneling through proxies [`http`, `https`, `socks4`, `socks5`]
- Round Robin for loaded proxy selection
- Integrates with [FastHTTP](https://github.com/valyala/fasthttp) ecosystem
- Ability to provide HTTP: accepts headers and general or specific headers
- CLI

### Usage
The usage, assumes you have common knowledge setting things up and understanding not to actually use this on your local machine but on dedicated servers...

Have the Go toolchain installed on your servers, clone this repo and run the following in the ``storm`` folder

- ``go build .``

For usage you can type ``./storm --help``

Example
- ``./storm lightup --target='https://simpaix.net' --timeout=15s``

> [!TIP]
> Do not forget to update ``/data/proxy.txt`` with actual high quality proxies, you can define the protocol using ``--proto`` flag, by default it assumes ``http``

#### License
See ``LICENSE`` file in this repo.

##### Disclaimer
> This software is provided "as is," without warranty of any kind, express or implied. The creator of this software shall not be held liable for any damages or consequences arising from its use or misuse. Users are solely responsible for how they choose to utilize this software.

### Author
z3ntl3, aka Efdal Sancak [my website](https://z3ntl3.com)
