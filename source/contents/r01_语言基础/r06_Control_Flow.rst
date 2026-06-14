============
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