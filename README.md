## «Утилита для управления ПВЗ»
В этом проекте реализована консольная утилита для менеджера ПВЗ.

Программа обладает командой help, благодаря которой можно получить список доступных команд с кратким описанием.

Список команд:

1. **Принять заказ от курьера**
   На вход принимается ID заказа, ID получателя и срок хранения. Заказ нельзя принять дважды. Если срок хранения в прошлом, приложение должно выдать ошибку. Список принятых заказов необходимо сохранять в файл. Формат файла остается на выбор автора.
2. **Вернуть заказ курьеру**
   На вход принимается ID заказа. Метод должен удалять заказ из вашего файла. Можно вернуть только те заказы, у которых вышел срок хранения и если заказы не были выданы клиенту.
3. **Выдать заказ клиенту**
   На вход принимается список ID заказов. Можно выдавать только те заказы, которые были приняты от курьера и чей срок хранения меньше текущей даты. Все ID заказов должны принадлежать только одному клиенту.
4. **Получить список заказов**
   На вход принимается ID пользователя как обязательный параметр и опциональные параметры.
   Параметры позволяют получать только последние N заказов или заказы клиента, находящиеся в нашем ПВЗ.
5. **Принять возврат от клиента**
   На вход принимается ID пользователя и ID заказа. Заказ может быть возвращен в течение двух дней с момента выдачи. Также необходимо проверить, что заказ выдавался с нашего ПВЗ.
6. **Получить список возвратов**
   Метод должен выдавать список пагинированно.

## Анализ запросов БД

Добавил два индекса:
```
CREATE INDEX idx_orders_recipient ON orders (recipient);
CREATE INDEX idx_orders_status ON orders (status);
```

Запрос:
`SELECT id, recipient, status, time_limit, delivered_at, returned_at FROM orders WHERE recipient = 2`
Результаты до создания индекса индекса:
```
"QUERY PLAN"
"Seq Scan on orders  (cost=0.00..21.00 rows=4 width=64) (actual time=0.009..0.010 rows=10 loops=1)"
"  Filter: (recipient = 2)"
"  Rows Removed by Filter: 10"
"Planning Time: 0.042 ms"
"Execution Time: 0.021 ms"
```
Результаты после создания индекса:
```
"QUERY PLAN"
"Seq Scan on orders  (cost=0.00..1.25 rows=1 width=64) (actual time=0.005..0.006 rows=10 loops=1)"
"  Filter: (recipient = 2)"
"  Rows Removed by Filter: 10"
"Planning Time: 0.108 ms"
"Execution Time: 0.013 ms"
```

Execution Time и actual time меньше

Запрос:
`SELECT id, recipient, status, time_limit, delivered_at, returned_at FROM orders WHERE status = 'delivered'`

Результат до создания индекса:
```
"QUERY PLAN"
"Seq Scan on orders  (cost=0.00..1.25 rows=1 width=64) (actual time=0.012..0.014 rows=4 loops=1)"
"  Filter: (status = 'delivered'::text)"
"  Rows Removed by Filter: 16"
"Planning Time: 0.080 ms"
"Execution Time: 0.031 ms"
```

Результат после создания индекса:
```
"QUERY PLAN"
"Seq Scan on orders  (cost=0.00..1.25 rows=1 width=64) (actual time=0.006..0.007 rows=4 loops=1)"
"  Filter: (status = 'delivered'::text)"
"  Rows Removed by Filter: 16"
"Planning Time: 0.103 ms"
"Execution Time: 0.015 ms"
```

Execution Time и actual time меньше