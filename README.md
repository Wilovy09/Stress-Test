# Stress tests

## Test general

```js
import http from "k6/http";
import { check } from "k6";

export const options = {
  stages: [
    { duration: "1m", target: 1000 },
    { duration: "1m", target: 2000 },
    { duration: "1m", target: 3000 },
    { duration: "1m", target: 0 },
  ],
};

export default function () {
  const url = "http://localhost:8080/login";
  const payload = JSON.stringify({ username: "Test", password: "MyAwesomePassword1234." });

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
    timeout: "60s",
  };

  const response = http.post(url, payload, params);

  check(response, {
    "success login": (r) => r.status === 200,
    "response contains data": (r) => r.body && r.body.includes("username"),
  });
}
```

## V

```v
module main

import vweb
import json

struct App {
	vweb.Context
}

fn main() {
	http_port := 8080
	app := &App{}
	vweb.run(app, http_port)
}

@['/login'; post]
fn (mut app App) login() vweb.Result {
	body := app.req.data
	data := json.encode(body)
	return app.json(data)
}
```

```txt
     ✗ success login
      ↳  98% — ✓ 182702 / ✗ 2927
     ✗ response contains data
      ↳  98% — ✓ 182702 / ✗ 2927

     checks.........................: 98.42% 365404 out of 371258
     data_received..................: 35 MB  139 kB/s
     data_sent......................: 36 MB  144 kB/s
     http_req_blocked...............: avg=226.16ms min=0s       med=299.76µs max=19.75s   p(90)=86.18ms p(95)=1.02s  
     http_req_connecting............: avg=225.23ms min=0s       med=210.14µs max=19.75s   p(90)=81.72ms p(95)=1.02s  
     http_req_duration..............: avg=1.57s    min=0s       med=487.22ms max=1m0s     p(90)=2.37s   p(95)=5.82s  
       { expected_response:true }...: avg=1.29s    min=185.32µs med=487.97ms max=53.16s   p(90)=2.18s   p(95)=5.7s   
     http_req_failed................: 1.57%  2927 out of 185629
     http_req_receiving.............: avg=4.88ms   min=0s       med=189.64µs max=631.21ms p(90)=11.36ms p(95)=26.7ms 
     http_req_sending...............: avg=4.36ms   min=0s       med=116.14µs max=612.49ms p(90)=10.07ms p(95)=24.29ms
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s       p(90)=0s      p(95)=0s     
     http_req_waiting...............: avg=1.57s    min=0s       med=478.95ms max=1m0s     p(90)=2.35s   p(95)=5.79s  
     http_reqs......................: 185629 749.14811/s
     iteration_duration.............: avg=2.07s    min=443.4µs  med=524.28ms max=1m0s     p(90)=3.51s   p(95)=11.16s 
     iterations.....................: 185629 749.14811/s
     vus............................: 10     min=2                max=3000
     vus_max........................: 3000   min=3000             max=3000


running (4m07.8s), 0000/3000 VUs, 185629 complete and 216 interrupted iterations
default ✓ [======================================] 0000/3000 VUs  4m0s
```

### Conclusiónes 

- La aplicación pudo manejar 3,000 usuarios virtuales concurrentes, logrando una alta tasa de éxito en las verificaciones y pocas fallas (1.57%).
- Sin embargo, hubo tiempos elevados en algunas métricas clave, como:
    - Tiempo promedio de bloqueo (226 ms).
    - Duración de la solicitud promedio (1.57 segundos).

Es posible que se requiera optimización en áreas como la conexión inicial y el tiempo de respuesta del servidor.
Estas métricas indican un buen desempeño general, pero con margen para mejorar la latencia y estabilidad bajo carga alta.

## Rust

```rust
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
```

```txt
     ✗ success login
      ↳  99% — ✓ 773296 / ✗ 3834
     ✗ response contains data
      ↳  99% — ✓ 773296 / ✗ 3834

     checks.........................: 99.50% 1546592 out of 1554260
     data_received..................: 126 MB 525 kB/s
     data_sent......................: 150 MB 624 kB/s
     http_req_blocked...............: avg=21.64ms  min=0s       med=5.41µs   max=21.27s p(90)=10.95µs  p(95)=17.23µs 
     http_req_connecting............: avg=21.52ms  min=0s       med=0s       max=20.78s p(90)=0s       p(95)=0s      
     http_req_duration..............: avg=349.25ms min=0s       med=116.02ms max=1m2s   p(90)=449.05ms p(95)=613.07ms
       { expected_response:true }...: avg=264.17ms min=275.78µs med=115.93ms max=59.48s p(90)=444.86ms p(95)=602.1ms 
     http_req_failed................: 0.49%  3837 out of 777133
     http_req_receiving.............: avg=698.12µs min=0s       med=36.41µs  max=2.83s  p(90)=105.14µs p(95)=391.17µs
     http_req_sending...............: avg=912.43µs min=0s       med=17.7µs   max=3.73s  p(90)=47.46µs  p(95)=163.26µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s     p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=347.64ms min=0s       med=115.13ms max=1m2s   p(90)=446.44ms p(95)=609.12ms
     http_reqs......................: 777133 3238.004066/s
     iteration_duration.............: avg=476.63ms min=395.21µs med=137.77ms max=1m5s   p(90)=517.12ms p(95)=756.4ms 
     iterations.....................: 777130 3237.991566/s
     vus............................: 49     min=0                  max=3000
     vus_max........................: 3000   min=1656               max=3000


running (4m00.0s), 0000/3000 VUs, 777130 complete and 56 interrupted iterations
default ✓ [======================================] 0000/3000 VUs  4m0s
```

### Conclusiónes

- La aplicación manejó con éxito hasta 3,000 usuarios virtuales concurrentes, alcanzando una tasa de éxito en las verificaciones de 99.50%, con solo un 0.49% de fallas.
- Las métricas clave muestran un rendimiento bastante sólido, aunque hay áreas que podrían optimizarse:
    - El tiempo promedio de bloqueo fue de 21.64 ms, con un máximo de hasta 21.27 segundos.
    - El tiempo promedio de duración de las solicitudes fue de 349.25 ms, con picos de hasta 1 minuto y 2 segundos.
    - El tiempo promedio de espera para las solicitudes fue de 347.64 ms, lo cual también muestra margen para mejora en la latencia.

Aunque la aplicación presentó un desempeño general muy bueno bajo la carga de 3,000 usuarios virtuales, las métricas de tiempo de respuesta y de conexión inicial sugieren que hay oportunidades de mejora en términos de optimización para reducir la latencia y mejorar la estabilidad cuando se enfrenta a altos volúmenes de tráfico.

## Go

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reqBody)
}

func main() {
	httpPort := 8080

	http.HandleFunc("/login", loginHandler)

	fmt.Printf("Servidor escuchando en el puerto %d...\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

```

```txt
     ✓ success login
     ✓ response contains data

     checks.........................: 100.00% 3909674 out of 3909674
     data_received..................: 321 MB  1.3 MB/s
     data_sent......................: 377 MB  1.6 MB/s
     http_req_blocked...............: avg=202.39µs min=2.05µs   med=4.92µs   max=745.16ms p(90)=7.91µs   p(95)=12.15µs 
     http_req_connecting............: avg=189.43µs min=0s       med=0s       max=744.99ms p(90)=0s       p(95)=0s      
     http_req_duration..............: avg=121.96ms min=118.94µs med=104.88ms max=1.02s    p(90)=212.04ms p(95)=288.74ms
       { expected_response:true }...: avg=121.96ms min=118.94µs med=104.88ms max=1.02s    p(90)=212.04ms p(95)=288.74ms
     http_req_failed................: 0.00%   0 out of 1954837
     http_req_receiving.............: avg=691.96µs min=11.87µs  med=31.9µs   max=664.22ms p(90)=108.25µs p(95)=347.99µs
     http_req_sending...............: avg=470.16µs min=6.65µs   med=15.68µs  max=967.9ms  p(90)=39.14µs  p(95)=154.94µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s       p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=120.8ms  min=74.05µs  med=104.11ms max=866.97ms p(90)=210.4ms  p(95)=284ms   
     http_reqs......................: 1954837 8145.014972/s
     iteration_duration.............: avg=179.67ms min=222.91µs med=141.14ms max=1.42s    p(90)=387.61ms p(95)=505.26ms
     iterations.....................: 1954837 8145.014972/s
     vus............................: 44      min=3                  max=3000
     vus_max........................: 3000    min=3000               max=3000


running (4m00.0s), 0000/3000 VUs, 1954837 complete and 0 interrupted iterations
default ✓ [======================================] 0000/3000 VUs  4m0s
```

### Conclusiónes

- La aplicación logró manejar de manera eficiente hasta 3,000 usuarios virtuales concurrentes, con una tasa de éxito del 100% en las verificaciones, sin registrar fallos en las solicitudes.
- El rendimiento general fue excelente, con métricas clave mostrando tiempos bajos en la mayoría de las operaciones:
    - El tiempo promedio de bloqueo fue de 202.39 µs, con un máximo de 745.16 ms.
    - El tiempo promedio de duración de las solicitudes fue de 121.96 ms, con picos de hasta 1.02 segundos.
    - El tiempo promedio de espera de la solicitud fue de 120.8 ms, sin superarse los 900 ms en los casos más extremos.

Con una tasa de fallos de 0.00% y un rendimiento general bastante estable, la aplicación demostró ser muy eficiente bajo carga alta. Sin embargo, aún existen algunos picos en las métricas de bloqueo y espera que podrían optimizarse para mantener la consistencia del rendimiento en escenarios de tráfico muy elevado.

