# Gamification Api

## Фейковая авторизация 
В сервисе доступна фейковая авторизация, если включить параметр в конфиге **`fake_auth_enabled: true`**

Фейковая авторизация работает путем подмены `claims`, информацию для которых она берет из заголовков запроса **`X-Auth-User-ID X-Auth-Role`**

Пример:
```
{
  "X-Auth-User-ID": "1",
  "X-Auth-Role": "super_admin"
}
```

Для работы с модификацией заголовков в хроме можно использовать [Mod Header](https://modheader.com/docs?product=ModHeader)

P.S. Лучше использовать ID пользователя, который есть в бд, иначе не гарантируется корректное выполнение всех запросов.
