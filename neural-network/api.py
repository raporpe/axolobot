# Neural network api that uses our bot to calculate the sentiment
# from tensorflow import keras
from flask import Flask
from flask import jsonify
from tensorflow import keras


app = Flask(__name__)


# model = create_model()
# model.load_weights('sentiments_model')
# classes = ["anger", "boredom", "empty", "enthusiasm", "fun", "happiness", "hate", "love",
#           "neutral", "relief", "sadness", "surprise", "worry"]


@ app.route("/v1/sentiment/<s>", methods=["GET"])
def get_sentiment(s):
    #prediction = model.predict([s])
    response = {}
    # for i in range(len(classes)):
    #    response.append({classes[i]: prediction[0][i]})
    return jsonify(response)


if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=False, port=8081)
