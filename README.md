
# Test task BackDev

Тестовое задание на позицию Junior Backend Developer


## How to run

Copy my project, then run:

```bash
  docker-compose up --build
```


## API

#### Register User

```http
  POST /api/user/register
```

```json
{
    "login": "string",
    "password": "string"
}
```

#### Sign In to get Access and Refresh Tokens

Первый маршрут выдает пару Access, Refresh токенов при входе в систему с логином и паролем. Проверяем пользователя в базе данных и с помощью идентификатора (GUID) указанным в DB вызываем функцию ```GenerateTokenPair(GUID uint, ip string) (map[string]string, error)``` в ```handler/jwt.go```

```http
  POST /api/user/sign_in
```
```json
{
    "login": "string",
    "password": "string"
}
```

#### Refresh Token Pair

Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

```http
  POST /api/user/refresh_tokens
```
```json
{
    "refresh_token": "string"
}
```
