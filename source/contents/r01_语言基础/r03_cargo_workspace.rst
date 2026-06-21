=============
工程化
=============

.. contents:: 目录
   :depth: 6
   :local:

Cargo
=============


cargo 是 Rust 的包管理工具和构建系统。它可以帮助你管理 Rust 项目的依赖、编译代码、运行测试等。与 Cargo 一起 还有 crates.io_，它为社区提供包注册服务，用户可以将自己的包发布到 crates.io_。

Cargo 负责:

- 项目创建
- 依赖管理
- 编译
- 测试
- 发布
- 文档生成
- 格式化
- 静态检查

.. csv-table:: cargo 常用命令
  :widths: 40, 70 
  :file: ./code/r02_cargo_usage/cargo_commands.csv
  :encoding: utf-8
  :align: left

cargo 上手使用
-----------------------

配置cargo
>>>>>>>>>>>>>

cargo 的配置文件位于 ``$CARGO_HOME/config.toml``，如果没有该文件，可以手动创建。 ``$CARGO_HOME`` 目录默认位于用户的 home 目录下的 ``.cargo`` 文件夹内。你可以通过设置环境变量 ``CARGO_HOME`` 来更改这个目录的位置。



.. code-block:: toml

  [source.crates-io]
  replace-with = 'ustc' # 使用 ustc 镜像替换 crates.io 源

  [source.ustc]
  registry = "https://mirrors.ustc.edu.cn/crates.io-index/" # ustc 镜像源
  
  [source.github]
  registry = "https://github.com/rust-lang/crates.io-index" # github 镜像源
  
  [build]
  target = "x86_64-unknown-linux-gnu" # 设置默认编译目标
  rustflags = ["-C", "target-cpu=native"] # 优化编译选项，针对本地 CPU 进行优化
  # target-cpu 可以设置为 native, x86-64, core2, corei7, skylake 等
  # 具体选项可以参考 https://doc.rust-lang.org/rustc/codegen-options/index.html#target-cpu
  target-dir = "target" # 设置构建输出目录
  # /Users/mac/.target
  # 设置构建输出目录，默认是 target，可以根据需要进行调整
  


创建新项目
>>>>>>>>>>>>>

.. code-block:: bash

  cargo new 02_hello_cargo --name hello_cargo
  cargo new hello_world --bin
  # 可以使用其他的VCS或不使用VCS：cargo new 的时候使用 --vcs这个 flag
  cargo new --vcs=git project_name
  cargo new my_project --vcs none

参数说明

- --vcs none：不创建任何版本控制文件夹
- 其他可选值：

  - git (默认)
  - hg (Mercurial)
  - pijul (Pijul)
  - fossil (Fossil)
  - none (不初始化版本控制)

切换到新的并行编译器前端
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

你可以在 ``Nightly`` 版本中，启用新的并行编译器前端。使用 ``-Z threads=8`` 选项运行 ``Nightly`` 编译器：

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
------------------------------

git 仓库作为依赖包
>>>>>>>>>>>>>>>>>>>>>>>>

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
>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: toml
  
  [dependencies]
  hello_utils = { path = "../hello_utils" }

根据平台引入依赖
>>>>>>>>>>>>>>>>>>>>>>>>

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


crate-type
----------------------

在 Rust 的 ``Cargo.toml`` 中:

.. code-block:: toml

  [lib]
  crate-type = ["cdylib", "rlib"]

表示这个库同时编译为：

.. code-block:: text

  1. rlib   （Rust库）
  2. cdylib （动态链接库）


rlib 是什么
>>>>>>>>>>>>>>>>>>>>>

Rust 默认生成的库格式。

例如：

.. code-block:: toml

  [lib]
  name = "my_lib"

执行：

>>> cargo build

生成：

.. code-block:: text

  target/debug/
  └── libmy_lib.rlib

rlib 只能被 Rust 使用：

.. code-block:: rust

  use my_lib::hello;

不能被：

.. code-block:: text

  Java
  Python
  C
  Go
  Node.js

直接调用。

cdylib 是什么
>>>>>>>>>>>>>>>>>>>>>

cdylib 是给其它语言调用的动态库。

生成：

Linux

.. code-block:: text

  libmy_lib.so

Windows

.. code-block:: text
  
  my_lib.dll

macOS

.. code-block:: text

  libmy_lib.dylib

例如：

.. code-block:: toml

  [lib]
  crate-type = ["cdylib"]

Rust：

.. code-block:: rust

  #[unsafe(no_mangle)]
  pub extern "C" fn add(a: i32, b: i32) -> i32 {
      a + b
  }

编译：

.. code-block:: shell

  cargo build --release

得到：

.. code-block:: text

  my_lib.dll

然后：

Python

.. code-block:: python

  import ctypes

  lib = ctypes.CDLL("./my_lib.dll")

  print(lib.add(1, 2))

为什么同时写

很多项目：

.. code-block:: toml

  [lib]
  crate-type = ["cdylib", "rlib"]

原因：

Rust内部使用

.. code-block:: text

  rlib

供其它 Rust crate 引用。

对外提供接口

.. code-block:: text

  cdylib

供：

.. code-block:: text

  Java JNI
  Python ctypes
  Node.js
  Go cgo
  C/C++

调用。

常见 crate-type
>>>>>>>>>>>>>>>>>>>>>


rlib
::::::::::::::

.. code-block:: toml

  crate-type = ["rlib"]

Rust静态库。

dylib
::::::::::::::

.. code-block:: toml

  crate-type = ["dylib"]

Rust动态库。

依赖 Rust Runtime。

较少使用。

cdylib
::::::::::::::

.. code-block:: toml

  crate-type = ["cdylib"]

给非 Rust 语言调用。

最常见。

staticlib
::::::::::::::

.. code-block:: toml

  crate-type = ["staticlib"]

生成：

Linux：

.. code-block:: text

  libxxx.a

Windows：

.. code-block:: text

  xxx.lib

完全静态链接。

适合：

.. code-block:: text

  C++
  Go
  嵌入式


proc-macro
::::::::::::::

.. code-block:: toml

  crate-type = ["proc-macro"]

过程宏库。

例如：

.. code-block:: rust

  #[derive(Serialize)]

背后就是 proc-macro。

实际项目例子
>>>>>>>>>>>>>>>>>>>>>

PyO3

Python扩展：

.. code-block:: toml

  [lib]
  crate-type = ["cdylib"]

生成：

.. code-block:: text

  xxx.pyd

供 Python 导入。

JNI

Java调用：

.. code-block:: toml

  [lib]
  crate-type = ["cdylib"]

生成：

.. code-block:: text

  xxx.dll

供 JNI 加载。

Rust SDK

既给 Rust 用又给外部用：

.. code-block:: toml

  [lib]
  crate-type = ["cdylib", "rlib"]

例如：

.. code-block:: text

  sdk
  ├── Rust调用
  └── Java调用

查看生成结果

执行：

.. code-block:: shell

  cargo build --release

然后：

.. code-block:: text

  ls target/release

Linux：

.. code-block:: text

  libmy_lib.rlib
  libmy_lib.so

Windows：

.. code-block:: text

  my_lib.dll
  my_lib.rlib

所以：

.. code-block:: toml

  crate-type = ["cdylib", "rlib"]

的含义就是：

同时生成：

.. code-block:: text

  ✓ Rust原生库（rlib）
  ✓ C兼容动态库（cdylib）

  既能给 Rust crate 使用，
  又能给 Java/Python/Go/C++ 等语言调用。



..  _crates.io: https://crates.io/

Package
=================

Package 就是：一个 Cargo 项目。由 Cargo.toml 定义。


关系
---------------

.. code-block:: text

  Package
   ↓
  Crate
   ↓
  Module


Crates
=================

Rust 最核心概念之一。Crate 就是：Rust编译单元。可以理解成：一个Jar，或者：一个Maven模块。

Crate 类型
---------------

Binary Crate
>>>>>>>>>>>>>>>>>>

生成可执行程序。


Library Crate
>>>>>>>>>>>>>>>>>>>

生成库


Module
=================

模块系统。类似 Java package。


Workspace
=================

Workspace类似：Maven Multi Module


Cargo Target
=================

Cargo 项目中包含有一些对象，它们包含的源代码文件可以被编译成相应的包，这些对象被称之为 Cargo Target。库对象 Library 、二进制对象 Binary、示例对象 Examples、测试对象 Tests 和 基准性能对象 Benches 都是 Cargo Target。

库对象(Library)
----------------------

库对象用于定义一个库，该库可以被其它的库或者可执行文件所链接。该对象包含的默认文件名是 src/lib.rs，且默认情况下，库对象的名称跟项目名是一致的，

一个工程只能有一个库对象，因此也只能有一个 src/lib.rs 文件，以下是一种自定义配置:

.. code-block:: toml

  # 一个简单的例子：在 Cargo.toml 中定制化库对象
  [lib]
  crate-type = ["cdylib"]
  bench = false


二进制对象(Binaries)
----------------------

二进制对象在被编译后可以生成可执行的文件，默认的文件名是 src/main.rs，二进制对象的名称跟项目名也是相同的。

大家应该还记得，一个项目拥有多个二进制文件，因此一个项目可以拥有多个二进制对象。当拥有多个对象时，对象的文件默认会被放在 src/bin/ 目录下。

二进制对象可以使用库对象提供的公共 API，也可以通过 [dependencies] 来引入外部的依赖库。

我们可以使用 cargo run --bin <bin-name> 的方式来运行指定的二进制对象，以下是二进制对象的配置示例：

.. code-block:: toml

  # Example of customizing binaries in Cargo.toml.
  [[bin]]
  name = "cool-tool"
  test = false
  bench = false

  [[bin]]
  name = "frobnicator"
  required-features = ["frobnicate"]

示例对象(Examples)
----------------------

示例对象的文件在根目录下的 examples 目录中。既然是示例，自然是使用项目中的库对象的功能进行演示。示例对象编译后的文件会存储在 target/debug/examples 目录下。

如上所示，示例对象可以使用库对象的公共 API，也可以通过 [dependencies] 来引入外部的依赖库。

默认情况下，示例对象都是可执行的二进制文件( 带有 fn main() 函数入口)，毕竟例子是用来测试和演示我们的库对象，是用来运行的。而你完全可以将示例对象改成库的类型:

.. code-block:: toml

  [[example]]
  name = "foo"
  crate-type = ["staticlib"]

如果想要指定运行某个示例对象，可以使用 ``cargo run --example <example-name>`` 命令。如果是库类型的示例对象，则可以使用 ``cargo build --example <example-name>`` 进行构建。

与此类似，还可以使用 ``cargo install --example <example-name>`` 来将示例对象编译出的可执行文件安装到默认的目录中，将该目录添加到 $PATH 环境变量中，就可以直接全局运行安装的可执行文件。

最后，cargo test 命令默认会对示例对象进行编译，以防止示例代码因为长久没运行，导致严重过期以至于无法运行。

测试对象(Tests)
----------------------

测试对象的文件位于根目录下的 tests 目录中，如果大家还有印象的话，就知道该目录是集成测试所使用的。

当运行 cargo test 时，里面的每个文件都会被编译成独立的包，然后被执行。

测试对象可以使用库对象提供的公共 API，也可以通过 [dependencies] 来引入外部的依赖库。

基准性能对象(Benches)
----------------------

该对象的文件位于 benches 目录下，可以通过 cargo bench 命令来运行，关于基准测试，

压力测试和基准测试。前者是针对接口 API，模拟大量用户去访问接口然后生成接口级别的性能数据；而后者是针对代码，可以用来测试某一段代码的运行速度，例如一个排序算法。

基准测试 ``benchmark``，在 Rust 中，有两种方式可以实现：

- 官方提供的 ``benchmark``
- 社区实现，例如 ``criterion.rs``

官方 benchmark
>>>>>>>>>>>>>>>>>>>>

官方提供的测试工具，目前最大的问题就是只能在非 stable 下使用，原因是需要在代码中引入 test 特性: ``#![feature(test)]``

设置 Rust 版本
:::::::::::::::::::

因此在开始之前，我们需要先将当前仓库中的 Rust 版本从 ``stable`` 切换为 ``nightly``:

1. 安装 nightly 版本：$ rustup install nightly
2. 使用以下命令确认版本已经安装成功:
   
   .. code-block:: console

      $ rustup toolchain list
      stable-aarch64-apple-darwin (default)
      nightly-aarch64-apple-darwin (override)

3. 进入项目的根目录，然后运行 ``rustup override set nightly``，将该项目使用的 rust 设置为 ``nightly``
很简单吧，其实只要一个命令就可以切换指定项目的 Rust 版本，例如你还能在基准测试后再使用 ``rustup override set stable`` 切换回 stable 版本。

使用 benchmark
:::::::::::::::::::

当完成版本切换后，就可以开始正式编写 ``benchmark`` 代码了。首先，将 ``src/lib.rs`` 中的内容替换成如下代码：

.. code-block:: rust

  #![feature(test)]

  extern crate test;

  pub fn add_two(a: i32) -> i32 {
      a + 2
  }

  #[cfg(test)]
  mod tests {
      use super::*;
      use test::Bencher;

      #[test]
      fn it_works() {
          assert_eq!(4, add_two(2));
      }

      #[bench]
      fn bench_add_two(b: &mut Bencher) {
          b.iter(|| add_two(2));
      }
  }

可以看出， ``benchmark`` 跟单元测试区别不大，最大的区别在于它是通过 ``#[bench]`` 标注，而单元测试是通过 ``#[test]`` 进行标注，这意味着 ``cargo test`` 将不会运行 ``benchmark`` 代码：

.. code-block:: console

  $ cargo test
  running 2 tests
  test tests::bench_add_two ... ok
  test tests::it_works ... ok

  test result: ok. 2 passed; 0 failed; 0 ignored; 0 measured; 0 filtered out; finished in 0.00s

cargo test 直接把我们的 benchmark 代码当作单元测试处理了，因此没有任何性能测试的结果产生。

对此，需要使用 cargo bench 命令：

.. code-block:: console

  $ cargo bench
  running 2 tests
  test tests::it_works ... ignored
  test tests::bench_add_two ... bench:           0 ns/iter (+/- 0)

  test result: ok. 0 passed; 0 failed; 1 ignored; 1 measured; 0 filtered out; finished in 0.29s

看到没，一个截然不同的结果，除此之外还能看出几点:

单元测试 ``it_works`` 被忽略，并没有执行: ``tests::it_works ... ignored``
``benchmark`` 的结果是 ``0 ns/iter``，表示每次迭代( ``b.iter`` )耗时 ``0 ns``，奇怪，怎么是 0 纳秒呢？别急，原因后面会讲

一些使用建议
:::::::::::::::::::

关于 benchmark，这里有一些使用建议值得大家关注:

- 将初始化代码移动到 ``b.iter`` 循环之外，否则每次循环迭代都会初始化一次，这里只应该存放需要精准测试的代码
- 让代码每次都做一样的事情，例如不要去做累加或状态更改的操作
- 最好让 ``iter`` 之外的代码也具有幂等性，因为它也可能被 ``benchmark`` 运行多次
- 循环内的代码应该尽量的短小快速，因为这样循环才能被尽可能多的执行，结果也会更加准确

谜一般的性能结果
::::::::::::::::::::::::

在写 ``benchmark`` 时，你可能会遇到一些很纳闷的棘手问题，例如以下代码:


.. code-block:: rust

  #![feature(test)]

  extern crate test;

  fn fibonacci_u64(number: u64) -> u64 {
      let mut last: u64 = 1;
      let mut current: u64 = 0;
      let mut buffer: u64;
      let mut position: u64 = 1;

      return loop {
          if position == number {
              break current;
          }

          buffer = last;
          last = current;
          current = buffer + current;
          position += 1;
      };
  }
  #[cfg(test)]
  mod tests {
      use super::*;
      use test::Bencher;

      #[test]
      fn it_works() {
        assert_eq!(fibonacci_u64(1), 0);
        assert_eq!(fibonacci_u64(2), 1);
        assert_eq!(fibonacci_u64(12), 89);
        assert_eq!(fibonacci_u64(30), 514229);
      }

      #[bench]
      fn bench_u64(b: &mut Bencher) {
          b.iter(|| {
              for i in 100..200 {
                  fibonacci_u64(i);
              }
          });
      }
  }

通过cargo bench运行后，得到一个难以置信的结果： ``test tests::bench_u64 ... bench: 0 ns/iter (+/- 0)``, 难道 Rust 已经到达量子计算机级别了？

其实，原因藏在LLVM中: LLVM认为 ``fibonacci_u64`` 函数调用的结果没有使用，同时也认为该函数没有任何副作用(造成其它的影响，例如修改外部变量、访问网络等), 因此它有理由把这个函数调用优化掉！

解决很简单，使用 Rust 标准库中的 ``black_box`` 函数:

.. code-block:: rust

  for i in 100..200 {
      test::black_box(fibonacci_u64(test::black_box(i)));
  }

通过这个函数，我们告诉编译器，让它尽量少做优化，此时 LLVM 就不会再自作主张了:)

.. code-block:: console

  $ cargo bench
  running 2 tests
  test tests::it_works ... ignored
  test tests::bench_u64 ... bench:       5,626 ns/iter (+/- 267)

  test result: ok. 0 passed; 0 failed; 1 ignored; 1 measured; 0 filtered out; finished in 0.67s

这次结果就明显正常了。

criterion.rs
>>>>>>>>>>>>>>>>>>>>

官方 ``benchmark`` 有两个问题，首先就是不支持 ``stable`` 版本的 Rust，其次是结果有些简单，缺少更详细的统计分布。

因此社区 ``benchmark`` 就应运而生，其中最有名的就是 `criterion.rs`_ ，它有几个重要特性:

统计分析，例如可以跟上一次运行的结果进行差异比对
图表，使用 ``gnuplots`` 展示详细的结果图表
首先，如果你需要图表，需要先安装 ``gnuplots``，其次，我们需要引入相关的包，在 ``Cargo.toml`` 文件中新增 :

.. code-block:: toml

  [dev-dependencies]
  criterion = "0.3"

  [[bench]]
  name = "my_benchmark"
  harness = false

接着，在项目中创建一个测试文件: ``$PROJECT/benches/my_benchmark.rs``，然后加入以下内容：

.. code-block:: rust

  use criterion::{black_box, criterion_group, criterion_main, Criterion};

  fn fibonacci(n: u64) -> u64 {
      match n {
          0 => 1,
          1 => 1,
          n => fibonacci(n-1) + fibonacci(n-2),
      }
  }

  fn criterion_benchmark(c: &mut Criterion) {
      c.bench_function("fib 20", |b| b.iter(|| fibonacci(black_box(20))));
  }

  criterion_group!(benches, criterion_benchmark);
  criterion_main!(benches);

最后，使用 cargo bench 运行并观察结果：

.. code-block:: console

  Running target/release/deps/example-423eedc43b2b3a93
  Benchmarking fib 20
  Benchmarking fib 20: Warming up for 3.0000 s
  Benchmarking fib 20: Collecting 100 samples in estimated 5.0658 s (188100 iterations)
  Benchmarking fib 20: Analyzing
  fib 20                  time:   [26.029 us 26.251 us 26.505 us]
  Found 11 outliers among 99 measurements (11.11%)
    6 (6.06%) high mild
    5 (5.05%) high severe
  slope  [26.029 us 26.505 us] R^2            [0.8745662 0.8728027]
  mean   [26.106 us 26.561 us] std. dev.      [808.98 ns 1.4722 us]
  median [25.733 us 25.988 us] med. abs. dev. [234.09 ns 544.07 ns]


.. _criterion.rs: https://bheisler.github.io/criterion.rs/book/criterion_rs.html
