use std::net::SocketAddr;

use axum::{Router, response::Html, routing::get, serve};
use tokio::net::TcpListener;


#[tokio::main]
async fn main() {
    let app = Router::new().route("/", get(handler));

    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));
    let listener = TcpListener::bind(addr).await.unwrap();

    serve(listener, app).await.unwrap();
}

async fn handler() -> Html<&'static str> {
    Html("<h1>Hello, world!</h1>")
}
