// Axum示例
use axum::{Router, extract::Path, response::Json, routing::get};
use serde_json::{Value, json};

async fn root() -> Json<Value> {
    Json(json!({
        "message": "Hello, World!"
    }))
}

async fn greet(Path(name): Path<String>) -> Json<Value> {
    Json(json!({
        "message": format!("Hello, {}!", name)
    }))
}

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/", get(root))
        .route("/{name}", get(greet));

    let listener = tokio::net::TcpListener::bind("127.0.0.1:8080")
        .await
        .unwrap();
    axum::serve(listener, app).await.unwrap();
}
