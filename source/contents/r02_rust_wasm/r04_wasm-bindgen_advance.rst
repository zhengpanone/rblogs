===================
wasm-bindgen进阶
===================

``wasm-bindgen`` 在 Rust 和 JavaScript 之间传递基本数据类型。但当面对结构体、集合类型等复杂数据时，这些基础的方法就显得力不从心了。

Rust 结构体到 JavaScript 类的完美转换
========================================

项目初始化
-------------------

.. code-block:: shell

  cargo new wasm-bindgen-advance --lib
  cd wasm-bindgen-advance

编辑Cargo.toml
>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: toml
  :caption: Cargo.toml

  [package]
  name = "wasm-bindgen-advance"
  version = "0.1.0"
  edition = "2024"

  [lib]
  crate-type = ['cdylib', "rlib"]

  [dependencies]
  wasm-bindgen = "0.2.104"

Rust结构体定义与导出
------------------------

wasm-bindgen 为 Rust 结构体提供了无缝的 JavaScript 类生成能力。

.. code-block:: rust
  :caption: src/lib.rs

  use wasm_bindgen::prelude::*;

  #[wasm_bindgen]
  pub struct User {
      // 私有字段需要getter/setter
      name: String,
      // 公开字段，这样JS才能访问
      pub age: u32,
  }

  #[wasm_bindgen]
  impl User {
      // 构造函数
      #[wasm_bindgen(constructor)]
      pub fn new(name: String, age: u32) -> User {
          User { name, age }
      }
      // 方法
      pub fn greet(&self) -> String {
          format!(
              "Hello, my name is {} and I'm {} year old.",
              self.name, self.age
          )
      }
      // getter
      #[wasm_bindgen(getter)]
      pub fn name(&self) -> String {
          self.name.clone()
      }
      // setter
      #[wasm_bindgen(setter)]
      pub fn set_name(&mut self, name: String) {
          self.name = name;
      }
  }

**关键点解析**：

- ``#[wasm_bindgen]``
    标记这个 struct 和 impl 块需要暴露给 JS。

- ``#[wasm_bindgen(constructor)]``
    将 ``new`` 函数指定为 JS 类的构造函数。

- 字段可见性
    默认情况下，Rust 的私有字段在 JS 中无法直接访问。那么可以选择将字段设为 ``pub``，或者提供 getter/setter 方法。

- 内存管理
    JS对象持有Rust内存指针，需手动释放。

Javascript 中使用导出的结构体
-------------------------------   

node版Demo
>>>>>>>>>>>>>>>>>

构建产物到 demos/node/pkg
::::::::::::::::::::::::::::::::::::::::::

.. code-block:: shell

  cd wasm-bindgen-advance
  wasm-pack build --target nodejs --out-dir demos/node/pkg --out-name wasm_bindgen_advance

demos/node/package.json
:::::::::::::::::::::::::::::::::::

.. code-block:: json
  :caption: demos/node/package.json

  {
    "name": "wasm-bindgen-advance",
    "version": "1.0.0",
    "description": "",
    "main": "index.mjs",
    "type": "module",
    "scripts": {
      "start": "node index.mjs"
    }
  }


demos/node/index.mjs
:::::::::::::::::::::::::::::::

.. code-block:: javascript
  :caption: demos/node/index.mjs

  import { User } from './pkg/wasm_bindgen_advance.js'

  async function run() {

    // 使用构造函数创建实例
    const user = new User('Alice', 28);

    // 访问公共字段
    console.log('Age:', user.age);

    // 调用方法
    console.log(user.greet());

    // 使用getter/setter
    console.log('Name:', user.name);
    user.name = 'Bob';
    console.log(user.greet());

    // 别忘了释放内存
    user.free();
  }
  run();

运行

.. code-block:: shell

  cd demos/node
  npm run start

Bundler版Demo
>>>>>>>>>>>>>>>>>

构建产物到 demos/bundler/pkg
:::::::::::::::::::::::::::::::

.. code-block:: shell

  cd wasm-bindgen-advance
  wasm-pack build --target bundler --out-dir demos/bundler/pkg --out-name wasm_bindgen_advance


demos/bundler/package.json
:::::::::::::::::::::::::::::::

.. code-block:: json
  :caption: demos/bundler/package.json

  {
    "name": "wasm-bundler-demo",
    "private": true,
    "type": "module",
    "scripts": {
      "dev": "vite",
      "build": "vite build",
      "preview": "vite preview"
    },
    "devDependencies": {
      "vite": "^7.1.0",
      "vite-plugin-top-level-await": "^1.6.0",
      "vite-plugin-wasm": "^3.5.0"
    }
  }


demos/bundler/index.html
:::::::::::::::::::::::::::::::

.. code-block:: html
  :caption: demos/bundler/index.html

  <!doctype html>
  <html>

  <head>
      <meta charset="utf-8" />
      <title>WASM Bundler Demo</title>
  </head>

  <body>
      <div id="app">Loading...</div>
      <script type="module" src="/main.js"></script>
  </body>
  </html>

demos/bundler/main.js
:::::::::::::::::::::::::::::::

.. code-block:: js
  :caption: demos/bundler/main.js

  // 不要用默认导入，改用命名空间导入，避免 “没有 default 导出” 的语法报错
  import * as mod from './pkg/wasm_bindgen_advance.js';
  // 用 ?url 让 Vite 处理 wasm 资源并返回 URL
  import wasmUrl from './pkg/wasm_bindgen_advance_bg.wasm?url';

  // 兼容不同 glue 导出名：default / init / __wbg_init / initSync
  const initFn = mod.default || mod.init || mod.__wbg_init || mod.initSync;
  // 有 init 就调用（传对象，避免 deprecated 提示）；有的 glue 在导入时已完成初始化，则此处会跳过
  if (typeof initFn === 'function') {
      await initFn({ url: wasmUrl });
  }

  const { User } = mod;

  const u = new User('Carol', 20);
  document.getElementById('app').textContent = [
      u.greet(),
      `name(before set): ${u.name}`,
  ].join('\n');

  u.name = 'Dave';
  const p = document.createElement('pre');
  p.textContent = [
      `name(after set): ${u.name}`,
      u.greet(),
  ].join('\n');
  document.body.appendChild(p);

demos/bundler/vite.config.js
:::::::::::::::::::::::::::::::

.. code-block:: js
  :caption: demos/bundler/vite.config.js

  import { defineConfig } from 'vite';
  import wasm from 'vite-plugin-wasm';
  import topLevelAwait from 'vite-plugin-top-level-await';

  export default defineConfig({
      plugins: [wasm(), topLevelAwait()],
  });

运行
:::::::::::::

Web版Demo
>>>>>>>>>>>>>>>>>

构建产物到 demos/web/pkg
:::::::::::::::::::::::::::::::

.. code-block:: shell

  cd wasm-bindgen-advance
  wasm-pack build --target web --out-dir demos/web/pkg --out-name wasm_bindgen_advance


demos/web/index.html
:::::::::::::::::::::::::::::::

.. code-block:: html
  :caption: demos/web/index.html

  <!doctype html>
  <html>

  <head>
    <meta charset="utf-8" />
    <title>WASM Web Demo</title>
  </head>

  <body>
    <div id="app">Loading...</div>
    <script type="module">
        // // 方式1 传URL
        // import init, { User } from './pkg/wasm_bindgen_advance.js';
        // // 纯浏览器：把 .wasm 的 URL 传给 init
        // await init({ url: new URL('./pkg/wasm_bindgen_advance_bg.wasm', import.meta.url) });

        // 方式2 传字节
        // 1. 导入 JS 胶水代码（不是 .wasm！）
        const wasm = await import('./pkg/wasm_bindgen_advance.js');
        const { default: init, User } = wasm;

        // 2. 加载 WASM 二进制文件
        const wasmPath = './pkg/wasm_bindgen_advance_bg.wasm';
        const response = await fetch(wasmPath);

        if (!response.ok) {
            throw new Error(`Failed to fetch WASM: ${response.status}`);
        }

        const wasmBytes = await response.arrayBuffer();
        await init({ url: wasmBytes });


        // 3. 使用 WASM 导出的功能
        const u = new User('Eve', 22);
        const pre = document.getElementById('app');
        pre.textContent = [
            u.greet(),
            `name(before set): ${u.name}`
        ].join('\n');

        u.name = 'Frank';
        const p = document.createElement('pre');
        p.textContent = [
            `name(after set): ${u.name}`,
            u.greet(),
        ].join('\n');
        document.body.appendChild(p);
    </script>
  </body>

  </html>

启动静态服务器
:::::::::::::::::::::::::::::::

.. code-block:: shell

  # Node
  npx serve -s .
  # 或者
  npx http-server .

  # Python
  python3 -m http.server

**注意**： ``wasm-bindgen`` 生成的 JS 对象持有一个指向 Rust 内存的指针。当不再需要这个对象时，最好调用 ``free()`` 方法来释放 Rust 分配的内存，避免内存泄漏。

Rust 接收 JavaScript 对象
========================================

灵活但类型不安全的JsValue方案
---------------------------------------

``JsValue`` 是一个万能类型，可以表示任何 JavaScript 值（对象、数组、字符串、数字等）。当不知道或不关心 JS 对象的具体结构时，它非常有用。

.. code-block:: rust
  :caption: src/lib.rs

  use js_sys::Reflect;
  use wasm_bindgen::prelude::*;

  #[wasm_bindgen]
  extern "C" {
      // 导入 JS 的 console.log
      #[wasm_bindgen(js_namespace = console)]
      fn log(s: &str);
  }

  #[wasm_bindgen]
  pub fn process_js_object(obj: &JsValue) -> Result<(), JsValue> {
      // 你可以使用 serde_wasm_bindgen 将其反序列化为 Rust 结构体
      log(&format!("Received JS value: {:?}", obj));
      let name_v = Reflect::get(obj, &JsValue::from_str("name"))?;
      if let Some(name) = name_v.as_string() {
          log(&format!("name = {}", name));
      }

      // age
      let age_v = Reflect::get(obj, &JsValue::from_str("age"))?;
      if let Some(age) = age_v.as_f64() {
          log(&format!("age = {}", age));
      }

      Ok(())
  }

注意:

- 需要安装 ``js-sys``, ``cargo add js-sys``

- ``Reflect::get(&JsValue, &JsValue)`` 返回 ``Result<JsValue, JsValue>``，所以函数签名用 -> ``Result<(), JsValue>``.

- ``as_string()`` / ``as_f64()`` 是 ``JsValue`` 自带的方法，用于做“JS → Rust”基础类型提取。

.. code-block:: js
  :caption: demos/node/index.mjs

  import { process_js_object } from './pkg/wasm_bindgen_advance.js'

  async function run() {
      const js_object = {
          id: 101,
          data: 'some payload',
          nested: { a: 1 },
          name: 'Tom',
          age: 20,
      };

      process_js_object(js_object);

      // 别忘了释放内存
      user.free();
  }
  run();

定义特定类型的安全方案
--------------------------------

如果 JS 对象的结构是固定的，我们可以使用 #[wasm_bindgen] 来定义一个类型，专门用来接收它。

.. code-block:: rust
  :caption: src/lib.rs

  // lib.rs
  use wasm_bindgen::prelude::*;

  // 使用 `typescript_type` 来告诉 wasm-bindgen 对应的 TS 类型
  #[wasm_bindgen(typescript_type = "MyJsObject")]
  pub extern "C" {
      // 定义一个类型来映射 JS 对象
      #[wasm_bindgen(extends = js_sys::Object)]
      #[derive(Debug, Clone)]
      type MyJsObject;

      // 定义 getter 方法来访问属性
      #[wasm_bindgen(method, getter)]
      fn id(this: &MyJsObject) -> u32;

      #[wasm_bindgen(method, getter)]
      fn data(this: &MyJsObject) -> String;
  }

  #[wasm_bindgen]
  pub fn process_typed_object(obj: &MyJsObject) {
      // 现在可以安全地访问属性了！
      log(&format!("Received typed object with id: {} and data: '{}'", obj.id(), obj.data()));
  }

.. code-block:: js
  :caption: demos/node/index.mjs

  import { process_typed_object } from './pkg/wasm_bindgen_advance.js'

  async function run() {
    const myObject = {
      id: 101,
      data: 'some payload',
    };

    process_typed_object(myObject);
    // 别忘了释放内存
    user.free();
  }
  run();

``serde`` + ``serde-wasm-bindgen``：终极解决方案
--------------------------------------------------------

``serde``（Rust 序列化标准库）+ ``serde-wasm-bindgen`` 可自动完成 Rust 结构体与 JS 对象的双向转换，兼顾灵活与安全。

``serde`` 是 Rust 生态中用于序列化和反序列化的标准库。 ``serde-wasm-bindgen`` 则是连接 serde 和 ``wasm-bindgen`` 的桥梁，可以自动将 Rust 结构体和 ``JsValue`` 进行相互转换。

.. code-block:: toml
  :caption: Cargo.toml

  [dependencies]
  js-sys = "0.3.81"
  serde = { version = "1.0.228", features = ["derive"] }
  serde-wasm-bindgen = "0.6.5"
  wasm-bindgen = "0.2.104"

结构体加上 ``#[derive(Serialize, Deserialize)]``，然后将函数的参数和返回值类型从具体结构体改为 ``JsValue`` 即可。  

.. code-block:: rust
  :caption: src/lib.rs
  
  use wasm_bindgen::prelude::*;
  use serde::{Serialize, Deserialize};


  #[derive(Serialize, Deserialize)]
  pub struct ComplexData {
      id: u32,
      name: String,
      tags: Vec<String>,
      active: bool,
  }

  // 接收 JS 对象，自动反序列化为 Rust 结构体
  #[wasm_bindgen]
  pub fn process_data_with_serde(val: JsValue) -> Result<JsValue, JsValue> {
      // 1. JsValue -> Rust struct
      let data: ComplexData = serde_wasm_bindgen::from_value(val)?;

      println!("Processed in Rust: {:?}", data.name);

      // 2. Rust struct -> JsValue
      let processed_data = ComplexData {
          id: data.id + 100,
          ..data // 使用struct update 语法
      };

      Ok(serde_wasm_bindgen::to_value(&processed_data)?)
  }

.. code-block:: js
  :caption: demos/node/index.mjs

  import { process_data_with_serde } from './pkg/wasm_bindgen_advance.js'

  async function run() {
      const myData = {
          id: 1,
          name: 'Wasm-Bindgen',
          tags: ['rust', 'webassembly', 'serde'],
          active: true
      };

      try {
          const result = process_data_with_serde(myData);
          console.log('Result from Rust:', result);
          // 输出: {id: 101, name: 'Wasm-Bindgen', tags: ['rust', 'webassembly', 'serde'], active: true}
      } catch (error) {
          console.error('Error from Rust:', error);
      }

      // 别忘了释放内存
      user.free();
  }
  run();

``serde-wasm-bindgen`` 几乎抹平了两种语言间数据结构的差异！











.. _wasm-bindgen_advance_Reference:

参考文档
================

- `Rust & WASM 之 wasm-bindgen 进阶：解锁 Rust 与 JS 的复杂数据交互秘籍`_

.. _`Rust & WASM 之 wasm-bindgen 进阶：解锁 Rust 与 JS 的复杂数据交互秘籍`: https://mp.weixin.qq.com/s?__biz=MzAwNzM0NDE3NA==&mid=2451927754&idx=1&sn=0a7da1231f103543d913488b2ace6ba2&chksm=8cae4b8bbbd9c29daa393c7d647dda045d8c95f393c32b06f62edc70608d86aa79facaa9bc19&cur_album_id=3982130130738102281&scene=189#wechat_redirect
