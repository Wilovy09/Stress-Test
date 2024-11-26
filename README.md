# Stress tests

## Specs del PC utilizado

- i5-4210U (4) @ 2.70GHz
- 8GB de RAM

## Test general

```js
import http from "k6/http";
import { check } from "k6";

function randomString(length) {
  const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  let result = "";
  for (let i = 0; i < length; i++) {
    result += charset.charAt(Math.floor(Math.random() * charset.length));
  }
  return result;
}

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

  const username = `${randomString(8)}`;
  const password = `${randomString(12)}`;

  const payload = JSON.stringify({ username, password });

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

> [!NOTE]
> Se esta usando la versión 0.4.8

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
      ↳  99% — ✓ 343939 / ✗ 2444
     ✗ response contains data
      ↳  99% — ✓ 343939 / ✗ 2444

     checks.........................: 99.29% 687878 out of 692766
     data_received..................: 63 MB  260 kB/s
     data_sent......................: 65 MB  266 kB/s
     http_req_blocked...............: avg=121.61ms min=0s       med=215.59µs max=19.84s   p(90)=27.27ms  p(95)=71.19ms
     http_req_connecting............: avg=121.45ms min=0s       med=151.44µs max=19.84s   p(90)=26.85ms  p(95)=70.12ms
     http_req_duration..............: avg=870.94ms min=0s       med=370.11ms max=59.01s   p(90)=650.4ms  p(95)=1.51s  
       { expected_response:true }...: avg=702.53ms min=212.07µs med=370.05ms max=57.8s    p(90)=647.12ms p(95)=1.4s   
     http_req_failed................: 0.70%  2444 out of 346383
     http_req_receiving.............: avg=2.37ms   min=0s       med=132.28µs max=561.67ms p(90)=5.76ms   p(95)=13.69ms
     http_req_sending...............: avg=2.1ms    min=0s       med=63.71µs  max=467.39ms p(90)=5.05ms   p(95)=12.54ms
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s       p(90)=0s       p(95)=0s     
     http_req_waiting...............: avg=866.45ms min=0s       med=366.67ms max=59s      p(90)=646.54ms p(95)=1.51s  
     http_reqs......................: 346383 1428.088601/s
     iteration_duration.............: avg=1.09s    min=495.33µs med=375.59ms max=1m0s     p(90)=732.61ms p(95)=3.49s  
     iterations.....................: 346383 1428.088601/s
     vus............................: 887    min=4                max=3000
     vus_max........................: 3000   min=3000             max=3000


running (4m02.6s), 0000/3000 VUs, 346383 complete and 128 interrupted iterations
default ✓ [======================================] 0000/3000 VUs  4m0s
```

### Conclusiónes 

- La aplicación soportó una carga de hasta 3,000 usuarios virtuales concurrentes, logrando una tasa de éxito del 99.29% en las verificaciones, con un 0.70% de fallos registrados (2,444 solicitudes fallidas de 346,383 totales).
- Aunque el desempeño fue sólido en general, algunas métricas clave indican áreas de mejora:
    - Tiempo promedio de bloqueo: 121.61 ms, con un máximo de 19.84 segundos.
    - Duración promedio de las solicitudes: 870.94 ms, con un máximo de 59.01 segundos.
    - Tiempo promedio de espera: 866.45 ms, mostrando latencias significativas en escenarios de alta demanda.
- Puntos destacados:
    - Los tiempos de recepción y envío fueron relativamente bajos, con promedios de 2.37 ms y 2.1 ms respectivamente.
    - La aplicación procesó 1,428 solicitudes por segundo, con una duración promedio de iteración de 1.09 segundos.

### Recomendaciones

- Reducir la latencia en el tiempo de bloqueo y de espera para mejorar la experiencia del usuario bajo condiciones de carga alta.
- Investigar las causas detrás de las solicitudes fallidas y los picos elevados en las métricas clave (e.g., duraciones superiores a 1 minuto).
- Implementar optimizaciones en la conexión inicial y el manejo de tráfico concurrente para mejorar la estabilidad y consistencia.

### Resumen

En resumen, la aplicación mostró un buen desempeño con una alta carga concurrente, aunque aún existen áreas de oportunidad en la latencia y el manejo de fallos para mejorar la experiencia bajo demanda extrema.

## Rust

> [!NOTE]
> Se esta usando la versión 1.81.0 

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
      ↳  99% — ✓ 1419455 / ✗ 3588
     ✗ response contains data
      ↳  99% — ✓ 1419455 / ✗ 3588

     checks.........................: 99.74%  2838910 out of 2846086
     data_received..................: 223 MB  929 kB/s
     data_sent......................: 266 MB  1.1 MB/s
     http_req_blocked...............: avg=6.07ms   min=0s       med=5.67µs   max=20.11s  p(90)=13.18µs  p(95)=19.58µs 
     http_req_connecting............: avg=6.06ms   min=0s       med=0s       max=20.11s  p(90)=0s       p(95)=0s      
     http_req_duration..............: avg=163.6ms  min=0s       med=64.62ms  max=1m0s    p(90)=182.34ms p(95)=253.81ms
       { expected_response:true }...: avg=101.36ms min=111.73µs med=64.64ms  max=59.67s  p(90)=181.38ms p(95)=251.78ms
     http_req_failed................: 0.25%   3593 out of 1423048
     http_req_receiving.............: avg=400.01µs min=0s       med=35.23µs  max=776.6ms p(90)=121.25µs p(95)=249.39µs
     http_req_sending...............: avg=688.94µs min=0s       med=17.66µs  max=1.19s   p(90)=64.8µs   p(95)=153.17µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s      p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=162.51ms min=0s       med=64.02ms  max=1m0s    p(90)=179.81ms p(95)=250.41ms
     http_reqs......................: 1423048 5929.146088/s
     iteration_duration.............: avg=260.37ms min=258.15µs med=100.95ms max=1m1s    p(90)=287.19ms p(95)=376.92ms
     iterations.....................: 1423041 5929.116922/s
     vus............................: 41      min=5                  max=3000
     vus_max........................: 3000    min=3000               max=3000


running (4m00.0s), 0000/3000 VUs, 1423041 complete and 435 interrupted iterations
default ✓ [======================================] 0000/3000 VUs  4m0s
```

### Conclusiónes

- La aplicación manejó exitosamente una carga de hasta 3,000 usuarios virtuales concurrentes, logrando una tasa de éxito del 99.74%, con solo un 0.25% de solicitudes fallidas (3,588 de un total de 1,423,041).
- Las métricas clave indican un buen rendimiento general, aunque hay áreas con oportunidad de optimización:
    - Duración promedio de las solicitudes: 163.6 ms, con un máximo registrado de 1 minuto.
    - Tiempo promedio de espera: 162.51 ms, mostrando un buen manejo en la mayoría de los casos, pero con picos notables.
    - Tiempo promedio de bloqueo: 6.07 ms, con un máximo de hasta 20.11 segundos.
- Puntos destacados:
    - Las métricas de recepción y envío fueron rápidas, con promedios de 400 µs y 689 µs respectivamente.
    - La aplicación procesó 5,929 solicitudes por segundo, con una duración promedio de iteración de 260.37 ms.

### Recomendaciones

- Analizar los picos elevados en la duración y el tiempo de espera de las solicitudes, especialmente aquellos que alcanzan hasta 1 minuto, para identificar posibles cuellos de botella.
- Optimizar el tiempo de bloqueo para reducir los valores máximos que podrían impactar el desempeño en condiciones de carga extrema.
- Monitorear y reducir las tasas de error, asegurándose de que las solicitudes fallidas no estén relacionadas con problemas críticos de la aplicación.

### Resumen

En general, la aplicación mostró un desempeño robusto, procesando un alto volumen de solicitudes concurrentes con una tasa de éxito muy elevada. Sin embargo, optimizar los tiempos máximos y analizar los errores podría mejorar aún más la experiencia bajo cargas intensivas.

## Go

> [!NOTE]
> Se esta usando la versión 1.23.3

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

     checks.........................: 100.00% 3178976 out of 3178976
     data_received..................: 251 MB  1.0 MB/s
     data_sent......................: 297 MB  1.2 MB/s
     http_req_blocked...............: avg=264.84µs min=2.3µs    med=5.28µs   max=758.85ms p(90)=8.46µs   p(95)=14.45µs 
     http_req_connecting............: avg=253.59µs min=0s       med=0s       max=758.57ms p(90)=0s       p(95)=0s      
     http_req_duration..............: avg=132.21ms min=114.56µs med=110.18ms max=1.16s    p(90)=248.87ms p(95)=316.3ms 
       { expected_response:true }...: avg=132.21ms min=114.56µs med=110.18ms max=1.16s    p(90)=248.87ms p(95)=316.3ms 
     http_req_failed................: 0.00%   0 out of 1589488
     http_req_receiving.............: avg=829.64µs min=11.97µs  med=32.69µs  max=897.8ms  p(90)=118.75µs p(95)=330.85µs
     http_req_sending...............: avg=675.44µs min=6.99µs   med=16.57µs  max=1.04s    p(90)=40.44µs  p(95)=150.62µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s       p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=130.7ms  min=82.39µs  med=109.27ms max=1.03s    p(90)=246.48ms p(95)=310.31ms
     http_reqs......................: 1589488 6622.764153/s
     iteration_duration.............: avg=216.94ms min=263.13µs med=165.18ms max=1.59s    p(90)=479.32ms p(95)=613.47ms
     iterations.....................: 1589488 6622.764153/s
     vus............................: 45      min=2                  max=3000
     vus_max........................: 3000    min=3000               max=3000


running (4m00.0s), 0000/3000 VUs, 1589488 complete and 0 interrupted iterations
default ✓ [======================================] 0000/3000 VUs  4m0s
```

### Conclusiónes

- La aplicación manejó exitosamente una carga de hasta 3,000 usuarios virtuales concurrentes, alcanzando una tasa de éxito del 100%, sin registrar fallos en las solicitudes (0 de 1,589,488).
- Las métricas clave indican un excelente rendimiento, con tiempos consistentes y bajos en la mayoría de los casos:
    - Duración promedio de las solicitudes: 132.21 ms, con un máximo de 1.16 segundos.
    - Tiempo promedio de espera: 130.7 ms, con una mediana de 109.27 ms, mostrando una experiencia de usuario fluida incluso bajo carga alta.
    - Tiempo promedio de bloqueo: 264.84 µs, con valores máximos bien controlados de hasta 758.85 ms.
- Puntos destacados:
    - Los tiempos de recepción y envío fueron rápidos, con promedios de 829 µs y 675 µs respectivamente.
    - La aplicación procesó 6,622 solicitudes por segundo, con una duración promedio de iteración de 216.94 ms.

### Recomendaciones

- Aunque el desempeño fue sobresaliente, sería beneficioso revisar los casos en los que las solicitudes alcanzaron duraciones máximas de hasta 1.16 segundos para garantizar que no haya latencias ocultas bajo carga extrema.
- Continuar monitoreando el tiempo de bloqueo, especialmente los valores en el percentil 95 y máximos, para asegurarse de que se mantenga consistente en escenarios futuros.

### Resumen

En resumen, la aplicación mostró un desempeño excepcional bajo una carga considerable, manejando un alto volumen de solicitudes con tiempos de respuesta rápidos y una tasa de éxito perfecta. Se encuentra en óptimas condiciones para manejar tráfico intenso.

