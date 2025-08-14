=======================
wasm-pack
=======================

wasm-pack介绍
=======================

wasm-pack 核心原理
------------------------------

``wasm-pack`` 是 Rust 生态中专门用于构建 WebAssembly 模块的高效工具，其核心在于整合 ``rustc``、 ``wasm-bindgen`` 等工具链，实现 Rust 代码到 Wasm 模块的无缝转换。

它通过解析 ``Cargo.toml`` 配置，自动生成 JavaScript 胶水代码（Glue Code），解决 Rust 与 JavaScript 之间的数据类型映射和函数调用问题，同时确保编译后的 Wasm 模块符合 Web 标准规范，支持在浏览器和 Node.js 环境中运行。

wasm-pack 核心功能
------------------------------

- 一键式构建
  
  - 通过 ``wasm-pack build`` 命令，自动完成 Rust 代码编译、Wasm 模块优化及 JavaScript 绑定生成，无需手动配置复杂的编译参数。
  
- 包管理集成
  
  - 自动生成 ``package.json``，支持通过 npm 进行依赖管理和版本发布，兼容现代前端工程化流程。
  
- 多环境适配
  
  - 支持生成不同环境（浏览器、Node.js、Web Workers）的构建产物，适配多样化的运行场景。
  
- 开发友好性
  
  - 提供清晰的错误提示和构建日志，集成 ``wasm-bindgen-test`` 实现浏览器端测试，提升开发调试效率。

wasm-pack 工作流程
--------------------------

当你运行 wasm-pack build 时，背后发生了什么？一张图看懂：

.. code-block:: text

              ┌──────────────────┐
              │ #[wasm_bindgen]  │
              │ fn my_func() {}  │
              └──────────────────┘
                      │
                      ▼
              ┌────────────────┐
              │ wasm-pack build│
              └────────────────┘
                      │
          ┌───────────┴────────────┐
          │                        │
          ▼                        ▼
  ┌─────────────────┐       ┌───────────────────┐
  │ cargo build     │       │   wasm-bindgen    │
  │(编译成 .wasm)    │──────▶│(生成 JS/TS胶水 )    │
  └─────────────────┘       └───────────────────┘
          │
          │
          ▼
  ┌──────────────────────────────────┐
  │    一个完整的 npm 包 （pkg目录）    │
  │                                  │
  │ - a_project_bg.wasm （核心 Wasm） │
  │ - a_project.js      （JS 接口）   │
  │ - a_project.d.ts    （TS 类型）   │
  │ - package.json      （包描述 ）    │
  │                                  │
  └──────────────────────────────────┘

整个过程全自动完成，只需敲下一行命令，就能得到一个可以直接在前端项目中 import 的现代化模块！

wasm-pack 核心命令
-----------------------

- 创建
  
  wasm-pack neww 命令用于创建一个新的 RustWasm 项目。它依赖cargo-generate命令，会自动生成一个 Rust 项目结构，包括 src/lib.rs、Cargo.toml 等文件。

  .. code-block:: bash

    wasm-pack new <name> --template <template> --mode <normal|noinstall|force>

  - ``<name>``: 项目名称
  - ``--template``: 可选，指定模板，默认模版：rustwasm/wasm-pack-template
  - ``--mode``: 可选，指定生成模式，如 ``normal``（默认）、 ``noinstall``、 ``force``。

- 构建
  
  ``wasm-pack build`` 命令用于编译 Rust 代码并生成 Wasm 模块和 JavaScript 绑定。它会自动处理 Rust 代码到 Wasm 的转换，同时生成 JavaScript 接口文件，使得 Rust 函数可以在前端代码中被调用。

  .. code-block:: bash

    wasm-pack build [--target <target>] [--release] [--dev]
  
  - ``--target``: 可选，指定目标平台，如 web、nodejs、bundler 等。
  - ``--release``: 可选，启用发布模式，进行优化编译。
  - ``--dev``: 可选，启用开发模式，生成调试信息。
  
- 测试
  
  ``wasm-pack test`` 命令用于运行测试用例。它会自动编译 Rust 代码并执行测试用例，确保代码的正确性。

  .. code-block:: bash

    wasm-pack test [--headless] [--firefox] [--chrome] [--safari] [--edge] [--node]

  - ``--headless``: 可选，启用无头模式，在无头浏览器中运行测试。
  - ``--firefox``、 ``--chrome``、 ``--safari``、 ``--edge``、 ``--node``:可选，指定要运行的浏览器或环境。

- 打包
  
  ``wasm-pack pack`` 命令用于将构建好的 Wasm 模块和 JavaScript 绑定打包成.tgz文件，tgz文件会生成在pkg目录下。用于发布到 npm 仓库。

  .. code-block:: bash

    wasm-pack pack [--pkg-dir <pkg-dir>]

  - ``--pkg-dir``:可选，指定输出目录，默认为 pkg。

- 发布
  
  ``wasm-pack publish`` 命令用于将构建好的 Wasm 模块和 JavaScript 绑定发布到 npm 仓库。

  .. code-block:: bash

    wasm-pack publish [--target <bundler|nodejs|web|no-modules>] [--access <public|restricted>] [--tag] [--pkg-dir <pkg-dir>]

  - ``--access``:可选，指定包的访问权限，默认为 ``public``；
  - ``--tag``:可选，指定标签，默认为 ``latest``；
  - ``--pkg-dir``:可选，指定输出目录，默认为 ``pkg``；
  - ``--target``:可选，指定目标平台，默认为 ``bundler``；

- 帮助
  
  ``wasm-pack --help`` 查看命令详细说明

wasm-pack 开发指南
=========================

- 使用 ``rustup target add wasm32-unknown-unknown`` 配置目标架构为 ``wasm32-unknown-unknown``，确保编译输出为 Wasm 格式。
- 使用 ``cargo install wasm-pack`` 安装，支持通过命令行快速调用核心功能。

初始化项目
----------------------------

.. code-block:: bash

  wasm-pack new hello-wasm
  wasm-pack new my_project --template rustwasm/wasm-pack-template

代码开发与接口导出
------------------------------

配置 Cargo.toml
>>>>>>>>>>>>>>>>>>>

指定 ``crate`` 类型为 ``cdylib``，引入 ``wasm-bindgen`` 依赖：

.. code-block:: toml

  [package]
  name = "hello_wasm"
  version = "0.1.0"
  edition = "2021"

  # cdylib 允许 Rust 生成一个动态链接库，这是 WebAssembly 所需要的格式。
  # rlib 允许 Rust 生成可供测试使用的库文件
  [lib]
  name = "hello_wasm"
  crate-type = ["cdylib", "rlib"]


  [dependencies]
  wasm-bindgen = "0.2.100"


  [dev-dependencies]
  wasm-bindgen-test = "0.3.50"

  [profile.release]
  # 优化 WebAssembly 包大小
  opt-level = "s"
  lto = true


编写 Rust 逻辑
>>>>>>>>>>>>>>>>>>>

在 ``src/lib.rs`` 中定义可导出的函数，通过 ``wasm-bindgen`` 宏声明对外接口：

.. code-block:: rust

  use wasm_bindgen::prelude::*;

  // 使用#[wasm_bindgen]宏将函数导出为JavaScript可调用的接口
  #[wasm_bindgen]
  pub fn greet(name: &str) -> String {
      format!(
          "Hello, {}! Welcome to the world of Rust and WebAssembly.",
          name
      )
  }

  #[wasm_bindgen]
  pub fn add(a: i32, b: i32) -> i32 {
      a + b
  }

  #[wasm_bindgen]
  pub struct MathUtils {
      value: i32,
  }

  #[wasm_bindgen]
  impl MathUtils {
      #[wasm_bindgen(constructor)]
      pub fn new(value: i32) -> Self {
          MathUtils { value }
      }

      pub fn multiply(&self, other: i32) -> i32 {
          self.value * other
      }
  }

测试驱动开发（TDD）
-----------------------------

单元测试
>>>>>>>>>>>>>>>>>>>>

使用 ``wasm-bindgen-test`` 编写浏览器端测试用例，保存于 ``tests/`` 目录：

.. code-block:: rust

  #[cfg(test)]
  mod tests {
      use wasm_bindgen_test::*;
      // 直接引用当前 crate 的函数
      use hello_wasm::{add,MathUtils};

      // 配置测试在浏览器中运行
      wasm_bindgen_test_configure!(run_in_browser);

      #[wasm_bindgen_test]
      fn test_add() {
          assert_eq!(add(2, 3), 5);
          assert_eq!(add(10, 15), 25);
          assert_eq!(add(100, 200), 300);
      }

      #[wasm_bindgen_test]
      fn test_math_utils() {
          let math_utils = MathUtils::new(4);
          // 通过 multiply(1) 来验证初始值
          assert_eq!(math_utils.multiply(3), 12);
      }
  }

运行测试
>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: bash

  wasm-pack test --headless --chrome # 在无头浏览器中执行测试  

构建、打包与发布
=========================

构建优化与产物生成
-----------------------------

开发模式构建
>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: bash

  wasm-pack build --dev # 生成未优化的调试版本，包含详细调试信息 

输出目录 `pkg/` 包含 `.wasm` 文件、JS 绑定文件及 `package.json`，支持直接引入前端项目。

生产模式构建
>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: bash

  wasm-pack build --release  # 启用编译器优化，减小文件体积并提升执行效率 
  wasm-pack build --scope zhengpanone

通过配置 ``Cargo.toml`` 中的 ``[profile.release]`` ，可进一步优化编译参数（如开启链接时优化 ``lto = true``, 即 ``Link Time Optimization`` （链接时优化））。 

发布到 npm 仓库
----------------------------

初始化 npm 包
>>>>>>>>>>>>>>>>>>>

.. code-block:: bash

  cd pkg  
  wasm-pack pack  # 生成 tgz 包 
  wasm-pack publish --access public # 发布到 npm 仓库










.. WebAssembly-pack_Reference:

参考文档
================

- `Wasm 与 Rust 生态初探（入门篇）`_

.. _`Wasm 与 Rust 生态初探（入门篇）`: https://mp.weixin.qq.com/s/QGdMpC_i-7hGTpwLyV2WOg
