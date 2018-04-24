from flask import Flask, jsonify
import json
import urllib.request

# This is a simple webserver written with Flask, designed to host a visualization of the data.
app = Flask(__name__)

# Default URL will display JSON raw.
@app.route('/')
def display_json():
	# URL of the blockchain node that returns JSON. Provided URL returns an example JSON for testing.
	blockchain_node_api_url = 'https://raw.githubusercontent.com/SMAPPNYU/ProgrammerGroup/master/LargeDataSets/sample-tweet.raw.json'
	# Python JSON object to manipulate
	blockchain_json = json.load(urllib.request.urlopen(blockchain_node_api_url))
	# jsonify returns formatted JSON
	return jsonify(blockchain_json["created_at"])

if __name__ == '__main__':
	# Default value for run() starts the server on localhost, good for debugging
	# app.run()

	# Running on 0.0.0.0 makes it visible using my machine's IP address over local network (i.e. BU network)
	# Test this with another machine close by
    app.run(host='0.0.0.0')