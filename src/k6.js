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

