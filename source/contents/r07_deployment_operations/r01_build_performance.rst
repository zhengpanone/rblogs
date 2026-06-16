编译构建与性能优化 (Build & Performance Optimization)
===========================================================

.. .. module:: r07_deployment_operations.r01_build_performance

编译速度和运行时性能直接影响开发体验和生产成本。
本章介绍 Rust 编译优化、构建配置、Profile 调优和二进制体积优化。

Cargo 构建配置 (Profile)
--------------------------

``Cargo.toml`` 中的 ``[profile.*]`` 配置控制编译器的优化级别和调试信息。

.. code-block:: toml

    # Cargo.toml

    [profile.dev]
    opt-level = 0          # 快速编译，不优化
    debug = true           # 包含调试信息
    overflow-checks = true # 整数溢出检查

    [profile.release]
    opt-level = 3          # 最高优化级别
    debug = false          # 不含调试信息
    lto = "fat"            # 链接时优化（增加编译时间）
    codegen-units = 1      # 单一代码生成单元（更好优化）
    panic = "abort"        # panic 时直接 abort（减小体积）
    strip = "symbols"      # 移除符号表（减小体积）
    overflow-checks = false

    # 自定义 Profile：兼顾性能与调试
    [profile.release-debug]
    inherits = "release"
    debug = true           # 保留调试符号
    strip = "none"

    [profile.dist]
    inherits = "release"
    lto = "fat"
    codegen-units = 1
    strip = "symbols"

.. list-table:: 常用 Profile 配置项
   :header-rows: 1

   * - 配置项
     - 说明
     - 推荐值 (release)
   * - ``opt-level``
     - 优化级别 0-3, "s", "z"
     - ``3`` 或 ``"z"`` （体积优先）
   * - ``lto``
     - 链接时优化
     - ``"fat"`` （最佳优化）
   * - ``codegen-units``
     - 并行代码生成单元数
     - ``1`` （最佳优化）
   * - ``panic``
     - panic 策略
     - ``"abort"`` （减小体积）
   * - ``strip``
     - 移除符号/debug 信息
     - ``"symbols"``
   * - ``overflow-checks``
     - 整数溢出检查
     - ``false``
   * - ``incremental``
     - 增量编译
     - dev 用 ``true``，release 用 ``false``

.. note::

   ``lto = "fat"`` 和 ``codegen-units = 1`` 会显著增加编译时间（可能 2-5 倍），
   但能带来 5%-20% 的运行时性能提升和更小的二进制体积。
   建议在 CI 的最终构建阶段使用，本地开发避免。

编译加速
----------

**sccache (共享编译缓存)：**

.. code-block:: bash

    # 安装
    cargo install sccache

    # 配置 Cargo 使用 sccache
    # ~/.cargo/config.toml
    [build]
    rustc-wrapper = "/path/to/sccache"

    # 查看缓存命中率
    sccache --show-stats

**mold (更快的链接器)：**

.. code-block:: bash

    # macOS
    brew install mold

    # Linux
    sudo apt install mold

.. code-block:: toml

    # .cargo/config.toml
    [target.x86_64-unknown-linux-gnu]
    linker = "clang"
    rustflags = ["-C", "link-arg=-fuse-ld=mold"]

    [target.x86_64-apple-darwin]
    rustflags = ["-C", "link-arg=-fuse-ld=lld"]

**cargo-binstall (预编译二进制安装)：**

.. code-block:: bash

    # 安装 cargo-binstall
    curl -L --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/cargo-bins/cargo-binstall/main/install-from-binstall-release.sh | bash

    # 使用预编译二进制安装 CLI 工具（跳过本地编译）
    cargo binstall cargo-audit cargo-deny cargo-watch

.. list-table:: 编译加速工具对比
   :header-rows: 1

   * - 工具
     - 加速方式
     - 效果
   * - ``sccache``
     - 编译缓存（跨项目共享）
     - CI 中 50-80% 命中率
   * - ``mold`` / ``lld``
     - 更快的链接器
     - 链接时间减少 5-10 倍
   * - ``cargo-binstall``
     - 跳过编译，直接安装二进制
     - 工具安装从分钟级到秒级
   * - ``cranelift`` 后端
     - 更快的 debug 编译
     - dev 编译快 30-50%
   * - ``cargo-nextest``
     - 更快的测试运行器
     - 测试运行快 2-3 倍

二进制体积优化
-----------------

.. code-block:: bash

    # 使用 cargo-bloat 分析体积
    cargo install cargo-bloat
    cargo bloat --release          # 按函数体积排序
    cargo bloat --release --crates # 按 Crate 体积排序

    # 使用 cargo-binutils 查看段信息
    cargo install cargo-binutils
    cargo size --release
    cargo size --release -- -A     # 详细段列表

.. code-block:: toml

    # 极致体积优化 Cargo.toml
    [profile.release]
    opt-level = "z"      # 体积优先优化
    lto = "fat"
    codegen-units = 1
    panic = "abort"
    strip = "symbols"

.. code-block:: bash

    # 使用 upx 进一步压缩（可选，可能有兼容性问题）
    upx --best --lzma target/release/my-app

.. note::

   ``opt-level = "z"`` 比 ``"s"`` 更激进地优化体积，但可能牺牲少量性能。
   对 CLI 工具和嵌入式场景推荐 ``"z"``，对服务器应用推荐 ``"s"`` 或 ``3``。

cargo-chef (Docker 层缓存优化)
----------------------------------

``cargo-chef`` 将依赖编译和应用编译分离，最大化 Docker 层缓存利用率。

.. code-block:: dockerfile

    # Dockerfile
    FROM rust:1.80-slim AS chef
    RUN cargo install cargo-chef
    WORKDIR /app

    # Step 1: 准备依赖清单
    FROM chef AS planner
    COPY . .
    RUN cargo chef prepare --recipe-path recipe.json

    # Step 2: 仅编译依赖（可被缓存）
    FROM chef AS builder
    COPY --from=planner /app/recipe.json recipe.json
    RUN cargo chef cook --release --recipe-path recipe.json

    # Step 3: 编译应用
    COPY . .
    RUN cargo build --release

    # Step 4: 最小运行时镜像
    FROM debian:bookworm-slim AS runtime
    RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
    COPY --from=builder /app/target/release/my-app /usr/local/bin/
    ENTRYPOINT ["/usr/local/bin/my-app"]

.. code-block:: bash

    # 构建
    docker build -t my-app:latest .

Feature Flags 优化
---------------------

.. code-block:: toml

    # Cargo.toml
    [features]
    default = ["json", "cli"]
    json = ["serde", "serde_json"]
    cli = ["clap"]
    server = ["tokio", "axum", "tower-http"]
    full = ["json", "cli", "server"]

    [dependencies]
    serde = { version = "1", optional = true }
    serde_json = { version = "1", optional = true }
    clap = { version = "4", optional = true }
    tokio = { version = "1", optional = true, features = ["full"] }
    axum = { version = "0.8", optional = true }
    tower-http = { version = "0.6", optional = true, features = ["cors"] }

.. code-block:: bash

    # 仅编译 CLI 功能
    cargo build --release --no-default-features --features cli

    # 检查未使用的 feature
    cargo install cargo-features-manager
    cargo features prune --check

Cross-Compilation (交叉编译)
------------------------------

``cross`` 提供了零配置的交叉编译体验。

.. code-block:: bash

    # 安装
    cargo install cross

    # 交叉编译到目标平台（自动使用 Docker）
    cross build --release --target aarch64-unknown-linux-gnu     # ARM64 Linux
    cross build --release --target x86_64-unknown-linux-musl     # 静态链接 Linux
    cross build --release --target x86_64-pc-windows-gnu          # Windows
    cross build --release --target aarch64-apple-darwin           # macOS ARM (需 macOS)

.. code-block:: toml

    # Cross.toml (可选配置)
    [target.aarch64-unknown-linux-gnu]
    pre-build = [
        "dpkg --add-architecture arm64",
        "apt-get update && apt-get install -y libssl-dev:arm64",
    ]

.. list-table:: 交叉编译常用 Target
   :header-rows: 1

   * - Target
     - 平台
     - libc
   * - ``x86_64-unknown-linux-gnu``
     - Linux x86_64
     - glibc
   * - ``x86_64-unknown-linux-musl``
     - Linux x86_64 静态
     - musl
   * - ``aarch64-unknown-linux-gnu``
     - Linux ARM64
     - glibc
   * - ``aarch64-unknown-linux-musl``
     - Linux ARM64 静态
     - musl
   * - ``x86_64-apple-darwin``
     - macOS x86_64
     - System
   * - ``aarch64-apple-darwin``
     - macOS ARM (Apple Silicon)
     - System
   * - ``x86_64-pc-windows-msvc``
     - Windows MSVC
     - MSVC

总结
-----

.. list-table:: 构建与优化 Crate/工具 总览
   :header-rows: 1

   * - 工具
     - 用途
     - 适用场景
   * - Profile 配置
     - 优化级别/体积控制
     - 所有项目
   * - ``sccache``
     - 编译缓存
     - CI 加速
   * - ``mold`` / ``lld``
     - 快速链接器
     - 本地开发/CI
   * - ``cargo-chef``
     - Docker 层缓存
     - 容器化部署
   * - ``cargo-bloat``
     - 体积分析
     - 二进制瘦身
   * - ``cross``
     - 交叉编译
     - 多平台发布
   * - ``cargo-binstall``
     - 预编译安装
     - 工具链安装
