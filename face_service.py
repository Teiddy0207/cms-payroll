import base64
import io
from flask import Flask, request, jsonify
import numpy as np
from PIL import Image
import face_recognition

app = Flask(__name__)

def decode_base64_image(base64_str):
	if "," in base64_str:
		base64_str = base64_str.split(",")[1]
	image_bytes = base64.b64decode(base64_str)
	image = Image.open(io.BytesIO(image_bytes))
	return np.array(image.convert("RGB"))

@app.route("/embed", methods=["POST"])
def get_embedding():
	data = request.get_json()
	if not data or "image" not in data:
		return jsonify({"status": "error", "message": "Missing image"}), 400
	
	try:
		img_np = decode_base64_image(data["image"])
		face_locations = face_recognition.face_locations(img_np)
		if len(face_locations) == 0:
			return jsonify({"status": "error", "message": "No face detected"}), 400
		
		encodings = face_recognition.face_encodings(img_np, face_locations)
		if len(encodings) == 0:
			return jsonify({"status": "error", "message": "Could not encode face"}), 400
		
		embedding = encodings[0].tolist()
		return jsonify({"status": "success", "embedding": embedding})
	except Exception as e:
		return jsonify({"status": "error", "message": str(e)}), 500

@app.route("/health", methods=["GET"])
def health():
	return jsonify({"status": "ok"})

if __name__ == "__main__":
	app.run(host="127.0.0.1", port=5050)
