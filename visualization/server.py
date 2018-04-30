from flask import Flask, jsonify, render_template, request
import json
import urllib.request
import urllib.parse
import os

# This is a simple webserver written with Flask, designed to host a visualization of the data.
app = Flask(__name__)

# Default URL will display JSON raw.
@app.route('/blockchainraw')
def raw_json():
	# jsonify returns formatted JSON
	return jsonify(get_blockchain_json())

# URL to display entire blockchain representation
@app.route('/')
def index():
	blockchain = get_blockchain_json()
	return render_template('blockchain.html', blockchain=blockchain)

# URL to display history of one item
@app.route('/getItemHistory', methods=['GET','POST'])
def get_item_history():
	if request.method == 'POST':
		blockchain = get_item_history_response(request.form['text'])
		return render_template('getItemHistoryResults.html', blockchain=blockchain)
	return render_template('getItemHistoryInput.html')

# URL to display all items a user owns
@app.route('/getUserItems', methods=['GET','POST'])
def get_user_items():
	if request.method == 'POST':
		item_list = get_items_of_owner(request.form['text'])
		name = {'name':request.form['text']}
		return render_template('getUserItemsResults.html', item_list=item_list, name=name)
	return render_template('getUserItemsInput.html')


# Helper function to grab JSON.
def get_blockchain_json():
	# URL of the blockchain node that returns JSON. 

	# Local URL for dummy testing JSON data.
	# blockchain_node_api_url = "file://" + os.getcwd() + "/sample.json"

	# Production URL, assuming you start your Go server on 8080. If different, change '8081' to the port you inputted plus 1.
	blockchain_node_api_url = "http://127.0.0.1:8081/joinGetBlock" 

	# Python JSON object to manipulate
	blockchain = json.load(urllib.request.urlopen(blockchain_node_api_url))

	# Parse block data to make it easier to put in HTML
	for block in blockchain:
		type = block['BlockTransaction']['TransactionType']
		transactionInfo = {}
		
		# If transaction type is empty, we assume it's a genesis block and continue
		if type == '':
			transactionInfo['type'] = "Genesis"
			transactionInfo['data'] = block['BlockTransaction']['Cr']
		else:
			transactionInfo['type'] = type
			transactionInfo["data"] = block['BlockTransaction'][type[:2]]
		block["ParsedBlockTransaction"] = transactionInfo

	return blockchain

def get_item_history_response(itemId):
	blockchain_node_api_url = "http://127.0.0.1:8081/getItemHistory"

	values = {"itemid": itemId}
	data = urllib.parse.urlencode(values)
	data = data.encode('ascii')
	req = urllib.request.Request(blockchain_node_api_url, data)
	response = urllib.request.urlopen(req)

	blockchain = json.load(response)
	# Parse block data to make it easier to put in HTML
	for block in blockchain:
		type = block['TransactionType']
		transactionInfo = {}
		
		# If transaction type is empty, we assume it's a genesis block and continue
		if type == '':
			transactionInfo['type'] = "Genesis"
			transactionInfo['data'] = block['Cr']
		else:
			transactionInfo['type'] = type
			transactionInfo["data"] = block[type[:2]]
		block["ParsedBlockTransaction"] = transactionInfo
	return blockchain


def get_items_of_owner(userId):
	blockchain_node_api_url = "http://127.0.0.1:8081/getItemsOfOwner"

	values = {"userid": userId}
	data = urllib.parse.urlencode(values)
	data = data.encode('ascii')
	req = urllib.request.Request(blockchain_node_api_url, data)
	response = urllib.request.urlopen(req)

	blockchain = json.load(response)
	return blockchain

if __name__ == '__main__':
	# Default value for run() starts the server on localhost, good for debugging
	app.run()

	# Running on 0.0.0.0 makes it visible using my machine's IP address over local network (i.e. BU network)
	# TODO: Test this with another machine close by
    # app.run(host='0.0.0.0')