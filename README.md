# darvazeh

Create MySQL Container:
```
docker run -d -p 3306:3306 --name pdns-mysql -v /home/mahdie/dockerv/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root  -d mysql:8.0.26
```
Create PowerDNS Container:
```
docker run -d -p 8080:53 -p 8080:53/udp -p 8081:8081 --name pdns \
  --link pdns-mysql:mysql \
  -e PDNS_gmysql_user=root \
  -e PDNS_gmysql_password=root \
  -e PDNS_gmysql_dbname=pdns \
  -e PDNS_api=yes \
  -e PDNS_api_key=gngjngogjrngoengong \
  -e PDNS_webserver=yes \
  -e PDNS_webserver_address=0.0.0.0 \
  -e PDNS_webserver_allow_from=172.0.0.0/8 \
  pschiffe/pdns-mysql\
  ```
  PowerDNS api on : <http://localhost:8081>\
  PowerDNS Swagger on : <http://localhost:8081\api\docs>\
  Swagger init changes, run in main project directory : 
  ```
  swag init -g cmd/api/main.go
  ```
  Checking records with nslookup:
  ```
 nslookup -port=PortDNS(8080) -type=TypeRecord Domain serverDNS
  ```

  Log level values for set env as LOGLEVEL : 
  ```
    - panic
    - fatal
    - error (*)
    - warn/warning
    - info(*)
    - debug
    - trace
  ```