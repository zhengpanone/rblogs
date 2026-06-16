=====================
Rust 基本语法
=====================

.. contents:: 目录
   :depth: 3
   :local:


数据类型(Date Types)
=====================

标量类型（Scalar Types）
-------------------------

代表单个值，Rust 有四类标量：

整数类型(Integer Types)
>>>>>>>>>>>>>>>>>>>>>>>>

.. list-table:: 整数类型表
   :widths: auto
   :header-rows: 1

   * - 长度
     - 有符号
     - 无符号
   * - 8 bit
     - i8
     - u8
   * - 16 bit
     - i16
     - u16
   * - 32 bit
     - i32
     - u32
   * - 64 bit
     - i64
     - u64
   * - 128 bit
     - i128
     - u128
   * - arch
     - isize
     - usize

**关键点**：

- 数字字面量可以加后缀： ``42u8``、 ``-5i32``
- 可以用下划线增强可读性：``1_000_000``
- 默认推断为 ``i32`` （除非上下文有其他要求）
- ``usize`` 常用于数组索引、集合容量（与平台指针宽度一致）

**溢出行为（Debug vs Release）**：

- Debug 模式下溢出会 **panic**
- Release 模式下默认 **环绕（wrapping）**，但可用显式方法控制：

.. code-block:: rust

  let x: u8 = 200;
  let y = x.wrapping_add(100);   // 44
  let z = x.saturating_add(100); // 255（卡在最大值）


浮点类型(Floating-Point Types)
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

- f32（单精度，约 7 位有效数字）
- f64（双精度，约 15 位有效数字，默认类型）

.. code-block:: rust

  let x = 2.0;      // f64
  let y: f32 = 3.0; // f32

字符类型(Character Types)
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

  let c = 'z';
  let heart_eyed_cat = '😻'; // Unicode，4 字节

注意：char是 **4 字节**，不是 ASCII。字符串中的字符是 UTF-8 编码的 u8序列。

布尔类型(Boolean Types)
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

大小：**1 字节**

.. code-block:: rust

  let t = true;
  let f: bool = false;

复合类型（Compound Types）
---------------------------------

将多个值组合成一个类型。

元组（Tuple）
>>>>>>>>>>>>>>>>>>>>>>>

固定长度，元素类型可以不同。

.. code-block:: rust

  fn main() {
    let person = ("Tom", 18);

    let (name, age) = person; // 解构(Destructuring)

    println!("{}", person.0);
    println!("{}", person.1); // 索引访问（从 0 开始）
  }

- 最大 12 个元素（标准库实现了 trait，超过则不能自动 derive）
- 空元组 ()称为 **unit 类型**，表示「没有值」，函数默认返回它
  

数组（Array）
>>>>>>>>>>>>>>>>>>>>>>>

固定长度，元素类型必须相同，**栈上分配**。

.. code-block:: rust

  let nums = [1, 2, 3, 4, 5]; // 类型推断为 [i32; 5]
  let arr: [i32; 5] = [0; 5]; // [0; 5] 表示 5 个元素都是 0

  println!("{}", nums[0]); // 索引访问

越界检查：​ 运行时 ``panic``，不会出现 C/C++ 那种缓冲区溢出。
  
  动态大小用 ``Vec<T>``。

特殊类型（Special Types）
-------------------------

newtype模式
>>>>>>>>>>>>>>>>>>>>>>

newtype 模式：用元组结构体包装现有类型，获得类型安全。

.. code-block:: rust

  // 没有类型安全
  fn login(username: String, password: String) {}
  fn register(username: String, password: String) {}

  fn main() {
      let user = String::from("larry");
      let pass = String::from("123456");
      
      // 参数顺序错了，但编译器不报错！
      login(pass, user);  // 密码当用户名传了...
  }

用newtype解决：

.. code-block:: rust

  // ✅ 类型安全
  struct Username(String);
  struct Password(String);

  fn login(username: Username, password: Password) {}
  fn register(username: Username, password: Password) {}

  fn main() {
      let user = Username(String::from("larry"));
      let pass = Password(String::from("123456"));
      
      // ❌ 现在顺序错了会编译错误！
      // login(pass, user);  // 类型不匹配！
      
      login(user, pass);  // ✅ 正确
  }

newtype 的优势：

  1. 类型安全：编译器帮你检查
  2. 零开销：编译后和直接用一个样
  #. 可以添加方法：给包装类型加专属方法
  #. 隐藏实现细节：外面不知道里面是啥

类型别名
>>>>>>>>>>>>>>>>>

类型别名：给现有类型起个"小名"，不是新类型。

.. code-block:: rust

  // 长名字
  type Result<T> = std::result::Result<T, std::io::Error>;

  // 用起来
  fn read_file() -> Result<String> {
      // ...
  }

类型别名 vs newtype：

.. list-table:: 特性对比表
   :widths: 20 30 30
   :header-rows: 1

   * - 特性
     - 类型别名
     - newtype
   * - 类型检查
     - 和原类型一样
     - 独立类型
   * - 开销
     - 零
     - 零
   * - 添加方法
     - 不能
     - 能
   * - 用途
     - 简化名字
     - 类型安全

.. code-block:: rust

  // 类型别名
  type Meters = f64;
  type Kilometers = f64;

  fn main() {
      let m: Meters = 1000.0;
      let k: Kilometers = 1.0;
      
      // 可以相加，但可能没意义
      let sum = m + k;  // 编译器不管
  }

  // newtype
  struct Meters(f64);
  struct Kilometers(f64);

  fn main() {
      let m = Meters(1000.0);
      let k = Kilometers(1.0);
      
      // 不能相加，类型不同
      // let sum = m + k;  // 编译错误！
  }

Never 类型 （Never Type）
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

never 类型： ``!`` ，表示"永远不会返回"。

用于永远不会返回的函数:

1. ``panic!``  程序崩溃
2. ``loop {}``  无限循环（没有 break）
3. ``std::process::exit()`` - 退出程序
4. ``continue``

.. code-block:: rust

  fn crash() -> ! {
      panic!("boom!");
  }

  fn forever() -> ! {
    loop {}
  }

编译器用它做类型统一，比如 match的分支：

.. code-block:: rust

  // let guess = match some_value {
  //  Ok(val) => val,
  //  Err(_) => continue,  // continue 的类型是 !
  // };

  fn main() {
    let numbers = vec!["10", "20", "abc", "30"];

    for s in numbers {
        let guess: i32 = match s.parse::<i32>() {
            Ok(val) => val,
            Err(_) => {
                println!("跳过非法输入: {}", s);
                continue; // 类型为 !
            }
        };

        println!("成功解析的数字: {}", guess);
    }

    println!("循环结束");
  }

never 类型的特点：

- 可以转换成任何类型
- 用在"不会执行到"的地方。
- 不能创建 ! 类型的值
- 只能作为返回值


``!`` 永远不会真正返回值，所以它可以假装是任何类型。

Sized特质
>>>>>>>>>>>>>>>>>>>

Sized 特质：一个类型在编译期是否拥有已知、固定的内存大小。

.. code-block:: rust

  pub trait Sized { }

- 没有方法
- 由编译器自动实现
- 表示：该类型的值可以在栈上分配

规则：

- 大多数类型都是 ``Sized`` 的（ ``i32``、 ``String``、结构体等）
- 有些类型不是 ``Sized`` 的（ ``str``、 ``[T]``、 ``dyn Trait``）

为什么需要 Sized

Rust 必须知道一个值占多少内存，才能：

- 在栈上分配空间
- 计算结构体字段偏移
- 做指针运算
- 生成正确的机器码
- 所以 几乎所有泛型都隐式要求 ``Sized``。

.. code-block:: rust

  // 不能这样
  let s: str = "hello";  // str 不是 Sized

  // 要这样
  let s: &str = "hello";  // &str 是 Sized（引用有固定大小）

Sized 作为泛型约束：

.. code-block:: rust

  // 泛型参数默认是 Sized 的
  fn generic<T>(x: T) {}  // 相当于 fn generic<T: Sized>(x: T) {}

  // 如果不是 Sized，要显式标注
  fn not_sized<T: ?Sized>(x: &T) {}  // ?Sized 表示可以是也可以不是

``?Sized`` 的含义：

- ``T: Sized`` → T 必须是 Sized 的
- ``T: ?Sized`` → T 可以是 Sized 的，也可以不是

.. code-block:: rust

  // 可以接受 str 和 String
  fn print_it<T: ?Sized>(x: &T) 
  where
      T: std::fmt::Display 
  {
      println!("{}", x);
  }

  fn main() {
      print_it("hello");      // &str
      print_it(&String::from("hello"));  // &String
  }

示例
>>>>>>>>>>>>>

newtype 模式基础

.. code-block:: rust

  // 包装类型
  struct Meters(f64);
  struct Seconds(f64);

  impl Meters {
      fn new(value: f64) -> Self {
          Meters(value)
      }
      
      fn to_kilometers(&self) -> Kilometers {
          Kilometers(self.0 / 1000.0)
      }
  }

  impl Seconds {
      fn new(value: f64) -> Self {
          Seconds(value)
      }
      
      fn to_minutes(&self) -> Minutes {
          Minutes(self.0 / 60.0)
      }
  }

  struct Kilometers(f64);
  struct Minutes(f64);

  fn main() {
      let distance = Meters::new(1500.0);
      let time = Seconds::new(120.0);
      
      println!("距离：{:?} 米", distance.0);
      println!("距离：{:?} 公里", distance.to_kilometers().0);
      println!("时间：{:?} 秒", time.0);
      println!("时间：{:?} 分钟", time.to_minutes().0);
      
      // 不能混用
      // let speed = distance.0 / time.0;  // 类型不匹配
  }

newtype 实现 Trait

.. code-block:: rust

  use std::fmt;

  struct Wrapper(Vec<String>);

  impl fmt::Display for Wrapper {
      fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
          write!(f, "[{}]", self.0.join(", "))
      }
  }

  fn main() {
      let w = Wrapper(vec![
          String::from("hello"),
          String::from("world"),
      ]);
      
      println!("w = {}", w);  // 输出：w = [hello, world]
  }

类型别名实战

.. code-block:: rust

  // 简化复杂类型
  type Thunk = Box<dyn Fn() + Send + 'static>;

  fn main() {
      let f: Thunk = Box::new(|| {
          println!("hello");
      });
      
      f();
  }

  // 带错误的 Result
  type IoResult<T> = Result<T, std::io::Error>;

  fn read_file(path: &str) -> IoResult<String> {
      std::fs::read_to_string(path)
  }

never 类型应用

.. code-block:: rust

  fn main() {
      // match 表达式中的 never 类型
      let result: Result<i32, &str> = Ok(5);
      
      let value = match result {
          Ok(v) => v,
          Err(_) => panic!("出错了！"),  // ! 可以当 i32 用
      };
      
      println!("值：{}", value);
      
      // 无限循环
      let mut count = 0;
      let x: i32 = loop {
          count += 1;
          if count == 10 {
              break count * 2;  // 返回 i32
          }
      };
      
      println!("x = {}", x);
  }

?Sized 约束

.. code-block:: rust

  use std::fmt::Display;

  // 可以接受 Sized 和 ?Sized 类型
  fn print_it<T: ?Sized>(x: &T) 
  where
      T: Display 
  {
      println!("{}", x);
  }

  // 只能接受 Sized 类型
  fn print_sized<T: Display>(x: T) {
      println!("{}", x);
  }

  fn main() {
      // &str 不是 Sized，但 &(&str) 是
      print_it("hello");  // 
      // print_sized("hello");  //  "hello" 是 &str，是 Sized 的
      
      // dyn Trait 不是 Sized
      let boxed: Box<dyn Display> = Box::new(42);
      print_it(&*boxed);  

      // print_sized(*boxed);  // 不对
  }

newtype 用于 API 设计

.. code-block:: rust

  // 用 newtype 防止 API 误用
  struct Email(String);
  struct Password(String);

  impl Email {
      fn new(email: &str) -> Result<Self, &'static str> {
          if email.contains('@') {
              Ok(Email(email.to_string()))
          } else {
              Err("无效的邮箱地址")
          }
      }
      
      fn as_str(&self) -> &str {
          &self.0
      }
  }

  impl Password {
      fn new(password: &str) -> Result<Self, &'static str> {
          if password.len() >= 8 {
              Ok(Password(password.to_string()))
          } else {
              Err("密码至少 8 位")
          }
      }
  }

  struct User {
      email: Email,
      password: Password,
  }

  fn main() {
      let email = Email::new("larry@example.com").unwrap();
      let password = Password::new("secure123").unwrap();
      
      let user = User {
          email,
          password,
      };
      
      println!("用户邮箱：{}", user.email.as_str());
      
      // 不能直接用字符串
      // let bad_user = User {
      //     email: "bad@email",  // 类型不匹配
      //     password: "123",     // 类型不匹配
      // };
  }

配置类型安全

.. code-block:: rust

  // 用 newtype 防止配置错误
  struct Port(u16);
  struct Host(String);
  struct Timeout(std::time::Duration);

  impl Port {
      fn new(port: u16) -> Result<Self, &'static str> {
          if port == 0 || port > 65535 {
              Err("无效端口")
          } else {
              Ok(Port(port))
          }
      }
  }

  struct DatabaseConfig {
      host: Host,
      port: Port,
      timeout: Timeout,
  }

  fn connect(config: &DatabaseConfig) {
      println!(
          "连接到 {}:{}，超时 {:?}",
          config.host.0,
          config.port.0,
          config.timeout.0
      );
  }

  fn main() {
      let config = DatabaseConfig {
          host: Host(String::from("localhost")),
          port: Port::new(5432).unwrap(),
          timeout: Timeout(std::time::Duration::from_secs(30)),
      };
      
      connect(&config);
      
      // 类型安全防止错误
      // let bad_config = DatabaseConfig {
      //     host: Host(String::from("localhost")),
      //     port: Port(99999),  // 超出范围，但构造函数会检查
      //     timeout: Timeout(std::time::Duration::from_secs(30)),
      // };
  }

类型转换（Type Casting）
-------------------------

Rust 非常严格，不会隐式转换，必须显式：

.. code-block:: rust

  let x = 5i32;
  let y = 3.0f64;

  // 错误：cannot add `f64` to `i32`
  // let sum = x + y;

  // 正确：显式转换
  let sum = x as f64 + y;       // as 转换
  let sum = (x as f64) + y;     // 更推荐加括号

as转换可能丢失精度（如 f64 → i32），安全转换用 ``TryFrom``：

.. code-block:: rust

  let x: i64 = 100;
  let y: i32 = x.try_into().unwrap(); // 返回 Result

类型推断（Type Inference）
------------------------------

Rust 能在大部分场景推断类型，但边界模糊时需要显式标注：

.. code-block:: rust

  let v = Vec::new();            // 无法推断元素类型 
  let v: Vec<i32> = Vec::new(); // 显式 
  let mut v = Vec::new();
  v.push(1);                     // 根据 push 推断为 Vec<i32> 


整数字面量写法一览
-------------------------

.. list-table:: Rust 整数字面量进制表示
   :widths: auto
   :header-rows: 1

   * - 进制
     - 示例
     - 说明
   * - 十进制
     - ``98_222``
     - 下划线分隔
   * - 十六进制
     - ``0xff``
     - 前缀 ``0x``
   * - 八进制
     - ``0o77``
     - 前缀 ``0o``
   * - 二进制
     - ``0b1111_0000``
     - 前缀 ``0b``
   * - 字节（仅 u8）
     - ``b'A'``
     - 单引号，ASCII


- **能用 i32就用 i32** ——CPU 通常对 32 位整数优化最好，除非你明确需要特定宽度或大量节省内存
- **数值运算默认不会溢出 panic** ，但建议在可能溢出的地方用 checked_add、overflowing_add等显式方法
- **char不等于 ASCII**，处理文本时多用 &str/ String而非 char数组
- **类型别名** 可以让复杂类型更清晰：
  
  .. code-block:: rust

    type Kilometers = i32;
    let distance: Kilometers = 50;


变量（Variable）
===================

默认不可变（Immutable by Default）

- 一个引用类型完整由 **四维度** 描述

 - **是否引用**: 有没有 ``&`` 拥有 vs 借用
 - **类型**: 是什么 ``str`` ``i32``
 - **生命周期**: 引用指向的数据活多久 仅引用才有
 - **可变性**: ``mut`` 可写 vs 只读
  
- 没有 ``&`` 的拥有型值 不涉及生命周期 数据归它自己管

.. code-block:: text

  letx : &'static str = "hi";
    ↑ ↑    ↑   ↑
    │ │    │   └── 值
    │ │    └───────── 类型 (str = 字符串)
    │ └───────────────── 生命周期 ('static = 活多久)
    └──────────────────── 是不是引用 (& = 是借用)

- ``'static`` 就是唯一一个有专属名字和内置含义的生命周期，表示活整个程序。
  
  - 其余生命周期都是 **泛型占位符** 名字自己起 惯例 ``'a`` ``'b`` ``'c``
  - 类比: ``'static`` 之于生命周期 = ``i32`` 之于类型(具体的); ``'a`` = ``T`` (泛型的)
  - ``'a`` 不是固定时长 是 **关系标记** 把多个引用的寿命绑在一起 具体多久由调用时决定
  - ``'_`` 匿名生命周期 = 有但让编译器推; 简单情况可直接省略编译器按省略规则自动补
  - 关键: ``'a`` 写在 ``<>`` 里像泛型 但占位的是 **存活区间** 而非类型 且 **零运行时开销** 编译完擦除



不可变变量（Immutable Variable）
-----------------------------------

.. code-block:: rust

  fn main() {
    let x = 5;
    println!("{}", x);
  }

可变变量（Mutable Variable）
-----------------------------------

.. code-block:: rust

  fn main() {
    let mut x = 5;
    println!("{}", x);
    x = 6;
    println!("{}", x);
  }


变量声明（Variable Declaration）
-----------------------------------

告诉编译器「我要引入一个名字」，并指定其初始值和类型（类型通常可推断）。

.. code-block:: rust

  // 最基本形式
  let x = 5;

  // 带类型注解
  let x: i32 = 5;

  // 可变声明
  let mut y = 10;

  // 先声明后赋值（必须保证后续一定赋值才能使用）
  let z: i32;
  // ... 一些条件逻辑 ...
  z = 42; // 必须在使用前完成赋值

关键点

- 使用关键字 let
- 声明时必须决定是否可变（mut）
- 类型可以在声明时标注，也可以由编译器推断
- 声明的变量必须初始化后才能使用——Rust 没有「未初始化变量」的概念（不像 C）

.. code-block:: rust

  let x: i32;
  println!("{}", x); // 编译错误：使用了未初始化的变量


变量绑定（Variable Binding）
-----------------------------------

将值与变量名关联起来的过程。Rust 里 let x = 5不叫「赋值」，而叫「绑定」。

绑定的本质

- 绑定将值与名称关联
- 绑定涉及所有权（ownership）的建立或转移
- 绑定一旦建立，变量名就拥有了该值的所有权（除非是引用或 Copy 类型）

.. code-block:: rust

  // 所有权绑定
  let s1 = String::from("hello"); // s1 拥有这个 String
  let s2 = s1;                    // 所有权从 s1 移动到 s2，s1 不再有效

  // 引用绑定（不转移所有权）
  let s3 = String::from("world");
  let r = &s3;                    // r 绑定了 s3 的引用，s3 仍拥有所有权

  // Copy 类型的绑定（实际上是按位拷贝）
  let a = 42;
  let b = a;                      // a 仍然有效，因为 i32 实现了 Copy

.. list-table:: 绑定与赋值操作
   :widths: auto
   :header-rows: 1

   * - 操作
     - 含义
     - 示例
   * - 绑定
     - 初次将值与名称关联
     - ``let x = 5;``
   * - 重新绑定
     - 用新值替换旧绑定（shadowing）
     - ``let x = x + 1;``
   * - 赋值
     - 修改已绑定可变变量的值
     - ``x = 10;`` （需要 ``let mut x``）
  
.. code-block:: rust

  let mut x = 5;  // 绑定
  x = 10;         // 赋值（修改已有绑定）
  let x = x + 1;  // 重新绑定（shadowing，创建新绑定）

变量作用域（Variable Scope）
-----------------------------------

变量在程序中有效的范围，从声明点到包含它的最近一对花括号 ``{}`` 的结束位置。

.. code-block:: rust

  {                          // 外层作用域开始
    let outer = "outside";
    
    {                      // 内层作用域开始
        let inner = "inside";
        println!("{} {}", outer, inner); // 可以访问外层变量
    }                      // inner 在此处销毁
    
    println!("{}", outer); // outer 仍然有效
    // println!("{}", inner); // inner 已超出作用域
  }                          // outer 在此处销毁

作用域与生命周期
>>>>>>>>>>>>>>>>>>>>>>

对于堆分配类型（如 String），离开作用域时会自动调用 drop释放内存：

.. code-block:: rust

  {
      let s = String::from("hello"); // s 进入作用域
      // 使用 s
  }                                  // 作用域结束，s.drop() 被自动调用，内存释放

作用域与借用
>>>>>>>>>>>>>>>>>>>>>>

借用不能超出被借用变量的作用域：

.. code-block:: rust

  let r;
  {
      let x = 5;
      r = &x; // 编译错误：x 的生命周期不够长
  }
  println!("{}", r);

作用域与 Shadowing
>>>>>>>>>>>>>>>>>>>>>>

Shadowing 发生在同一作用域或嵌套作用域中：

.. code-block:: rust

  let x = 1;

  {
      let x = 2;        // 遮蔽外层的 x
      println!("{}", x); // 2
  }

  println!("{}", x);    // 1（内层遮蔽已结束）


变量遮蔽（Variable Shadowing）
-----------------------------------

同一个作用域内可以用同名变量覆盖之前的绑定：

.. code-block:: rust

  let x = 5;
  let x = x + 1;   // 新绑定，遮蔽旧的 x
  {
      let x = x * 2; // 内部作用域再次遮蔽
      println!("{}", x); // 12
  }
  println!("{}", x); // 6（内部作用域的遮蔽已结束）

与 mut的区别：

- mut是修改同一块内存的值
- shadowing 是创建新的绑定，旧变量被遮蔽但仍可能存在（所有权未转移的话）
- shadowing 可以改变类型，mut不行

.. code-block:: rust

  let spaces = "   ";
  let spaces = spaces.len(); //  类型从 &str 变成 usize

  let mut spaces = "   ";
  spaces = spaces.len(); // 不能改变类型

常量
=================

.. code-block:: rust

  const MAX_POINTS: u32 = 100_000;
  const THREE_HOURS_IN_SECONDS: u32 = 60 * 60 * 3; // 编译期常量表达式


.. list-table:: let 与 const 对比
   :widths: auto
   :header-rows: 1

   * - 特性
     - let
     - const
   * - 可变性
     - 默认不可变，可加 ``mut``
     - 永远不可变
   * - 类型注解
     - 通常可省略
     - 必须显式标注
   * - 求值时机
     - 运行时
     - 编译期
   * - 作用域
     - 块级作用域
     - 全局有效（任意作用域）
   * - 内存位置
     - 栈上（或堆上）
     - 每次使用时内联展开，无固定内存地址
