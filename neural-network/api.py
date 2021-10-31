# Neural network api that uses our bot to calculate the sentiment

from tensorflow import keras
import numpy as np


model = keras.models.load_model('sentiments_model')


classes = ["anger", "boredom", "empty", "enthusiasm", "fun", "happiness", "hate", "love",
 "neutral", "relief", "sadness", "surprise", "worry"]

while True:
    frase = input("\nWrite a Tweet:")
    prediction = model.predict([frase])
    for sentiment in range(len(classes)):
        print("Percentage of",classes[sentiment],":",prediction[0][sentiment])
    testPred = np.argmax(prediction, axis=1)[0]
    classPred = classes[testPred]
    print("\nSentimiento del Tweet:",classPred)

