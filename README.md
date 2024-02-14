# aTES
[Ссылка на схему в Miro](https://miro.com/app/board/uXjVNtPkTMM=/?share_link_id=710319574736)

## Сервисы
- Task (создание задач, ассайн)
- Accounting (калькуляция заработка)

Не стал выделять сервис аналитики, так как по сути там ничего не происходит, данные запрашиваются из аккаунтинга

## Коммуникации
### Async
* Task Created (Task -> Accounting) - при создании задачи (стриминг данных задачи)
* Task Completed (Task -> Accounting) - при выполнении задачи 
* Task Assigned (Task -> Accounting) - при назначении исполнителя на задачу