=================
引用和值转换
=================

引用转换为值
-----------------

.. list-table:: Option 转换方法对照表
   :widths: 25 15 15 45
   :header-rows: 1

   * - 方法
     - 输入
     - 输出
     - 要求
   * - ``copied()``
     - ``Option<&T>``
     - ``Option<T>``
     - ``T: Copy``
   * - ``cloned()``
     - ``Option<&T>``
     - ``Option<T>``
     - ``T: Clone``
   * - ``to_owned()``
     - ``&T``
     - ``T``
     - ``T: ToOwned``
   * - ``clone()``
     - ``&T``
     - ``T``
     - ``T: Clone``
   * - ``map(|x| x.clone())``
     - ``Option<&T>``
     - ``Option<T>``
     - ``T: Clone``
   * - ``map(Clone::clone)``
     - ``Option<&T>``
     - ``Option<T>``
     - ``T: Clone``

copied()
================

**Option<&T> 转换为 Option<T>**

.. code-block:: rust

  let x = Some(&100);
  let y = x.copied();
  println!("{:?}", y);

结果

>>> Some(100)

等价于

.. code-block:: rust

  let y = x.map(|v| *v);

**Iterator 转换为 Option<T>**

.. code-block:: rust

  let nums = vec![1, 2, 3];

  let v: Vec<i32> = nums.iter().copied().collect();


类型变换：

.. code-block:: text

  Iter<&i32>
    ↓
  copied()
      ↓
  Iter<i32>

相当于

.. code-block:: rust

  nums.iter().map(|x| *x)


cloned()
================

**Option<&T> 转换为 Option<T>**

.. code-block:: rust

  let x = Some(&String::from("hello"));
  let y = x.cloned();

结果

>>> Some(String::from("hello"))

等价于

.. code-block:: rust

  let y = x.map(|v| v.clone());

**Iterator 转换为 Option<T>**

.. code-block:: rust

  let strs = vec![
      String::from("a"),
      String::from("b")
  ];

  let v: Vec<String> =
      strs.iter().cloned().collect();


等价于

.. code-block:: rust

  strs.iter().map(|x| x.clone());

值转换为引用
-----------------

as_ref()
===============

.. code-block:: rust

  let s = Some(String::from("abc"));
  let r = s.as_ref();

类型转换：

.. code-block:: text

  Option<String>
    ↓
  as_ref()
    ↓
  Option<&String>

as_deref()
===============
.. code-block:: rust

  let s = Some(String::from("abc"));
  let r = s.as_deref();

类型转换：

.. code-block:: text

  Option<String>
    ↓
  as_deref()
    ↓
  Option<&str>

等价于:

.. code-block:: rust

  s.as_ref().map(|v| v.as_str())

as_mut()
===============

.. code-block:: rust

  let mut x = Some(100);
  if let Some(v) = x.as_mut() {
      *v += 1;
  }

类型转换：

.. code-block:: text

  Option<i32>
    ↓
  as_mut()
    ↓
  Option<&mut i32>

Option/Result的转置
-----------------------

transpose()
===============
	
**交换两层包装**

- **Option<Result<T, E>> 转换为 Result<Option<T>, E>**
- **Result<Option<T>, E> 转换为 Option<Result<T, E>>**

定义:

.. code-block:: rust

  // Option 上的 transpose
  impl<T, E> Option<Result<T, E>> {
      pub fn transpose(self) -> Result<Option<T>, E>;
  }

  // Result 上的 transpose
  impl<T, E> Result<Option<T>, E> {
      pub fn transpose(self) -> Option<Result<T, E>>;
  }

作用:

在 ``Option<Result<T, E>>`` 和 ``Result<Option<T>, E>`` 之间相互转换

示例:

.. code-block:: rust

  // Option<Result> → Result<Option>
  let x: Option<Result<i32, &str>> = Some(Ok(42));
  let y: Result<Option<i32>, &str> = x.transpose();
  assert_eq!(y, Ok(Some(42)));

  let x: Option<Result<i32, &str>> = Some(Err("error"));
  let y: Result<Option<i32>, &str> = x.transpose();
  assert_eq!(y, Err("error"));

  let x: Option<Result<i32, &str>> = None;
  let y: Result<Option<i32>, &str> = x.transpose();
  assert_eq!(y, Ok(None));

  // Result<Option> → Option<Result>
  let x: Result<Option<i32>, &str> = Ok(Some(42));
  let y: Option<Result<i32, &str>> = x.transpose();
  assert_eq!(y, Some(Ok(42)));

记忆口诀:

    **外层变成内层，内层变成外层**

- ``Some(Ok(v))↔ Ok(Some(v))``
- ``Some(Err(e))↔ Err(e)``
- ``None↔ Ok(None)``

典型应用场景

处理可能失败的查找操作：

.. code-block:: rust

  fn find_user(id: u32) -> Option<Result<User, Error>> { ... }

  // 传统写法
  let result: Result<Option<User>, Error> = match find_user(42) {
      Some(Ok(user)) => Ok(Some(user)),
      Some(Err(e)) => Err(e),
      None => Ok(None),
  };

  // 用 transpose
  let result: Result<Option<User>, Error> = find_user(42).transpose();


inspect()
---------------

**在不改变值的情况下插入调试或副作用操作**

定义

.. code-block:: rust

  // Iterator 上的 inspect
  impl<I: Iterator> Iterator for Inspect<I, F> {
      // ...
  }

  fn inspect<F>(self, f: F) -> Inspect<Self, F>
  where
      F: FnMut(&Self::Item);

作用

在迭代器的每个元素上调用闭包，但 **不改变元素本身**，类似于 forEach但保持惰性求值。

示例

.. code-block:: rust

  let v = vec![1, 2, 3, 4, 5];

  let result: Vec<i32> = v.iter()
    .inspect(|x| println!("before map: {}", x))
    .map(|x| x * 2)
    .inspect(|x| println!("after map: {}", x))
    .filter(|x| x > 5)
    .collect();

  // 输出：
  // before map: 1
  // after map: 2
  // before map: 2
  // after map: 4
  // before map: 3
  // after map: 6
  // before map: 4
  // after map: 8
  // before map: 5
  // after map: 10

调试用途

.. code-block:: rust

  // 调试过滤条件
  let nums = vec![1, 2, 3, 4, 5];
  let evens: Vec<_> = nums.iter()
      .inspect(|x| println!("checking {}...", x))
      .filter(|x| *x % 2 == 0)
      .inspect(|x| println!("{} passed filter", x))
      .collect();


与 for_each的区别

.. code-block:: rust

  // inspect — 惰性，只在消费时执行
  let iter = vec![1, 2, 3].iter().inspect(|x| println!("{}", x));
  // 此时还没打印任何东西
  let sum: i32 = iter.sum(); // 现在才打印

  // for_each — 立即消费整个迭代器
  vec![1, 2, 3].iter().for_each(|x| println!("{}", x));
  // 立即打印


在 Option/Result上也有 inspect

.. code-block:: rust

  // Option::inspect
  let x: Option<i32> = Some(42);
  x.inspect(|v| println!("value is {}", v)); // 打印 "value is 42"

  let x: Option<i32> = None;
  x.inspect(|v| println!("value is {}", v)); // 不打印

  // Result::inspect / inspect_err
  let x: Result<i32, &str> = Ok(42);
  x.inspect(|v| println!("ok: {}", v));      // 打印 "ok: 42"
  x.inspect_err(|e| println!("err: {}", e)); // 不打印

  let x: Result<i32, &str> = Err("fail");
  x.inspect(|v| println!("ok: {}", v));      // 不打印
  x.inspect_err(|e| println!("err: {}", e)); // 打印 "err: fail"


常见组合
-----------

最大值

.. code-block:: rust

  let max = 
      nums.iter()
        .max()
        .copied()
        .unwrap_or(0);


字符串

.. code-block:: rust

  let name =
    map.get("name")
       .cloned()
       .unwrap_or_default();

Option<String> → Option<&str>

.. code-block:: rust

  let name =
    config.name
          .as_deref()
          .unwrap_or("guest");


