# api


## POST /api/users/:id/tasks

parameter | required |description
--- | --- | ---
task | true | raw text in tweet format
mentions | false | list of user names or lists
hashtags | false | list of hash tags
urls | false | list of urls which will be shortened on the server side 

errors | description
--- | ---
400 | parameters
404 | user not found

## POST /api/users/:id/follows

##

