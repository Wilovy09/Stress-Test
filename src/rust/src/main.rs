use actix_web::{post, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize)]
struct LoginRequest {
    username: String,
    password: String,
}

#[post("/login")]
async fn login(req_body: web::Json<LoginRequest>) -> impl Responder {
    HttpResponse::Ok().json(req_body.into_inner())
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let http_port = 8080;

    HttpServer::new(|| {
        App::new()
            .service(login)
    })
    .bind(("127.0.0.1", http_port))?
    .run()
    .await
}

