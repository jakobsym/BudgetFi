Backend:
    - API (Amazon ECS (container service) || digital ocean?)
    - DB: Any repo, currently using Microsoft SQL server

Frontend:
    - WebUI (Firebase/Vercel/Netlify/GHPages)

Example cURL request:
    curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"test user","email":"testuser@email.com", "google_id": "123456789"}' \
  http://localhost:8080/register


  