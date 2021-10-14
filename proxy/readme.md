# nginx proxy setting

brew로 설치했을시 `/usr/local/etc/nginx/`에 nginx.conf가 존재합니다.

nginx로 서버를 배포했을때 403에러가 발생하는 이유는 파일에 접근 가능한 권한이 부족해서 생기는 문제다. 프로젝트 폴더에 `httpd_sys_rw_content_t`권한을 부여하면 해결된다.

이때 모든 경로에 접근 가능한 권한이 설정되어있어야된다. `/a/b/c` 라면 `/a`, `/a/b`, `/a/b/c`에 모두 권한이 있어야 한다.