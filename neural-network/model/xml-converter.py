import xml.dom.minidom as xml
import pandas as pd

doc = xml.parse("dataset.xml")

df = pd.DataFrame(columns=['sentiment', 'text'])

sent_converter = {
    "NEU": 0.5,
    "P+": 1,
    "P": 0.75,
    "N": 0.25,
    "N+": 0,
    "NONE": -1,
}

for tweet in doc.getElementsByTagName("tweets"):
    try:
        sentiment = tweet.getElementsByTagName("sentiments")[0].getElementsByTagName("polarity")[0].getElementsByTagName("value")[0].firstChild.data
        text = tweet.getElementsByTagName("content")[0].firstChild.data
        print("Sentiment -> {} - {}".format(sentiment, text))
        sentiment = sent_converter[sentiment]

        # Discard sentiments with NONE value
        if sentiment == -1: continue

        # Add sentiment to dataframe
        df = df.append({"text": text, "sentiment": sentiment}, ignore_index=True)

    except Exception as e:
        print(e)
        continue

df.to_csv("./dataset1.csv")
print(len(df), "in total.")