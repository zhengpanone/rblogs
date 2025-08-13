const http = require('http');
const fs = require('fs');
const path = require('path');
const url = require('url');

const port = 3000;
const publicDir = path.join(__dirname, '/');

// 定义MIME类型映射
const mimeTypes = {
  ".html": "text/html",
  ".js": "application/javascript",
  ".wasm": "application/wasm",
  ".css": "text/css",
  ".json": "application/json",
  ".txt": "text/plain",
  ".png": "image/png",
  ".jpg": "image/jpeg",
};

const server = http.createServer((req, res) => {
  // 解析请求路径
  const parsedUrl = url.parse(req.url);
  let pathname = path.normalize(parsedUrl.pathname);

  // 处理根路径
  if(pathname === '/'){
    pathname = '/index.html';
  }
  // 构造完整文件路径
  const filePath = path.join(publicDir, pathname);

  // 安全检查：防止目录遍历
  if(!filePath.startsWith(publicDir)){
    res.writeHead(403, { "Content-Type": "text/plain" });
    res.end("Forbidden");
    return;
  }

  // 获取文件扩展名并确定Content-Type
  const ext = path.extname(filePath);
  const contentType = mimeTypes[ext] || "application/octet-stream";

  fs.readFile(filePath, (err, data) => {
    if(err){
      if(err.code === 'ENOENT'){
        // 文件不存在
        res.writeHead(404, { "Content-Type": "text/plain" });
        res.end("Not Found");
        return;
      }else{
        // 其他错误
        res.writeHead(500, { "Content-Type": "text/plain" });
        res.end("Internal Server Error");
        return;
      }
    }else{
      res.writeHead(200, { "Content-Type": contentType });
      res.end(data);
    }
  });
});

server.listen(port, () => {
  console.log(`Server is listening on http://localhost:${port}`);
});
