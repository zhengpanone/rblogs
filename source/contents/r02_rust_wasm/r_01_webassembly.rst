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

.. WebAssembly_Reference:

参考文档
================

- `Wasm 与 Rust 生态初探（入门篇）`_

.. _`Wasm 与 Rust 生态初探（入门篇）`: https://mp.weixin.qq.com/s/QGdMpC_i-7hGTpwLyV2WOg



