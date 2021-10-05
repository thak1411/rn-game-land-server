# rn-game-land-serer

건모의 보드게임 서버

golang 공부용 토이 프로젝트

---

# Documentation

## Architecture

~~~
-config  : Config, Option, Password, ... Private Constants
-handler : End Point Handler
-model   : Basement Business Object
-router  : Rest API Router
-usecase : Business Logic
~~~

위 아키텍처에서 DB같은 저장소를 관리하기 위한 레이어가 추가될 수도 있음 -> usecase 하위 레이어로 삽입