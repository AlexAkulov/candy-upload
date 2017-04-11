Candy Upload
===

Very simple backend for upload and processing files


upload via curl
```
curl -F 'uploadfile=@rabbitmq-server-3.4.3-1.noarch.rpm' "http://example.com:9090/pkg/CentOS/7/Packages"
```

upload.html
```
<html>
  <head><title>Upload file</title></head>
  <h1>pkg/CentOS/7</h1>
  <body>
    <form enctype="multipart/form-data" action="../Packages" method="post">
      <input type="file" name="uploadfile" />
      <input type="submit" value="upload" />
    </form>
  </body>
</html>

```


nginx config
```
server {
    listen 80;
    client_max_body_size 256m;
    location / {
        if ($request_method = POST) {
                proxy_pass http://localhost:9090;
        }
        root /srv/www/repo;
        autoindex on;
    }
}

```

