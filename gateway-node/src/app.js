const http = require("http");

const port = Number(process.env.PORT || 3001);
const redisHost = process.env.REDIS_HOST || "redis";

function sendJson(response, statusCode, payload) {
  response.writeHead(statusCode, { "Content-Type": "application/json" });
  response.end(JSON.stringify(payload));
}

const server = http.createServer((request, response) => {
  const url = new URL(request.url, `http://${request.headers.host || "localhost"}`);

  if (url.pathname === "/api" || url.pathname === "/api/") {
    sendJson(response, 200, {
      service: "node-gateway",
      status: "ok",
      redisHost,
      path: url.pathname,
    });
    return;
  }

  if (url.pathname === "/api/health") {
    sendJson(response, 200, {
      service: "node-gateway",
      status: "healthy",
    });
    return;
  }

  sendJson(response, 404, {
    error: "not_found",
    service: "node-gateway",
    path: url.pathname,
  });
});

server.listen(port, "0.0.0.0", () => {
  console.log(`node-gateway listening on ${port}`);
});
