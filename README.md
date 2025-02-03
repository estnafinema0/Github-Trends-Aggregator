npm install -g wscat

wscat -c ws://localhost:8080/ws

Как запускать клиента:

npm install npm start

(если будет ругаться, то нужно export NODE_OPTIONS=--openssl-legacy-provider, чтобы разрешить есть старые версии openssl)

Далее он сразу перекинет на порт 3000, где будет запущен красивый клиент.