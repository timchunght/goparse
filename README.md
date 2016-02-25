# GoParse
---

Make sure mongodb is running

	PORT=8080 ./goparse

This will connect to the local mongodb server at port 27017

You can also run it as:

	PORT=8080 MONGO_URL=mongodb://myurl ./goparse

Visit PORT ``8080``.

To make sure that it is working, try:

	curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"score":1337,"playerName":"Sean Plott","cheatMode":false}' \
	http://localhost:8080/classes/GameScore

That's it!