===================
测试 & 开发
===================

Rust 生态中用于基准测试、Mock、属性测试和开发调试的核心 crate。

.. contents:: 目录
   :depth: 3
   :local:

criterion
==========

统计驱动的基准测试框架。提供精确的性能测量、回归检测和可视化报告。

.. code-block:: toml

   [dev-dependencies]
   criterion = { version = "0.5", features = ["html_reports"] }

   [[bench]]
   name = "my_benchmark"
   harness = false

基准测试文件 ``benches/my_benchmark.rs``：

.. code-block:: rust

   use criterion::{black_box, criterion_group, criterion_main, Criterion, BenchmarkId};

   fn fibonacci(n: u64) -> u64 {
       match n {
           0 => 0,
           1 => 1,
           _ => fibonacci(n - 1) + fibonacci(n - 2),
       }
   }

   fn fibonacci_iter(n: u64) -> u64 {
       let mut a = 0;
       let mut b = 1;
       for _ in 0..n {
           let tmp = a + b;
           a = b;
           b = tmp;
       }
       a
   }

   fn bench_fibonacci(c: &mut Criterion) {
       let mut group = c.benchmark_group("fibonacci");
       for i in [10, 20, 30].iter() {
           group.bench_with_input(BenchmarkId::new("递归", i), i, |b, i| {
               b.iter(|| fibonacci(*i))
           });
           group.bench_with_input(BenchmarkId::new("迭代", i), i, |b, i| {
               b.iter(|| fibonacci_iter(*i))
           });
       }
       group.finish();
   }

   criterion_group!(benches, bench_fibonacci);
   criterion_main!(benches);

参数化基准测试：

.. code-block:: rust

   use criterion::{black_box, Criterion};

   fn bench_sort(c: &mut Criterion) {
       let mut group = c.benchmark_group("sort");

       // 不同大小
       for size in [100, 1000, 10000].iter() {
           group.bench_with_input(
               criterion::BenchmarkId::new("vec_sort", size),
               size,
               |b, &size| {
                   let mut vec: Vec<i32> = (0..size).rev().collect();
                   b.iter(|| {
                       vec.sort();
                       black_box(&vec); // 防止编译器优化掉
                   })
               },
           );
       }

       group.finish();
   }

   criterion_group!(benches, bench_sort);
   criterion_main!(benches);

对比两个实现：

.. code-block:: rust

   use criterion::{Criterion, black_box};

   fn hashmap_lookup() -> std::collections::HashMap<String, i32> {
       let mut map = std::collections::HashMap::new();
       for i in 0..1000 {
           map.insert(format!("key_{}", i), i);
       }
       map
   }

   fn btreemap_lookup() -> std::collections::BTreeMap<String, i32> {
       let mut map = std::collections::BTreeMap::new();
       for i in 0..1000 {
           map.insert(format!("key_{}", i), i);
       }
       map
   }

   fn bench_lookup(c: &mut Criterion) {
       let hashmap = hashmap_lookup();
       let btreemap = btreemap_lookup();

       let mut group = c.benchmark_group("map_lookup");

       group.bench_function("HashMap", |b| {
           b.iter(|| hashmap.get(&"key_500".to_string()))
       });

       group.bench_function("BTreeMap", |b| {
           b.iter(|| btreemap.get(&"key_500".to_string()))
       });

       group.finish();
   }

   criterion_group!(benches, bench_lookup);
   criterion_main!(benches);

criterion 常用功能：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 功能
     - 说明
   * - ``bench_function``
     - 单个基准测试函数
   * - ``benchmark_group``
     - 分组基准测试，可对比多个实现
   * - ``BenchmarkId``
     - 参数化测试标识
   * - ``bench_with_input``
     - 带参数的基准测试
   * - ``black_box``
     - 防止编译器优化掉被测代码
   * - ``html_reports``
     - 生成 HTML 可视化报告（target/criterion/report/）
   * - ``throughput``
     - 设置吞吐量单位（Bytes, Elements）
   * - ``sample_size``
     - 自定义采样次数
   * - ``measurement_time``
     - 自定义测量时间

mockall
==========

强大的 Mock 库，基于自动生成的 trait 实现。用于单元测试中模拟外部依赖。

.. code-block:: toml

   [dev-dependencies]
   mockall = "0.12"

Mock trait：

.. code-block:: rust

   use mockall::*;
   use mockall::predicate::*;

   #[automock]
   #[async_trait]
   pub trait UserRepository {
       async fn find_by_id(&self, id: u64) -> Option<String>;
       async fn save(&self, name: &str) -> Result<u64, String>;
       fn count(&self) -> usize;
   }

   #[tokio::test]
   async fn test_with_mock_trait() {
       let mut mock = MockUserRepository::new();

       // 设置期望：find_by_id(1) 被调用时返回 "Alice"
       mock.expect_find_by_id()
           .with(eq(1))
           .times(1)
           .returning(|_| Some("Alice".to_string()));

       // 设置期望：save 被调用任意次，返回 Ok
       mock.expect_save()
           .returning(|name| Ok(name.len() as u64));

       // 设置期望：count 返回固定值
       mock.expect_count()
           .return_const(42);

       // 使用 mock
       let user = mock.find_by_id(1).await;
       assert_eq!(user, Some("Alice".to_string()));

       let id = mock.save("Bob").await.unwrap();
       assert_eq!(id, 3);

       assert_eq!(mock.count(), 42);
   }

   #[tokio::test]
   async fn test_mock_with_side_effects() {
       let mut mock = MockUserRepository::new();

       // 多次调用返回不同值
       mock.expect_find_by_id()
           .times(3)
           .returning(|id| {
               match id {
                   1 => Some("Alice".to_string()),
                   2 => Some("Bob".to_string()),
                   _ => None,
               }
           });

       assert_eq!(mock.find_by_id(1).await, Some("Alice".to_string()));
       assert_eq!(mock.find_by_id(2).await, Some("Bob".to_string()));
       assert_eq!(mock.find_by_id(3).await, None);
   }

Mock 普通结构体的方法：

.. code-block:: rust

   use mockall::*;

   mock! {
       pub Database {
           pub fn connect(&self, url: &str) -> Result<(), String>;
           pub fn query(&self, sql: &str) -> Vec<String>;
           pub async fn query_async(&self, sql: &str) -> Vec<String>;
       }
   }

   #[test]
   fn test_mock_struct() {
       let mut mock_db = MockDatabase::new();

       mock_db.expect_connect()
           .with(eq("postgres://localhost/test"))
           .times(1)
           .returning(|_| Ok(()));

       mock_db.expect_query()
           .returning(|sql| {
               if sql.contains("users") {
                   vec!["Alice".to_string(), "Bob".to_string()]
               } else {
                   vec![]
               }
           });

       assert!(mock_db.connect("postgres://localhost/test").is_ok());
       assert_eq!(mock_db.query("SELECT * FROM users").len(), 2);
       assert!(mock_db.query("SELECT * FROM orders").is_empty());
   }

   #[tokio::test]
   async fn test_async_mock_struct() {
       let mut mock_db = MockDatabase::new();

       mock_db.expect_query_async()
           .returning(|_| vec!["result1".to_string(), "result2".to_string()]);

       let results = mock_db.query_async("SELECT * FROM table").await;
       assert_eq!(results.len(), 2);
   }

验证调用顺序：

.. code-block:: rust

   use mockall::Sequence;

   #[test]
   fn test_call_sequence() {
       let mut seq = Sequence::new();
       let mut mock = MockUserRepository::new();

       mock.expect_find_by_id()
           .times(1)
           .in_sequence(&mut seq)
           .returning(|_| Some("Alice".to_string()));

       mock.expect_save()
           .times(1)
           .in_sequence(&mut seq)
           .returning(|_| Ok(1));

       // 必须按顺序调用
       let _ = mock.find_by_id(1);
       let _ = mock.save("Bob");
   }

常用 Predicate（匹配器）：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - Predicate
     - 说明
   * - ``eq(x)``
     - 等于 x
   * - ``ne(x)``
     - 不等于 x
   * - ``gt(x)`` / ``ge(x)``
     - 大于 / 大于等于 x
   * - ``lt(x)`` / ``le(x)``
     - 小于 / 小于等于 x
   * - ``always()``
     - 始终匹配
   * - ``never()``
     - 从不匹配
   * - ``function(f)``
     - 自定义匹配函数 ``f(x) -> bool``
   * - ``in_iter(iter)``
     - 在迭代器内
   * - ``str::contains(s)``
     - 字符串包含 s
   * - ``str::starts_with(s)``
     - 字符串以 s 开头

Mock 最佳实践：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 实践
     - 说明
   * - Mock trait 而非具体类型
     - 使用 ``#[automock]`` 标注 trait，依赖注入
   * - 明确期望次数
     - 使用 ``.times(n)`` 精确控制调用次数
   * - 使用 predicate
     - 精确匹配参数，避免 ``always()`` 滥用
   * - 验证调用顺序
     - 使用 ``Sequence`` 确保关键操作顺序
   * - 检查点 (checkpoint)
     - 在测试中间调用 ``mock.checkpoint()`` 验证已设置的期望
   * - 清理期望
     - 使用完毕后 ``drop(mock)`` 自动验证所有期望是否满足

proptest
==========

基于性质的测试（Property-Based Testing）框架。自动生成大量随机输入，验证代码在所有输入下满足的性质。

.. code-block:: toml

   [dev-dependencies]
   proptest = "1"

基础用法：

.. code-block:: rust

   use proptest::prelude::*;

   // 性质：反转两次等于原值
   proptest! {
       #[test]
       fn test_reverse_twice(v in ".*") {
           let reversed: String = v.chars().rev().collect();
           let twice: String = reversed.chars().rev().collect();
           assert_eq!(v, twice);
       }

       #[test]
       fn test_sort_len(a in prop::collection::vec(0i32..100, 0..100)) {
           let mut sorted = a.clone();
           sorted.sort();
           // 排序后长度不变
           prop_assert_eq!(a.len(), sorted.len());
           // 排序后首元素 <= 尾元素
           if !sorted.is_empty() {
               prop_assert!(sorted.first() <= sorted.last());
           }
       }
   }

自定义策略（Strategy）：

.. code-block:: rust

   use proptest::prelude::*;
   use proptest::strategy::{Strategy, Just};
   use proptest::arbitrary::Arbitrary;

   // 自定义类型实现 Arbitrary
   #[derive(Debug, Clone)]
   struct UserId(u64);

   impl Arbitrary for UserId {
       type Parameters = ();
       type Strategy = proptest::strategy::Map<
           proptest::num::u64::Any,
           fn(u64) -> UserId,
       >;

       fn arbitrary_with(_: Self::Parameters) -> Self::Strategy {
           (1u64..1000).prop_map(UserId)
       }
   }

   proptest! {
       #[test]
       fn test_user_id_range(id: UserId) {
           prop_assert!(id.0 >= 1);
           prop_assert!(id.0 < 1000);
       }
   }

复杂策略组合：

.. code-block:: rust

   use proptest::prelude::*;
   use proptest::collection::{vec, btree_map};
   use proptest::string::string_regex;

   proptest! {
       // 生成有效的 email 格式字符串
       #[test]
       fn test_email_parsing(email in string_regex("[a-z]+@[a-z]+\\.[a-z]{2,}").unwrap()) {
           let parts: Vec<&str> = email.split('@').collect();
           prop_assert_eq!(parts.len(), 2);
           prop_assert!(parts[1].contains('.'));
       }

       // 生成任意 Vec<i32> 并验证排序性质
       #[test]
       fn test_sort_properties(
           mut v in vec(0i32..1000, 0..50)
       ) {
           let original_len = v.len();
           v.sort();
           prop_assert_eq!(v.len(), original_len);

           // 验证有序性
           for window in v.windows(2) {
               prop_assert!(window[0] <= window[1]);
           }
       }

       // 验证 HashMap 插入性质
       #[test]
       fn test_hashmap_insert(
           items in vec((any::<String>(), 0i32..100), 0..100)
       ) {
           let mut map = std::collections::HashMap::new();
           for (k, v) in &items {
               map.insert(k.clone(), *v);
           }
           // 去重后 map 长度 <= 插入项数
           prop_assert!(map.len() <= items.len());
       }

       // 验证解析-序列化往返性质
       #[test]
       fn test_roundtrip(n in 0i64..1_000_000) {
           let s = n.to_string();
           let parsed: i64 = s.parse().unwrap();
           prop_assert_eq!(n, parsed);
       }
   }

条件假设与过滤：

.. code-block:: rust

   use proptest::prelude::*;

   proptest! {
       #[test]
       fn test_division(a in 0i32..100, b in 0i32..100) {
           // 跳过 b == 0 的情况
           prop_assume!(b != 0);

           let result = a / b;
           // 商 * 除数 + 余数 == 被除数
           prop_assert_eq!(result * b + a % b, a);
       }

       #[test]
       fn test_non_empty_vec(
           v in prop::collection::vec(any::<i32>(), 1..100)  // 长度 >= 1
       ) {
           prop_assert!(!v.is_empty());
           let max = v.iter().max().unwrap();
           prop_assert!(v.iter().all(|x| x <= max));
       }
   }

回归测试 —— 发现失败用例后复现：

.. code-block:: rust

   use proptest::prelude::*;

   proptest! {
       #![proptest_config(ProptestConfig {
           // 失败后保存用例到文件
           failure_persistence: Some(Box::new(
               proptest::test_runner::FileFailurePersistence::WithSource("regression")
           )),
           ..ProptestConfig::default()
       })]

       #[test]
       fn test_buggy_function(a: i32, b: i32) {
           let result = a.checked_add(b);
           // proptest 会自动保存失败的 (a, b) 并复现
           prop_assert!(result.is_some());
       }
   }

proptest 常用策略：

.. list-table::
   :header-rows: 1
   :widths: 30 50

   * - 策略
     - 说明
   * - ``any::<T>()``
     - 生成任意 T 类型值（需实现 Arbitrary）
   * - ``a..b`` (range)
     - 在范围 [a, b) 内生成值
   * - ``prop::bool::ANY``
     - 随机布尔值
   * - ``prop::char::range(a, b)``
     - 在字符范围内生成
   * - ``vec(strategy, size_range)``
     - 生成 Vec，指定元素策略和长度范围
   * - ``btree_map(key, val, size)``
     - 生成 BTreeMap
   * - ``string_regex(pattern)``
     - 根据正则表达式生成字符串
   * - ``Just(x)``
     - 始终返回固定值 x
   * - ``prop_oneof![s1, s2, ...]``
     - 从多个策略中随机选择
   * - ``prop::option::of(s)``
     - 生成 Option<T>
   * - ``s.prop_map(f)``
     - 对策略生成的值应用映射函数
   * - ``s.prop_filter(pred, s)``
     - 过滤满足条件的值（优先用 prop_assume!）

测试工具扩展
==============

其他常用测试辅助工具：

pretty_assertions
------------------

更可读的断言输出：

.. code-block:: toml

   [dev-dependencies]
   pretty_assertions = "1"

.. code-block:: rust

   use pretty_assertions::{assert_eq, assert_ne};

   #[test]
   fn test_pretty_assertions() {
       let expected = vec![
           ("Alice", 30),
           ("Bob", 25),
           ("Charlie", 35),
       ];
       let actual = vec![
           ("Alice", 30),
           ("Bob", 26),  // 差异
           ("Charlie", 35),
       ];

       // 输出彩色 diff，而不是普通断言
       assert_eq!(expected, actual);
   }

rstest
-------

参数化测试框架：

.. code-block:: toml

   [dev-dependencies]
   rstest = "0.19"

.. code-block:: rust

   use rstest::*;

   #[rstest]
   #[case(0, 0)]
   #[case(1, 1)]
   #[case(2, 1)]
   #[case(3, 2)]
   #[case(10, 55)]
   fn test_fibonacci(#[case] input: u64, #[case] expected: u64) {
       assert_eq!(fibonacci(input), expected);
   }

   #[rstest]
   fn test_with_fixture(#[values(1, 2, 3)] base: i32) {
       assert!(base > 0);
   }

   // fixture
   #[fixture]
   fn user_data() -> Vec<String> {
       vec!["Alice".into(), "Bob".into(), "Charlie".into()]
   }

   #[rstest]
   fn test_with_data(user_data: Vec<String>) {
       assert_eq!(user_data.len(), 3);
   }

   fn fibonacci(n: u64) -> u64 {
       match n {
           0 => 0,
           1 => 1,
           _ => fibonacci(n - 1) + fibonacci(n - 2),
       }
   }

tempfile
--------

临时文件和目录：

.. code-block:: toml

   [dev-dependencies]
   tempfile = "3"

.. code-block:: rust

   use tempfile::{tempdir, tempfile, NamedTempFile};
   use std::io::Write;

   #[test]
   fn test_temp_file() -> std::io::Result<()> {
       // 临时文件（自动删除）
       let mut file = NamedTempFile::new()?;
       writeln!(file, "Hello, temp!")?;

       let path = file.path().to_path_buf();
       let content = std::fs::read_to_string(&path)?;
       assert!(content.contains("Hello, temp!"));

       // 关闭时自动删除
       file.close()?;
       Ok(())
   }

   #[test]
   fn test_temp_dir() -> std::io::Result<()> {
       // 临时目录（自动删除）
       let dir = tempdir()?;
       let file_path = dir.path().join("test.txt");

       std::fs::write(&file_path, "content")?;
       assert!(file_path.exists());

       // dir 离开作用域时自动删除目录及内容
       Ok(())
   }

serial_test
------------

串行执行测试（避免共享状态冲突）：

.. code-block:: toml

   [dev-dependencies]
   serial_test = "3"

.. code-block:: rust

   use serial_test::serial;

   #[test]
   #[serial]
   fn test_with_shared_resource_a() {
       // 此测试不会与其他标记 #[serial] 的测试并行
       let mut guard = SHARED_RESOURCE.lock().unwrap();
       *guard = 42;
       assert_eq!(*guard, 42);
   }

   #[test]
   #[serial]
   fn test_with_shared_resource_b() {
       let mut guard = SHARED_RESOURCE.lock().unwrap();
       *guard = 100;
       assert_eq!(*guard, 100);
   }

   use std::sync::Mutex;
   lazy_static::lazy_static! {
       static ref SHARED_RESOURCE: Mutex<i32> = Mutex::new(0);
   }

总结
==========

测试 & 开发 Crate 总览：

.. list-table::
   :header-rows: 1
   :widths: 18 15 45 22

   * - Crate
     - 定位
     - 核心能力
     - 使用场景
   * - ``criterion``
     - 基准测试
     - 统计分析、参数化、对比测试、HTML 报告
     - 性能测量、回归检测
   * - ``mockall``
     - Mock 框架
     - 自动生成 Mock、期望设置、异步支持、调用顺序
     - 隔离外部依赖的单元测试
   * - ``proptest``
     - 属性测试
     - 随机输入生成、性质验证、回归用例保存
     - 发现边界情况、验证不变量
   * - ``pretty_assertions``
     - 断言增强
     - 彩色 diff 输出
     - 复杂数据结构的测试调试
   * - ``rstest``
     - 参数化测试
     - 用例参数化、fixture 注入
     - 减少重复测试代码
   * - ``tempfile``
     - 临时文件
     - 自动清理的临时文件/目录
     - 文件 I/O 测试
   * - ``serial_test``
     - 串行测试
     - 标记测试串行执行
     - 共享全局状态的测试
