from flask import Flask, jsonify, render_template
import json
import urllib.request
import os

# This is a simple webserver written with Flask, designed to host a visualization of the data.
app = Flask(__name__)

# Default URL will display JSON raw.
@app.route('/')
def raw_json():
	# jsonify returns formatted JSON
	return jsonify(get_blockchain_json())

@app.route('/index')
def index():
	blockchain = get_blockchain_json()
	user = {'username': 'Mark'}
	return render_template('index.html', user=user, blockchain=blockchain)

# Helper function to grab JSON.
def get_blockchain_json():
	# URL of the blockchain node that returns JSON. Provided URL returns an example JSON for testing.
	#blockchain_node_api_url = "file://" + os.getcwd() + "/sample.json"


	# Assuming you start your Go server on 8080. If different, change '8081' to the port you inputted plus 1.
	blockchain_node_api_url = "http://127.0.0.1:8081/joinGetBlock" 


	# Python JSON object to manipulate
	blockchain_json = json.load(urllib.request.urlopen(blockchain_node_api_url))

	return blockchain_json

if __name__ == '__main__':
	# Default value for run() starts the server on localhost, good for debugging
	app.run()

	# Running on 0.0.0.0 makes it visible using my machine's IP address over local network (i.e. BU network)
	# TODO: Test this with another machine close by
    # app.run(host='0.0.0.0')