import pandas as pd

df = pd.read_csv("text_emotion.csv")

print(df['sentiment'].value_counts())

#define values
values = ["neutral"]

#drop rows that contain any value in the list
df = df[df.sentiment.isin(values) == False]

df.loc[df['sentiment'] == "anger", 'sentiment'] = "hate"
df.loc[df['sentiment'] == "empty", 'sentiment'] = "sadness"
df.loc[df['sentiment'] == "boredom", 'sentiment'] = "sadness"
df.loc[df['sentiment'] == "worry", 'sentiment'] = "sadness"
df.loc[df['sentiment'] == "hate", 'sentiment'] = "sadness"
df.loc[df['sentiment'] == "enthusiasm", 'sentiment'] = "happiness"
df.loc[df['sentiment'] == "fun", 'sentiment'] = "happiness"
df.loc[df['sentiment'] == "relief", 'sentiment'] = "happiness"
df.loc[df['sentiment'] == "love", 'sentiment'] = "happiness"
df.loc[df['sentiment'] == "surprise", 'sentiment'] = "happiness"



print(df['sentiment'].value_counts())


# Test model

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


