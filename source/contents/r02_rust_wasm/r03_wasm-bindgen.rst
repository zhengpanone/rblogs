=======================
wasm-bindgen
=======================

在 Web 开发中， ``WASM`` 打破了语言壁垒，而 ``wasm-bindgen`` 则架起了 Rust 与 JavaScript 之间的高速桥梁。作为 Rust 与 JavaScript 交互的「翻译官」， ``wasm-bindgen`` 让两种语言突破边界，实现了「双向调用」。

wasm-bindgen简介
=======================

``wasm-bindgen`` 是由 Mozilla Research 开源的一个 Rust 库和 CLI 工具，它可以允许你在 Rust 中调用 JavaScript 中的函数和方法，也可以将 Rust 中的函数方法暴露给 JavaScript 调用。旨在简化 Rust 生成的 WebAssembly 模块与 JavaScript 能够方便、安全地交换复杂数据类型（如字符串、对象、函数）和调用彼此的函数。它通过生成胶水代码（glue code），使得 Rust 函数可以像普通 JavaScript 函数一样被调用，反之亦然。 通过 wasm-bindgen，开发者可以轻松地将 Rust 编写的高性能逻辑与 JavaScript 的灵活性相结合，从而充分利用两者的优点。 ``wasm-bindgen`` 处理了复杂的数据类型转换、内存管理和异步操作，使得开发者可以专注于业务逻辑，而无需担心底层的实现细节。

核心定位
-------------------------------

- 跨语言交互枢纽
    专为 Rust 和 JavaScript 设计，支持 Rust 函数导出到 JS，也允许 Rust 调用 JS 函数。
- 高层级抽象
    无需手动处理内存分配（如线性内存），自动处理数据类型映射（字符串、对象、数组等）。
- 生态集成
    与 ``wasm-pack`` 配合，一键生成 JS 绑定代码，兼容 Webpack、Vite 等前端构建工具。

核心优势
---------------------------

- 类型安全
    严格校验跨语言数据类型，避免运行时错误。
- 零运行时开销
    生成的胶水代码（glue code）轻量高效，不引入额外性能损耗。
- 渐进式集成
    支持从简单函数调用到复杂类结构的交互，适配不同项目规模。

底层原理
=======================

核心原理
------------------------------

``wasm-bindgen`` 的核心原理是通过在 Rust 和 JavaScript 之间生成绑定代码，从而实现两者的交互。这些绑定代码会处理类型转换、内存管理等复杂问题，使得开发者可以专注于业务逻辑的实现，而无需担心底层的细节。

1. 代码标注
    通过 ``#[wasm_bindgen]`` 宏标记需交互的函数、结构体或枚。 ``wasm-bindgen`` 会分析 Rust 代码中的 ``#[wasm_bindgen]`` 标记。
2. 类型转换
    自动生成 Rust 端和 JavaScript 端的绑定代码，处理复杂类型（字符串、数组、对象等）在 WASM 线性内存与 JavaScript 堆之间的转换。
#. 函数映射
    将被标记的 Rust 函数转换为 JavaScript 可调用的函数；反之，也可以将 JavaScript 函数或 Web API 包装成 Rust 可调用的形式。
#. 错误桥接
    将 Rust 的 ``Result`` 转换为 JavaScript 的异常，或将 JavaScript 异常转换为 Rust 的 Result。
#. 胶水代码
    ``wasm-bindgen`` 工具根据元数据生成 JS 接口代码，屏蔽底层 Wasm 细节。

数据类型映射规则
-----------------------------

.. list-table:: Rust 与 JavaScript 类型对照
   :header-rows: 1
   :widths: 25 25 50

   * - Rust 类型
     - JavaScript 类型
     - 示例场景
   * - 基本类型（i32、f64、bool）
     - 对应原始类型
     - 数值计算、逻辑判断
   * - &str / String
     - String
     - 文本处理、日志输出
   * - &[T] / Vec<T>
     - TypedArray 或 Array
     - 数组数据传递（如图片像素）
   * - Rust 结构体 / 类
     - JS 对象
     - 复杂数据结构交互（如配置项）

核心功能与基础用法
========================

环境准备
-----------------------------------

.. code-block:: toml
  :caption: Cargo.toml

  [dependencies]
  wasm-bindgen = "0.2"# 核心库
  web-sys = "0.3"# 浏览器 API 绑定（可选，需调用 JS 原生接口时使用）


从 Rust 到 JavaScript：导出函数供 JS 调用
------------------------------------------------

.. code-block:: rust
  :caption: src/lib.rs

  use wasm_bindgen::prelude::*;

  // 导出函数到 JS，支持基本类型参数和返回值
  #[wasm_bindgen]
  pub fn add(a: i32, b: i32) -> i32 {
    a + b
  }

  // 导出结构体及方法
  pub struct MathUtils{
    base: i32,
  }

  #[wasm_bindgen]
  impl MathUtils {
    // 构造函数
    #[wasm_bindgen(constructor)]
    pub fn new(base: i32) -> MathUtils {
        MathUtils { base }
    }

    // 导出方法到 JS
    pub fn multiply(&self, factor: i32) -> i32 {
        self.base * factor
    }
  }

编译rust代码

.. code-block:: shell

  wasm-pack build --target web

JS 调用 Rust 函数：调用 Rust 函数
----------------------------------

.. code-block:: javascript
  :caption: index.js

  import init, { add, MathUtils } from './pkg/hello_wasm-bindgen.js';

  async function run() {
    // 初始化 WASM 模块
    await init();

    // 调用导出的 Rust 函数
    console.log(add(2, 3)); // 输出：5

    // 使用导出的结构体
    const mathUtils = new MathUtils(10);
    console.log(mathUtils.multiply(5)); // 输出：50
  }

  run();

从 JavaScript 到 Rust：调用 JS 函数
------------------------------------------------

在 Rust 中调用 JS 的 console.log

.. code-block:: rust
  :caption: src/lib.rs

  use wasm_bindgen::prelude::*;

  // 导入 JS 函数
  #[wasm_bindgen]
  extern "C" {
    // 从全局作用域导入，等价于调用 window.console.log
    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);
  }

  // Rust 函数中调用 JS 函数
  #[wasm_bindgen]
  pub fn greet(name: &str) {
    log(&format!("Hello, {}!", name)); // 在浏览器控制台输出
  }

动态导入JS模块
---------------------------------

JavaScript代码，导出一个函数和一个类。

.. code-block:: javascript
  :caption: defined-in-js.js

  export function jsFunction() {
    return "Hello, from JS Function";
  }

  export class JavaScriptClass {
    constructor() {
      this._text = "This is from a JS class.";
    }
    getText() {
      return this._text;
    }
    setText(text) {
      this._text = text;
    }
    render() {
      return "This is render method." + this._text;
    }
  }

在 Rust 中指定这个js文件，声明外部的函数和类型，然后就可以在 Rust 中使用了。

.. code-block:: rust
  :caption: src/lib.rs

  use wasm_bindgen::prelude::*;

  // 导入 JS 函数
  #[wasm_bindgen(module = "./defined-in-js.js")]
  extern "C" {
    fn jsFunction() -> String;

    // 导入 JS 类
    type JavaScriptClass;

    #[wasm_bindgen(constructor)]
    fn new() -> JavaScriptClass;

    #[wasm_bindgen(method, getter)]
    fn getText(this: &JavaScriptClass) -> String;

    #[wasm_bindgen(method, setter)]
    fn setText(this: &JavaScriptClass, text: &str);

    #[wasm_bindgen(method)]
    fn render(this: &JavaScriptClass) -> String;
  }

  #[wasm_bindgen]
  extern "C" {
      #[wasm_bindgen(js_namespace = console)]
      fn log(s: &str);
  }

  // Rust 函数调用 JS 函数和类
  #[wasm_bindgen]
  pub fn call_js_function() -> String {
      jsFunction()
  }

  #[wasm_bindgen]
  pub fn use_js_class() -> String {
      let js_instance = JavaScriptClass::new();
      js_instance.setText("Hello from Rust!");
      js_instance.render()
  }

  // wasm模块初始化后调用
  #[wasm_bindgen(start)]
  fn run() {
      log(&format!("Hello from {}!", name())); // 输出 "Hello from Rust!"

      let x = MyClass::new();
      assert_eq!(x.number(), 42);
      x.set_number(10);
      log(&x.render());
  }









.. wasm-bindgen_Reference:

参考文档
================

- `Rust & WASM 之 wasm-bindgen 基础：让 Rust 与 JavaScript 无缝对话`_

.. _`Rust & WASM 之 wasm-bindgen 基础：让 Rust 与 JavaScript 无缝对话`: https://mp.weixin.qq.com/s?__biz=MzAwNzM0NDE3NA==&mid=2451927750&idx=1&sn=39de88cf70015f2fb54f2a4b360ea333&chksm=8cae4b87bbd9c291b8928517577c8b4cc8f73d40171ef8596bff9e71ba4cafbbf79f7ee530bd&cur_album_id=3982130130738102281&scene=190#rd

