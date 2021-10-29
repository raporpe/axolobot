import pandas as pd

df = pd.read_csv("text_emotion.csv")

print(df['sentiment'].value_counts())

#define values
#values = ["empty","enthusiasm","boredom","anger"]

#drop rows that contain any value in the list
#df = df[df.sentiment.isin(values) == False]

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