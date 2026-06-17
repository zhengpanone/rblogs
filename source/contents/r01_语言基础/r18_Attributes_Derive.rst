================================
Attributes and Derive
================================

Attribute（属性）和Derive（派生）——让编译器自动生成样板代码的魔法。

Attributes
=================

Attribute 是给编译器或工具的元数据，告诉它们怎么处理你的代码。属性是应用于 Rust 代码项的元数据，用 #[...]语法表示。


属性的两种形式

- 外部属性​ #[attr]— 应用于后面的项
- 内部属性​ #![attr]— 应用于包含它的整个容器（如 crate、模块）

**条件编译 Attribute**

.. code-block:: rust

  #[cfg(target_os = "linux")]
  fn linux_only() {
      println!("只在 Linux 上运行");
  }

  #[cfg(target_os = "windows")]
  fn windows_only() {
      println!("只在 Windows 上运行");
  }

  #[cfg(debug_assertions)]
  fn debug_only() {
      println!("只在 debug 模式运行");
  }

  #[cfg(not(debug_assertions))]
  fn release_only() {
      println!("只在 release 模式运行");
  }

  fn main() {
      // 根据平台调用不同的函数
      #[cfg(target_os = "linux")]
      linux_only();
      
      #[cfg(target_os = "windows")]
      windows_only();
  }


**Lint 控制 Attribute**

.. code-block:: rust

  // 允许死代码警告
  #[allow(dead_code)]
  fn unused_function() {
      println!("这个函数没被调用，但别警告我");
  }

  // 禁止某个 lint
  #[deny(clippy::all)]
  fn strict_function() {
      // 这里 Clippy 警告会变成错误
  }

  // 全局配置
  #![warn(missing_docs)]// 警告缺少文档

  /// 这个函数有文档
  pub fn documented_function() {}


**Attribute 的位置**

.. code-block:: rust

  // 外部属性（作用于下一个项）
  #[derive(Debug)]
  struct Person {
      // 内部属性（作用于当前项）
      #![allow(dead_code)]
      
      name: String,
      age: u32,
  }

  // 模块级别
  #![warn(missing_docs)]

  mod my_module {
      // 模块内容
  }

  // 函数级别
  fn my_function() {
      #![allow(unused_variables)]
      
      let x = 5;  //  unused 也不警告
  }

Derive Macros（派生宏）
============================

derive属性允许编译器自动为自定义类型实现 trait。这是 Rust 中最常用的元编程特性之一。

