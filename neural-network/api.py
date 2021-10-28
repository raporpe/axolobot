# Neural network api that uses our bot to calculate the sentiment
import pandas as pd
import tensorflow as tf
from tensorflow import keras
from collections import Counter
from tensorflow.keras.preprocessing.sequence import pad_sequences
from tensorflow.keras.preprocessing.text import Tokenizer

df = pd.read_csv("corpus.csv")

df = df[df['Tweet'].notnull()]

# Obtener el numero de palabras únicas de nuestro corpus

def counter_word(text_col):
    count = Counter()
    for text in text_col.values:
        for word in text.split():
            count[word] += 1
    return count

counter = counter_word(df.Tweet)
num_unique_words = len(counter)

max_length = 20 # Longitud maxima arbitraria

model = keras.models.load_model('sentiments.h5')
tokenizer = Tokenizer(num_words=num_unique_words)
tokenizer.fit_on_texts(df.Tweet.to_numpy())

while True:
    frase = input("escribe un tweet:")
    frase_sequence=tokenizer.texts_to_sequences([frase])
    print(frase_sequence)
    frase_padded = pad_sequences(frase_sequence,maxlen=max_length,padding="post", truncating="post")
    print(frase_padded)
    print("negative" if model.predict(frase_padded)<0.5 else "positive")

