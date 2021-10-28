from tensorflow.keras import layers
from tensorflow.keras.preprocessing.sequence import pad_sequences
from tensorflow.keras.preprocessing.text import Tokenizer
from nltk.corpus import stopwords
from collections import Counter
from tensorflow import keras
import pandas as pd
import re


df = pd.read_csv("train.csv")

# Preprocesado

hashtags = re.compile(r"^#\S+|\s#\S+")
mentions = re.compile(r"^@\S+|\s@\S+")
urls = re.compile(r"https?://\S+")

def process_text(text):
    text = re.sub(r'http\S+', '', text)
    text = hashtags.sub(' hashtag', text)
    text = mentions.sub(' entity', text)
    return text.strip().lower()

df['Tweet'] = df.tweet.apply(process_text)

# Eliminando las stop words

stop = set(stopwords.words("english"))

def remove_stopwords(text):
    filtered_words = [word.lower()
                      for word in text.split() if word.lower() not in stop]
    return " ".join(filtered_words)

df["Tweet"] = df.tweet.map(remove_stopwords)

df.to_csv("corpus.csv")

# Obtener el numero de palabras únicas de nuestro corpus

def counter_word(text_col):
    count = Counter()
    for text in text_col.values:
        for word in text.split():
            count[word] += 1
    return count

counter = counter_word(df.tweet)
num_unique_words = len(counter)

# Separar datos de entramiento y validacion

train_size = int(df.shape[0]*0.8)

train_df = df[:train_size]
val_df = df[train_size:]

train_sentences = train_df.tweet.to_numpy()
train_labels = train_df.label.to_numpy()
val_sentences = val_df.tweet.to_numpy()
val_labels = val_df.label.to_numpy()

# Crear tokenizador con el corpus
tokenizer = Tokenizer(num_words=num_unique_words)
tokenizer.fit_on_texts(train_sentences)

# Tokenizar nuestras instancias de entrenamiento y validación
train_sequences = tokenizer.texts_to_sequences(train_sentences)
val_sequences = tokenizer.texts_to_sequences(val_sentences)

# Necesitamos a todas las instances con la misma longitud por lo que usamos pad

max_length = 20 # Longitud maxima arbitraria

train_padded = pad_sequences(
    train_sequences, maxlen=max_length, padding="post", truncating="post")
val_padded = pad_sequences(
    val_sequences, maxlen=max_length, padding="post", truncating="post")

# Creamos el modelo con keras

model = keras.models.Sequential()
model.add(layers.Embedding(num_unique_words, 32, input_length=max_length))
model.add(layers.LSTM(64, dropout=0.1))
model.add(layers.Dense(1, activation="sigmoid"))

print(model.summary())

# Parametros para el entramiento

loss = keras.losses.BinaryCrossentropy(from_logits=False)
optim = keras.optimizers.Adam(learning_rate=0.001)
metrics = ["accuracy"]

model.compile(loss=loss, optimizer=optim, metrics=metrics)

model.fit(train_padded, train_labels, epochs=10,
          validation_data=(val_padded, val_labels),batch_size=800, verbose=1)

model.save("sentiments.h5")

# Comprobamos el modelo

val_predictions = model.predict(val_padded)
val_predictions = ["positive" if p > 0.7 else "neutral" if p >
                   0.5 else "negative" for p in val_predictions]
print(val_sentences[50:55])
print(val_labels[50:55])
print(val_predictions[50:55])