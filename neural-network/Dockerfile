FROM python:3.8

WORKDIR /usr/src/app

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# Dropbox S3-like-but-cheaper hosting of the file
RUN wget https://www.dropbox.com/s/gop1dbpi7f2s2lj/sentiments_en.h5 -P ./model/

COPY . .

CMD [ "python", "api.py" ]