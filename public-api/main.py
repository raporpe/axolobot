# Public api for getting information about axolobot
from logging import exception
from flask import Flask
from flask import jsonify
import os
import mysql.connector
import time
from flask_cors import CORS, cross_origin


app = Flask(__name__)
cors = CORS(app)

class DataBase:
    
    db = None
    last_update_time = time.time()
    last_total_requests = 0
    db_password = os.getenv("DB_PASSWORD")

    def __init__(self):
        # Connect to the database
        while self.db == None:
            try:
                self.db = mysql.connector.connect(host="database",
                                            database="axolobot",
                                            user= "root",
                                            password=self.db_password)
            except:
                print("❌ Error when connecting to database, retrying...")
                time.sleep(1)
                continue
            break

    def get_total_requests(self):
        if time.time() - self.last_update_time > 5:
            sql_cursor = self.db.cursor()
            sql_cursor.execute("SELECT * FROM mention")
            self.last_total_requests = len(sql_cursor.fetchall())
            self.db.commit()
            self.last_update_time = time.time()
        return self.last_total_requests


@app.route("/v1/info", methods=["GET"])
def get_total_requests():
    response = {"mentions": db.get_total_requests()}
    return jsonify(response)


if __name__ == "__main__":
    db = DataBase()
    app.run(host="0.0.0.0", debug=False, port=8080)
