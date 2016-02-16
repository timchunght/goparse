# Modernplanit
	
### Development:

	go get github.com/codegangsta/gin
	gin -p 5000

This above will run the application at port 5000

### Build and Deploy:

If you wish the build the application in the traditional "go way",

	go build

To run the binary file, run:

	export PORT=8080 && ./modernplanit

For the app to work correctly, please ensure you fill out the ``.env`` with the correct environment variable. It will use the following by default if no environment variable given.

	MONGO_URL=mongodb://localhost:27017

You can change the port to the one your prefer.

TODOS:

Add ability to configure more Mongodb settings.

### Usage:

Currently, only one endpoint functions and you can make a sample request to see if everything is running

	POST   http://localhost:5000/events?name=chicago

Sample Response:

	{
	  "id": "56ae2b503d10891dced23e05",
	  "event_date": "0001-01-01T00:00:00Z",
	  "city": "",
	  "name": "chicago",
	  "country": "",
	  "weather_url": "",
	  "updated_at": "2016-01-31T10:42:08.44106907-05:00",
	  "created_at": "2016-01-31T10:42:08.44106907-05:00"
	}

### Dependency Vendoring

We are currently using ``godep`` to manage dependencies. All dependencies are tracked in ``Godeps.json`` and copied into the ``Godeps`` directory. If the project directory is changed, say from ``tim`` to ``timothy``, run the following. However, if any of the source files are referencing same-project package(s), you need to change those import paths first.

For example, I have a package in the same project named ``connection`` and my old project name is ``tim``, I will have to reference ``connection`` pacakge in my ``main.go`` as ``tim/connection``. Now the directory name is ``timothy`` and I will have to change it to ``timothy/connection`` and then run the following,

	godep save -r ./...
		
``godep`` is not smart enough yet to distinguish whether the package belongs to the same project or is an external dependency.

Go Debugging
---

If you have used Rails, you will miss using binding.pry. However, Go does have similar. Go has a package called "Godebug"

Install Godebug by running:

	go get github.com/mailgun/godebug

Insert a breakpoint anywhere in a source file you want to debug:
	
	_ = "breakpoint"

Replace ``<pkgs>`` with the packages we will be debugging if it is not the ``main`` package

	godebug build -instrument <pkgs>

For example:
	
	godebug build -instrument goparse/connection

godebug will generate a binary named ``yourprojectname.debug``, run that binary with the necessary arguments or environment variables and use it as you would binding.pry

For example,

	PORT=8080 ./modernplanit.debug

### Dealing with Mongodb

Enter the console by typing ``mongo`` in terminal

The following command allows you to rename field/column

	db.events.update({},{ $rename: { 'current_field_name': 'new_name'}}, { multi: true })

We will later write a script to execute these changes automatically.

### Associations:

ActivityGroup has_many ActivitySubgroup(s)

ActivitySubgroup has_many Activitie(s)

FoodGroup has_many FoodVenue(s)

Event has_many ActivityGroup(s)

Event has_many FoodGroup(s)

Event has_many local_info(s)

Event has_many Map(s)

Event has_many Schedule(s)

Event has_many Trivia(s)

Event has_many WelcomeMessage(s)

Event has_many Contact(s)

Event has_many DiscussionTopic(s)

Event has_many User(s)

DiscussionTopic has_many DiscussionMessage(s)

Installation (optional; used to track deviceId for push notifications)

Models
---
User
Event
Trivia
Activity
ActivityGroup
ActivitySubGroup
FoodGroup
FoodVenue
WelcomeMessage
Contact
DiscussionTopic
DiscussionMessage
Schedule
Map
LocalInfo

### Schema Details:

Implemented: Finished both testing and CRUD endpoint except Query endpoint

Completed: CRUD + Query + Testing

WelcomeMessage (Completed)
------------------------------------
Id            				: string :mandatory
EventId        				: string :mandatory
HeaderImageUrl    		: string :optional (url to image)
Message        				: string :mandatory

Trivia (Completed)
------------------------------------
Id             			  : string :mandatory
EventId        				: string :mandatory
Name        					: string :mandatory
Description    				: string :mandatory


Schedule (Completed)
------------------------------------
Id            				: string :mandatory
EventId        				: string :mandatory
Details        				: string :mandatory
EventName    					: string :mandatory
EventTime        			: string :mandatory
EventDate        			: date   :mandatory
DisplayImageUrl    		: string :mandatory (url to image)

Map (Completed)
------------------------------------
Id             			  : string :mandatory
EventId        				: string :mandatory
Name        					: string :mandatory
Description    				: string :mandatory
HasFile        				: bool 	 :mandatory
Source        				: string :optional (url to image)
Latitude        			: string :optional
Longitude        			: string :optional


Contact (Completed)
------------------------------------
Id            				: string :mandatory
EventId        				: string :mandatory
Name        					: string :mandatory
Number        				: string :mandatory
DisplayImageUrl    		: string :mandatory (url to image)

FoodGroup (Completed)
------------------------------------
Id             			  : string :mandatory
EventId        				: string :mandatory
Name        					: string :mandatory
VenueCount    				: int 	 :mandatory (default 0)

FoodVenue (Completed)
------------------------------------
Id             			  : string :mandatory
FoodGroupId    				: string :mandatory
Name        					: string :mandatory
Address        				: string :mandatory
DisplayImageUrl    		: string :mandatory (url to image)
Link        					: string :mandatory
Description   				: string :mandatory
Latitude        			: string :mandatory
Longitude        			: string :mandatory

Event (Implemented) (Query by user_id)
------------------------------------
Id             			  : string :mandatory
EventUserId						: string :mandatory
Name        					: string :mandatory
City        					: string :mandatory
Country        				: string :mandatory
MainBackgroundUrl   	: string :mandatory (url to image)
WeatherUrl    				: string :mandatory
Latitude        			: string :mandatory
Longitude        			: string :mandatory
UserIds               : []string :mandatory

DiscussionTopic (Completed)
------------------------------------
Id            			 	: string :mandatory
EventId        			 	: string :mandatory
Name        				 	: string :mandatory
Title        				 	: string :mandatory
TotalMessages    		 	: int 	 :mandatory

DiscussionMessage (Completed)
------------------------------------
Id             		    : string :mandatory
DiscussionTopicId     : string :mandatory
Name        					: string :mandatory
Message        				: string :mandatory
Avatar        				: string :mandatory

ActivityGroup (Completed)
------------------------------------
Id             			  : string :mandatory
EventId        				: string :mandatory
Name        					: string :mandatory
SubGroupCount    			: int    :mandatory (default 0)

ActivitySubGroup (Completed)
------------------------------------
Id             			  : string :mandatory
ActivityGroupId    		: string :mandatory
Name        					: string :mandatory
ActivityCount    			: int    :mandatory (default 0)

Activity (Completed)
------------------------------------
Id             			  : string :mandatory
ActivitySubGroupId    : string :mandatory
Name        					: string :mandatory
Description    				: string :mandatory
DisplayImageUrl    		: string :mandatory (url to image)

User (Completed, must encrypt password using salt, Query with event_id)
------------------------------------
Id             			  : string :mandatory
EventUserId						: string :mandatory
Username        			: string :mandatory
Password        			: string :mandatory
EventIds              : []string :mandatory

AccessInfo (Completed)
------------------------------------
Id             			  : string :mandatory
EventId        				: string :mandatory
Username        			: string :mandatory
Password        			: string :mandatory

Photo (Completed)
------------------------------------
FileName 							:string :mandatory
DownloadUrl 					:string :mandatory
UploadedBy					 	:string :mandatory
EventId 							:string :mandatory

### REST API Endpoints

Create:
	
	POST /<under_score_pluralised_class_name>

	request_body:

		{"required_param_1": "value1",
		 "required_param_2": "value2",
		 "required_param_3": "value3"
		} 

Retrieve by Id:
	
	GET /<under_score_pluralised_class_name>/:id

Update:
	
	PUT /<under_score_pluralised_class_name>/:id

		request_body:

		{"required_param_1": "value1",
		 "required_param_2": "value2",
		 "required_param_3": "value3"
		}

Destroy/Delete Record:
	
	DELETE /<under_score_pluralised_class_name>/:id

Query (TO BE IMPLEMENTED) Foreign key is the only ``query_param`` accepted

	GET /<under_score_pluralised_class_name>?<query_param>=value

### Setup the Image Conversion API 

Run the Docker Image:

	docker rm -f imageproxy && docker run --restart=always -p 80:8080 --name imageproxy -d timchunght/imageproxy -addr 0.0.0.0:8080 -cache /tmp/imageproxy -baseURL http://s3.amazonaws.com/modernplanit/

Get photo (provide ``s3_file_id`` and it will be ``http://s3.amazonaws.com/modernplanit/``):
  
  http://img.modernplanit.com/{{width}}x{{height}}/{{s3_file_id}}

### Testing

Run 
	
	go test -v ./...

-----

Parse Schema JSON

"results":[{"className":"Contact","fields":{"ACL":{"type":"ACL"},"DisplayImage":{"type":"File"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"Number":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Schedule","fields":{"ACL":{"type":"ACL"},"Details":{"type":"String"},"DisplayImage":{"type":"File"},"EventDate":{"type":"Date"},"EventDateStr":{"type":"String"},"EventId":{"type":"Pointer","targetClass":"Event"},"EventName":{"type":"String"},"EventTime":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"FoodGroup","fields":{"ACL":{"type":"ACL"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"VenueCount":{"type":"Number"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"FoodVenue","fields":{"ACL":{"type":"ACL"},"Address":{"type":"String"},"Coordinates":{"type":"GeoPoint"},"Description":{"type":"String"},"DisplayImage":{"type":"File"},"Group":{"type":"Pointer","targetClass":"FoodGroup"},"Link":{"type":"String"},"Name":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"ActivityGroup","fields":{"ACL":{"type":"ACL"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"SubGroupCount":{"type":"Number"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"_User","fields":{"ACL":{"type":"ACL"},"EventId":{"type":"Pointer","targetClass":"Event"},"authData":{"type":"Object"},"createdAt":{"type":"Date"},"email":{"type":"String"},"emailVerified":{"type":"Boolean"},"objectId":{"type":"String"},"password":{"type":"String"},"updatedAt":{"type":"Date"},"username":{"type":"String"}}},{"className":"_Role","fields":{"ACL":{"type":"ACL"},"createdAt":{"type":"Date"},"name":{"type":"String"},"objectId":{"type":"String"},"roles":{"type":"Relation","targetClass":"_Role"},"updatedAt":{"type":"Date"},"users":{"type":"Relation","targetClass":"_User"}}},{"className":"ActivitySubGroup","fields":{"ACL":{"type":"ACL"},"ActivityCount":{"type":"Number"},"Group":{"type":"Pointer","targetClass":"ActivityGroup"},"Name":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Activity","fields":{"ACL":{"type":"ACL"},"Address":{"type":"String"},"Description":{"type":"String"},"DisplayImage":{"type":"File"},"Group":{"type":"Pointer","targetClass":"ActivitySubGroup"},"Link":{"type":"String"},"Name":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Map","fields":{"ACL":{"type":"ACL"},"Coordinates":{"type":"GeoPoint"},"Description":{"type":"String"},"EventId":{"type":"Pointer","targetClass":"Event"},"HasFile":{"type":"Boolean"},"Name":{"type":"String"},"Source":{"type":"File"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"LocalInfo","fields":{"ACL":{"type":"ACL"},"Description":{"type":"String"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Trivia","fields":{"ACL":{"type":"ACL"},"Description":{"type":"String"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"DiscussionTopic","fields":{"ACL":{"type":"ACL"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"Title":{"type":"String"},"TotalMessages":{"type":"Number"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"DiscussionMessage","fields":{"ACL":{"type":"ACL"},"Avatar":{"type":"String"},"Message":{"type":"String"},"Name":{"type":"String"},"Topic":{"type":"Pointer","targetClass":"DiscussionTopic"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"_Installation","fields":{"ACL":{"type":"ACL"},"GCMSenderId":{"type":"String"},"appIdentifier":{"type":"String"},"appName":{"type":"String"},"appVersion":{"type":"String"},"badge":{"type":"Number"},"channels":{"type":"Array"},"createdAt":{"type":"Date"},"deviceToken":{"type":"String"},"deviceType":{"type":"String"},"installationId":{"type":"String"},"localeIdentifier":{"type":"String"},"objectId":{"type":"String"},"parseVersion":{"type":"String"},"pushType":{"type":"String"},"timeZone":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Event","fields":{"ACL":{"type":"ACL"},"City":{"type":"String"},"Country":{"type":"String"},"EventDate":{"type":"Date"},"EventDateStr":{"type":"String"},"Latitude":{"type":"String"},"Longitude":{"type":"String"},"MainBackground":{"type":"File"},"Name":{"type":"String"},"WeatherURL":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Welcome","fields":{"ACL":{"type":"ACL"},"EventId":{"type":"Pointer","targetClass":"Event"},"HeaderImage":{"type":"File"},"Message":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"_Session","fields":{"ACL":{"type":"ACL"},"createdAt":{"type":"Date"},"createdWith":{"type":"Object"},"expiresAt":{"type":"Date"},"installationId":{"type":"String"},"objectId":{"type":"String"},"restricted":{"type":"Boolean"},"sessionToken":{"type":"String"},"updatedAt":{"type":"Date"},"user":{"type":"Pointer","targetClass":"_User"}}},{"className":"Photo","fields":{"ACL":{"type":"ACL"},"DownloadURL":{"type":"String"},"EventId":{"type":"Pointer","targetClass":"Event"},"Filename":{"type":"String"},"UploadedBy":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}}]}
