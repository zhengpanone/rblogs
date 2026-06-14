===================================
数据结构 & 集合 Crate
===================================

标准库的集合（``Vec``、``HashMap``、``HashSet`` 等）已经很强，但以下 crate 提供了更丰富的迭代器能力、并行计算和特殊数据结构。

.. contents:: 目录
   :depth: 3
   :local:

itertools
==============

迭代器增强工具，提供标准库没有的迭代器适配器和组合方法。

.. code-block:: rust

   use itertools::Itertools;

   fn main() {
       // chunk：分组
       let data = vec![1, 2, 3, 4, 5, 6];
       for chunk in &data.into_iter().chunks(3) {
           println!("{:?}", chunk.collect::<Vec<_>>());
       }
       // [1, 2, 3]
       // [4, 5, 6]

       // tuple_windows：滑动窗口
       let v = vec![1, 2, 3, 4, 5];
       let diffs: Vec<i32> = v.iter().tuple_windows().map(|(&a, &b)| b - a).collect();
       println!("相邻差: {:?}", diffs); // [1, 1, 1, 1]

       // combinations：组合
       let items = ['a', 'b', 'c'];
       for combo in items.iter().combinations(2) {
           println!("{:?}", combo);
       }
       // ['a', 'b']
       // ['a', 'c']
       // ['b', 'c']

       // permutations：排列
       let perms: Vec<_> = (0..3).permutations(2).collect();
       println!("排列: {:?}", perms);
       // [[0, 1], [0, 2], [1, 0], [1, 2], [2, 0], [2, 1]]

       // sorted / unique
       let v = vec![3, 1, 2, 1, 3, 4];
       let sorted_unique: Vec<_> = v.into_iter().sorted().unique().collect();
       println!("排序去重: {:?}", sorted_unique); // [1, 2, 3, 4]

       // group_by：按 key 分组
       let words = vec!["apple", "banana", "avocado", "blueberry", "apricot"];
       let groups: Vec<_> = words
           .into_iter()
           .group_by(|w| w.chars().next().unwrap())
           .into_iter()
           .map(|(key, group)| (key, group.collect::<Vec<_>>()))
           .collect();
       println!("分组: {:?}", groups);
       // [('a', ["apple", "avocado", "apricot"]), ('b', ["banana", "blueberry"])]

       // join：拼接
       let names = vec!["Alice", "Bob", "Charlie"];
       println!("{}", names.iter().join(", ")); // Alice, Bob, Charlie

       // cartesian_product：笛卡尔积
       let xs = 0..3;
       let ys = 10..12;
       let product: Vec<_> = xs.cartesian_product(ys).collect();
       println!("笛卡尔积: {:?}", product);
       // [(0, 10), (0, 11), (1, 10), (1, 11), (2, 10), (2, 11)]

       // zip_longest：不等长 zip
       let a = [1, 2, 3];
       let b = ["x", "y"];
       let zipped: Vec<_> = a.iter().zip_longest(b.iter()).collect();
       println!("不等长 zip: {:?}", zipped);
       // [Both(&1, &"x"), Both(&2, &"y"), Left(&3)]
   }

常用方法：

.. list-table:: itertools 常用方法
   :header-rows: 1
   :widths: 30 70

   * - 方法
     - 说明
   * - ``chunks(n)``
     - 按固定大小分组
   * - ``tuple_windows()``
     - 滑动窗口（(a,b), (b,c), (c,d)...）
   * - ``combinations(k)``
     - 组合（C(n, k)）
   * - ``permutations(k)``
     - 排列（P(n, k)）
   * - ``sorted()``
     - 排序（返回新迭代器）
   * - ``unique()``
     - 去重（连续重复）
   * - ``group_by(key_fn)``
     - 按 key 分组
   * - ``join(sep)``
     - 用分隔符拼接
   * - ``cartesian_product(other)``
     - 笛卡尔积
   * - ``zip_longest(other)``
     - 不等长 zip
   * - ``duplicates()``
     - 找出重复元素
   * - ``minmax()``
     - 同时求最小最大值

rayon
==========

数据并行库，将迭代器计算自动并行化，几乎不需要改代码。

基本使用：

.. code-block:: rust

   use rayon::prelude::*;

   fn main() {
       let numbers: Vec<i32> = (0..1_000_000).collect();

       // 并行求和
       let sum: i32 = numbers.par_iter().sum();
       println!("并行求和: {}", sum);

       // 并行 map
       let squares: Vec<i32> = numbers.par_iter().map(|x| x * x).collect();
       println!("前 5 个平方: {:?}", &squares[..5]);

       // 并行 filter
       let evens: Vec<i32> = numbers.par_iter().filter(|&&x| x % 2 == 0).cloned().collect();
       println!("偶数个数: {}", evens.len());

       // 并行 find
       let found = numbers.par_iter().find_any(|&&x| x > 900_000);
       println!("找到: {:?}", found);

       // 并行 fold
       let total: i32 = numbers.par_iter().fold(|| 0, |acc, &x| acc + x).sum();
       println!("并行 fold 求和: {}", total);
   }

并行可变操作：

.. code-block:: rust

   use rayon::prelude::*;

   fn main() {
       let mut numbers: Vec<i32> = (0..100).collect();

       // 并行修改
       numbers.par_iter_mut().for_each(|x| *x = (*x) * 2);

       // 并行排序
       numbers.par_sort();
       println!("排序后前 5 个: {:?}", &numbers[..5]);
   }

CPU 密集型计算：

.. code-block:: rust

   use rayon::prelude::*;

   fn is_prime(n: u64) -> bool {
       if n < 2 { return false; }
       (2..=((n as f64).sqrt() as u64)).all(|i| n % i != 0)
   }

   fn main() {
       let range = 1_000_000u64..1_000_100;
       let primes: Vec<u64> = range
           .into_par_iter()
           .filter(|&n| is_prime(n))
           .collect();
       println!("质数: {:?}", primes);
   }

核心 API：

.. list-table:: rayon 核心 API
   :header-rows: 1
   :widths: 30 70

   * - 方法
     - 说明
   * - ``par_iter()``
     - 并行不可变迭代器
   * - ``par_iter_mut()``
     - 并行可变迭代器
   * - ``into_par_iter()``
     - 并行消费迭代器（转移所有权）
   * - ``par_sort()``
     - 并行排序
   * - ``par_chunks(n)``
     - 并行分块
   * - ``find_any(predicate)``
     - 并行查找（任意一个匹配即返回）
   * - ``reduce(identity, op)``
     - 并行归约
   * - ``fold(init_fn, op)``
     - 并行折叠

全局线程池配置：

.. code-block:: rust

   // 设置线程数
   rayon::ThreadPoolBuilder::new()
       .num_threads(8)
       .build_global()
       .unwrap();

注意事项：

.. list-table:: rayon 使用注意
   :header-rows: 1
   :widths: 25 75

   * - 注意点
     - 说明
   * - 数据量小时不划算
     - 并行有调度开销，小数据串行更快
   * - 避免在并行闭包中使用 ``Rc``
     - ``Rc`` 不是 ``Send``，应使用 ``Arc``
   * - ``collect`` 保证顺序
     - ``par_iter().map().collect()`` 保持元素顺序
   * - ``find_any`` 不保证顺序
     - 返回最先被某线程找到的匹配

dashmap
==========

并发哈希表，无锁读取，接近 ``HashMap`` 的 API，无需 ``Mutex`` 包裹。

.. code-block:: rust

   use dashmap::DashMap;
   use std::sync::Arc;
   use std::thread;

   fn main() {
       let map = Arc::new(DashMap::new());

       // 并发写入
       let mut handles = vec![];
       for i in 0..10 {
           let map = Arc::clone(&map);
           handles.push(thread::spawn(move || {
               for j in 0..100 {
                   map.insert(i * 100 + j, format!("value-{}", i * 100 + j));
               }
           }));
       }

       for handle in handles {
           handle.join().unwrap();
       }

       println!("map 大小: {}", map.len()); // 1000

       // 读取
       if let Some(v) = map.get(&42) {
           println!("key=42: {}", *v);
       }

       // 遍历
       for entry in map.iter() {
           println!("{} -> {}", entry.key(), entry.value());
       }

       // 原子更新
       map.entry(42).and_modify(|v| *v = "updated".to_string()).or_insert("default".to_string());
   }

常用 API：

.. list-table:: dashmap 常用 API
   :header-rows: 1
   :widths: 30 70

   * - 方法
     - 说明
   * - ``insert(key, value)``
     - 插入键值对
   * - ``get(&key)``
     - 获取值（返回 ``Option<Ref<K, V>>``）
   * - ``remove(&key)``
     - 删除键
   * - ``entry(key)``
     - 原子 entry API（and_modify / or_insert）
   * - ``iter()``
     - 遍历所有条目
   * - ``len()``
     - 获取条目数
   * - ``clear()``
     - 清空
   * - ``contains_key(&key)``
     - 是否包含键

DashMap vs Arc\<Mutex\<HashMap\>\>：

.. list-table:: DashMap vs Arc\<Mutex\<HashMap\>\>
   :header-rows: 1
   :widths: 25 40 35

   * - 特性
     - ``DashMap``
     - ``Arc<Mutex<HashMap>>``
   * - 并发读
     - 无锁，高并发
     - 需要获取锁
   * - 并发写
     - 分段锁，细粒度
     - 互斥锁，串行
   * - API 兼容
     - 接近 HashMap
     - 完全兼容 HashMap
   * - 内存开销
     - 略高（分片结构）
     - 低
   * - 适用场景
     - 读多写少、高并发
     - 简单场景、锁持有时间短

indexmap
==========

保持插入顺序的哈希表 / 集合。

.. code-block:: rust

   use indexmap::IndexMap;

   fn main() {
       let mut map = IndexMap::new();

       map.insert("c", 3);
       map.insert("a", 1);
       map.insert("b", 2);

       // 遍历顺序 = 插入顺序
       for (key, value) in &map {
           println!("{}: {}", key, value);
       }
       // c: 3
       // a: 1
       // b: 2

       // 按索引访问
       if let Some((key, value)) = map.get_index(0) {
           println!("第一个: {} -> {}", key, value); // c -> 3
       }

       // 获取键的索引
       if let Some(index) = map.get_index_of("a") {
           println!("'a' 的位置: {}", index); // 1
       }

       // 交换位置
       map.swap_indices(0, 1);
       println!("交换后: {:?}", map.keys().collect::<Vec<_>>()); // ["a", "c", "b"]

       // 按索引插入
       map.shift_insert(0, "first", 0);
       println!("插入后: {:?}", map.keys().collect::<Vec<_>>());
       // ["first", "a", "c", "b"]
   }

IndexSet：

.. code-block:: rust

   use indexmap::IndexSet;

   fn main() {
       let mut set = IndexSet::new();
       set.insert("c");
       set.insert("a");
       set.insert("b");

       // 保持插入顺序
       println!("{:?}", set.iter().collect::<Vec<_>>()); // ["c", "a", "b"]

       // 按索引访问
       println!("第一个: {}", set.get_index(0).unwrap()); // "c"
   }

``HashMap`` vs ``IndexMap`` vs ``BTreeMap``：

.. list-table:: Map 类型对比
   :header-rows: 1
   :widths: 20 25 30 25

   * - 类型
     - 顺序
     - 查找性能
     - 适用场景
   * - ``HashMap``
     - 无序
     - O(1) 平均
     - 通用场景，不关心顺序
   * - ``IndexMap``
     - 插入顺序
     - O(1) 平均
     - 需要保持插入顺序
   * - ``BTreeMap``
     - 排序顺序
     - O(log n)
     - 需要按 key 排序遍历

总结
=====

.. list-table:: 数据结构 & 集合 Crate 总览
   :header-rows: 1
   :widths: 20 30 50

   * - Crate
     - 用途
     - 典型场景
   * - ``itertools``
     - 迭代器增强
     - 组合/排列、分组、滑动窗口、笛卡尔积
   * - ``rayon``
     - 数据并行
     - CPU 密集型批量计算、并行 map/filter/sort
   * - ``dashmap``
     - 并发哈希表
     - 多线程共享缓存、计数器、全局状态
   * - ``indexmap``
     - 有序哈希表/集合
     - 保持插入顺序的场景、LRU 缓存、去重但保序
