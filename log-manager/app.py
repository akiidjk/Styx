import json
import os
from datetime import datetime

import pandas as pd
from flask import Flask, jsonify, render_template, request
from werkzeug.utils import secure_filename

app = Flask(__name__, template_folder=".")
app.config["UPLOAD_FOLDER"] = "./uploads"

def load_logs(file_path):
    """Load structured logs from a file."""
    try:
        with open(file_path, "r") as file:
            data = []
            for line in file:
                log = json.loads(line)
                log["time"] = datetime.fromtimestamp(log["time"]).strftime(
                    "%Y-%m-%d %H:%M:%S"
                )
                data.append(log)
            return pd.DataFrame(data)
    except Exception as e:
        print(f"Error loading logs: {e}")
        return None

@app.route("/logs", methods=["POST", "GET"]) # type: ignore
def get_logs():
    """Endpoint to retrieve logs."""
    if request.method == "POST":
        if "file" not in request.files:
            return jsonify({"error": "No file provided"}), 400
        file = request.files["file"]
        if file and file.filename:
            filename = secure_filename(file.filename)
            file.save(os.path.join(app.config["UPLOAD_FOLDER"], filename))
        else:
            return jsonify({"error": "No file provided"}), 400

    elif request.method == "GET":
        if "fileName" not in request.args:
            return jsonify({"error": "No file provided"}), 400
        fileName = request.args.get("fileName")
        if fileName:
            filename = secure_filename(fileName)
            if filename not in os.listdir(app.config["UPLOAD_FOLDER"]):
                return jsonify({"error": "File not found"}), 404
        else:
            return jsonify({"error": "No file provided"}), 400

    else:
        return jsonify({"error": "Invalid request method"}), 405

    df = load_logs(os.path.join(app.config["UPLOAD_FOLDER"], filename))
    if df is None:
        return jsonify({"error": "Could not load logs"}), 500

    return df.to_json(orient="records"), 200

@app.route("/")
def index():
    """Render the index page for the frontend."""
    return render_template("index.html")


# if __name__ == "__main__":
    # app.run(debug=False, host="127.0.0.1", port=5000)
