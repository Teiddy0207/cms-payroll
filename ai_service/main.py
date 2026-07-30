"""
main.py — Entry point của AI Service.

Đăng ký tất cả Blueprint và khởi chạy Flask server.

Routes:
    GET  /health              → Kiểm tra trạng thái service
    POST /embed               → Face recognition (face_service)
    POST /analyze-meeting     → Meeting AI summarization (meeting_service)

Chạy:
    python main.py
    # hoặc production:
    gunicorn -w 2 -b 0.0.0.0:5050 main:app
"""

from flask import Flask, jsonify

from face_service import face_bp
from meeting_service import meeting_bp

# ──────────────────────────────────────────────────────────────
# App khởi tạo
# ──────────────────────────────────────────────────────────────

app = Flask(__name__)

# Tăng giới hạn upload (file audio cuộc họp có thể lớn)
app.config["MAX_CONTENT_LENGTH"] = 500 * 1024 * 1024  # 500 MB

# ──────────────────────────────────────────────────────────────
# Đăng ký Blueprints
# ──────────────────────────────────────────────────────────────

app.register_blueprint(face_bp)
app.register_blueprint(meeting_bp)

# ──────────────────────────────────────────────────────────────
# Health check
# ──────────────────────────────────────────────────────────────

@app.route("/health", methods=["GET"])
def health():
    return jsonify({
        "status":   "ok",
        "services": ["face-recognition", "meeting-ai"],
        "routes": {
            "face":    "POST /embed",
            "meeting": "POST /analyze-meeting",
        },
    })


# ──────────────────────────────────────────────────────────────
# Entry point
# ──────────────────────────────────────────────────────────────

if __name__ == "__main__":
    app.run(host="127.0.0.1", port=5050, debug=False)
