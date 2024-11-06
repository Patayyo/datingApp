# Dating APP


If u wanna launch this project, use next steps:
1. clone this repository, to do this use the following command: git clone https://github.com/Patayyo/datingApp.git
2. To deploy a project via Docker use the following command: Docker-compose up --build


CURL Requests:
1. for healtcheck curl -X GET http://localhost:8080(or your host)/ping
2. for register curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\", \"username\":\"test\", \"password\":\"pass123\"}" you can use any email, username and password
3. for login curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\", \"password\":\"pass123\"}" you can use any email, username and password
4. for check user profile curl -X GET http://localhost:8080/user/profile -H "Authorization: " you need to use the token received after authorization
5. for check matches curl -X GET http://localhost:8080/match/matches -H "Authorization: " you need to use the token received after authorization
6. for update user profile curl -X PUT http://localhost:8080/user/profile -H "Content-Type: application/json" -H "Authorization: " -d "{\"gender\":\"male\"}"(for example) You can use any of the suggested themes: age, username, gender, location, interests, bio
7. for match curl -X POST http://localhost:8080/match/like/2 -H "Authorization: " you need to use the token received after authorization
