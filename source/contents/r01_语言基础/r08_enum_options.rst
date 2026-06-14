=====================
枚举&选项
=====================

枚举（Enum）
========================

枚举是 Rust 中极其强大的类型系统特性，远不止 C 语言那种简单的整数常量。

1. 基本定义
-----------------------

.. code-block:: rust

  enum IpAddrKind {
      V4,
      V6,
  }

  let four = IpAddrKind::V4;
  let six = IpAddrKind::V6;

1. 枚举可以关联数据
------------------------------------

每个变体可以携带不同类型、不同数量的数据：

.. code-block:: rust

  enum IpAddr {
      V4(String),
      V6(String),
  }

  let home = IpAddr::V4(String::from("127.0.0.1"));
  let loopback = IpAddr::V6(String::from("::1"));

甚至每个变体的数据类型可以完全不同：

.. code-block:: rust

  enum Message {
      Quit,                       // 没有关联数据
      Move { x: i32, y: i32 },   // 具名字段（类似结构体）
      Write(String),             // 单个字符串
      ChangeColor(i32, i32, i32), // 三个 i32
  }

3. 枚举也可以有方法
-----------------------

.. code-block:: rust

  impl Message {
      fn call(&self) {
          match self {
              Message::Quit => println!("Quitting"),
              Message::Move { x, y } => println!("Moving to ({}, {})", x, y),
              Message::Write(text) => println!("Writing: {}", text),
              Message::ChangeColor(r, g, b) => println!("Changing color to ({}, {}, {})", r, g, b),
          }
      }
  }

  let msg = Message::Write(String::from("hello"));
  msg.call();

4. 枚举与 match 的黄金搭档
------------------------------

.. code-block:: rust

  enum Coin {
      Penny,
      Nickel,
      Dime,
      Quarter(UsState),
  }

  fn value_in_cents(coin: Coin) -> u8 {
      match coin {
          Coin::Penny => {
              println!("Lucky penny!");
              1
          }
          Coin::Nickel => 5,
          Coin::Dime => 10,
          Coin::Quarter(state) => {
              println!("State quarter from {:?}!", state);
              25
          }
      }
  }

5. 枚举的内存布局
-----------------------

.. code-block:: rust

  enum Message {
      Quit,
      Move { x: i32, y: i32 },
      Write(String),
      ChangeColor(i32, i32, i32),
  }

Rust 会在编译时计算所有变体中最大的那个尺寸，然后分配足够大的空间。还有一个标签（discriminant）来标识当前是哪个变体。
可以通过 ``std::mem::size_of``查看：

.. code-block:: rust

  println!("{}", std::mem::size_of::<Message>());
  // 通常是 32 字节（8 字节标签 + 24 字节最大变体对齐后的尺寸）

6. 空指针优化（Null Pointer Optimization）
------------------------------------------------

Rust 利用枚举的标签位来做内存优化。常见的 Option<&T>就是典型例子：

.. code-block:: rust

  // Option<&T> 实际上和 &T 大小相同！
  // None 用空指针表示，Some 用有效指针表示
  println!("{}", std::mem::size_of::<Option<&i32>>()); // 8 字节（64位系统）
  println!("{}", std::mem::size_of::<&i32>());          // 也是 8 字节

同理 Box<T>、Vec<T>、String等智能指针类型与 Option组合时也有此优化。

Option 枚举
========================

Option是 Rust 标准库中最重要、使用最频繁的枚举。它取代了其他语言中的 null/ nil/ NULL。

1. 定义
-----------------------

.. code-block:: rust

  enum Option<T> {
      None,     // 没有值
      Some(T),  // 有值，类型为 T
  }

Option被 prelude 自动导入，所以可以直接用 Some和 None，无需写 Option::Some。

2. 基本用法
-----------------------

.. code-block:: rust

  let some_number = Some(5);          // Option<i32>
  let some_string = Some("hello");    // Option<&str>

  let absent_number: Option<i32> = None; // 必须指定类型

3. 为什么 Option 比 null 安全
----------------------------------

在其他语言中，你可以对 null 值调用方法，导致 NullPointerException：

.. code-block:: java

  // Java
  String s = null;
  int length = s.length(); // NullPointerException

但在 Rust 中，Option<T>和 T是完全不同的类型，不能混用：

.. code-block:: rust

  let x: i32 = 5;
  let y: Option<i32> = Some(5);

  let sum = x + y; // 编译错误：cannot add `Option<i32>` to `i32`
  
你必须显式处理​ None的情况：

.. code-block:: rust

  let sum = x + y.unwrap_or(0); // 如果 y 是 None，用 0 代替

4. Option 的常用方法
-----------------------

取值（谨慎使用）

.. code-block:: rust

  let opt = Some(42);

  // 不安全：如果是 None 会 panic
  let val = opt.unwrap();        // 42
  let val = opt.expect("custom message"); // 42，panic 时可自定义消息

  // 安全取值
  let val = opt.unwrap_or(0);    // 42，如果是 None 返回 0
  let val = opt.unwrap_or_else(|| compute_default()); // 延迟计算默认值
  let val = opt.unwrap_or_default(); // 使用类型的 Default 值

转换与操作

.. code-block:: rust

  let opt = Some(5);

  // map：对 Some 内的值做变换
  let doubled = opt.map(|x| x * 2); // Some(10)

  // and_then：链式操作（flat_map）
  let result = opt.and_then(|x| {
      if x > 0 {
          Some(x * 2)
      } else {
          None
      }
  });

  // filter：过滤
  let filtered = opt.filter(|x| *x > 10); // None（因为 5 > 10 为 false）

  // or：提供备选 Option
  let backup = Some(0);
  let result = opt.or(backup); // Some(5)，如果 opt 是 None 则返回 backup

判断与提取

.. code-block:: rust

  let opt = Some(42);

  // 判断
  if opt.is_some() { /* ... */ }
  if opt.is_none() { /* ... */ }

  // 模式匹配
  if let Some(val) = opt {
      println!("Value: {}", val);
  }

  // while let
  let mut opt = Some(0);
  while let Some(i) = opt {
      if i > 5 {
          opt = None;
      } else {
          println!("{}", i);
          opt = Some(i + 1);
      }
  }

5. ? 运算符（错误传播）
-----------------------

?是处理 Option和 Result的语法糖：

.. code-block:: rust

  fn find_user(id: u32) -> Option<String> {
      // 假设从数据库查找
      if id == 1 {
          Some(String::from("Alice"))
      } else {
          None
      }
  }

  fn get_user_email(id: u32) -> Option<String> {
      let user = find_user(id)?; // 如果是 None，立即返回 None
      // 继续处理 user...
      Some(format!("{}@example.com", user))
  }

等价于：

.. code-block:: rust

  fn get_user_email(id: u32) -> Option<String> {
      let user = match find_user(id) {
          Some(u) => u,
          None => return None,
      };
      Some(format!("{}@example.com", user))
  }

6. Option 与迭代器
-----------------------

.. code-block:: rust

  let items = vec![Some(1), None, Some(2), None, Some(3)];

  // 过滤掉 None
  let valid: Vec<_> = items.iter().flatten().collect(); // [1, 2, 3]

  // 或者用 filter_map
  let valid: Vec<_> = items.iter().filter_map(|x| *x).collect();

  // 在 Option 上迭代
  let opt = Some(5);
  for x in opt {
      println!("{}", x); // 只会打印一次
  }

7. 常见模式：组合多个 Option
------------------------------

.. code-block:: rust

  fn combine(a: Option<i32>, b: Option<i32>) -> Option<i32> {
      a.zip(b).map(|(x, y)| x + y)
  }

  let result = combine(Some(3), Some(4)); // Some(7)
  let result = combine(Some(3), None);    // None

实用建议
========================

- 永远不要用 unwrap()除非你 100% 确定不会是 None——在生产代码中用 expect()并提供有意义的信息
- 优先用 ?运算符——比手动 match 更简洁
- 用 if let处理单分支匹配——比 match 更简洁
- 链式调用 map、and_then、or_else——写出流畅的管道式代码
- Option与 Result的区别：Option表示可能有值也可能没有，Result表示操作可能成功也可能失败（附带错误信息）
