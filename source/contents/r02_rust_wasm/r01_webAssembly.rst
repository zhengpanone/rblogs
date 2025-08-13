================
WebAssembly
================

什么是 WebAssembly（Wasm）
====================================

WebAssembly，常简称为 Wasm，是一种为 Web 优化的二进制字节码格式，专为现代浏览器设计。它能以接近原生的速度运行，并与 JavaScript 无缝协作。

简单来说，WebAssembly 是一种可移植、体积小、加载快并且兼容 Web 的全新格式。它不是要取代 JavaScript，而是要与 JavaScript 并肩作战！

WebAssembly 的核心优势在于其 **“沙箱安全环境”** 和 **“跨平台可移植性”**。

- 在沙箱环境中，WebAssembly 代码被隔离运行，避免了对系统的潜在风险，就像在一个安全的独立容器里运行程序，不会影响到外面的世界。

- 跨平台可移植性则意味着，编写一次 WebAssembly 代码，可以在不同的浏览器、Node.js 环境，甚至是嵌入式设备中运行，极大地提高了代码的复用性和通用性。

三大核心优势，解锁 Web 新可能
--------------------------------------------

- **高性能计算**：相比 JavaScript，WebAssembly 省去了解释执行和动态类型检查的开销，以接近原生的速度运行，在数值计算、图形渲染等密集型任务中性能提升显著。使得网页运行 3D 游戏、视频处理等复杂应用成为可能。

- **代码复用革命**：WebAssembly 实现了一次编写，多端运行。开发者可以复用现有的 Rust 或 C++ 代码库，通过编译生成 Wasm 模块，轻松嵌入 Web 应用。

- **多语言协同生态**：作为 “语言中立” 的中间格式，WebAssembly 支持多种编程语言编译接入。这形成了以 JavaScript 为 “胶水”、Wasm 为 “引擎” 的混合开发模式。在这种模式下，JavaScript 负责处理 DOM 操作、用户交互等上层逻辑，而 WebAssembly 则专注于执行计算密集型任务，两者相互协作，发挥各自的优势。

应用场景：从浏览器到全平台
------------------------------------------

- **Web 端**：在游戏领域，无论是 2D 还是 3D 游戏，Wasm 都可以让游戏引擎直接在浏览器中运行，提供流畅的游戏体验，无需插件；在音视频领域，实现浏览器内的实时音频分析、效果器、音乐制作软件、视频剪辑、转码、特效处理等；在科学计算与数据可视化方面，处理和渲染海量数据点，能够加速图表的渲染，使数据展示更加流畅和实时。

- **跨平台**：通过 Electron/Tauri 等框架，WebAssembly 可以构建高性能桌面应用。开发者可以用 Rust 编写核心逻辑，然后编译成 Wasm 模块，再结合前端技术，实现一套代码同时运行在 Web 和桌面平台。

- **边缘计算**：在资源受限的嵌入式设备中，WebAssembly 也能发挥重要作用。运行轻量级的 Wasm 模块，既能够满足设备的性能需求，又能保证安全性。例如，在智能摄像头中，WebAssembly 可以运行图像识别算法，实现实时的目标检测和分析 ，而无需依赖强大的云端计算资源。

一站式工具链
=====================

Rust 拥有一套完善的工具链，为开发者提供了全方位的支持，极大地提升了开发体验。

- **rustup**：这是 Rust 的安装程序和版本管理工具，就像是一个贴心的管家，帮助开发者轻松管理 Rust 版本和 Wasm 编译目标。只需要一行命令 ``rustup target add wasm32-unknown-unknown``，就可以添加 Wasm 编译目标，让开发者能够快速开始 Rust 和 WebAssembly 的开发。

- **cargo**：作为 Rust 的构建系统和包管理工具，cargo 原生支持 Wasm 项目。通过 Cargo.toml 文件，开发者可以方便地配置项目依赖。比如，要使用 ``wasm-bindgen`` 库，只需要在 Cargo.toml 中添加 ``wasm-bindgen = "0.2"`` ，就可以轻松引入这个库，为后续的开发做好准备。

- **wasm-pack**：这是一个将 Rust 代码编译、打包为 Wasm 模块，并自动生成 JavaScript 调用接口的工具。它彻底简化了跨语言交互流程，让开发者无需担心 Rust 和 JavaScript 之间的通信问题。使用 ``wasm-pack build --target web`` 命令，就可以一键完成编译和打包，生成的 JavaScript 文件可以直接在浏览器中引入使用 ，大大提高了开发效率。

wasm-bindgen
---------------------

在 Rust 与 WebAssembly 的开发中，wasm-bindgen 是一个非常重要的工具，它就像是一位 “翻译官”，打通了 Rust 与 JavaScript 之间的交流障碍，让两者能够轻松地进行交互。

通过#[wasm_bindgen]宏，Rust 函数可以轻松地导出为 JavaScript 可调用的接口。例如，我们在 Rust 代码中定义一个简单的函数：

.. code-block:: rust

  use wasm_bindgen::prelude::*;

  #[wasm_bindgen]
  pub fn add(a: i32, b: i32) -> i32 {
      a + b
  }

在这段代码中，``#[wasm_bindgen]`` 宏将add函数标记为可以被 JavaScript 调用。编译后，在 JavaScript 中就可以通过导入的方式来调用这个函数，实现无缝交互。

wasm-bindgen 也能让我们在 Rust 中调用浏览器 API。比如，我们想要在 Rust 中操作 DOM 元素，通过 wasm-bindgen 就可以实现。首先，在 Rust 代码中导入相关的模块：

.. code-block:: rust

  use wasm_bindgen::prelude::*;
  use web_sys::window;
  
  #[wasm_bindgen]
  pub fn set_body_text(text: &str) {
    if let Some(window) = window(){
      let Some(document) = window.document(){
        let body = document.body().unwrap();
        body.set_inner_html(text);
    }  
  }

在这段代码中，我们通过wasm_bindgen和web_sys库，在 Rust 中调用了浏览器的window、document和body等 API，实现了设置网页正文内容的功能。

wasm-pack
------------------------

``wasm-pack`` 是另一个不可或缺的工具，它就像是一个 “智能打包机”，能够帮助我们将 Rust 代码编译、打包为 Wasm 模块，并自动生成 JavaScript 调用接口，极大地简化了开发流程。

只需运行 ``wasm-pack build --target web``， ``wasm-pack`` 就会自动完成一系列操作：

- **Rust 代码编译为 Wasm 字节码**： ``wasm-pack`` 会调用 Rust 编译器，将我们编写的 Rust 代码编译成 WebAssembly 字节码，这是在浏览器中运行的核心代码。

- **生成 JS 胶水代码（.js 文件）和类型声明（.d.ts）**：为了让 JavaScript 能够方便地调用 Wasm 模块， ``wasm-pack`` 会生成相应的 JavaScript 胶水代码，这些代码就像是桥梁，连接了 Rust 和 JavaScript。同时，还会生成类型声明文件（.d.ts），可以为 TypeScript 项目提供类型检查和代码提示，提高代码的可维护性。

- **配置 package.json**： ``wasm-pack`` 会读取项目的 Cargo.toml 文件，并生成相应的 package.json 文件。这个 package.json 文件包含了项目的元数据、依赖信息等，支持通过 npm 发布或本地引用。

最终， ``wasm-pack`` 会生成一个 ``pkg/`` 目录，里面包含了编译后的 Wasm 文件、生成的 JavaScript 文件和类型声明文件等。

这个 ``pkg/`` 目录可以直接嵌入 Web 项目，无需手动配置复杂的编译流程。例如，我们在一个 Webpack 项目中使用 wasm-pack 生成的 Wasm 模块，只需要在项目中安装相关依赖，然后在 JavaScript 文件中导入生成的模块即可使用，非常方便。


搭建第一个 Rust Wasm 开发环境
=======================================

安装核心工具 wasm-pack
---------------------------------

``wasm-pack``，它是将 Rust 代码编译、打包为 Wasm 模块，并自动生成 JavaScript 调用接口的核心工具。使用 ``cargo`` 安装 ``wasm-pack`` 非常简单，在终端中执行以下命令：

.. code-block:: shell

  cargo install wasm-pack

``cargo`` 会自动从 crates.io_ 下载 ``wasm-pack`` 并安装到本地。安装完成后，执行 ``wasm-pack --version`` 命令，验证是否安装成功。如果能显示 ``wasm-pack`` 的版本号，说明安装成功。

创建第一个 Rust Wasm 项目
--------------------------------------

使用 ``cargo new`` 命令创建一个新的 Rust 项目：

.. code-block:: shell

  cargo new hello_wasm --lib
  cd hello_wasm

这里我们使用 ``--lib`` 标志，表示创建一个库项目，因为我们要编写的是供 JavaScript 调用的库代码。

配置 Cargo.toml
---------------------------------

在项目根目录下的 ``Cargo.toml`` 文件中，添加对 ``wasm-bindgen`` 的依赖：

.. code-block:: toml

  [package]
  name = "hello_wasm"
  version = "0.1.0"
  edition = "2021"

  [lib]
  crate-type = ["cdylib"]

  [dependencies]
  wasm-bindgen = "0.2"

[lib]部分告诉 Cargo 我们要生成一个动态链接库（cdylib），这是 WebAssembly 所需要的格式。

[dependencies]部分添加了 ``wasm-bindgen`` 依赖，它是 Rust 与 JavaScript 之间进行交互的重要工具，版本号0.2表示我们使用的是 0.2 版本的 ``wasm-bindgen``。

添加依赖后，Cargo 会自动下载并管理这些依赖。

编写 Rust 逻辑（src/lib.rs）
-----------------------------------------

.. code-block:: rust

  use wasm_bindgen::prelude::*;

  // 使用#[wasm_bindgen]宏将函数导出为JavaScript可调用的接口
  #[wasm_bindgen]
  pub fn greet(name: &str) -> String {
      format!("Hello, {}!", name)
  }

在这段代码中，我们首先导入了 ``wasm_bindgen::prelude::*`` ，这是使用 ``wasm_bindgen`` 的必要步骤，它包含了一些常用的宏和类型定义。

然后，我们定义了一个 ``greet`` 函数，这个函数接受一个字符串类型的参数 ``name`` ，并返回一个格式化后的问候语字符串。

``#[wasm_bindgen]`` 宏将这个函数标记为可以被 JavaScript 调用的接口，这样在编译后，我们就可以在 JavaScript 中调用这个函数了。

编译为 Wasm 模块
----------------------------------------

编写好 Rust 代码后，接下来我们使用 ``wasm-pack`` 将其编译为 Wasm 模块，并生成 JavaScript 调用接口。在项目根目录下的终端中执行以下命令：

.. code-block:: shell

  wasm-pack build --target web

``wasm-pack build`` 命令会将 Rust 代码编译为 Wasm 模块，并生成相应的 JavaScript 胶水代码和类型声明文件。

``--target web`` 参数表示我们生成的 Wasm 模块是用于 Web 环境的。

编译完成后，项目目录下会生成一个 ``pkg`` 目录，里面包含了编译后的文件：

.. code-block:: console

  pkg
  ├── hello_wasm_bg.wasm  // WebAssembly字节码文件
  ├── hello_wasm_bg.wasm.d.ts  // Wasm模块的类型声明文件（TypeScript）
  ├── hello_wasm.d.ts  // JavaScript调用接口的类型声明文件（TypeScript）
  ├── hello_wasm.js  // JavaScript调用接口文件，用于在JavaScript中调用Wasm模块
  └── package.json  // 项目的元数据和依赖信息文件，支持通过npm发布或本地引用

其中， ``hello_wasm_bg.wasm`` 是 WebAssembly 字节码文件，它包含了我们编写的 Rust 代码编译后的机器码，是在浏览器中实际执行的代码。

``hello_wasm.js`` 是 JavaScript 调用接口文件，它提供了在 JavaScript 中调用 Wasm 模块的方法，通过这个文件，可以在 JavaScript 中轻松地调用 Rust 编写的函数。

在JavaScript中调用
-----------------------------------------

接下来，我们创建一个 HTML 页面来调用编译好的 Wasm 模块。在项目根目录下创建一个index.html文件，编写以下代码：

.. literalinclude:: code/r01_webAssembly/index.html
  :caption: index.html
  :language: html

在这段 HTML 代码中，我们首先通过 ``<script type="module">`` 标签引入了 ``pkg/hello_wasm.js`` 文件，这是wasm-pack生成的 JavaScript 调用接口文件。

然后，我们调用 ``init()`` 函数来初始化 Wasm 模块，这个函数会加载并实例化 ``hello_wasm_bg.wasm`` 文件。

初始化完成后，我们调用 ``greet('WebAssembly')`` 函数，传入参数 ``WebAssembly``，并将返回的问候语字符串打印到控制台，同时也通过创建段落元素的方式将问候语显示在页面上。


运行项目
------------------------------------------

完成以上步骤后，我们就可以在浏览器中运行这个 Rust+Wasm 的 “Hello World” 程序了。

在项目根目录下创建server.js文件，用于启动一个web服务，实现在浏览器中访问页面：

.. literalinclude:: code/r01_webAssembly/server.js
  :caption: server.js
  :language: javascript

然后在项目根目录下的终端中执行以下命令：

.. code-block:: shell

  node server.js


打开浏览器，访问 http://localhost:3000/ ，你会在浏览器的控制台中看到输出的问候语，同时页面上也会显示出问候语：

.. code-block:: console

  Hello, WebAssembly! Welcome to the world of Rust and WebAssembly.

这样，我们就成功地在浏览器中运行了一个用 Rust 编写的 “Hello, World!” 程序，


.. WebAssembly_Reference:

参考文档
================

- `Wasm 与 Rust 生态初探（入门篇）`_

.. _`Wasm 与 Rust 生态初探（入门篇）`: https://mp.weixin.qq.com/s/QGdMpC_i-7hGTpwLyV2WOg


..  _crates.io: https://crates.io/



