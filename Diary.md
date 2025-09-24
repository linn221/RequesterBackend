# About
This is the file that will contains my instructions. I will be adding this in the context, and will point to the specific header you are supposed to read from. And after finishing the task, I want the read portion to be placed under Archives header, at the top most place, a h2 header with the current date and time via command `date '+%B %d %Y (%I:%M %p)'`
If this is the first time reading this file, store the significane of this file in the project memory or in the context file or something else.
## Tips
- When you need further context, you can look through the archives, most recent tasks.
## Rules
- Do not make any modifications to this file apart from moving the read portion under Archives and labeling it.
- If I add this file in context without providing a sentence, a dot, you are to do the task below.

# Tasks

# Archives
## September 24 2025 (11:48 PM)
## Update DTO or the spec to show tags
I have received no errors related to gorm. Now, you must respond tags (TagDTO) in Endpoint,Program,Note, etc (both listing and detail). I believe the spec file is already updated (include the fields) but could be wrong.
## September 24 2025 (11:31 PM)
## Polymorphic many-to-many relationships of tags
I have seen you have been doing the association of Tag with Endpoint,Program,Note, wrong. Let me make it clear. You should use a join table taggable with taggable_type(enum of endpoints, programs,notes,...) and taggable_id, and tag_id. Make the association work for gorm so that I can simply do .Preload("Tags") on Endpoint,Note, etc. In total, there will be 2 tables: tags, taggable. make sure it works with gorm because it can be difficult with its conventions
## September 24 2025 (10:57 PM)
## Logging Middleware
Create middleware that will log the request IP address, latency in millisecond, url, method, status code, etc by sticking to the standard library. I believe the conventional approach is to create a custom type for Request embedding http.Request, and modifying the http Writer interface or something like that.
## July 2 2025 (12:05 AM)
## Use sqlite for db
Use sqlite for database instead of MySQL, the file should be named app.db. Do not delete the existing code of connecting to MySQL database, instead, move it in a function, then create a new function for connecting to sqlite, and call it in the app. If I need to use MySQL, I should be able to do so by simply calling a different function.
## July 2 2025 (11:55 PM)
## Tag Further Implementations
Add the routes in routes.go. ApplyTag method should get the referenceId, referenceType from request PathValue function instead of manually spliting strings. ({referencId})
Add the model in migration function.
I have seen that you have not updated the service files of Note,Endpoint,Program,etc. You must preload the model to resolve Tags, as I have seen you using in types.go.
## July 2 2025 (11:45 PM)
## Adding Tag Resource
I want to add Tag resource, that have just a name, id and priority, and it will have many to many polymorhpic relationships with Request, Endpoint, Program, Vuln, Note. Create the model for me and add the association in existing models. Then, add endpoints for creating a new tag, and renaming it. (no delete), along with listing tags. do POST /apply-tags/{tagId}/{referenceType}/{referenceId} for adding a tag to a reference type(Request, Endpoint, Program, etc). Create a DTO for tag containing just id and name. Each reference type's both Detail and Listing DTO will have array of the tag DTO, and input DTO of these types will have tag_ids field(int array). TagService will have a method for connecting with a reference type, and reuse it in Creating Endpoint,Program,Note after the respective service has completed creating. Do not make modifications to existing models that I have not specifically said. Update the openapi spec file.
## July 2 2025 (10:10 PM)
## the header
i am getting invalid image id when going to the url of the image. i think the path images conflict with the api endpoint. use the same directory when you store the attachement files, then you no longer have to worry about the conflict. 
the newly genearated openapi spec portion does not fully match with the existing ones. (missing enum of reference_type, for example). fix the newly added ones
Add Image DTO array when returning the detail of Request, Endpoint, Program, etc. Just like Request detail DTO will return array of Notes and Attachments associated with it. Refactor both the spec and the code.

## July 1 2025 (10:10 PM)
## Active
the listing vuln dto must have parent name. add it.
modify parentId validation to check only if parentId is greater than 0, and makes sure if id is not equal to parentId (to prevent existing vuln to update as its own parent)
