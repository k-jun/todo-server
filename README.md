# todo-server

## Run

```bash
go run .

curl -XPOST -d '{"id": "d4f97b8f-a1dd-4806-8aa7-46fd9499f96e","title": "hello world", "done": false}' localhost:8080/todos
curl -XGET localhost:8080/todos
curl -XGET localhost:8080/todos/d4f97b8f-a1dd-4806-8aa7-46fd9499f96e
curl -XPUT -d '{"id": "d4f97b8f-a1dd-4806-8aa7-46fd9499f96e","title": "hello world", "done": true}' localhost:8080/todos/d4f97b8f-a1dd-4806-8aa7-46fd9499f96e
curl -XDELETE localhost:8080/todos/d4f97b8f-a1dd-4806-8aa7-46fd9499f96e
```

## Deploy

```bash
export GCP_PROJECT='' SERVICE_NAME=''
gcloud --project $GCP_PROJECT run deploy $SERVICE_NAME --source ./
```

