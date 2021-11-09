# Public api for getting information about axolobot
from flask import Flask
from flask import jsonify
import os
import mysql.connector


app = Flask(__name__)
db_password = os.getenv("DB_PASSWORD")
db = mysql.connector.connect(host="database",
                             database="axolobot",
                             user="root",
                             password=db_password)


@app.route("/v1/info", methods=["GET"])
def get_total_requests():
    sql_cursor = db.cursor()
    sql_cursor.execute("SELECT * FROM mention")
    total_requests = len(sql_cursor.fetchall())
    response = {"mentions": total_requests}
    return jsonify(response)


if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=False, port=8080)
