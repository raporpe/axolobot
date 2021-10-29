from sklearn.preprocessing import LabelBinarizer
from nltk.corpus import stopwords
from tensorflow import keras
import pandas as pd
import numpy as np
import re


df = pd.read_csv("text_emotion.csv")

#combinar Classes

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


# Preprocesado

hashtags = re.compile(r"^#\S+|\s#\S+")
mentions = re.compile(r"^@\S+|\s@\S+")
urls = re.compile(r"https?://\S+")

def process_text(text):
    text = re.sub(r'http\S+', '', text)
    text = hashtags.sub(' hashtag', text)
    text = mentions.sub(' entity', text)
    return text.strip().lower()

df['content'] = df.content.apply(process_text)

# Eliminando las stop words

stop = set(stopwords.words("english"))

def remove_stopwords(text):
    filtered_words = [word.lower()
                      for word in text.split() if word.lower() not in stop]
    return " ".join(filtered_words)

df["content"] = df.content.map(remove_stopwords)


# Separar datos de entramiento y validacion

encoder = LabelBinarizer()

train_size = int(df.shape[0]*0.8)

train_df = df[:train_size]
val_df = df[train_size:]

train_sentences = train_df.content.to_numpy()
train_labels = encoder.fit_transform(train_df.sentiment)
val_sentences = val_df.content.to_numpy()
val_labels = encoder.fit_transform(val_df.sentiment)

# Creamos el modelo con keras

VOCAB_SIZE = 1000
encoderLayer = keras.layers.TextVectorization(max_tokens=VOCAB_SIZE)
encoderLayer.adapt(train_sentences)

model = keras.Sequential([
    encoderLayer,
    keras.layers.Embedding(input_dim=len(encoderLayer.get_vocabulary()),output_dim=64),
    keras.layers.LSTM(64),
    keras.layers.Dense(64, activation='relu'),
    keras.layers.Dense(64, activation='relu'),
    keras.layers.Dense(3, activation='softmax')
])

# Parametros para el entramiento

optim = keras.optimizers.Adam(learning_rate=0.001)

model.compile(loss="categorical_crossentropy", optimizer=optim, metrics="accuracy")

model.fit(train_sentences , train_labels, epochs=1, validation_data=(val_sentences, val_labels),batch_size=10, verbose=1)

model.save("sentiments_model")

# Comprobamos el modelo

prediction = model.predict(val_sentences)
testPred = np.argmax(prediction, axis=1)
classPred = encoder.classes_[testPred]
print(val_sentences[20:30])
print(classPred[20:30])