=============
cargo使用
=============

cargo 是 Rust 的包管理工具和构建系统。它可以帮助你管理 Rust 项目的依赖、编译代码、运行测试等。与 Cargo 一起 还有 crate.io_，它为社区提供包注册服务，用户可以将自己的包发布到 crate.io_。

.. csv-table:: cargo 常用命令
  :widths: 50, 70 
  :file: ./code/r02_cargo_usage/cargo_commands.csv
  :encoding: utf-8
  :align: left

cargo 上手使用
===========================

.. code-block:: bash

  cargo new 02_hello_cargo --name hello_cargo
  cargo new hello_world --bin
  # 可以使用其他的VCS或不使用VCS：cargo new 的时候使用 --vcs这个 flag
  cargo new --vcs=git projectName


切换到新的并行编译器前端
----------------------------

你可以在 Nightly 版本中，启用新的并行编译器前端。使用 -Z threads=8 选项运行 Nightly 编译器：

.. code-block:: shell
  
  RUSTFLAGS="-Z threads=8" cargo +nightly build

也可以通过添加 -Z threads=8到~/.cargo/config.toml文件中将其设为默认值：

.. code-block:: toml

  [build]
  rustflags = ["-Z", "threads=8"]

  # 或者在项目的 Cargo.toml 中添加
  [profile.dev]
  rustflags = ["-Z", "threads=8"]

  [profile.release]
  rustflags = ["-Z", "threads=8"]

还可以在 shell 的配置文件中设置别名（例如~/.bashrc或~/.zshrc）

.. code-block:: shell

  # bash, zsh, fish
  alias cargo="RUSTFLAGS='-Z threads=8' cargo +nightly"

- 移除没有的依赖项
  
  - 删除未使用的依赖，减少构建时间和资源消耗及减小项目体积。
  
.. code-block:: shell

  cargo install cargo-machete && cargo machete


- 找出代码库中编译缓慢的 crate
  
  - 运行 ``cargo build --timings`` 命令，这会提供关于每个 crate 编译所花费的时间信息。


Cargo 缓存
----------------------------

Cargo使用缓存来提高构建效率，当执行构建命令时，它会把下载的依赖包存放在CARGO_HOME目录下。该目录默认位于用户的home目录下的.cargo文件夹内。
你可以通过设置环境变量CARGO_HOME来更改这个目录的位置。

.. code-block:: shell

  echo $CARGO_HOME
  echo $HOME/.cargo/


Cargo.toml 文件
===========================

git 仓库作为依赖包
------------------------------

1. 默认不指定版本，从主分支拉去最新 commit
   
.. code-block:: toml

  [dependencies]
  regex = { git = "https://github.com/rust-lang/regex" }

2. 指定分支

.. code-block:: toml

  [dependencies]
  regex = { git = "https://github.com/rust-lang/regex", branch = "next" }

3. 根据tag 拉取指定版本的代码
  
.. code-block:: toml

  [dependencies]
  regex = { git = "https://github.com/rust-lang/regex", tag = "v0.1.0" }

4. 根据 commit hash 拉取指定版本的代码
    
.. code-block:: toml
  
  [dependencies]
  regex = { git = "https://github.com/rust-lang/regex", rev = "c8480030aa6b1ef330874f83ad31e693480c008e" }


任何非 tag 和 branch 的类型都可以通过 rev 来引入 例如 rev= “hash”

通过路径引入本地依赖包
------------------------------

.. code-block:: toml
  
  [dependencies]
  hello_utils = { path = "../hello_utils" }

根据平台引入依赖
------------------------------

.. code-block:: toml

  [target.'cfg(windows)'.dependencies]
  winapi = "0.3"
  winhttp = "0.4.0"

  [target.'cfg(unix)'.dependencies]
  libc = "0.2"
  openssl = "1.0.1"

  [target.'cfg(target_arch = "x86")'.dependencies]
  native = { path = "native/i686" }

  [target.'cfg(target_arch = "x86_64")'.dependencies]
  native = { path = "native/x86_64" }




..  _crate.io: https://crates.io/


