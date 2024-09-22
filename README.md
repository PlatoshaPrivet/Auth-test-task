# Auth-test-task

Это программа предоставляет реализацию простого сервиса аутентификации. Она включает два маршрута:

## Маршруты

### 1. Создание токенов

Этот маршрут отвечает за создание пары токенов: access_token и refresh_token. 
- **Access Token**: JWT токен, который содержит информацию о клиенте (включая IP адрес) и используется для аутентификации в дальнейшем. Прогорает через 15 минут после создания.
- **Refresh Token**: Токен, закодированный в формате base64, который позволяет получать новый access и refresh токены. Хэшированный вариант refresh токена сохраняется в базе данных для обеспечения безопасности.

### 2. Операция Refresh

Этот маршрут отвечает за обновление access и refresh токена на основе предоставленного refresh токена. При этом происходит проверка:
- Если IP адрес клиента отличается от того, который был указан при создании access токена, программа отправляет предупреждающее сообщение на указанный email адрес.

## Технологии

- **Go**
- **JWT**
- **PostgreSQL**
- **Docker**
