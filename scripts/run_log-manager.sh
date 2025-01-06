cd log-viewer
source .venv/bin/activate
gunicorn -w 1 -b 127.0.0.1:5000 app:app
