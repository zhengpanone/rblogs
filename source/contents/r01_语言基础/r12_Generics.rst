==============================
泛型（Generics）
==============================

泛型是 Rust 中实现代码复用和类型安全的强大机制。通过泛型，可以用同一段代码处理不同类型的数据，避免重复编写相似逻辑。

在 Rust 中，泛型广泛应用于：

- 函数泛型
- 结构体泛型
- 枚举泛型
- 方法泛型
- Trait 约束（Trait Bound）

泛型在编译期通过**单态化（Monomorphization）**为每个具体类型生成专门代码，运行时**零成本抽象**。

.. contents:: 目录
   :depth: 3
   :local:

为什么需要泛型
==================

没有泛型时，需要为每个类型写一份重复代码：

.. code-block:: rust

   fn largest_i32(list: &[i32]) -> &i32 {
       let mut largest = &list[0];
       for item in list {
           if item > largest {
               largest = item;
           }
       }
       largest
   }

   fn largest_f64(list: &[f64]) -> &f64 {
       let mut largest = &list[0];
       for item in list {
           if item > largest {
               largest = item;
           }
       }
       largest
   }

   fn largest_char(list: &[char]) -> &char {
       let mut largest = &list[0];
       for item in list {
           if item > largest {
               largest = item;
           }
       }
       largest
   }

可以看到，除了类型不同，逻辑完全一样。泛型解决的就是这个问题。

泛型函数
============

使用泛型将上面的代码统一：

.. code-block:: rust

   fn largest<T: PartialOrd>(list: &[T]) -> &T {
       let mut largest = &list[0];
       for item in list {
           if item > largest {
               largest = item;
           }
       }
       largest
   }

   fn main() {
       let numbers = vec![10, 20, 5, 30, 15];
       let chars = vec!['a', 'z', 'b', 'm'];

       println!("最大数字: {}", largest(&numbers));
       println!("最大字符: {}", largest(&chars));
   }

其中 ``<T>`` 声明泛型参数，``PartialOrd`` 是 Trait 约束，表示 ``T`` 必须支持 ``>`` 比较。

泛型结构体
==============

单个泛型参数：

.. code-block:: rust

   struct Point<T> {
       x: T,
       y: T,
   }

   fn main() {
       let integer_point = Point { x: 5, y: 10 };
       let float_point = Point { x: 1.2, y: 3.4 };
   }

多个泛型参数：

.. code-block:: rust

   struct Point<T, U> {
       x: T,
       y: U,
   }

   fn main() {
       let p = Point { x: 5, y: 3.14 };
       let q = Point { x: "hello", y: 'c' };
   }

泛型枚举
============

Rust 标准库中最经典的泛型枚举：

.. code-block:: rust

   enum Option<T> {
       Some(T),
       None,
   }

   enum Result<T, E> {
       Ok(T),
       Err(E),
   }

使用：

.. code-block:: rust

   fn main() {
       let some_number: Option<i32> = Some(42);
       let some_string: Option<String> = Some(String::from("hello"));

       let ok_result: Result<i32, String> = Ok(100);
       let err_result: Result<i32, String> = Err(String::from("出错了"));
   }

泛型方法
============

在 ``impl`` 块中为泛型结构体实现方法：

.. code-block:: rust

   struct Point<T> {
       x: T,
       y: T,
   }

   impl<T> Point<T> {
       fn x(&self) -> &T {
           &self.x
       }

       fn y(&self) -> &T {
           &self.y
       }
   }

   fn main() {
       let p = Point { x: 5, y: 10 };
       println!("p.x = {}", p.x());
   }

为特定类型实现方法
======================

可以对泛型参数的具体类型单独实现方法：

.. code-block:: rust

   struct Point<T> {
       x: T,
       y: T,
   }

   // 只对 f64 类型实现 distance_from_origin 方法
   impl Point<f64> {
       fn distance_from_origin(&self) -> f64 {
           (self.x.powi(2) + self.y.powi(2)).sqrt()
       }
   }

   fn main() {
       let p = Point { x: 3.0, y: 4.0 };
       println!("距离原点: {}", p.distance_from_origin());

       let q = Point { x: 1, y: 2 };
       // q.distance_from_origin(); // 编译错误：i32 没有此方法
   }

Trait 约束（Trait Bound）
=============================

泛型不是万能的——需要对类型行为进行约束。

基本语法：

.. code-block:: rust

   // 方式一：在泛型声明后
   fn notify<T: Summary>(item: &T) {
       println!("{}", item.summarize());
   }

   // 方式二：where 子句（推荐，清晰）
   fn notify<T>(item: &T)
   where
       T: Summary,
   {
       println!("{}", item.summarize());
   }

多个 Trait 约束：

.. code-block:: rust

   fn process<T>(item: &T)
   where
       T: Summary + Display,
   {
       println!("{}", item);
       println!("{}", item.summarize());
   }

通过 where 简化复杂约束
============================

当约束很多时，``where`` 子句比内联写法清晰得多：

.. code-block:: rust

   use std::fmt::Display;
   use std::ops::Add;

   fn complex_function<T, U>(a: T, b: U) -> T
   where
       T: Add<U, Output = T> + Display + Clone,
       U: Display,
   {
       let result = a + b;
       println!("{} + {} = {}", a, b, result);
       result
   }

   fn main() {
       complex_function(10, 20);
   }

对比内联写法（不推荐）：

.. code-block:: rust

   fn complex_function<T: Add<U, Output = T> + Display + Clone, U: Display>(a: T, b: U) -> T {
       // ...
   }

泛型 + Trait 综合示例
==========================

定义一个 ``Summary`` Trait，然后用泛型约束使用它：

.. code-block:: rust

   pub trait Summary {
       fn summarize(&self) -> String;
   }

   struct Article {
       pub title: String,
       pub content: String,
   }

   impl Summary for Article {
       fn summarize(&self) -> String {
           format!("《{}》: {}", self.title, &self.content[..20.min(self.content.len())])
       }
   }

   struct Tweet {
       pub username: String,
       pub text: String,
   }

   impl Summary for Tweet {
       fn summarize(&self) -> String {
           format!("@{}: {}", self.username, &self.text[..15.min(self.text.len())])
       }
   }

   // 泛型函数：接受任何实现了 Summary 的类型
   fn notify<T: Summary>(item: &T) {
       println!("通知: {}", item.summarize());
   }

   fn main() {
       let article = Article {
           title: String::from("Rust 泛型详解"),
           content: String::from("泛型是 Rust 中实现代码复用的核心机制..."),
       };

       let tweet = Tweet {
           username: String::from("rustacean"),
           text: String::from("Rust 的泛型零成本抽象太棒了！"),
       };

       notify(&article);
       notify(&tweet);
   }

泛型的性能：单态化（Monomorphization）
===========================================

Rust 泛型在编译时展开为具体类型，运行时零开销：

.. code-block:: rust

   fn add<T: Add<Output = T>>(a: T, b: T) -> T {
       a + b
   }

   fn main() {
       let a = add(1, 2);     // 编译器生成 add_i32
       let b = add(1.5, 2.5); // 编译器生成 add_f64
   }

编译器实际生成的代码：

.. code-block:: rust

   fn add_i32(a: i32, b: i32) -> i32 {
       a + b
   }

   fn add_f64(a: f64, b: f64) -> f64 {
       a + b
   }

这就是**零成本抽象（Zero-Cost Abstraction）**：

.. list-table:: 泛型 vs 动态派发对比
   :header-rows: 1
   :widths: 30 35 35

   * - 特性
     - 静态派发（泛型/单态化）
     - 动态派发（dyn Trait）
   * - 运行时开销
     - 无（零成本抽象）
     - 有（虚函数表查找）
   * - 二进制大小
     - 可能增大（多份代码）
     - 较小（一份代码）
   * - 灵活性
     - 编译期确定类型
     - 运行时多态
   * - 典型用法
     - ``<T: Trait>``
     - ``&dyn Trait`` / ``Box<dyn Trait>``

生命周期与泛型
====================

泛型经常与生命周期参数一起使用：

.. code-block:: rust

   struct Excerpt<'a, T> {
       content: &'a T,
       start: usize,
       end: usize,
   }

   impl<'a, T> Excerpt<'a, T>
   where
       T: Display,
   {
       fn display(&self) {
           println!("{}", self.content);
       }
   }

泛型与关联类型
====================

关联类型是与 Trait 绑定的类型占位符，常与泛型配合：

.. code-block:: rust

   trait Container {
       type Item;
       fn get(&self, index: usize) -> Option<&Self::Item>;
   }

   impl<T> Container for Vec<T> {
       type Item = T;
       fn get(&self, index: usize) -> Option<&T> {
           self.as_slice().get(index)
       }
   }

关联类型 vs 泛型的区别：

.. list-table:: 关联类型 vs 泛型参数
   :header-rows: 1
   :widths: 30 35 35

   * - 特性
     - 关联类型
     - 泛型参数
   * - 声明位置
     - trait 内部 ``type Item;``
     - trait 名后 ``trait Foo<T>``
   * - 每个实现
     - 只能指定一次
     - 可多次实现不同 T
   * - 使用方式
     - ``T::Item``
     - ``T`` 直接作为类型
   * - 典型例子
     - ``Iterator::Item``
     - ``From<T>`` / ``Into<T>``

泛型与常量泛型（Const Generics）
=====================================

通过常量值参数化类型（Rust 1.51+ 稳定）：

.. code-block:: rust

   struct Array<T, const N: usize> {
       data: [T; N],
   }

   impl<T, const N: usize> Array<T, N> {
       fn new(data: [T; N]) -> Self {
           Array { data }
       }

       fn len(&self) -> usize {
           N
       }
   }

   fn main() {
       let arr1: Array<i32, 3> = Array::new([1, 2, 3]);
       let arr2: Array<i32, 5> = Array::new([1, 2, 3, 4, 5]);

       println!("arr1 长度: {}", arr1.len()); // 3
       println!("arr2 长度: {}", arr2.len()); // 5
   }

常见应用场景：

.. code-block:: rust

   // 标准库中的数组
   let a: [i32; 3] = [1, 2, 3];

   // 实际上 [T; N] 就是常量泛型

   // 自定义固定大小的向量运算
   fn dot_product<const N: usize>(a: [f64; N], b: [f64; N]) -> f64 {
       a.iter().zip(b.iter()).map(|(x, y)| x * y).sum()
   }

泛型参数的默认类型
========================

可以为泛型参数指定默认类型：

.. code-block:: rust

   struct Container<T = i32> {
       value: T,
   }

   fn main() {
       let c1 = Container { value: 42 };        // T 默认为 i32
       let c2 = Container { value: "hello" };    // T 推断为 &str
       let c3: Container<f64> = Container { value: 3.14 }; // 显式指定
   }

标准库中的典型例子——``Add`` Trait：

.. code-block:: rust

   // std::ops::Add 的定义（简化）
   pub trait Add<Rhs = Self> {
       type Output;
       fn add(self, rhs: Rhs) -> Self::Output;
   }

``Rhs = Self`` 表示默认右操作数类型与自身相同，所以 ``1 + 2`` 中两个操作数都是 ``i32``。

常见泛型模式
==================

.. list-table:: Rust 常见泛型模式
   :header-rows: 1
   :widths: 30 50 20

   * - 模式
     - 说明
     - 示例
   * - 类型占位符
     - 用 ``T``, ``U`` 等表示任意类型
     - ``struct Point<T>``
   * - Trait Bound
     - 约束泛型必须实现某 Trait
     - ``fn foo<T: Display>(x: T)``
   * - 多重约束
     - 要求同时满足多个 Trait
     - ``T: Clone + Debug``
   * - 泛型返回值
     - 返回类型也是泛型
     - ``fn parse<T: FromStr>(s: &str) -> Result<T, T::Err>``
   * - 泛型 + 生命周期
     - 同时约束类型和引用有效期
     - ``fn longest<'a, T>(x: &'a T, y: &'a T) -> &'a T``
   * - 泛型 impl 块
     - 为泛型结构体实现方法
     - ``impl<T> Point<T> { ... }``
   * - 条件实现
     - 仅在满足约束时实现
     - ``impl<T: Display> Point<T> { ... }``
   * - 常量泛型
     - 用编译期常量参数化
     - ``fn foo<const N: usize>(arr: [i32; N])``

条件实现（Conditional Implementation）
===========================================

仅在泛型参数满足某些条件时才实现 Trait：

.. code-block:: rust

   use std::fmt::Display;

   struct Pair<T> {
       x: T,
       y: T,
   }

   // 总是可以 new
   impl<T> Pair<T> {
       fn new(x: T, y: T) -> Self {
           Pair { x, y }
       }
   }

   // 只有当 T 实现了 Display + PartialOrd 时才有 cmp_display
   impl<T: Display + PartialOrd> Pair<T> {
       fn cmp_display(&self) {
           if self.x >= self.y {
               println!("较大的是 x = {}", self.x);
           } else {
               println!("较大的是 y = {}", self.y);
           }
       }
   }

   fn main() {
       let p1 = Pair::new(10, 20);
       p1.cmp_display(); // OK

       let p2 = Pair::new(vec![1], vec![2]);
       // p2.cmp_display(); // 编译错误：Vec 没有实现 Display
   }

标准库中的经典条件实现——Blanket Implementation：

.. code-block:: rust

   // 任何实现了 Display 的类型，自动实现 ToString
   impl<T: Display> ToString for T {
       // ...
   }

泛型与高级 Trait Bound
============================

.. code-block:: rust

   use std::fmt::Debug;
   use std::ops::Add;

   // 要求 T 能相加，相加结果也是 T，且 T 可以 Debug
   fn sum_and_debug<T>(a: T, b: T) -> T
   where
       T: Add<Output = T> + Debug,
   {
       let result = a + b;
       println!("{:?} + {:?} = {:?}", a, b, result);
       result
   }

   fn main() {
       sum_and_debug(10, 20);
       sum_and_debug(1.5, 2.5);
   }

更复杂的例子——返回实现了某 Trait 的类型：

.. code-block:: rust

   // impl Trait 语法：返回某种实现了 Summary 的类型
   fn make_summarizable() -> impl Summary {
       Article {
           title: String::from("默认文章"),
           content: String::from("这是一篇默认文章的内容..."),
       }
   }

泛型最佳实践
==================

.. list-table:: 泛型使用建议
   :header-rows: 1
   :widths: 30 70

   * - 原则
     - 说明
   * - 从具体到抽象
     - 先写具体类型代码，跑通后再抽象为泛型
   * - Trait Bound 最小化
     - 只约束真正需要的 Trait，不过度约束
   * - 优先用 where 子句
     - 当约束超过一个时，where 比内联更清晰
   * - 利用编译器推断
     - 不需要每次都显式标注类型，让编译器推断
   * - 避免过早泛型化
     - 不要一开始就写成泛型，确认有复用需求后再抽象
   * - 考虑 impl Trait
     - 返回值类型复杂或不便命名时，用 ``-> impl Trait``

总结
=====

.. code-block:: text

   泛型
   │
   ├── 函数泛型：fn foo<T>(x: T)
   │
   ├── 结构体泛型：struct Point<T>
   │
   ├── 枚举泛型：enum Option<T>
   │
   ├── 方法泛型：impl<T> Point<T> { ... }
   │
   ├── Trait 约束：<T: Display>
   │
   ├── where 子句：where T: Display + Clone
   │
   ├── 常量泛型：<const N: usize>
   │
   ├── 关联类型：type Item;
   │
   ├── 条件实现：impl<T: Display> for ...
   │
   └── 单态化：编译期零成本展开

泛型是 Rust 类型系统的基石，与 Trait、生命周期共同构成了 Rust 强大的抽象能力，同时保持零运行时开销。
