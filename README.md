Manually Setup
1. Clone the repository at here
2. cd admin-dashboard
3. Copy 2 json file and paste it to goravel/storage folder
4. Copy .env and paste it to goravel folder
5. For Goravel (go to goravel folder)
   Run:
      1.  go run . artisan key:generate
      2.  go run . artisan migrate
      3.  go run . artisan db:seed
      4.  go run main.go
6. For React (go to my-app folder)
   Run:
      1. npm install
      2. npm run dev
8. After successfully, click http://localhost:5173/ to run the website


















Docker Setup

1. Clone the repository at here
2. cd admin-dashboard
3. Copy 2 json file and paste it to goravel/storage folder
4. Copy .env.docker and paste it to goravel folder
5. Run:
   1. docker-compose up --build
6. After successfully, click http://localhost:5173/ to run the website


-------
Sample dummy data 

User
- Manager (manager@gmail.com): Super Admin
- Admin (admin@gmail.com): Reader
- Developer (developer@gmail.com): Developer
- Client (client@gmail.com)
  
-> all user password is "qwer1234"

----------
Roles and Permissions

Super Admin:
- View Roles: Allowed to GET on /roles (View the roles list)
- Create Role: Allowed to POST to /roles (Create a new role)
- Edit Role: Allowed to POST to /roles/:id (Edit an existing role)

Developer:
- View Roles: Allowed to GET on /roles (View the roles list)
- Create Role: Allowed to POST to /roles (Create a new role)

Reader:
- View Roles: Allowed to GET on /roles (View the roles list)
