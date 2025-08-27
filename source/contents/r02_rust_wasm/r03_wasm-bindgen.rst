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

Rust 与 JavaScript 交互示例
================================

Node.js环境
-------------------------------

.. _create_wasm_bindgen_node_project:

创建Rust项目
>>>>>>>>>>>>>>>>>

.. code-block:: shell
  
  # 创建一个新的 Rust 库项目
  cargo new hello_wasm_bindgen_node --lib --vcs none

  # 进入项目目录
  cd hello_wasm_bindgen_node
  
  # 添加 wasm-bindgen、web-sys 依赖
  cargo add wasm-bindgen web-sys


.. code-block:: toml
  :caption: Cargo.toml

  [dependencies]
  wasm-bindgen = "0.2"# 核心库
  web-sys = "0.3"# 浏览器 API 绑定（可选，需调用 JS 原生接口时使用）

.. _edit_wasm_bindgen_node_cargo_toml_project:

编辑Cargo.toml
>>>>>>>>>>>>>>>>>

添加 ``[lib]`` 段

- ``crate-type = ["cdylib"]``:指定编译为动态库（供 WASM 使用）;
  
- ``crate-type = ["rlib"]``: 仍然能作为普通 Rust 库依赖（方便单元测试或共享逻辑）。


.. code-block:: toml
  :caption: Cargo.toml

  [lib]
  crate-type = ["cdylib", "rlib"] # 编译为动态库，供 WASM 使用

.. _edit_wasm_bindgen_node_project_lib_rs:

编辑lib.rs
>>>>>>>>>>>>>>>>>

.. code-block:: rust
  :caption: src/lib.rs

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

  #[cfg(test)]
  mod tests {
      use super::*;
      use wasm_bindgen_test::*;

      #[wasm_bindgen_test]
      fn test_greet() {
          let result = greet("Alice");
          assert_eq!(
              result,
              "Hello, Alice! Welcome to the world of Rust and WebAssembly."
          );
      }

      #[wasm_bindgen_test]
      fn test_add() {
          assert_eq!(add(2, 3), 5);
      }

      #[wasm_bindgen_test]
      fn test_math_utils_multiply() {
          let math = MathUtils::new(10);
          assert_eq!(math.multiply(5), 50);
      }
  }

测试
>>>>>>>>>>>>>>>>>

.. code-block:: shell

  # 安装 wasm-pack
  cargo install wasm-pack
  # 运行测试
  wasm-pack test --node

  # 编译为 WASM
  wasm-pack build --target nodejs


使用
>>>>>>>>>>>>>>>>>

编辑package.json

.. code-block:: json
  :caption: package.json

  {
    "name": "hello_wasm_bindgen_node",
    "version": "1.0.0",
    "type": "module",
    "scripts": {
      "start": "node index.js"
    }
  }


编辑index.js

.. code-block:: js
  :caption: index.js

  // Node.js 版本 wasm-bindgen 直接导出函数和类，不需要 init()
  import { add, MathUtils } from './pkg/hello_wasm_bindgen_node.js';

  function run() {
    console.log(add(2, 3)); // 输出：5

    const mathUtils = new MathUtils(10);
    console.log(mathUtils.multiply(5)); // 输出：50
  }

  run();

.. run_wasm-bindgen_node_project_index_js:

运行项目
>>>>>>>>>>>>>>>>>

.. code-block:: shell

  node index.js


Web环境
-------------------------------

.. _create_wasm_bindgen_web_project:

创建Rust项目
>>>>>>>>>>>>>>>>>

.. code-block:: shell
  
  # 创建一个新的 Rust 库项目
  cargo new hello_wasm_bindgen_web --lib --vcs none

  # 进入项目目录
  cd hello_wasm_bindgen_web
  
  # 添加 wasm-bindgen、web-sys 依赖
  cargo add wasm-bindgen web-sys


.. code-block:: toml
  :caption: Cargo.toml

  [dependencies]
  wasm-bindgen = "0.2"# 核心库
  web-sys = "0.3"# 浏览器 API 绑定（可选，需调用 JS 原生接口时使用）

.. _edit_wasm_bindgen_web_cargo_toml_project:

编辑Cargo.toml
>>>>>>>>>>>>>>>>>

添加 ``[lib]`` 段

- ``crate-type = ["cdylib"]``:指定编译为动态库（供 WASM 使用）;
  
- ``crate-type = ["rlib"]``: 仍然能作为普通 Rust 库依赖（方便单元测试或共享逻辑）。


.. code-block:: toml
  :caption: Cargo.toml

  [lib]
  crate-type = ["cdylib", "rlib"] # 编译为动态库，供 WASM 使用

.. _edit_wasm_bindgen_web_project_lib_rs:

编辑lib.rs
>>>>>>>>>>>>>>>>>

.. code-block:: rust
  :caption: src/lib.rs

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

  #[cfg(test)]
  mod tests {
      use super::*;
      use wasm_bindgen_test::*;
      wasm_bindgen_test::wasm_bindgen_test_configure!(run_in_browser);

      #[wasm_bindgen_test]
      fn test_greet() {
          let result = greet("Alice");
          assert_eq!(
              result,
              "Hello, Alice! Welcome to the world of Rust and WebAssembly."
          );
      }

      #[wasm_bindgen_test]
      fn test_add() {
          assert_eq!(add(2, 3), 5);
      }

      #[wasm_bindgen_test]
      fn test_math_utils_multiply() {
          let math = MathUtils::new(10);
          assert_eq!(math.multiply(5), 50);
      }
  }

.. _test_wasm_bindgen_web_project:

测试
>>>>>>>>>>>>>>>>>

.. code-block:: shell

  # 安装 wasm-pack
  cargo install wasm-pack
  # 运行测试
  wasm-pack test --headless --chrome

  # 编译为 WASM
  wasm-pack build --target web


.. _use_wasm-bindgen_web_project_wasm:

使用
>>>>>>>>>>>>>>>>>

编辑index.html

.. code-block:: html
  :caption: index.html

  <!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Rust + WASM Demo</title>
  </head>
  <body>
    <script type="module">
      import init, { add, MathUtils, greet } from "./pkg/hello_wasm_bindgen_web.js";


      async function run() {
        await init(); // 必须初始化 wasm

        console.log(add(2, 3)); // 5
        // 调用Rust导出的greet函数，并传入参数"WebAssembly"
        const message = greet("WebAssembly");
        console.log(message); // Hello, WebAssembly! Welcome to the world of Rust and WebAssembly.

        // 也可以将问候语显示在页面上，例如创建一个段落元素并添加到页面中
        const p = document.createElement("p");
        p.textContent = message;
        document.body.appendChild(p);
        const math = new MathUtils(10);
        console.log(math.multiply(5)); // 50
      }

      run();
    </script>
  </body>
  </html>


JavaScript与Rust交互示例 
================================

从 JavaScript 到 Rust：调用 JS 函数
------------------------------------------------

在 Rust 中调用 JS 的 console.log


创建Rust项目
>>>>>>>>>>>>>>>>>

.. code-block:: shell
  
  # 创建一个新的 Rust 库项目
  cargo new rust_call_js --lib --vcs none

  # 进入项目目录
  cd rust_call_js
  
  # 添加 wasm-bindgen、web-sys 依赖
  cargo add wasm-bindgen web-sys
  cargo add --dev wasm-bindgen-test

编辑Cargo.toml
>>>>>>>>>>>>>>>>>

添加 ``[lib]`` 段

- ``crate-type = ["cdylib"]``:指定编译为动态库（供 WASM 使用）;
  
- ``crate-type = ["rlib"]``: 仍然能作为普通 Rust 库依赖（方便单元测试或共享逻辑）。


.. code-block:: toml
  :caption: Cargo.toml

  [lib]
  crate-type = ["cdylib", "rlib"] # 编译为动态库，供 WASM 使用

编辑lib.rs
>>>>>>>>>>>>>>>>>

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

  #[cfg(test)]
  mod tests {
    use super::*;
    use wasm_bindgen_test::*;
    #[wasm_bindgen_test]
    fn test_greet(){
        greet("World");

    }
  }
  
测试
>>>>>>>>>>>>>>>>>

.. code-block:: shell
  
  # 运行测试
  wasm-pack test --node
  wasm-pack test --headless --chrome

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



.. _wasm-bindgen_Reference:

参考文档
================

- `Rust & WASM 之 wasm-bindgen 基础：让 Rust 与 JavaScript 无缝对话`_

.. _`Rust & WASM 之 wasm-bindgen 基础：让 Rust 与 JavaScript 无缝对话`: https://mp.weixin.qq.com/s?__biz=MzAwNzM0NDE3NA==&mid=2451927750&idx=1&sn=39de88cf70015f2fb54f2a4b360ea333&chksm=8cae4b87bbd9c291b8928517577c8b4cc8f73d40171ef8596bff9e71ba4cafbbf79f7ee530bd&cur_album_id=3982130130738102281&scene=190#rd

