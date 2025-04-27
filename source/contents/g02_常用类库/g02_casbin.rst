===================
Casbin 
===================

Casbin介绍
=================

Casbin 是一个开源的访问控制库，用于实现权限管理和访问控制模型。它提供了一种简单而灵活的方式来定义和强制应用程序中的访问控制规则。

Casbin 主要解决的是如何控制用户对资源的访问权限，确保只有具有合适权限的用户可以执行指定操作。


Casbin 的核心思想是基于两种常见的访问控制模型：

- **访问控制列表（Access Control List, ACL）**：通过列出每个资源的访问权限来管理用户对资源的访问。每个资源都有一个列表，列出哪些用户可以执行哪些操作。
  
- **角色访问控制（Role-Based Access Control，RBAC）**：通过角色来管理用户权限，每个用户被赋予一个或多个角色，每个角色拥有一组权限。这种方式使得权限管理更加集中和高效。


Casbin 允许通过编程的方式定义资源、操作和角色之间的关系，并在运行时根据这些规则进行验证和授权。

访问控制模型的核心概念
-------------------------------------------------------

Casbin 的访问控制模型由三个主要概念组成：

**模型规则（Model Rule）**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

  - 定义资源、操作和角色之间的关系。模型规则是描述权限管理逻辑的蓝图，它定义了哪些请求应该被允许，哪些请求应该被拒绝。
  - 使用类似于自然语言的策略语法来描述访问控制规则。例如，可以定义规则来指定用户、角色和资源之间的关系。

线编辑器: `Casbin Editor <https://casbin.org/editor/>`_

.. literalinclude:: ./code/2_casbin/example_rbac_model.pml
  :encoding: utf-8
  :language: ini
  :linenos:

模型配置主要有五个部分：[request_definition]，[policy_definition]，[role_definition]，[policy_effect] 和 [matchers]，分别表示请求定义、策略定义、角色定义、策略效果定义、匹配器，其中 [role_definition] 角色定义是用于基于角色的模型（RBAC），支持用 # 开头表示注释。

请求定义
:::::::::::::::

- 定义请求的格式，`r` 表示请求，`sub` 表示用户，`obj` 表示资源，`act` 表示操作。
- 例如，`r = sub, obj, act` 表示请求包含用户、资源和操作三个部分。 

.. code-block:: ini
   :linenos:

    [request_definition]
    r = sub, obj, act

请求的定义，定义了决策器决策 enforce() 时传入的参数。sub, obj, act 就表示 enforce() 方法需传入 3 个参数，支持增加或删除参数数量，保证顺序和数量一一对应即可。

定义的参数名也可以修改，例如，要表示用户对接接口地址的权限，

则请求可以定义为：

.. code-block:: ini
  :linenos:

    [request_definition]
    r = uid, uri, method

那么决策时可以这样使用：

.. code-block:: go
  :linenos:

    enforce(888,"/api/post/add", "POST");

意思是需要决策 uid 为 888 的用户对接口 /api/post/add 是否具有 POST 请求的权限。

策略定义
:::::::::::::::

.. code-block:: ini
  :linenos:

    [policy_definition]
    p = sub, obj, act

策略定义，它定义了策略中值的含义及顺序，参数名和数量同样支持修改。

假设，我们的策略表有如下的策略记录：

.. csv-table:: Casbin 策略示例
   :header: "p", "uid", "uri", "method", "created_at"

   p, 886, /api/post/add, POST, 2025-01-10 08:16:01
   p, 887, /api/post/delete, DELETE, 2025-01-10 08:20:11
   p, 888, /api/post/list, GET, 2025-01-10 08:30:22


那么，我们的策略可以定义为：

.. code-block:: ini
  :linenos:

    [policy_definition]
    p = uid, uri, method, created_at

这个四个参数就依次表示了策略中的 v1~v4 列，v0 可以理解为策略类型或分组，就是定义中的 p 。

那 Casbin 在加载策略的时候，就会按照这个对应关系加载，使用某条策略中的某个字段时可以用p.uid的方式。


角色定义
:::::::::::::::

.. code-block:: ini
  :linenos:
  
    [role_definition]
    g = _, _

角色定义，表明了用户和角色的继承关系，通常用于基于角色的权限模型。其中 ``_, _``，两个下划线，可以称作前向和后项，表示前向继承后项，通常用于表示用户属于某个角色，或者角色继承另一个角色。

如果为 ``_, _, _`` 三个下划线，用于多租户模型，则最后一项表示多租户里的域。

策略效果定义
:::::::::::::::

.. code-block:: ini
  :linenos:

    [policy_effect]
    e = some(where (p.eft == allow))

策略的定义的取值是固定值，并且是系统内置的硬编码，不支持自定义。目前支持以下几种配置：

- ``some(where (p.eft == allow))`` ，表示在所有命中的策略中，只要有一条 ``allow`` 的策略，那么结果就是 ``true``，通俗来讲，就是 ``一票赞成制`` 。
- ``!some(where (p.eft == deny))`` ，表示在命中的策略中，如果没有 ``deny`` 的，那最终的决策结果就是true，否则就是 ``false``，即 ``一票否决制`` 。
- ``some(where (p.eft == allow)) && !some(where (p.eft == deny))`` ，表示在所有命中的策略中，要有一条 ``allow`` 的策略，并且不能有 ``deny`` 的策略，则最终决策结果就是 ``true`` ，例如：只要有人赞成且没有人反对即通过。
- ``priority(p.eft) || deny``，用于隐式/显式优先级模型，以命中的策略的 ``eft`` 的值作为决策结果，如果没有命中，或者 ``eft`` 没有明确的结果，那么就是 ``false`` 。
- ``subjectPriority(p.eft)``，用于基于用户和角色的层级关系的优先级模型。

匹配器
:::::::::::::::

.. code-block:: ini
  :linenos:
  
    [matchers]
    m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act

匹配器，是对匹配规则的定义，m 的值是一个表达式。在调用 ``enforce(...)`` 方法进行决策执行的时候，会根据表达式中的变量带入运算执行表达式，依据表达式的执行结果来决定策略是否命中或未命中。

匹配器支持常用的算术运算符如 ``+, -, *, /``，以及一些逻辑运算符如 ``&&, ||, !``，并且还支持内置函数和自定义函数，像例子中的 ``g(r.sub, p.sub)`` 就是一个内置函数 ``(g(...))``。

这个表达式中的变量名就是前面请求定义、策略定义中的定义的名称， **支持自行修改变量名，保证上下文一致即可**。比如前文我们演示了 ``uid, uri, method`` 这样的定义，那么表达式中可以这样使用 ``r.uri == p.uri && r.method == p.method``。

**策略规则（Policy Rule）**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

  - 存储了用户、角色、资源和操作之间的映射关系。策略规则决定了具体的权限控制配置，定义了用户、角色和权限之间的映射。
  - 策略规则可以从外部数据源（如数据库、配置文件等）加载，也可以通过代码进行动态配置。

.. literalinclude:: ./code/2_casbin/example_casbin_policy.go
  :encoding: utf-8
  :language: golang
  :lines: 1-8
  :linenos:

**执行器（Enforcer）**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

  - 执行器是 Casbin 的核心组件，负责验证访问请求是否符合访问控制规则。执行器根据定义的模型规则和策略规则，判断某个用户是否有权访问某个资源。
  - 执行器根据规则的配置判断是否允许或拒绝访问。


.. literalinclude:: ./code/2_casbin/example_casbin_policy.go
  :encoding: utf-8
  :language: golang
  :lines: 10-19
  :linenos:

4. **模型加载**


.. literalinclude:: ./code/2_casbin/example_casbin_policy.go
  :encoding: utf-8
  :language: golang
  :lines: 22-50
  :linenos:

当初始化Gorm Mysql Adpter并加载策略时，适配器会自动使用或创建一个默认表名为 ``casbin_rule`` ；也可以使用自定义表名初始化适配器，但无论使用默认表还是自定义表，Casbin都会按照如下结构创建表。

.. literalinclude:: ./code/2_casbin/casbin_rule.sql
  :encoding: utf-8
  :language: sql
  :linenos:

添加/保存/删除策略

.. code-block:: ini
  :linenos:

  res, err := enforcer.AddPolicy("alice","data1","read")
  if err !=nil{
    panic(err)
  }

  err = enforcer.SavePolicy()
  if err !=nil{
    panic(err)
  }

  res, err := enforcer.RemovePolicy("alice","data1","read")
  if err !=nil{
    panic(err)
  }

验证规则

.. code-block:: ini
  :linenos:

  // 权限检查，创建请求
  sub:="alice"
  obj :="data1"
  act :="read"
  ok, err := enforcer.Enforce(sub, obj, act)
  if err !=nil{
      log.Println("err:", err)
  }
  if ok ==true{
      log.Println("true")
  }else{
      log.Println("false")
  }

5. **中间件定义**

.. literalinclude:: ./code/2_casbin/example_casbin_policy.go
  :encoding: utf-8
  :language: golang
  :lines: 52-86
  :linenos:

支持的访问控制模型
-------------------------------------------------------
Casbin 支持多种访问控制模型，您可以根据需求选择合适的模型进行实现。常见的模型有：

1. **基于角色的访问控制（RBAC）**  
    - 通过定义用户角色，并将权限分配给角色，管理不同角色的权限。每个用户可以有多个角色，每个角色拥有特定的访问权限。

2. **基于属性的访问控制（ABAC）**  
    - 通过定义用户和资源的属性来管理权限。例如，权限可以基于用户的部门、资源的类型等属性进行动态控制。

3. **多租户访问控制**  
    - 在多租户环境中，Casbin 支持根据租户 ID 进行访问控制。不同租户可以有不同的资源和权限控制。

实现细节与集成
-------------------------------------------------------
Casbin 提供了与多个编程语言和框架的集成，包括 Go、Java、Python、Node.js 等。它还提供了与常见存储（如文件、数据库等）的集成，使得权限配置可以灵活存储和管理。

数据存储与策略管理
>>>>>>>>>>>>>>>>>>>>>>

Casbin 提供了多种存储方式，包括：

- **内存存储**：适用于简单的应用场景，适合快速开发和原型设计。
- **数据库存储**：将策略规则保存在数据库中，适用于需要持久化和跨会话的场景。
- **文件存储**：将策略规则保存在文件中，适用于配置较为简单的应用。

通过与数据库（如 MySQL、PostgreSQL）或缓存（如 Redis）集成，Casbin 可以灵活地管理权限，并动态更新策略。

可扩展性与灵活性
>>>>>>>>>>>>>>>>>>>>>>

Casbin 是高度可扩展的，您可以根据需要定义自己的访问控制模型和策略规则，甚至可以在运行时动态修改权限。Casbin 还提供了丰富的中间件和 API，用于与 Web 框架（如 Gin、Express）集成，简化权限验证过程。

优势
-------------------------------------------------------
1. **灵活性**：Casbin 支持多种访问控制模型，可以根据具体业务需求选择合适的模型。
2. **可扩展性**：支持动态加载和修改策略，满足不同场景下的权限控制需求。
3. **易于集成**：Casbin 提供了与多种语言和框架的集成方式，开发者可以快速将 Casbin 集成到现有系统中。

总结
-------------------------------------------------------
Casbin 提供了一种强大、灵活的方式来实现复杂的权限管理。无论是基于角色的访问控制（RBAC）、基于属性的访问控制（ABAC），还是多租户的场景，Casbin 都能够满足不同的需求。通过简洁的配置和强大的执行器，Casbin 帮助开发者实现了精细化的权限控制，提升了应用的安全性和可维护性。
