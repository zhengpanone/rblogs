容器化与 CI/CD (Containerization & CI/CD)
============================================

.. .. module:: r07_deployment_operations.r02_containerization_cicd

将 Rust 应用部署到生产环境需要可靠的容器化方案和自动化 CI/CD 流程。
本章涵盖 Docker 最佳实践、CI/CD 配置和云部署方案。

Docker 最佳实践
-----------------

**多阶段构建 (Multi-stage Build)：**

.. code-block:: dockerfile

    # ============ 构建阶段 ============
    FROM rust:1.80-slim-bookworm AS builder

    # 创建空项目用于缓存依赖
    RUN USER=root cargo init --bin /app
    WORKDIR /app

    # 先复制 Cargo.toml 和 Cargo.lock
    COPY Cargo.toml Cargo.lock ./

    # 构建依赖（利用 Docker 层缓存）
    RUN mkdir src && echo "fn main() {}" > src/main.rs
    RUN cargo build --release
    RUN rm -rf src

    # 复制真实源码并构建
    COPY src ./src
    # 强制重新编译（touch 使缓存失效）
    RUN touch src/main.rs
    RUN cargo build --release

    # ============ 运行阶段 ============
    FROM debian:bookworm-slim

    # 安装运行时依赖
    RUN apt-get update && \
        apt-get install -y --no-install-recommends \
            ca-certificates \
            libssl3 \
        && rm -rf /var/lib/apt/lists/*

    # 创建非 root 用户
    RUN useradd -m -u 1000 appuser
    USER appuser
    WORKDIR /home/appuser

    # 复制二进制
    COPY --from=builder /app/target/release/my-app /usr/local/bin/my-app

    # 健康检查
    HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
        CMD curl -f http://localhost:8080/health || exit 1

    EXPOSE 8080
    ENTRYPOINT ["/usr/local/bin/my-app"]

**静态链接 (musl) 最小镜像：**

.. code-block:: dockerfile

    FROM rust:1.80-alpine AS builder
    RUN apk add --no-cache musl-dev
    WORKDIR /app
    COPY . .
    RUN cargo build --release --target x86_64-unknown-linux-musl

    FROM scratch
    COPY --from=builder /app/target/x86_64-unknown-linux-musl/release/my-app /my-app
    ENTRYPOINT ["/my-app"]

.. list-table:: Docker 基础镜像选择
   :header-rows: 1

   * - 镜像
     - 大小
     - 适用场景
   * - ``scratch``
     - ~0 MB
     - 纯静态链接应用
   * - ``alpine``
     - ~7 MB
     - 需要少量系统工具
   * - ``debian:bookworm-slim``
     - ~80 MB
     - 需要 glibc / OpenSSL
   * - ``distroless``
     - ~20 MB
     - 安全优先（无 shell/包管理器）
   * - ``ubuntu:jammy``
     - ~77 MB
     - 需要完整系统环境

.dockerignore
^^^^^^^^^^^^^^^^

.. code-block:: text

    # .dockerignore
    target/
    .git/
    **/*.rs.bk
    *.pdb
    **/target/
    Dockerfile
    .dockerignore
    .env
    *.log

Docker Compose 编排
---------------------

.. code-block:: yaml

    # docker-compose.yml
    version: "3.8"

    services:
      app:
        build:
          context: .
          dockerfile: Dockerfile
        ports:
          - "8080:8080"
        environment:
          - DATABASE_URL=postgres://user:pass@db:5432/mydb
          - RUST_LOG=info
        depends_on:
          db:
            condition: service_healthy
        restart: unless-stopped
        healthcheck:
          test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
          interval: 30s
          timeout: 3s
          retries: 3

      db:
        image: postgres:16-alpine
        environment:
          POSTGRES_USER: user
          POSTGRES_PASSWORD: pass
          POSTGRES_DB: mydb
        volumes:
          - pgdata:/var/lib/postgresql/data
        healthcheck:
          test: ["CMD-SHELL", "pg_isready -U user"]
          interval: 5s
          timeout: 5s
          retries: 5

      redis:
        image: redis:7-alpine
        volumes:
          - redisdata:/data

    volumes:
      pgdata:
      redisdata:

GitHub Actions CI/CD
---------------------

.. code-block:: yaml

    # .github/workflows/ci.yml
    name: CI

    on:
      push:
        branches: [main]
      pull_request:
        branches: [main]

    env:
      CARGO_TERM_COLOR: always
      RUSTFLAGS: "-D warnings"

    jobs:
      # ============ 代码质量 ============
      lint:
        runs-on: ubuntu-latest
        steps:
          - uses: actions/checkout@v4
          - uses: dtolnay/rust-toolchain@stable
            with:
              components: clippy, rustfmt
          - uses: Swatinem/rust-cache@v2

          - name: Check formatting
            run: cargo fmt --all -- --check

          - name: Clippy
            run: cargo clippy --all-targets --all-features -- -D warnings

      # ============ 测试 ============
      test:
        runs-on: ubuntu-latest
        needs: lint
        services:
          postgres:
            image: postgres:16-alpine
            env:
              POSTGRES_USER: test
              POSTGRES_PASSWORD: test
              POSTGRES_DB: testdb
            ports:
              - 5432:5432
            options: >-
              --health-cmd pg_isready
              --health-interval 10s
              --health-timeout 5s
              --health-retries 5

        steps:
          - uses: actions/checkout@v4
          - uses: dtolnay/rust-toolchain@stable
          - uses: Swatinem/rust-cache@v2

          - name: Run tests
            run: cargo test --all-features
            env:
              DATABASE_URL: postgres://test:test@localhost:5432/testdb

      # ============ 安全审计 ============
      security:
        runs-on: ubuntu-latest
        steps:
          - uses: actions/checkout@v4
          - uses: actions-rust-lang/audit@v1

      # ============ 构建 & 推送镜像 ============
      build-and-push:
        runs-on: ubuntu-latest
        needs: [test, security]
        if: github.ref == 'refs/heads/main'
        steps:
          - uses: actions/checkout@v4

          - name: Set up Docker Buildx
            uses: docker/setup-buildx-action@v3

          - name: Login to Container Registry
            uses: docker/login-action@v3
            with:
              registry: ghcr.io
              username: ${{ github.actor }}
              password: ${{ secrets.GITHUB_TOKEN }}

          - name: Build and push
            uses: docker/build-push-action@v5
            with:
              context: .
              push: true
              tags: |
                ghcr.io/${{ github.repository }}:latest
                ghcr.io/${{ github.repository }}:${{ github.sha }}
              cache-from: type=gha
              cache-to: type=gha,mode=max

.. list-table:: GitHub Actions 关键 Actions
   :header-rows: 1

   * - Action
     - 用途
   * - ``dtolnay/rust-toolchain``
     - 安装 Rust 工具链（推荐）
   * - ``Swatinem/rust-cache``
     - Rust 编译缓存
   * - ``actions-rust-lang/audit``
     - 依赖安全审计
   * - ``docker/build-push-action``
     - Docker 构建与推送
   * - ``docker/setup-buildx-action``
     - 启用 Docker BuildKit

shuttle (零配置 Rust 云部署)
-------------------------------

``shuttle`` 是 Rust 原生的 Serverless 平台，通过注解式部署。

.. code-block:: rust

    use axum::{routing::get, Router};
    use shuttle_runtime::SecretStore;

    async fn hello() -> &'static str {
        "Hello from Shuttle!"
    }

    #[shuttle_runtime::main]
    async fn main(
        #[shuttle_shared_db::Postgres] pool: sqlx::PgPool,
        #[shuttle_secrets::Secrets] secrets: SecretStore,
    ) -> shuttle_axum::ShuttleAxum {
        // 自动运行数据库迁移
        sqlx::migrate!().run(&pool).await.unwrap();

        let app = Router::new().route("/", get(hello));
        Ok(app.into())
    }

.. code-block:: bash

    # 安装 CLI
    cargo install cargo-shuttle

    # 本地运行
    cargo shuttle run

    # 部署到 Shuttle
    cargo shuttle deploy

.. note::

   Shuttle 支持 Axum / Actix-web / Rocket / Poem / Salvo 等主流框架，
   以及 PostgreSQL / Redis / MongoDB / Secrets 等资源自动配置。

Cloud Run / Kubernetes 部署
----------------------------

**Cloud Run (GCP)：**

.. code-block:: yaml

    # cloudbuild.yaml
    steps:
      - name: 'gcr.io/cloud-builders/docker'
        args:
          - 'build'
          - '-t'
          - 'gcr.io/$PROJECT_ID/my-app:$COMMIT_SHA'
          - '.'
      - name: 'gcr.io/cloud-builders/docker'
        args:
          - 'push'
          - 'gcr.io/$PROJECT_ID/my-app:$COMMIT_SHA'
      - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
        args:
          - 'gcloud'
          - 'run'
          - 'deploy'
          - 'my-app'
          - '--image=gcr.io/$PROJECT_ID/my-app:$COMMIT_SHA'
          - '--region=us-central1'
          - '--platform=managed'
          - '--allow-unauthenticated'

**Kubernetes Deployment：**

.. code-block:: yaml

    # k8s/deployment.yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: rust-app
    spec:
      replicas: 3
      selector:
        matchLabels:
          app: rust-app
      template:
        metadata:
          labels:
            app: rust-app
        spec:
          containers:
            - name: rust-app
              image: ghcr.io/myorg/rust-app:latest
              ports:
                - containerPort: 8080
              env:
                - name: DATABASE_URL
                  valueFrom:
                    secretKeyRef:
                      name: app-secrets
                      key: database-url
              resources:
                requests:
                  memory: "64Mi"
                  cpu: "100m"
                limits:
                  memory: "256Mi"
                  cpu: "500m"
              livenessProbe:
                httpGet:
                  path: /health
                  port: 8080
                initialDelaySeconds: 10
                periodSeconds: 15
              readinessProbe:
                httpGet:
                  path: /ready
                  port: 8080
                initialDelaySeconds: 5
                periodSeconds: 10

.. list-table:: 部署平台对比
   :header-rows: 1

   * - 平台
     - 类型
     - 适合场景
   * - Shuttle
     - Rust 原生 Serverless
     - 快速原型、个人项目
   * - Fly.io
     - 边缘计算
     - 全球分布式部署
   * - Cloud Run
     - 托管容器
     - 按需扩缩、零运维
   * - Kubernetes
     - 容器编排
     - 大规模微服务、企业级
   * - Docker Compose
     - 单机编排
     - 开发环境、小规模部署

cargo-release (版本发布自动化)
--------------------------------

.. code-block:: toml

    # release.toml 或 Cargo.toml [package.metadata.release]
    pre-release-replacements = [
      { file="README.md", search="my-app = \"[0-9.]+\"", replace="my-app = \"{{version}}\"" },
      { file="CHANGELOG.md", search="Unreleased", replace="{{version}}", exactly=1 },
      { file="CHANGELOG.md", search="\\.\\.\\.HEAD", replace="...{{tag_name}}", exactly=1 },
    ]
    pre-release-commit-message = "chore: release {{version}}"
    tag-message = "{{tag_name}}"
    tag-name = "v{{version}}"
    publish = true
    push = true

.. code-block:: bash

    # 安装
    cargo install cargo-release

    # Dry run（预览会做什么）
    cargo release patch --dry-run

    # 正式发布（自动 bump 版本、提交、打 tag、推送）
    cargo release patch --execute

    # 手动指定版本
    cargo release 1.2.0 --execute

总结
-----

.. list-table:: 容器化与 CI/CD 工具总览
   :header-rows: 1

   * - 工具
     - 用途
     - 适用场景
   * - Docker 多阶段构建
     - 最小化镜像
     - 所有容器化项目
   * - ``cargo-chef``
     - 依赖缓存加速 Docker 构建
     - 大型项目
   * - GitHub Actions
     - CI/CD 自动化
     - 开源项目、团队协作
   * - Shuttle
     - Rust 原生 Serverless
     - 快速部署
   * - Kubernetes
     - 容器编排
     - 企业级微服务
   * - ``cargo-release``
     - 版本发布自动化
     - Crate 发布流程
