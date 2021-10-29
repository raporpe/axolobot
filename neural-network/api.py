# Neural network api that uses our bot to calculate the sentiment

from tensorflow import keras
import numpy as np


model = keras.models.load_model('sentiments_model')


classes = ["Happy","Neutral","Sad"]

while True:
    frase = input("\nWrite a Tweet:")
    prediction = model.predict([frase])
    print("\nPercentage of Happiness:",prediction[0][0])
    print("Percentage of Neutrality:",prediction[0][1])
    print("Percentage of Sadness:",prediction[0][2])
    testPred = np.argmax(prediction, axis=1)[0]
    classPred = classes[testPred]
    print(classPred)

