# Neural network api that uses our bot to calculate the sentiment
#from tensorflow import keras
import numpy as np
from flask import Flask, json
from flask import jsonify

#model = keras.models.load_model('sentiments_model')


app = Flask(__name__)

@app.route("/api/v1/sentiment/<s>", methods=["GET"])
def get_sentiment(s):
    response = {"sentiment":s}
    return jsonify(response)


if __name__ == "__main__":
    app.run(debug=True, port=8080)