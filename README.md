# comment-system
Система для добавления и чтения постов и комментариев аналогичная Reddit или Хабр

## Пример запросов

Получение постов
```
curl \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{ "query": "{posts {id, name}}" }' \
  http://localhost:8080/query
```
```
curl -X POST http://localhost:8080/query \
-H "Content-Type: application/json" \
-d '{
  "query": "mutation {createPost(input: {userId: 1, commentsEnabled: true, name: New Post , text: Hello, world!}}",
  "variables": {}
}'
```

```
curl -X POST http://localhost:8080/query \
-H "Content-Type: application/json" \
-d '{"query":"{comments(postId: 1) {id, text}}"}'
```

```
curl -X POST http://localhost:8080/query \
-H "Content-Type: application/json" \
-d '{
  "query": "mutation {createComment(input: {parentId: null, userId: 1, postId: 1, text: \"Great post\"}){id, text}}",
  "variables": {}
}'
```