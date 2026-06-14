========================
所有权和借用
========================

所有权（Ownership）
==========================

三条核心规则

- Rust 中每个值都有一个所有者（owner）
- 同一时刻只能有一个所有者
- 当所有者离开作用域，值会被丢弃（drop）

栈 vs 堆
--------------------

理解所有权之前，先分清两种内存：

.. code-block:: rust
    
  // 栈上分配（固定大小，编译期已知）
  let x = 5;          // i32，栈上
  let b = true;       // bool，栈上
  let t = (1, 2);     // 元组，栈上

  // 堆上分配（动态大小，运行时确定）
  let s = String::from("hello"); // String，数据在堆上，指针在栈上
  let v = vec![1, 2, 3];         // Vec<i32>，数据在堆上

栈上的数据：大小固定、拷贝成本低、离开作用域直接弹出
堆上的数据：大小可变、拷贝成本高、需要手动管理释放（Rust 通过所有权自动管理）

1. 移动（Move）
--------------------

当把一个堆分配的值赋给另一个变量时，所有权发生转移：

.. code-block:: rust

  let s1 = String::from("hello");
  let s2 = s1; // s1 的所有权移动到 s2

  println!("{}", s1); // 编译错误：s1 已被移动
  println!("{}", s2); // hello

为什么这样设计？

.. code-block:: rust

  // 如果 s1 和 s2 都持有指针，离开作用域时会 double free
  // Rust 通过移动语义避免了这个问题

内存变化示意：

.. code-block:: rust

  移动前：
  s1 → 栈: (ptr, len, capacity) → 堆: "hello"

  移动后：
  s1 → 已失效（编译器阻止访问）
  s2 → 栈: (ptr, len, capacity) → 堆: "hello"

4. 克隆（Clone）
--------------------

如果你真的需要深度拷贝堆上的数据：

.. code-block:: rust

  let s1 = String::from("hello");
  let s2 = s1.clone(); // 深拷贝，堆上的数据也被复制

  println!("s1 = {}, s2 = {}", s1, s2); // 两者都有效

5. 拷贝（Copy）
--------------------

栈上类型实现了 Copytrait，赋值是拷贝而非移动：

.. code-block:: rust

  let x = 5;
  let y = x; // i32 实现了 Copy，这里是拷贝

  println!("x = {}, y = {}", x, y); // 两者都有效

常见实现了 Copy 的类型：

- 所有整数类型：i32, u32, i64等
- 布尔类型：bool
- 浮点类型：f32, f64
- 字符类型：char
- 元组（仅当所有元素都实现 Copy）：(i32, i32)是 Copy，(i32, String)不是

6. 函数传参时的所有权
--------------------

.. code-block:: rust

  fn main() {
      let s = String::from("hello");
      takes_ownership(s); // s 的所有权被移入函数
      // println!("{}", s); // s 已失效

      let x = 5;
      makes_copy(x); // i32 是 Copy，x 仍然有效
      println!("{}", x); // 
  }

  fn takes_ownership(some_string: String) {
      println!("{}", some_string);
  } // 这里 some_string 被 drop

  fn makes_copy(some_integer: i32) {
      println!("{}", some_integer);
  } // 这里 nothing special

7. 返回值与所有权
--------------------

.. code-block:: rust

  fn gives_ownership() -> String {
      let s = String::from("hello");
      s // 所有权返回给调用者
  }

  fn takes_and_gives_back(s: String) -> String {
      s // 接收所有权，再返回
  }

  fn main() {
      let s1 = gives_ownership();
      let s2 = String::from("world");
      let s3 = takes_and_gives_back(s2);
      // s2 已失效，s3 获得了所有权
  }

借用（Borrowing）
==========================

每次都转移所有权太麻烦，借用让你临时使用值而不获取所有权。

不可变引用（&T）
--------------------

.. code-block:: rust

  fn calculate_length(s: &String) -> usize {
      s.len()
  } // s 离开作用域，但它只是引用，不会被 drop

  fn main() {
      let s = String::from("hello");
      let len = calculate_length(&s); // 借出引用
      println!("The length of '{}' is {}.", s, len); // ✅ s 仍然可用
  }

内存示意：

.. code-block:: rust

  s → 栈: (ptr, len, cap) → 堆: "hello"
          ↑
          |
  &s → 指向 s 的栈上指针

可变引用（&mut T）
--------------------

.. code-block:: rust

  fn change(s: &mut String) {
      s.push_str(", world");
  }

  fn main() {
      let mut s = String::from("hello");
      change(&mut s);
      println!("{}", s); // "hello, world"
  }

借用规则
--------------------

同一时间只能满足以下之一：

- 一个或多个不可变引用（&T）
- 恰好一个可变引用（&mut T）

.. code-block:: rust

  let mut s = String::from("hello");

  let r1 = &s;      //  不可变引用
  let r2 = &s;      //  多个不可变引用没问题
  println!("{} {}", r1, r2); // 最后一次使用不可变引用

  let r3 = &mut s;  // 前面不可变引用已用完
  println!("{}", r3);

错误的例子：

.. code-block:: rust

  let mut s = String::from("hello");

  let r1 = &s;       //
  let r2 = &s;       //
  let r3 = &mut s;   // 不能同时有不可变引用和可变引用
  println!("{}, {}, {}", r1, r2, r3);
  let mut s = String::from("hello");

  let r1 = &mut s;   // 
  let r2 = &mut s;   //  不能同时有两个可变引用
  println!("{}, {}", r1, r2);

悬垂引用（Dangling References）
------------------------------------

Rust 编译器在编译时就杜绝了悬垂引用：

.. code-block:: rust

  fn dangle() -> &String {
      let s = String::from("hello");
      &s // 编译错误：返回局部变量的引用
  } // s 被 drop，返回的引用指向无效内存

  // 正确做法：返回所有权
  fn no_dangle() -> String {
      let s = String::from("hello");
      s // 所有权转移出去
  }

引用的作用域
--------------------

引用的作用域从声明开始，到最后一次使用结束：

.. code-block:: rust

  let mut s = String::from("hello");

  let r1 = &s;       // r1 作用域开始
  let r2 = &s;       // r2 作用域开始
  println!("{} {}", r1, r2); // r1 和 r2 的作用域在这里结束

  let r3 = &mut s;   // 因为前面的引用已经不再使用
  println!("{}", r3);

内部可变性（Interior Mutability）
-----------------------------------------

有些情况下需要在不可变引用下修改值，可以用 RefCell<T>或 Mutex<T>：

.. code-block:: rust

  use std::cell::RefCell;

  let data = RefCell::new(5);

  let mut ref_mut = data.borrow_mut();
  *ref_mut += 1;

  println!("{}", data.borrow()); // 6

但这会在运行时检查借用规则，违反时会 panic。

对比总结
==========================

.. list-table:: 所有权与引用对比
   :widths: auto
   :header-rows: 1

   * - 概念
     - 语法
     - 所有权
     - 能否修改
     - 数量限制
   * - 所有权转移
     - ``let s2 = s1``
     - 转移
     - 取决于 ``mut``
     - 唯一所有者
   * - 不可变引用
     - ``&s``
     - 不转移
     - 否
     - 多个
   * - 可变引用
     - ``&mut s``
     - 不转移
     - 是
     - 唯一
  
记忆口诀

所有权三规则：一值一主，主离值消
借用两规则：读写不共存，引用不悬垂

实用建议
==========================

- 优先用引用而不是转移所有权——除非你真的想让函数接管数据
- 优先用不可变引用——只在需要修改时才用可变引用
- 尽早释放引用——缩小引用的作用域，方便后续使用可变引用
- 理解 NLL（Non-Lexical Lifetimes）——引用的作用域到最后一次使用，不是到花括号结束
- Copy类型不用操心所有权——它们是按位复制的
