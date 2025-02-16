# Merch Store

## Запуск

1. Установите Docker и Docker Compose.
2. Выполните команду `docker-compose up --build`.
3. Сервис будет доступен по адресу `http://localhost:8080`.

## WARNING
Бд поднимается автоматически по миграциям, нужно лишь в файле .env поменять логин и пароль от PostgrSQL

### Тесты E2E
Что проверяется в тестах:
    TestAuth: Проверяет, что пользователь может зарегистрироваться и получить JWT-токен.
    TestSendCoins: Проверяет, что пользователь может отправить монеты другому пользователю, и балансы обновляются корректно.
    TestBuyMerch: Проверяет, что пользователь может купить товар, и инвентарь обновляется правильно.


При проверке тестового задания, прошу учесть, что на Go я пишу впервые, ни в универе, ни самостоятельно я его не изучал. Всё было сделано в рамках одной недели. Выбрал его осознанно, зная что он является предпочтительным. Этим я хотел показать свое желание попасть на стажировку, ведь она будет на Go.
Основной мой стек: C#, .Net, ASP.NET



Не знаю как приложить тестирование, поэтому пускай прямо тут будет (на предупреждение 
"2025/02/16 18:19:41 C:/avito/merch-store/internal/repository/user_repository.go:35 record not found
[1.922ms] [rows:0] SELECT * FROM "users" WHERE username = 'user2' ORDER BY "users"."id" LIMIT 1 "  смотреть не надо т.к. всё работает[это как исключение при отсутствии записи в бд - всё пашет], просто не знаю как поправить это):


PS C:\avito\merch-store> go test ./test/e2e -v
=== RUN   TestAuth
2025/02/16 18:19:40 Database merch_test does not exist, creating...
2025/02/16 18:19:40 Database merch_test created successfully
2025/02/16 18:19:41 Connected to the database successfully
2025/02/16 18:19:41 Migrations applied successfully
2025/02/16 18:19:41 Authenticating or creating user: testuser
2025/02/16 18:19:41 Fetching user by username: testuser      

2025/02/16 18:19:41 C:/avito/merch-store/internal/repository/user_repository.go:35 record not found
[3.126ms] [rows:0] SELECT * FROM "users" WHERE username = 'testuser' ORDER BY "users"."id" LIMIT 1
2025/02/16 18:19:41 User testuser not found
2025/02/16 18:19:41 User testuser not found, creating...
2025/02/16 18:19:41 Creating new user: &{ID:0 Username:testuser Password:$2a$10$00Yo5ZH/hdOoJcyrt0qyyutLxftjEQfX0A3C1DJsIYkLWTUifjCo6 Coins:1000 CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC}
2025/02/16 18:19:41 User created successfully: &{ID:1 Username:testuser Password:$2a$10$00Yo5ZH/hdOoJcyrt0qyyutLxftjEQfX0A3C1DJsIYkLWTUifjCo6 Coins:1000 CreatedAt:2025-02-16 18:19:41.2408584 +0300 MSK UpdatedAt:2025-02-16 18:19:41.2408584 +0300 MSK}
2025/02/16 18:19:41 User testuser created successfully with ID: 1
--- PASS: TestAuth (0.87s)
=== RUN   TestSendCoins
2025/02/16 18:19:41 Connected to the database successfully
2025/02/16 18:19:41 Database is already migrated
2025/02/16 18:19:41 Authenticating or creating user: user1
2025/02/16 18:19:41 Fetching user by username: user1

2025/02/16 18:19:41 C:/avito/merch-store/internal/repository/user_repository.go:35 record not found
[3.719ms] [rows:0] SELECT * FROM "users" WHERE username = 'user1' ORDER BY "users"."id" LIMIT 1
2025/02/16 18:19:41 User user1 not found
2025/02/16 18:19:41 User user1 not found, creating...
2025/02/16 18:19:41 Creating new user: &{ID:0 Username:user1 Password:$2a$10$QbkjKoPFQQ28jwhlZ6FQlOcHLwxekLJ7dUkm3GK8AUyriwa552eDy Coins:1000 CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC}
2025/02/16 18:19:41 User created successfully: &{ID:2 Username:user1 Password:$2a$10$QbkjKoPFQQ28jwhlZ6FQlOcHLwxekLJ7dUkm3GK8AUyriwa552eDy Coins:1000 CreatedAt:2025-02-16 
18:19:41.6094527 +0300 MSK UpdatedAt:2025-02-16 18:19:41.6094527 +0300 MSK}
2025/02/16 18:19:41 User user1 created successfully with ID: 2
2025/02/16 18:19:41 Authenticating or creating user: user2
2025/02/16 18:19:41 Fetching user by username: user2

2025/02/16 18:19:41 C:/avito/merch-store/internal/repository/user_repository.go:35 record not found
[1.922ms] [rows:0] SELECT * FROM "users" WHERE username = 'user2' ORDER BY "users"."id" LIMIT 1
2025/02/16 18:19:41 User user2 not found
2025/02/16 18:19:41 User user2 not found, creating...
2025/02/16 18:19:41 Creating new user: &{ID:0 Username:user2 Password:$2a$10$lvKT.inbOKbdbyiZ9Hpoje48r3IWhGp5nq4S1K9Znp.Q2h3wVjn1G Coins:1000 CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC}
2025/02/16 18:19:41 User created successfully: &{ID:3 Username:user2 Password:$2a$10$lvKT.inbOKbdbyiZ9Hpoje48r3IWhGp5nq4S1K9Znp.Q2h3wVjn1G Coins:1000 CreatedAt:2025-02-16 
18:19:41.7521461 +0300 MSK UpdatedAt:2025-02-16 18:19:41.7521461 +0300 MSK}
2025/02/16 18:19:41 User user2 created successfully with ID: 3
2025/02/16 18:19:41 Authenticating or creating user: user1
2025/02/16 18:19:41 Fetching user by username: user1
2025/02/16 18:19:41 Getting user ID for username: user2
2025/02/16 18:19:41 Fetching user by username: user2
2025/02/16 18:19:41 User ID for username user2 is 3
2025/02/16 18:19:41 Getting user balance for userID: 2
2025/02/16 18:19:41 User balance for userID: 2 is 1000 coins
2025/02/16 18:19:41 Processing GetUserInfo request for userID: 2
2025/02/16 18:19:41 Getting user info for userID: 2
2025/02/16 18:19:41 Getting user balance for userID: 2
2025/02/16 18:19:41 User balance for userID: 2 is 900 coins
2025/02/16 18:19:42 User info for userID: 2 is map[coinHistory:map[received:[] sent:[map[amount:100 toUser:user2]]] coins:900 inventory:[]]
2025/02/16 18:19:42 GetUserInfo request processed successfully for userID: 2. Info: map[coinHistory:map[received:[] sent:[map[amount:100 toUser:user2]]] coins:900 inventory:[]]
--- PASS: TestSendCoins (0.80s)
=== RUN   TestBuyMerch
2025/02/16 18:19:42 Connected to the database successfully
2025/02/16 18:19:42 Database is already migrated
2025/02/16 18:19:42 Authenticating or creating user: testuser
2025/02/16 18:19:42 Fetching user by username: testuser
2025/02/16 18:19:42 Processing BuyMerch request for userID: 1, itemName: t-shirt
2025/02/16 18:19:42 Starting BuyMerch for userID: 1, itemName: t-shirt
2025/02/16 18:19:42 Getting item price for itemName: t-shirt
2025/02/16 18:19:42 Item price for t-shirt is 80 coins
2025/02/16 18:19:42 Item price for t-shirt is 80 coins
2025/02/16 18:19:42 Getting user balance for userID: 1
2025/02/16 18:19:42 User balance for userID: 1 is 1000 coins
2025/02/16 18:19:42 User balance for userID: 1 is 1000 coins
2025/02/16 18:19:42 Deducted 80 coins from userID: 1
2025/02/16 18:19:42 Logged purchase of t-shirt for userID: 1
2025/02/16 18:19:42 Committing transaction
2025/02/16 18:19:42 Processing GetUserInfo request for userID: 1
2025/02/16 18:19:42 Getting user info for userID: 1
2025/02/16 18:19:42 Getting user balance for userID: 1
2025/02/16 18:19:42 User balance for userID: 1 is 920 coins
2025/02/16 18:19:42 User info for userID: 1 is map[coinHistory:map[received:[] sent:[]] coins:920 inventory:[map[t-shirt:1]]]
2025/02/16 18:19:42 GetUserInfo request processed successfully for userID: 1. Info: map[coinHistory:map[received:[] sent:[]] coins:920 inventory:[map[t-shirt:1]]]
    e2e_test.go:206:
                Error Trace:    C:/avito/merch-store/test/e2e/e2e_test.go:206
                Error:          Not equal:
                                expected: string("t-shirt")
                                actual  : <nil>(<nil>)
                Error:          Not equal:
                                expected: string("t-shirt")
                                actual  : <nil>(<nil>)
                Test:           TestBuyMerch
    e2e_test.go:207:
                Error Trace:    C:/avito/merch-store/test/e2e/e2e_test.go:207
                Error:          Not equal:
                                expected: float64(1)
                                actual  : <nil>(<nil>)
                Test:           TestBuyMerch
--- FAIL: TestBuyMerch (0.51s)
FAIL
FAIL    github.com/KsenoTech/merch-store/test/e2e       2.358s
FAIL
PS C:\avito\merch-store>
                Error:          Not equal:
                                expected: string("t-shirt")
                                actual  : <nil>(<nil>)
                Test:           TestBuyMerch
    e2e_test.go:207:
                Error Trace:    C:/avito/merch-store/test/e2e/e2e_test.go:207
                Error:          Not equal:
                                expected: float64(1)
                                actual  : <nil>(<nil>)
                Test:           TestBuyMerch
--- FAIL: TestBuyMerch (0.51s)
FAIL
FAIL    github.com/KsenoTech/merch-store/test/e2e       2.358s
                Error:          Not equal:
                                expected: string("t-shirt")
                                actual  : <nil>(<nil>)
                Test:           TestBuyMerch
    e2e_test.go:207:
                Error Trace:    C:/avito/merch-store/test/e2e/e2e_test.go:207
                Error:          Not equal:
                                expected: float64(1)
                                actual  : <nil>(<nil>)
                Test:           TestBuyMerch
--- FAIL: TestBuyMerch (0.51s)
                Error:          Not equal:
                                expected: string("t-shirt")
                                actual  : <nil>(<nil>)
                Test:           TestBuyMerch
    e2e_test.go:207:
                Error Trace:    C:/avito/merch-store/test/e2e/e2e_test.go:207
                Error:          Not equal:
                                expected: float64(1)
                                actual  : <nil>(<nil>)
                Error:          Not equal:
                                expected: string("t-shirt")
                                actual  : <nil>(<nil>)
                Test:           TestBuyMerch
    e2e_test.go:207:
                Error Trace:    C:/avito/merch-store/test/e2e/e2e_test.go:207
                Error:          Not equal:
                                expected: float64(1)
                Test:           TestBuyMerch
    e2e_test.go:207:
                Error Trace:    C:/avito/merch-store/test/e2e/e2e_test.go:207
    e2e_test.go:207:
                Error Trace:    C:/avito/merch-store/test/e2e/e2e_test.go:207
                Error Trace:    C:/avito/merch-store/test/e2e/e2e_test.go:207
                Error:          Not equal:
                                expected: float64(1)
                                actual  : <nil>(<nil>)
                Test:           TestBuyMerch
--- FAIL: TestBuyMerch (0.51s)
FAIL
FAIL    github.com/KsenoTech/merch-store/test/e2e       2.358s
FAIL
PS C:\avito\merch-store>
