# Neural network api that uses our bot to calculate the sentiment
# from tensorflow import keras
from flask import Flask
from flask import jsonify
import tensorflow as tf
import pickle
from tensorflow.keras.preprocessing.sequence import pad_sequences

app = Flask(__name__)

model = tf.keras.models.load_model("./model/sentiments.h5")
with open("./model/tokenizer.pkl", 'rb') as handle:
    tokenizer = pickle.load(handle)

def decode_sentiment(score, include_neutral=True):
    if include_neutral:        
        label = "NEUTRAL"
        if score <= 0.4:
            label = "NEGATIVE"
        elif score >= 0.7:
            label = "POSITIVE"

        return label
    else:
        return "NEGATIVE" if score < 0.5 else "POSITIVE"

def predict(text, include_neutral=True):
    SEQUENCE_LENGTH = 50
    # Tokenize text
    x_test = pad_sequences(tokenizer.texts_to_sequences([text]), maxlen=SEQUENCE_LENGTH)
    # Predict
    score = model.predict([x_test])[0]
    # Decode sentiment
    label = decode_sentiment(score, include_neutral=include_neutral)

    return {"sentiment": label, "score": float(score)} 


@ app.route("/v1/sentiment/<sentiment>", methods=["GET"])
def get_sentiment(sentiment):
    #prediction = model.predict([s])
    response = predict(sentiment)
    # for i in range(len(classes)):
    #    response.append({classes[i]: prediction[0][i]})
    return jsonify(response)


if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=False, port=8081)
