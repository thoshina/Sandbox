実行コマンドメモ

- Welcome to the HomePage! と表示
curl http://127.0.0.1:8081/

- Articles（固定）を表示
curl http://127.0.0.1:8081/articles

- Articles（固定）をデータベースに書き込み
curl http://127.0.0.1:8081/write

- データベースに保存されている内容を表示
  - すべて
curl http://127.0.0.1:8081/fetch
  - ID指定
curl -X post -H "Content-type: application/json" -d '{"id":"1"}' http://127.0.0.1:8081/fetch

- データベースに保存されている内容を削除して削除件数を表示
  - すべて
curl http://127.0.0.1:8081/delete
  - ID指定
curl -X post -H "Content-type: application/json" -d '{"id":"1"}' http://127.0.0.1:8081/delete

- データベースに指定した内容を書き込み
curl -X post -H "Content-type: application/json" -d '{"Title":"Greeting","Description":"Konichiwa!","Content":"HELLO"}' http://127.0.0.1:8081/postart
