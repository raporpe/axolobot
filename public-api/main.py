# Public api for getting information about axolobot
from flask import Flask
from flask import jsonify
import os
import mysql.connector


app = Flask(__name__)
global sql_cursor

@app.route("/api/v1/total-requests", methods=["GET"])
def get_total_requests():
    sql_cursor.execute("SELECT * FROM mention")
    total_requests = len(sql_cursor.fetchall())
    response = {"total-requests":total_requests}
    return jsonify(response)


if __name__ == "__main__":
    db_password = os.getenv("DB_PASSWORD")
    db = mysql.connector.connect(host="localhost",
                                database = "axolobot",
                                user = "root",
                                password = db_password)

    sql_cursor = db.cursor()
    
    app.run(debug=True, port=8080)




