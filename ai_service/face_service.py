import base64
import io

import numpy as np
from flask import Blueprint, jsonify, request
from PIL import Image

try:
    import face_recognition
    _FACE_AVAILABLE = True
except ImportError:
    _FACE_AVAILABLE = False

face_bp = Blueprint("face", __name__)


# ──────────────────────────────────────────────────────────────
# Helpers
# ──────────────────────────────────────────────────────────────

def _decode_base64_image(base64_str: str) -> np.ndarray:
    """Giải mã base64 string → numpy RGB array."""
    if "," in base64_str:
        base64_str = base64_str.split(",")[1]
    image_bytes = base64.b64decode(base64_str)
    image = Image.open(io.BytesIO(image_bytes))
    return np.array(image.convert("RGB"))


# ──────────────────────────────────────────────────────────────
# Routes
# ──────────────────────────────────────────────────────────────

@face_bp.route("/embed", methods=["POST"])
def get_embedding():
    """
    Nhận ảnh base64 → trích xuất face embedding 128-chiều.

    Request JSON:
        { "image": "<base64_string>" }

    Response JSON:
        { "status": "success", "embedding": [...] }
    """
    if not _FACE_AVAILABLE:
        return jsonify({"status": "error", "message": "face_recognition chưa được cài đặt"}), 503

    data = request.get_json()
    if not data or "image" not in data:
        return jsonify({"status": "error", "message": "Thiếu trường 'image'"}), 400

    try:
        img_np = _decode_base64_image(data["image"])

        face_locations = face_recognition.face_locations(img_np)
        if not face_locations:
            return jsonify({"status": "error", "message": "Không phát hiện khuôn mặt"}), 400

        encodings = face_recognition.face_encodings(img_np, face_locations)
        if not encodings:
            return jsonify({"status": "error", "message": "Không trích xuất được embedding"}), 400

        return jsonify({"status": "success", "embedding": encodings[0].tolist()})

    except Exception as e:
        return jsonify({"status": "error", "message": str(e)}), 500
