=======================
Rust原生跨平台编译
=======================

**交叉编译** 也就是大家说的 **跨平台编译**。

在 Rust 中,跨平台编译有以下主要优势:

1. **无需依赖虚拟机** 不同于 Java 和 .NET 等需要虚拟机的语言,Rust 编译器 **直接将代码编译为机器码**,因此可以直接在目标平台上运行,无需额外的运行时环境,提高了性能。
2. **静态链接** Rust 默认静态链接所有依赖库,生成的可执行文件是独立的,无需依赖共享库即可运行,便于部署和分发。
3. **LLVM 支持** Rust 使用 LLVM 作为编译器后端,LLVM 提供了强大的跨平台支持,能为多种 CPU 架构生成高质量的机器码。
4. **标准库的跨平台支持** Rust 的标准库就设计为跨平台的,它利用了一些跨平台的抽象层,如跨平台系统调用接口,从而使标准库能够在不同操作系统上运行。
5. **编译时单元测试** Rust 的单元测试在编译时就运行,可以确保在发布时,程序在不同平台上的行为是一致的。
   
需要说明的是,虽然 Rust 为跨平台编译提供了很好的支持,但由于不同平台的差异,仍然可能需要一些平台特定的代码。不过相比其他语言,Rust 的跨平台编译支持无疑更加方便和高效。

Rust 目标三元组
====================



要进行跨平台编译，我们需要知道我们要构建的平台的「目标三元组」（target triple）。Rust使用与LLVM[1]相同的格式。格式为 ``<arch><sub>-<vendor>-<sys>-<env>``。

例如:

- ``x86_64-unknown-linux-gnu`` 代表一个64位Linux机器
- ``x86_64-pc-windows-gnu`` 代表一个64位的Windows机器
  
我们可以运行 ``rustc --print target-list`` 将打印出Rust支持的所有目标。这是一段又臭又长的数据信息。

确定我们关心的平台的目标三元组的两种最佳方法是：

1. 在该平台上运行 ``rustc -vV`` ，并查找以 ``host:`` 开头的行——该行的其余部分将是目标三元组
2. 或者在 ``rust platform-support``  [#]_ 页面中查找下面一些比较常见的目标三元组

下面一些比较常见的目标三元组

.. list-table::
  :header-rows: 1
  :widths: 20 60

  * - 目标三元组名
    - 描述
  * - x86_64-unknown-linux-gnu
    - 64位Linux（内核3.2+，glibc 2.17+）
  * - x86_64-pc-windows-gnu
    - 64位MinGW（Windows 7+）
  * - x86_64-pc-windows-msvc
    - 64位MSVC（Windows 7+）
  * - x86_64-apple-darwin
    - 64位macOS（10.7+，Lion+）
  * - aarch64-unknown-linux-gnu
    - ARM64 Linux（内核4.1，glibc 2.17+）
  * - aarch64-apple-darwin
    - ARM64 macOS（11.0+，Big Sur+）	
  * - aarch64-apple-ios	
    - ARM64 iOS
  * - aarch64-apple-ios-sim
    - ARM64上的Apple iOS模拟器
  * - armv7-linux-androideabi
    - ARMv7a Android
	
跨平台编译（Cross-Compilation）
===============================

在 Rust 中实现跨平台编译（Cross-Compilation）通常涉及几个步骤，包括安装目标工具链、配置交叉编译器（如果需要）、以及使用 cargo 编译。下面我给你整理一个完整、通用的步骤说明：

安装 Rust 的目标平台工具链
---------------------------

Rust 使用 rustup 管理工具链，可以为每个目标平台安装对应的 target

.. code-block:: shell

  # 查看支持的目标平台列表
  rustup target list

  # 添加目标平台，例如编译 Windows 64 位
  rustup target add x86_64-pc-windows-gnu
  rustup target add aarch64-unknown-linux-gnu  # ARM Linux
  rustup target add x86_64-apple-darwin       # macOS

.. note::

  不同平台可能有 GNU 或 MSVC 版本，比如 Windows 就有 gnu 和 msvc。

安装目标平台的 C 编译器和工具链
----------------------------------

对于非纯 Rust 项目（有 C/C++ 依赖，如 openssl 或 ring），通常需要目标平台的编译器工具链。例如：

**Linux → Windows**

.. code-block:: shell

  sudo apt install mingw-w64

**macOS → Windows**

.. code-block:: shell

  brew install mingw-w64

使用 Cargo 进行跨平台编译
------------------------------

.. code-block:: shell

  # 编译指定目标
  cargo build --target x86_64-pc-windows-gnu --release

编译完成后，可在：

.. code-block:: shell

  target/x86_64-pc-windows-gnu/release/

找到可执行文件。

使用 cross 工具（推荐）
---------------------------

`Cross`_ 是官方推荐的跨平台编译工具，底层使用 Docker 容器封装环境，免去手动安装交叉编译工具链。

.. code-block:: shell

  cargo install cross

  # 使用 cross 编译
  cross build --target x86_64-pc-windows-gnu --release

.. note::

  cross 会自动下载对应平台的 Docker 镜像，处理依赖，非常适合 Linux 主机编译 Windows/macOS/ARM 可执行文件。


注意事项
============

1. 动态库依赖 
   
    跨平台编译时，要注意目标系统上的动态库依赖，比如 Windows 需要 **msvcrt**，Linux 需要 *libc*。

2. 不同架构的二进制

   - x86_64 → 64 位
   - i686 → 32 位
   - aarch64 → ARM 64 位

.. [#] https://doc.rust-lang.org/nightly/rustc/platform-support.html


.. _`Cross`: https://github.com/cross-rs/cross.git





