# Neural network api that uses our bot to calculate the sentiment
# from tensorflow import keras
from flask import Flask, request
from flask import jsonify
import tensorflow as tf
import pickle
from tensorflow.keras.preprocessing.sequence import pad_sequences
import re
import base64

app = Flask(__name__)

# Load the model and tokenizer for english version
model = tf.keras.models.load_model("./model/sentiments.h5")
with open("./model/tokenizer.pkl", 'rb') as handle:
    tokenizer = pickle.load(handle)


# Uses the model to predict the sentiment
def predict(text):
    # Clean text from url's and tags
    text = re.sub("((http|https)\:\/\/)?[a-zA-Z0-9\.\/\?\:@\-_=#]+\.([a-zA-Z]){2,6}([a-zA-Z0-9\.\&\/\?\:@\-_=#])*", "", text)
    text = re.sub("@(\w){1,15}", "", text)

    # Padding for the neural network
    SEQUENCE_LENGTH = 50

    # Tokenize text
    x_test = pad_sequences(tokenizer.texts_to_sequences([text]), maxlen=SEQUENCE_LENGTH)

    # Predict the sentiment
    score = model.predict([x_test])[0]

    return {"score": str(int(score[0]*100))} 


# This function decodes the headers in which the sentiment text is passed encoded in base64
def get_sentiment_from_request(request):
    return base64.b64decode(request.headers.get("sentiment")).decode('utf-8')


# Sentiment analysis for english text
@ app.route("/v1/sentiment/en", methods=["GET"])
def get_sentiment_en():
    sentiment = get_sentiment_from_request(request)
    response = predict(sentiment)
    return jsonify(response)

# Sentiment analysis for spanish text
@ app.route("/v1/sentiment/es", methods=["GET"])
def get_sentiment_es():
    sentiment = get_sentiment_from_request(request)
    response = predict(sentiment)
    return jsonify(response)


if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=False, port=8081)
