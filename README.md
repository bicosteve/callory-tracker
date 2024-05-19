# Callorytracker Application.

The application provides **CRUD** functionality for users to be able to track their callory intake.
It allos users to create, retrieve, update and delete the created records.

## 1. Prerequisites

To be able to run this application in your local marchine, you will need the following;

- Golang
- MySQL

## 2. Installation

1. Clone the repository to local marchine:

- git clone git@github.com:bicosteve/callory-tracker.git

2. Navigate to the project dir

- cd callory-tracker

3. Install dependancies

- go mod download
- go mod verify

4. Run the application
   I have provided an environment variables examples on .env.example.
   Use the tables dir to find the scripts of created tables for the application.
   After creating the tables and providing the envrionment vars;
   run **make run** this will start the server on your localhost

## 3. Usage

Once the application is running your provided port, access the home page by going to http://localhost:{port}/
These are the end points on the application:

- "/food/add" -> loading 'add food' page and submitting the page
- "/food/day" -> generate the consumption for the day
- "/food/edit" -> load the edit form and submit the form
- "/food/delete" -> remove specific 'food sample'
- "/food/total" -> calculate the total calorry for the day.
- "/user/me" -> get users details
- "/" -> gets home page
- "/user/register" -> load and submit register form
- "/user/login" -> load and submit login form
- "/user/logout" -> destroy user session
- router.NotFound(app.pageNotFound) -> 404 not found
