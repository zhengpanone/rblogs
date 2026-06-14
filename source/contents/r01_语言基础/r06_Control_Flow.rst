============
控制流
============

if / else if / else
=============================

if是一个表达式（expression），可以返回值：

.. code-block:: rust

  let number = 6;

  if number % 4 == 0 {
      println!("number is divisible by 4");
  } else if number % 3 == 0 {
      println!("number is divisible by 3");
  } else if number % 2 == 0 {
      println!("number is divisible by 2");
  } else {
      println!("number is not divisible by 4, 3, or 2");
  }

if 作为表达式使用
-----------------------------

.. code-block:: rust

  let condition = true;
  let number = if condition { 5 } else { 6 }; // 类型必须一致

  println!("The value of number is: {}", number); // 5

重要限制：​ 条件必须是 bool类型，不会像 C 那样自动把非零值当作 true。

.. code-block:: rust

  let x = 1;
  if x { /* 编译错误：expected bool, found integer */ }

循环（Loops）
=============================

Rust 提供了三种循环结构：loop、while 和 for。

loop —— 无限循环
-----------------------------

.. code-block:: rust

  loop {
      println!("again!");
  }

可以用 break退出，break也可以返回值：

.. code-block:: rust

  let mut counter = 0;

  let result = loop {
      counter += 1;

      if counter == 10 {
          break counter * 2; // 退出循环并返回 20
      }
  };

  println!("The result is {}", result); // 20

while —— 条件循环
-----------------------------

.. code-block:: rust

  let mut number = 3;

  while number != 0 {
      println!("{}!", number);
      number -= 1;
  }

  println!("LIFTOFF!!!");


for —— 遍历迭代器（最常用）
-----------------------------

.. code-block:: rust

  let arr = [10, 20, 30, 40, 50];

  for element in arr {
      println!("the value is: {}", element);
  }

Range 语法（左闭右开）
>>>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

  for number in 1..4 {
      println!("{}!", number); // 1, 2, 3
  }

  for number in (1..4).rev() {
      println!("{}!", number); // 3, 2, 1
  }

包含右端点的 Range
>>>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

  for number in 1..=4 {
      println!("{}!", number); // 1, 2, 3, 4
  }

**推荐优先用 for：比 while+ 索引更安全（不会越界），而且性能通常一样好。**

循环标签（Loop Labels）
=============================

嵌套循环时，可以用标签指定 break或 continue作用于哪个循环：

.. code-block:: rust

  'outer: for i in 1..=3 {
      'inner: for j in 1..=3 {
          if i == 2 && j == 2 {
              break 'outer; // 跳出外层循环
          }
          println!("i={}, j={}", i, j);
      }
  }
  // 输出：
  // i=1, j=1
  // i=1, j=2
  // i=1, j=3
  // i=2, j=1

标签名以单引号开头（类似生命周期语法），命名惯例是 snake_case。

match —— 模式匹配
=============================

match是 Rust 中最强大的控制流结构，必须穷尽所有可能性（exhaustive）：

.. code-block:: rust

  enum Coin {
      Penny,
      Nickel,
      Dime,
      Quarter,
  }

  fn value_in_cents(coin: Coin) -> u8 {
      match coin {
          Coin::Penny => 1,
          Coin::Nickel => 5,
          Coin::Dime => 10,
          Coin::Quarter => 25,
      }
  }

绑定匹配的值
--------------------

.. code-block:: rust

  #[derive(Debug)]
  enum UsState {
      Alabama,
      Alaska,
      // ...
  }

  enum Coin {
      Penny,
      Nickel,
      Dime,
      Quarter(UsState),
  }

  fn value_in_cents(coin: Coin) -> u8 {
      match coin {
          Coin::Penny => 1,
          Coin::Nickel => 5,
          Coin::Dime => 10,
          Coin::Quarter(state) => {
              println!("State quarter from {:?}!", state);
              25
          }
      }
  }

通配符 _和 other
--------------------

.. code-block:: rust

  let dice_roll = 9;
  match dice_roll {
      3 => add_fancy_hat(),
      7 => remove_fancy_hat(),
      other => move_player(other), // 捕获剩余所有值
      // _ => reroll(),            // 忽略剩余所有值
      // _ => (),                  // 忽略且什么都不做
  }

if let —— 简洁的模式匹配
=============================

当你只关心一种模式而忽略其他情况时，if let比 match更简洁：

.. code-block:: rust

  let config_max = Some(3u8);

  // match 写法
  match config_max {
      Some(max) => println!("The maximum is configured to be {}", max),
      _ => (),
  }

  // if let 写法（等价）
  if let Some(max) = config_max {
      println!("The maximum is configured to be {}", max);
  }

可以加 else：

.. code-block:: rust

  let mut stack = Vec::new();
  stack.push(1);
  stack.push(2);
  stack.push(3);

  while let Some(top) = stack.pop() {
      println!("{}", top);
  }


while let —— 条件模式匹配循环
=============================

.. code-block:: rust

  let mut optional = Some(0);

  while let Some(i) = optional {
      if i > 9 {
          println!("Greater than 9, quit!");
          optional = None;
      } else {
          println!("`i` is `{:?}`. Try again.", i);
          optional = Some(i + 1);
      }
  }

let else —— 解构失败时提前返回（Rust 1.65+）
===============================================

.. code-block:: rust

  fn get_first_item(s: &str) -> Option<&str> {
    let Some(first) = s.split(',').next() else {
        return None; // 如果 pattern 不匹配，执行 else 分支
    };
    Some(first.trim())
  }

等价于：

.. code-block:: rust

  fn get_first_item(s: &str) -> Option<&str> {
      let first = match s.split(',').next() {
          Some(v) => v,
          None => return None,
      };
      Some(first.trim())
  }

return —— 提前返回
=============================

.. code-block:: rust

  fn early_return(x: i32) -> i32 {
    if x < 0 {
        return 0; // 提前返回
    }
    x + 1 // 最后一个表达式作为返回值（不加分号）
  }


实用建议
=============================

- 优先用 for而不是 while+ 索引——更安全、更符合习惯
- match必须穷举——这迫使你考虑所有边界情况，减少运行时 bug
- if let用于单分支匹配——代码更紧凑
- let else用于「不匹配就提前返回」——减少嵌套层级
- 循环标签只在多层嵌套时使用——滥用反而降低可读性
- loop+ break返回值——适合需要「找到结果就退出」的场景


iter
============

iter()返回：值的不可变引用

.. code-block:: rust

  let numbers: Vec<i32> = (0..1000).collect();

  let sum: i32 = numbers
      .iter()
      .filter(|&&x| x % 2 == 0)
      .map(|&x| x * x)
      .sum();

  let sum: i32 = numbers
    .iter()
    .copied()              // &i32 → i32
    .filter(|x| x % 2 == 0)
    .map(|x| x * x)
    .sum();