# GoParse
	
### Development:

Make sure to use the mongo database that you want to connect to, it will default to ``localhost:27017`` if no ``MONGO_URL`` environment variable exists

	go get github.com/codegangsta/gin
	MONGO_URL=mongodb://localhost:27017/parse_sample_db gin -p 5000

This above will run the application at port 5000

### Build and Deploy:

If you wish the build the application in the traditional "go way",

	go build

To run the binary file, run:

	export PORT=8080 && ./goparse

You can change the port to the one your prefer since the application reads from ``PORT`` environment variable

### Implemented:

Create (Done, tests pending)
Retrieve (Done, tests pending)
Update (Done, tests pending)
Delete (Done, tests pending)

Queries (Not implemented)


# Objects

## Object Format

Storing data through the Parse REST API is built around a JSON encoding of the object's data. This data is schemaless, which means that you don't need to specify ahead of time what keys exist on each object. You simply set whatever key-value pairs you want, and the backend will store it.

For example, let's say you're tracking high scores for a game. A single object could contain:

```json
{
  "score": 1337,
  "playerName": "Sean Plott",
  "cheatMode": false
}
```

Keys must be alphanumeric strings. Values can be anything that can be JSON-encoded.

Each object has a class name that you can use to distinguish different sorts of data. For example, we could call the high score object a `GameScore`. We recommend that you NameYourClassesLikeThis and nameYourKeysLikeThis, just to keep your code looking pretty.

When you retrieve objects from Parse, some fields are automatically added: `createdAt`, `updatedAt`, and `objectId`. These field names are reserved, so you cannot set them yourself. The object above could look like this when retrieved:

```json
{
  "score": 1337,
  "playerName": "Sean Plott",
  "cheatMode": false,
  "createdAt": "2011-08-20T02:06:57.931Z",
  "updatedAt": "2011-08-20T02:06:57.931Z",
  "objectId": "Ed1nuqPvcm"
}
```

`createdAt` and `updatedAt` are UTC timestamps stored in ISO 8601 format with millisecond precision: `YYYY-MM-DDTHH:MM:SS.MMMZ`. `objectId` is a string unique to this class that identifies this object.

In the REST API, the class-level operations operate on a resource based on just the class name. For example, if the class name is `GameScore`, the class URL is:

```js
https://api.parse.com/1/classes/GameScore
```

Users have a special class-level url:

```js
https://api.parse.com/1/users
```

The operations specific to a single object are available a nested URL. For example, operations specific to the `GameScore` above with `objectId` equal to `Ed1nuqPvcm` would use the object URL:

```js
https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```


## Creating Objects

To create a new object on Parse, send a POST request to the class URL containing the contents of the object. For example, to create the object described above:

```bash
  curl -X POST \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"score":1337,"playerName":"Sean Plott","cheatMode":false}' \
  https://api.parse.com/1/classes/GameScore
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('POST', '/1/classes/GameScore', json.dumps({
       "score": 1337,
       "playerName": "Sean Plott",
       "cheatMode": False
     }), {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}",
       "Content-Type": "application/json"
     })
results = json.loads(connection.getresponse().read())
print results
```

When the creation is successful, the HTTP response is a `201 Created` and the `Location` header contains the object URL for the new object:

```js
Status: 201 Created
Location: https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```

The response body is a JSON object containing the `objectId` and the `createdAt` timestamp of the newly-created object:

```json
{
  "createdAt": "2011-08-20T02:06:57.931Z",
  "objectId": "Ed1nuqPvcm"
}
```

## Retrieving Objects

Once you've created an object, you can retrieve its contents by sending a GET request to the object URL returned in the location header. For example, to retrieve the object we created above:

```bash
curl -X GET \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('GET', '/1/classes/GameScore/Ed1nuqPvcm', '', {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}"
     })
result = json.loads(connection.getresponse().read())
print result
```

The response body is a JSON object containing all the user-provided fields, plus the `createdAt`, `updatedAt`, and `objectId` fields:

```json
{
  "score": 1337,
  "playerName": "Sean Plott",
  "cheatMode": false,
  "skills": [
    "pwnage",
    "flying"
  ],
  "createdAt": "2011-08-20T02:06:57.931Z",
  "updatedAt": "2011-08-20T02:06:57.931Z",
  "objectId": "Ed1nuqPvcm"
}
```

When retrieving objects that have pointers to children, you can fetch child objects by using the `include` option. For instance, to fetch the object pointed to by the "game" key:

```bash
curl -X GET \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -G \
  --data-urlencode 'include=game' \
  https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```
```python
import json,httplib,urllib
connection = httplib.HTTPSConnection('api.parse.com', 443)
params = urllib.urlencode({"include":"game"})
connection.connect()
connection.request('GET', '/1/classes/GameScore/Ed1nuqPvcm?%s' % params, '', {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}"
     })
result = json.loads(connection.getresponse().read())
print result
```

## Updating Objects

To change the data on an object that already exists, send a PUT request to the object URL. Any keys you don't specify will remain unchanged, so you can update just a subset of the object's data. For example, if we wanted to change the score field of our object:

```bash
curl -X PUT \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"score":73453}' \
  https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('PUT', '/1/classes/GameScore/Ed1nuqPvcm', json.dumps({
       "score": 73453
     }), {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}",
       "Content-Type": "application/json"
     })
result = json.loads(connection.getresponse().read())
print result
```

The response body is a JSON object containing just an `updatedAt` field with the timestamp of the update.

```json
{
  "updatedAt": "2011-08-21T18:02:52.248Z"
}
```

### Counters

To help with storing counter-type data, Parse provides the ability to atomically increment (or decrement) any number field. So, we can increment the score field like so:

```bash
curl -X PUT \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"score":{"__op":"Increment","amount":1}}' \
  https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('PUT', '/1/classes/GameScore/Ed1nuqPvcm', json.dumps({
       "score": {
         "__op": "Increment",
         "amount": 1
       }
     }), {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}",
       "Content-Type": "application/json"
     })
result = json.loads(connection.getresponse().read())
print result
```

To decrement the counter, use the `Increment` operator with a negative number:

```bash
curl -X PUT \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"score":{"__op":"Increment","amount":-1}}' \
  https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('PUT', '/1/classes/GameScore/Ed1nuqPvcm', json.dumps({
       "score": {
         "__op": "Increment",
         "amount": -1
       }
     }), {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}",
       "Content-Type": "application/json"
     })
result = json.loads(connection.getresponse().read())
print result
```

### Arrays

To help with storing array data, there are three operations that can be used to atomically change an array field:

*   `Add` appends the given array of objects to the end of an array field.
*   `AddUnique` adds only the given objects which aren't already contained in an array field to that field. The position of the insert is not guaranteed.
*   `Remove` removes all instances of each given object from an array field.

Each method takes an array of objects to add or remove in the "objects" key. For example, we can add items to the set-like "skills" field like so:

```bash
curl -X PUT \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"skills":{"__op":"AddUnique","objects":["flying","kungfu"]}}' \
  https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('PUT', '/1/classes/GameScore/Ed1nuqPvcm', json.dumps({
       "skills": {
         "__op": "AddUnique",
         "objects": [
           "flying",
           "kungfu"
         ]
       }
     }), {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}",
       "Content-Type": "application/json"
     })
result = json.loads(connection.getresponse().read())
print result
```

### Relations

 In order to update `Relation` types, Parse provides special operators to atomically add and remove objects to a relation.  So, we can add an object to a relation like so:

```bash
curl -X PUT \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"opponents":{"__op":"AddRelation","objects":[{"__type":"Pointer","className":"Player","objectId":"Vx4nudeWn"}]}}' \
  https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('PUT', '/1/classes/GameScore/Ed1nuqPvcm', json.dumps({
       "opponents": {
         "__op": "AddRelation",
         "objects": [
           {
             "__type": "Pointer",
             "className": "Player",
             "objectId": "Vx4nudeWn"
           }
         ]
       }
     }), {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}",
       "Content-Type": "application/json"
     })
result = json.loads(connection.getresponse().read())
print result
```

To remove an object from a relation, you can do:

```bash
curl -X PUT \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"opponents":{"__op":"RemoveRelation","objects":[{"__type":"Pointer","className":"Player","objectId":"Vx4nudeWn"}]}}' \
  https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('PUT', '/1/classes/GameScore/Ed1nuqPvcm', json.dumps({
       "opponents": {
         "__op": "RemoveRelation",
         "objects": [
           {
             "__type": "Pointer",
             "className": "Player",
             "objectId": "Vx4nudeWn"
           }
         ]
       }
     }), {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}",
       "Content-Type": "application/json"
     })
result = json.loads(connection.getresponse().read())
print result
```

## Deleting Objects

To delete an object from the Parse Cloud, send a DELETE request to its object URL. For example:

```bash
curl -X DELETE \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('DELETE', '/1/classes/GameScore/Ed1nuqPvcm', '', {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}"
     })
result = json.loads(connection.getresponse().read())
print result
```

You can delete a single field from an object by using the `Delete` operation:

```bash
curl -X PUT \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"opponents":{"__op":"Delete"}}' \
  https://api.parse.com/1/classes/GameScore/Ed1nuqPvcm
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('PUT', '/1/classes/GameScore/Ed1nuqPvcm', json.dumps({
       "opponents": {
         "__op": "Delete"
       }
     }), {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}",
       "Content-Type": "application/json"
     })
result = json.loads(connection.getresponse().read())
print result
```

## Batch Operations

To reduce the amount of time spent on network round trips, you can create, update, or delete up to 50 objects in one call, using the batch endpoint.

Each command in a batch has `method`, `path`, and `body` parameters that specify the HTTP command that would normally be used for that command. The commands are run in the order they are given. For example, to create a couple of `GameScore` objects:

```bash
curl -X POST \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
        "requests": [
          {
            "method": "POST",
            "path": "/1/classes/GameScore",
            "body": {
              "score": 1337,
              "playerName": "Sean Plott"
            }
          },
          {
            "method": "POST",
            "path": "/1/classes/GameScore",
            "body": {
              "score": 1338,
              "playerName": "ZeroCool"
            }
          }
        ]
      }' \
  https://api.parse.com/1/batch
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('POST', '/1/batch', json.dumps({
       "requests": [
         {
           "method": "POST",
           "path": "/1/classes/GameScore",
           "body": {
             "score": 1337,
             "playerName": "Sean Plott"
           }
         },
         {
           "method": "POST",
           "path": "/1/classes/GameScore",
           "body": {
             "score": 1338,
             "playerName": "ZeroCool"
           }
         }
       ]
     }), {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}",
       "Content-Type": "application/json"
     })
result = json.loads(connection.getresponse().read())
print result
```

The response from batch will be a list with the same number of elements as the input list. Each item in the list with be a dictionary with either the `success` or `error` field set. The value of `success` will be the normal response to the equivalent REST command:

```json
{
  "success": {
    "createdAt": "2012-06-15T16:59:11.276Z",
    "objectId": "YAfSAWwXbL"
  }
}
```

The value of `error` will be an object with a numeric `code` and `error` string:

```json
{
  "error": {
    "code": 101,
    "error": "object not found for delete"
  }
}
```

Other commands that work in a batch are `update` and `delete`.

```bash
curl -X POST \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
        "requests": [
          {
            "method": "PUT",
            "path": "/1/classes/GameScore/Ed1nuqPvcm",
            "body": {
              "score": 999999
            }
          },
          {
            "method": "DELETE",
            "path": "/1/classes/GameScore/Cpl9lrueY5"
          }
        ]
      }' \
  https://api.parse.com/1/batch
```
```python
import json,httplib
connection = httplib.HTTPSConnection('api.parse.com', 443)
connection.connect()
connection.request('POST', '/1/batch', json.dumps({
       "requests": [
         {
           "method": "PUT",
           "path": "/1/classes/GameScore/Ed1nuqPvcm",
           "body": {
             "score": 999999
           }
         },
         {
           "method": "DELETE",
           "path": "/1/classes/GameScore/Cpl9lrueY5"
         }
       ]
     }), {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}",
       "Content-Type": "application/json"
     })
result = json.loads(connection.getresponse().read())
print result
```

Note that N requests sent in a batch will still count toward your request limit as N requests.

## Data Types

So far we have only used values that can be encoded with standard JSON. The Parse mobile client libraries also support dates, geolocations, and relational data. In the REST API, these values are encoded as JSON hashes with the `__type` field set to indicate their type, so you can read or write these fields if you use the correct encoding. Overall, the following types are allowed for each field in your object:

* String
* Number
* Boolean
* Arrays
* JSON Objects
* DateTime
* File
* Pointer to another Parse Object
* Relation to another Parse Class
* Null

The `Date` type contains a field `iso` which contains a UTC timestamp stored in ISO 8601 format with millisecond precision: `YYYY-MM-DDTHH:MM:SS.MMMZ`.

```json
{
  "__type": "Date",
  "iso": "2011-08-21T18:02:52.249Z"
}
```

Dates are useful in combination with the built-in `createdAt` and `updatedAt` fields. For example, to retrieve objects created since a particular time, just encode a Date in a comparison query:

```bash
curl -X GET \
  -H "X-Parse-Application-Id: ${APPLICATION_ID}" \
  -H "X-Parse-REST-API-Key: ${REST_API_KEY}" \
  -G \
  --data-urlencode 'where={"createdAt":{"$gte":{"__type":"Date","iso":"2011-08-21T18:02:52.249Z"}}}' \
  https://api.parse.com/1/classes/GameScore
```
```python
import json,httplib,urllib
connection = httplib.HTTPSConnection('api.parse.com', 443)
params = urllib.urlencode({"where":json.dumps({
       "createdAt": {
         "$gte": {
           "__type": "Date",
           "iso": "2011-08-21T18:02:52.249Z"
         }
       }
     })})
connection.connect()
connection.request('GET', '/1/classes/GameScore?%s' % params, '', {
       "X-Parse-Application-Id": "${APPLICATION_ID}",
       "X-Parse-REST-API-Key": "${REST_API_KEY}"
     })
result = json.loads(connection.getresponse().read())
print result
```

The `Pointer` type is used when mobile code sets a `%{ParseObject}` as the value of another object. It contains the `className` and `objectId` of the referred-to value.

```json
{
  "__type": "Pointer",
  "className": "GameScore",
  "objectId": "Ed1nuqPvc"
}
```

Note that the bult-in User, Role, and Installation classes are prefixed by an underscore. For example, pointers to user objects have a `className` of `_User`. Prefixing with an underscore is forbidden for developer-defined classes and signifies the class is a special built-in.

The `Relation` type is used for many-to-many relations when the mobile uses `PFRelation` or `%{ParseRelation}` as a value.  It has a `className` that is the class name of the target objects.

```json
{
  "__type": "Relation",
  "className": "GameScore"
}
```

When querying, `Relation` objects behave like arrays of Pointers. Any operation that is valid for arrays of pointers (other than `include`) works for `Relation` objects.

We do not recommend storing large pieces of binary data like images or documents on a Parse object. Parse objects should not exceed 128 kilobytes in size. To store more, we recommend you use `File`. You may associate a [previously uploaded file](#files) using the `File` type.

```json
{
  "__type": "File",
  "name": "...profile.png"
}
```

When more data types are added, they will also be represented as hashes with a `__type` field set, so you may not use this field yourself on JSON objects. For more information about how Parse handles data, check out our documentation on [Data](#data).



### TODOS:
* Implement test suite
* Implement object query features
* Point, Relation type
* Middlewares and Authentication

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

	PORT=8080 ./goparse.debug

### Dealing with Mongodb

Enter the console by typing ``mongo`` in terminal

The following command allows you to rename field/column

	db.events.update({},{ $rename: { 'current_field_name': 'new_name'}}, { multi: true })

### Testing

Run 
	
	go test -v ./...

-----

Parse Schema JSON

"results":[{"className":"Contact","fields":{"ACL":{"type":"ACL"},"DisplayImage":{"type":"File"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"Number":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Schedule","fields":{"ACL":{"type":"ACL"},"Details":{"type":"String"},"DisplayImage":{"type":"File"},"EventDate":{"type":"Date"},"EventDateStr":{"type":"String"},"EventId":{"type":"Pointer","targetClass":"Event"},"EventName":{"type":"String"},"EventTime":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"FoodGroup","fields":{"ACL":{"type":"ACL"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"VenueCount":{"type":"Number"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"FoodVenue","fields":{"ACL":{"type":"ACL"},"Address":{"type":"String"},"Coordinates":{"type":"GeoPoint"},"Description":{"type":"String"},"DisplayImage":{"type":"File"},"Group":{"type":"Pointer","targetClass":"FoodGroup"},"Link":{"type":"String"},"Name":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"ActivityGroup","fields":{"ACL":{"type":"ACL"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"SubGroupCount":{"type":"Number"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"_User","fields":{"ACL":{"type":"ACL"},"EventId":{"type":"Pointer","targetClass":"Event"},"authData":{"type":"Object"},"createdAt":{"type":"Date"},"email":{"type":"String"},"emailVerified":{"type":"Boolean"},"objectId":{"type":"String"},"password":{"type":"String"},"updatedAt":{"type":"Date"},"username":{"type":"String"}}},{"className":"_Role","fields":{"ACL":{"type":"ACL"},"createdAt":{"type":"Date"},"name":{"type":"String"},"objectId":{"type":"String"},"roles":{"type":"Relation","targetClass":"_Role"},"updatedAt":{"type":"Date"},"users":{"type":"Relation","targetClass":"_User"}}},{"className":"ActivitySubGroup","fields":{"ACL":{"type":"ACL"},"ActivityCount":{"type":"Number"},"Group":{"type":"Pointer","targetClass":"ActivityGroup"},"Name":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Activity","fields":{"ACL":{"type":"ACL"},"Address":{"type":"String"},"Description":{"type":"String"},"DisplayImage":{"type":"File"},"Group":{"type":"Pointer","targetClass":"ActivitySubGroup"},"Link":{"type":"String"},"Name":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Map","fields":{"ACL":{"type":"ACL"},"Coordinates":{"type":"GeoPoint"},"Description":{"type":"String"},"EventId":{"type":"Pointer","targetClass":"Event"},"HasFile":{"type":"Boolean"},"Name":{"type":"String"},"Source":{"type":"File"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"LocalInfo","fields":{"ACL":{"type":"ACL"},"Description":{"type":"String"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Trivia","fields":{"ACL":{"type":"ACL"},"Description":{"type":"String"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"DiscussionTopic","fields":{"ACL":{"type":"ACL"},"EventId":{"type":"Pointer","targetClass":"Event"},"Name":{"type":"String"},"Title":{"type":"String"},"TotalMessages":{"type":"Number"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"DiscussionMessage","fields":{"ACL":{"type":"ACL"},"Avatar":{"type":"String"},"Message":{"type":"String"},"Name":{"type":"String"},"Topic":{"type":"Pointer","targetClass":"DiscussionTopic"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"_Installation","fields":{"ACL":{"type":"ACL"},"GCMSenderId":{"type":"String"},"appIdentifier":{"type":"String"},"appName":{"type":"String"},"appVersion":{"type":"String"},"badge":{"type":"Number"},"channels":{"type":"Array"},"createdAt":{"type":"Date"},"deviceToken":{"type":"String"},"deviceType":{"type":"String"},"installationId":{"type":"String"},"localeIdentifier":{"type":"String"},"objectId":{"type":"String"},"parseVersion":{"type":"String"},"pushType":{"type":"String"},"timeZone":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Event","fields":{"ACL":{"type":"ACL"},"City":{"type":"String"},"Country":{"type":"String"},"EventDate":{"type":"Date"},"EventDateStr":{"type":"String"},"Latitude":{"type":"String"},"Longitude":{"type":"String"},"MainBackground":{"type":"File"},"Name":{"type":"String"},"WeatherURL":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"Welcome","fields":{"ACL":{"type":"ACL"},"EventId":{"type":"Pointer","targetClass":"Event"},"HeaderImage":{"type":"File"},"Message":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}},{"className":"_Session","fields":{"ACL":{"type":"ACL"},"createdAt":{"type":"Date"},"createdWith":{"type":"Object"},"expiresAt":{"type":"Date"},"installationId":{"type":"String"},"objectId":{"type":"String"},"restricted":{"type":"Boolean"},"sessionToken":{"type":"String"},"updatedAt":{"type":"Date"},"user":{"type":"Pointer","targetClass":"_User"}}},{"className":"Photo","fields":{"ACL":{"type":"ACL"},"DownloadURL":{"type":"String"},"EventId":{"type":"Pointer","targetClass":"Event"},"Filename":{"type":"String"},"UploadedBy":{"type":"String"},"createdAt":{"type":"Date"},"objectId":{"type":"String"},"updatedAt":{"type":"Date"}}}]}
