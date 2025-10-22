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








.. _wasm-bindgen_advance_Reference:

参考文档
================

- `Rust & WASM 之 wasm-bindgen 进阶：解锁 Rust 与 JS 的复杂数据交互秘籍`_

.. _`Rust & WASM 之 wasm-bindgen 进阶：解锁 Rust 与 JS 的复杂数据交互秘籍`: https://mp.weixin.qq.com/s?__biz=MzAwNzM0NDE3NA==&mid=2451927754&idx=1&sn=0a7da1231f103543d913488b2ace6ba2&chksm=8cae4b8bbbd9c29daa393c7d647dda045d8c95f393c32b06f62edc70608d86aa79facaa9bc19&cur_album_id=3982130130738102281&scene=189#wechat_redirect
