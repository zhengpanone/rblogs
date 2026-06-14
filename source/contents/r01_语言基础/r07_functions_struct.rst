=====================
函数和结构体
=====================

函数（Functions）
========================

基本定义
--------------

.. code-block:: rust

  fn function_name(param1: Type1, param2: Type2) -> ReturnType {
      // 函数体
      return_value // 最后一个表达式作为返回值（不加分号）
  }

完整示例：

.. code-block:: rust

  fn greet(name: &str) -> String {
      format!("Hello, {}!", name)
  }

  fn add(x: i32, y: i32) -> i32 {
      x + y // 不加分号，作为返回值
  }

  fn print_sum(a: i32, b: i32) {
      println!("Sum: {}", a + b); // 没有返回值，返回 ()
  }

语句 vs 表达式
---------------------

- 语句（Statement）：执行操作但不返回值，结尾有分号
- 表达式（Expression）：计算结果并返回值，结尾没有分号
  
.. code-block:: rust

  fn example() -> i32 {
      let y = 6;          // 语句
      let x = {           // 代码块也是表达式
          let z = 5;
          z + 1           // 这个表达式的值是 6
      };
      x + y               // 函数的返回值：12
  }

参数传递方式
---------------------

.. code-block:: rust

  // 1. 传值（所有权转移）
  fn take_ownership(s: String) {
      println!("{}", s);
  } // s 在这里被 drop

  // 2. 不可变引用（借用）
  fn borrow(s: &String) -> usize {
      s.len()
  } // 不获取所有权

  // 3. 可变引用
  fn modify(s: &mut String) {
      s.push_str(" world");
  }

泛型函数
---------------------

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

  // 使用
  let numbers = vec![34, 50, 25, 100, 65];
  let result = largest(&numbers);

高阶函数与闭包
---------------------

函数可以作为参数传递：

.. code-block:: rust

  fn apply_twice(f: fn(i32) -> i32, x: i32) -> i32 {
      f(f(x))
  }

  fn double(x: i32) -> i32 { x * 2 }

  let result = apply_twice(double, 5); // 20

闭包（匿名函数）：

.. code-block:: rust

  let add_one = |x: i32| -> i32 { x + 1 };
  let add_two = |x| x + 2; // 类型可推断

  let result = add_one(5); // 6

函数指针 vs 闭包
---------------------

.. code-block:: rust
  
  fn add_one(x: i32) -> i32 { x + 1 }

  // 函数指针
  let f: fn(i32) -> i32 = add_one;

  // 闭包（不捕获环境）
  let c: fn(i32) -> i32 = |x| x + 1;

  // 闭包（捕获环境）— 不能转为 fn 指针
  let y = 5;
  let c2 = |x| x + y; // 捕获了 y

结构体（Structs）
========================

结构体是 Rust 自定义数据类型的主要方式，有三种形式。

具名字段结构体（Named Fields Struct）
--------------------------------------
  
最常用的形式：

.. code-block:: rust

  struct User {
      active: bool,
      username: String,
      email: String,
      sign_in_count: u64,
  }

  // 创建实例
  let user1 = User {
      email: String::from("someone@example.com"),
      username: String::from("someusername123"),
      active: true,
      sign_in_count: 1,
  };

  // 修改（整个实例必须可变）
  let mut user2 = User {
      email: String::from("another@example.com"),
      username: String::from("anotheruser567"),
      active: true,
      sign_in_count: 1,
  };
  user2.email = String::from("anotheremail@example.com");

  // 结构体更新语法（从其他实例继承字段）
  let user3 = User {
      email: String::from("third@example.com"),
      ..user1 // 其余字段从 user1 复制
  };
  // 注意：username 所有权从 user1 移到了 user3，user1 不能再使用

元组结构体（Tuple Structs）
--------------------------------------
   
字段没有名字，通过索引访问：

.. code-block:: rust

  struct Color(i32, i32, i32);
  struct Point(i32, i32, i32);

  let black = Color(0, 0, 0);
  let origin = Point(0, 0, 0);

  // 访问
  println!("Red: {}", black.0);

  // 尽管字段类型相同，Color 和 Point 是不同的类型
  // let p: Point = black; // 类型不匹配


单元结构体（Unit-Like Structs）
--------------------------------------
   
没有字段，类似于 ()类型，常用于标记或实现 trait：

.. code-block:: rust

  struct AlwaysEqual;

  let subject = AlwaysEqual;

  // 常用来实现 trait
  impl SomeTrait for AlwaysEqual {
      // ...
  }

方法（Methods）
--------------------------------------

在 impl块中定义，第一个参数总是 self：

.. code-block:: rust

  struct Rectangle {
      width: u32,
      height: u32,
  }

  impl Rectangle {
      // 不可变引用
      fn area(&self) -> u32 {
          self.width * self.height
      }

      // 可变引用
      fn set_width(&mut self, width: u32) {
          self.width = width;
      }

      // 获取所有权（很少见）
      fn consume(self) -> u32 {
          self.width * self.height
      }

      // 关联函数（不是方法，没有 self）
      fn square(size: u32) -> Self {
          Self {
              width: size,
              height: size,
          }
      }
  }

  let rect = Rectangle { width: 30, height: 50 };
  println!("Area: {}", rect.area());

  let square = Rectangle::square(10); // 关联函数用 :: 调用

多个 impl 块
--------------------------------------

一个结构体可以有多个 impl块：

.. code-block:: rust

  impl Rectangle {
      fn can_hold(&self, other: &Rectangle) -> bool {
          self.width > other.width && self.height > other.height
      }
  }

  impl Rectangle {
      fn is_square(&self) -> bool {
          self.width == self.height
      }
  }

派生 Trait（Derive）
--------------------------------------

自动实现一些常用 trait：

.. code-block:: rust

  #[derive(Debug, Clone, Copy, PartialEq, Eq)]
  struct Point {
      x: i32,
      y: i32,
  }

  let p1 = Point { x: 1, y: 2 };
  let p2 = p1; // Copy 允许这样
  println!("{:?}", p1); // Debug 允许打印
  println!("{}", p1 == p2); // PartialEq 允许比较

常用 derive trait：

- Debug：格式化输出 {:?}
- Clone：显式克隆 .clone()
- Copy：按位复制（要求所有字段都实现 Copy）
- PartialEq/ Eq：相等比较 ==
- PartialOrd/ Ord：排序比较 <, >
- Hash：哈希计算

所有权与结构体
--------------------------------------

.. code-block:: rust

  struct User {
      username: String, // 拥有所有权的字段
      email: String,
  }

  // 错误：使用字符串切片会导致生命周期问题
  // struct User<'a> {
  //     username: &'a str, // 需要生命周期标注
  //     email: &'a str,
  // }


如果结构体包含引用，必须标注生命周期：

.. code-block:: rust

  struct User<'a> {
      username: &'a str,
      email: &'a str,
  }

  let name = String::from("Alice");
  let email = String::from("alice@example.com");
  let user = User {
      username: &name,
      email: &email,
  };


实用建议
========================

- 函数命名用 snake_case，结构体用 PascalCase
- 优先用 &self而不是 &mut self，除非必须修改
- 结构体字段尽量用拥有所有权的类型（如 String而非 &str），避免生命周期困扰
- 善用 #[derive(Debug)]​ 方便调试打印
- 结构体更新语法 ..other会移动所有权，小心原变量后续使用
- 方法名可以和字段名相同——Rust 会自动区分
