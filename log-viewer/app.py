import json
import os
import time
import tempfile
import datetime

import pandas as pd
from flask import Flask, jsonify, render_template, request, Response
from werkzeug.utils import secure_filename
import requests_toolbelt.multipart.encoder as mp_encoder

app = Flask(__name__, template_folder=".")
app.config["UPLOAD_FOLDER"] = "./uploads"
app.config["LOGS_FOLDER"] = os.environ.get("LOGS_FOLDER", "/var/log/styx")
MAX_LINES = 1_000_000


def load_logs(file_path, max_lines=MAX_LINES):
    """Load structured logs from a file with improved performance."""
    try:
        with open(file_path, "r") as file:
            lines = (file.readline() for _ in range(max_lines))
            data = (json.loads(line) for line in lines if line.strip())
            logs = list(data)
        df = pd.DataFrame(logs)
        if "time" in df.columns:
            df["time"] = pd.to_datetime(df["time"], unit="s").dt.strftime(
                "%Y-%m-%d %H:%M:%S"
            )

        df.sort_values(by="time", ascending=False, inplace=True)

        return df
    except Exception as e:
        print(f"Error loading logs: {e}")
        return None


@app.route("/logs", methods=["GET"])  # type: ignore
def get_logs():
    """Endpoint to retrieve logs."""
    if "fileName" not in request.args:
        return jsonify({"error": "No file provided"}), 400
    fileName = request.args.get("fileName")
    if fileName:
        filename = secure_filename(fileName)
        if filename not in os.listdir(
            app.config["UPLOAD_FOLDER"]
        ) and filename not in os.listdir(app.config["LOGS_FOLDER"]):
            return jsonify({"error": "File not found"}), 404
        else:
            print(f"{datetime.datetime.now()} | Loading logs from {filename}")
            start_time = time.time()
            if filename in os.listdir(app.config["UPLOAD_FOLDER"]):
                df = load_logs(os.path.join(app.config["UPLOAD_FOLDER"], filename))
            else:
                df = load_logs(os.path.join(app.config["LOGS_FOLDER"], filename))
            print(
                f"{datetime.datetime.now()} | Loaded {len(df) if df is not None else 0} logs in {time.time() - start_time}"
            )
    else:
        return jsonify({"error": "No file provided"}), 400

    if df is None:
        return jsonify({"error": "Could not load logs"}), 500

    try:
        temp_file = tempfile.NamedTemporaryFile(delete=False, suffix=".json")
        df.to_json(temp_file.name, orient="records")
        temp_file.close()
        encoder = mp_encoder.MultipartEncoder(
            fields={
                "file": (
                    f"{filename}.json",
                    open(temp_file.name, "rb"),
                    "application/json",
                ),
                "columns": json.dumps(df.columns.tolist()),
            }
        )
        return Response(encoder.to_string(), mimetype=encoder.content_type, status=201)
    except Exception as e:
        return jsonify({"error": f"Error generating file: {e}"}), 500


@app.route("/logs", methods=["POST"])  # type: ignore
def upload_logs():
    if "file" not in request.files:
        return jsonify({"error": "No file provided"}), 400
    file = request.files["file"]
    if file and file.filename:
        filename = secure_filename(file.filename)
        file.save(os.path.join(app.config["UPLOAD_FOLDER"], filename))
        df = load_logs(os.path.join(app.config["UPLOAD_FOLDER"], filename))
    else:
        return jsonify({"error": "No file provided"}), 400

    if df is None:
        return jsonify({"error": "Could not load logs"}), 500

    try:
        temp_file = tempfile.NamedTemporaryFile(delete=False, suffix=".json")
        df.to_json(temp_file.name, orient="records")
        temp_file.close()
        print(json.dumps(df.columns.tolist()))
        encoder = mp_encoder.MultipartEncoder(
            fields={
                "file": (
                    f"{filename}.json",
                    open(temp_file.name, "rb"),
                    "application/json",
                ),
                "columns": json.dumps(df.columns.tolist()),
            }
        )
        return Response(encoder.to_string(), mimetype=encoder.content_type, status=201)
    except Exception as e:
        return jsonify({"error": f"Error generating file: {e}"}), 500


@app.route("/")
def index():
    """Render the index page for the frontend."""
    file_list = os.listdir(app.config["LOGS_FOLDER"])
    return render_template("index.html", file_list=file_list)


if __name__ == "__main__":
    app.run(debug=True, host="127.0.0.1", port=5000)
