PROJECT DESCRIPTION

The target audience for our project is students in universities that enjoy movies and want to know what their peers or other people their age think about movies. Our audience wants an easy and fast way to find details about movies, their ratings, reviews, and other information. The target audience is also interested in finding new movies to watch from different genres and make a watchlist! Right now, there is a huge industry just for critiquing movies, however, there isn’t something specifically for college students who may have a unique perspective relative to professional critics. 

Our project allows for students to see what their peers enjoy as well as take their perspective on their favorite movies without input from other professional critics or people. We all love movies and we want to know what our peers think are great movies as their opinions may differ from what society deems “classics” and “good”. It’s cool to know what other people our age enjoy.
 
TECHNICAL DESCRIPTION
ARCHITECTURAL DIAGRAM
 
 
 
 
 
 
SUMMARY TABLE
Priority 
User
Description
P0
As a guest
Use search to find movies


P0
As a guest
Create an account
P0
As a user
Login and logout


P0
As a user
Create and write reviews


P1
As a user
Edit user-written reviews


P2
As a user
Leave comments on reviews


P2
As a user
Favorite movies
P1
As an editor
Create and delete movie pages


P1
As an editor
Edit movie information and description


P2
As an editor
Create and delete actor/director pages
P2
As an editor
Edit movie actor/director information and description
 
 


IMPLEMENTATION STRATEGY
As a guest,
I want to create an account
This requires capturing the user registration credentials and storing it in a new row in the mySQL DB USERS table and then start a new session with Redis for the new user.
I want to use search and find movies
This would require capturing the user search text and then pulling our list of movies from a mysql of movies that we have stored in the MOVIES table 
As a user,
I want to login and logout
This requires authenticating the user credentials with mySQL then afterwards, if they are authenticated, we will create a new session for them with Redis
I want to create and write a review
This requires taking the comment and then storing the UserID,  MovieID, and Comment in the mySQL REVIEW table
I want to edit user reviews:
This requires authenticating the user to see if it's the user that wrote the review with their credentials against the mySQL USERS table. Then, if verified, we will allow them to update the comment and then edit the mySQL COMMENT table. We will replace the comments entry but leave UserID and MovieID the same..
I want to favorite movies
This requires saving the movie into the mySQL FAVORITED_MOVIE table that stores the UserID and the MovieID in a row.
As a editor,
I want to create and delete movie pages
This requires editor authorization by checking mySQL USERS table
I want to edit movie page information and description
This requires editor authorization by checking mySQL USERS table
I want to 
 
I want to edit movie director and actor details
This requires editor authorization by checking mySQL USERS table and then if authorized, we can save their comments in the MOVIE table that holds Title, Desc, ActorID,  Date, and DirectorID.
ENDPOINTS
PUT /create
Guests will be able to create an account
New UserID will be created within SQL database
 
POST /login
Users will be able to log in with their credentials
 
GET /movie/{MovieID}
Any user will be able to view the movie page		
 
GET /movie/{MovieID}/reviews
Any user will be able to view the reviews for a given movie
 
PUT /movie/{MovieID}/reviews/create/{ReviewID}
Logged-in users will be able to create a user-written review for a given movie
 
PUT /movie/{MovieID}/reviews/edit/{ReviewID}
Logged-in users will be able to edit their own user-written review for a given movie
 
DEL /movie/{MovieID}/reviews/delete/{ReviewID}
Logged-in users will be able to delete their own user-written review for a given movie

PUT /movie/create/
Editor level users will be able to create movie listing pages
New MovieID will be created within SQL database
 
PUT /movie/edit/{MovieID}
Editor level users will be able to edit movie listing pages
 
DEL /movie/delete/{MovieID}
Editor level users will be able to delete movie listing pages
 
GET /name/{PersonID}
Any user will be able to view the page a person (director, producer, actor, etc).
 
APPENDIX


For ERD and Pictures: 
https://docs.google.com/document/d/1izy4wkWV9OE95jnoSNXvE7fvJqGpKXupAmpwStWXP-c/edit#



