Backend:
    - API (Amazon ECS (container service))
    - Amazon RDS (MySQL)

Frontend:
    - WebUI (Firebase/Vercel/Netlify/GHPages)

Example cURL request:
    curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"username":"xyz","password":"xyz"}' \
  http://localhost:8080/register