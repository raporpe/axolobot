import pandas as pd

df = pd.read_csv("text_emotion.csv")

print(df['sentiment'].value_counts())

#define values
values = ["empty","enthusiasm","boredom","anger"]

#drop rows that contain any value in the list
df = df[df.sentiment.isin(values) == False]

print(df['sentiment'].value_counts())