======================
常见的Trait和用法
======================

Newtype 模式
==================

NewType它可以解决：

- 原始类型污染（String 到处飞）

- ID 混用（UserId 和 OrderId 都是 i64）

- 无法在类型层表达领域语义

- 违反 DDD 的“显式建模”

Newtype 本质就是：用 struct 包一层单字段类型，创建一个“新类型”

.. code-block:: rust

  struct UserId(i64);

虽然底层是 i64，但：

- 它 不等于 i64

- 不能和 OrderId 混用

- 是一个独立类型

这就是 DDD 想要的：领域语义通过类型表达

Newtype 应该主要出现在：

- Value Object

- Entity Id

- 强语义字段（Email / Password / Username）

Newtype 的使用场景
---------------------------

.. code-block:: rust

  #[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
  pub struct UserId(i64);

  #[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
  pub struct OrderId(i64);


Path 参数转换

.. code-block:: rust

  use axum::{extract::Path};
  use std::str::FromStr;

  async fn get_user(Path(id): Path<UserId>) {
  }

  impl FromStr for UserId {
    type Err = std::num::ParseIntError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(UserId(s.parse()?))
    }
  }