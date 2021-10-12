# rn-game-land-serer

건모의 보드게임 서버

golang 공부용 토이 프로젝트

[동시 개발중인 웹 클라이언트](https://github.com/thak1411/rn-game-land-web)

---

# Documentation

## Architecture

~~~
-config     : Config, Option, Password, ... Private Constants
-database   : DB Handler
-handler    : End Point Handler
-middleware : Preprocessor in End Point Handler
-model      : Basement Business Object
-router     : Rest API Router
-temp       : Legacy Test Code
-usecase    : Business Logic
-util       : Util Functions
~~~

## Custom Http Status Code

| Code | Message |
| ---- | ------- |
| 909  | 중복된 이름, 회원 가입 실패 |
| 910  | 유저 생성 성공 |