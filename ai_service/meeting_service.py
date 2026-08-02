"""
meeting_service.py — Blueprint tóm tắt cuộc họp bằng AI.

Pipeline:
    File audio (.wav/.mp3/.m4a)
        │
        ▼ Bước 2 — Speech-to-Text
    Faster-Whisper (small, INT8, CPU)
        │
        ▼ Bước 3 — Vectorization
    TF-IDF + positional feature → State vector S_t
        │
        ▼ Bước 4 — RL Optimization Loop (REINFORCE)
    Policy Network (Linear → ReLU → Softmax)
    Reward = Jaccard(selected, all_sentences) − penalty(ratio)
        │
        ▼ Bước 5 — Output
    { summary, key_decisions, action_items, sentiment, efficiency_score }
"""

import math
import os
import re
import tempfile

import numpy as np
from flask import Blueprint, jsonify, request

meeting_bp = Blueprint("meeting", __name__)

# ──────────────────────────────────────────────────────────────
# Whisper — lazy load để không nặng khi khởi động
# ──────────────────────────────────────────────────────────────

_whisper_model = None


def _get_whisper():
    """Load Faster-Whisper model (base, INT8, CPU) lần đầu gọi."""
    global _whisper_model
    if _whisper_model is None:
        from faster_whisper import WhisperModel
        _whisper_model = WhisperModel("base", device="cpu", compute_type="int8")
    return _whisper_model


# ──────────────────────────────────────────────────────────────
# Module 1 — TF-IDF Vectorizer (thuần NumPy, không cần sklearn)
# ──────────────────────────────────────────────────────────────

def _tokenize(text: str) -> list[str]:
    return re.findall(r"\w+", text.lower())


def _tfidf_matrix(sentences: list[str]) -> np.ndarray:
    """
    Tính ma trận TF-IDF.

    Returns:
        np.ndarray shape (n_sentences, vocab_size)
    """
    tokenized = [_tokenize(s) for s in sentences]
    vocab = sorted({tok for doc in tokenized for tok in doc})
    vocab_idx = {w: i for i, w in enumerate(vocab)}
    n, v = len(sentences), len(vocab)

    # Term Frequency
    tf = np.zeros((n, v), dtype=np.float32)
    for i, doc in enumerate(tokenized):
        if not doc:
            continue
        for tok in doc:
            tf[i, vocab_idx[tok]] += 1
        tf[i] /= len(doc)

    # Inverse Document Frequency (smooth)
    df = np.count_nonzero(tf, axis=0).astype(np.float32)
    idf = np.log((n + 1) / (df + 1)) + 1.0

    return tf * idf  # (n, vocab_size)


# ──────────────────────────────────────────────────────────────
# Module 2 — Policy Network (REINFORCE)
#
# Kiến trúc: input_dim → Linear → ReLU → Linear → Softmax(2)
# Tham số  : W1 (input_dim×64), b1 (64,), W2 (64×2), b2 (2,)
# ──────────────────────────────────────────────────────────────

class PolicyNetwork:
    """
    Mạng chính sách π_θ(a_t | S_t).

    Forward pass:
        h₁ = ReLU(S_t · W₁ + b₁)          # (64,)
        logits = h₁ · W₂ + b₂              # (2,)
        probs  = Softmax(logits)            # P(0)=loại, P(1)=chọn

    Backward pass (REINFORCE):
        L(θ) = −Σ log π_θ(aₜ|Sₜ) · R
        Gradient descent on W1, b1, W2, b2
    """

    def __init__(self, input_dim: int, lr: float = 0.02):
        self.lr = lr

        he_in  = np.sqrt(2.0 / input_dim)
        he_hid = np.sqrt(2.0 / 64)

        self.W1 = np.random.randn(input_dim, 64).astype(np.float32) * he_in
        self.b1 = np.zeros(64, dtype=np.float32)
        self.W2 = np.random.randn(64, 2).astype(np.float32) * he_hid
        self.b2 = np.zeros(2, dtype=np.float32)

    # ── Forward ──────────────────────────────────────────────

    def forward(self, x: np.ndarray) -> tuple[np.ndarray, np.ndarray]:
        """
        Args:
            x: state vector (input_dim,)
        Returns:
            probs: (2,)  — [P(loại), P(chọn)]
            h1:   (64,) — activations lớp ẩn (dùng cho backward)
        """
        h1 = np.maximum(0.0, x @ self.W1 + self.b1)   # ReLU
        logits = h1 @ self.W2 + self.b2
        logits -= logits.max()                          # numerical stability
        exp = np.exp(logits)
        probs = exp / exp.sum()                         # Softmax
        return probs, h1

    # ── Backward (REINFORCE gradient step) ───────────────────

    def update(self, states: np.ndarray, actions: list[int], reward: float):
        """
        Cập nhật tham số θ theo Policy Gradient.

        L(θ) = −Σₜ log π_θ(aₜ|Sₜ) · R

        Gradient của CrossEntropy-Softmax:
            ∂L/∂logits[a] = probs[a] − 1   (nếu a là hành động được chọn)
            ∂L/∂logits[¬a] = probs[¬a]
        Nhân thêm −reward để biến thành policy gradient.
        """
        for state, action in zip(states, actions):
            probs, h1 = self.forward(state)

            # Gradient tại output
            d_logits = probs.copy()
            d_logits[action] -= 1.0           # ∂CrossEntropy/∂logits
            d_logits *= -reward * self.lr     # scale by reward & lr

            # Update W2, b2
            self.W2 -= np.outer(h1, d_logits)
            self.b2 -= d_logits

            # Backprop qua ReLU
            d_h1 = d_logits @ self.W2.T
            d_h1 *= (h1 > 0).astype(np.float32)

            # Update W1, b1
            self.W1 -= np.outer(state, d_h1)
            self.b1 -= d_h1


# ──────────────────────────────────────────────────────────────
# Module 3 — Reward Function: Jaccard Coverage
# ──────────────────────────────────────────────────────────────

def _jaccard(selected: list[str], all_sents: list[str]) -> float:
    """
    R = |tokens(selected) ∩ tokens(all)| / |tokens(selected) ∪ tokens(all)|

    Đo mức độ bao phủ thông tin của tập câu được chọn so với toàn văn bản.
    Thay thế ROUGE-1 trong môi trường demo (không cần reference summary).
    """
    sel_tokens = {tok for s in selected for tok in _tokenize(s)}
    all_tokens = {tok for s in all_sents for tok in _tokenize(s)}
    if not all_tokens:
        return 0.0
    return len(sel_tokens & all_tokens) / len(sel_tokens | all_tokens)


# ──────────────────────────────────────────────────────────────
# Module 4 — RL Extractive Summarizer (Online Fine-tuning)
# ──────────────────────────────────────────────────────────────

def _rl_summarize(
    sentences: list[str],
    n_episodes: int = 40,
    target_ratio: float = 0.35,
) -> list[str]:
    """
    Chạy REINFORCE để học chính sách chọn câu tóm tắt tối ưu.

    MDP:
        Environment : D = {s₁, …, sₙ}
        State S_t   : TF-IDF(sₜ) ⊕ position(t)     ← vector số
        Action aₜ  : {0=loại, 1=chọn}
        Policy π_θ  : PolicyNetwork
        Reward R    : Jaccard(selected, D) − penalty(|ratio − target|)

    Args:
        sentences    : Danh sách câu từ Whisper.
        n_episodes   : Số episode huấn luyện.
        target_ratio : Tỉ lệ câu muốn chọn (~35% tổng).

    Returns:
        Danh sách câu được chọn (giữ thứ tự gốc).
    """
    n = len(sentences)
    if n == 0:
        return []
    if n == 1:
        return sentences[:]

    # Bước 3: Xây dựng State vector S_t
    tfidf     = _tfidf_matrix(sentences)                                 # (n, vocab)
    positions = np.linspace(0, 1, n, dtype=np.float32).reshape(-1, 1)   # (n, 1)
    states    = np.hstack([tfidf, positions])                            # (n, vocab+1)

    policy       = PolicyNetwork(states.shape[1], lr=0.02)
    best_selected: list[str] = []
    best_reward   = -1.0

    # Bước 4: Vòng lặp Episode
    for _ in range(n_episodes):
        # Agent duyệt từng câu → sample hành động
        actions = [
            int(np.random.choice(2, p=policy.forward(states[t])[0]))
            for t in range(n)
        ]

        selected = [sentences[t] for t, a in enumerate(actions) if a == 1]

        # Fallback: nếu không chọn câu nào, lấy câu TF-IDF tổng cao nhất
        if not selected:
            selected = [sentences[int(np.argmax(tfidf.sum(axis=1)))]]

        # Terminal Reward
        reward = _jaccard(selected, sentences)
        ratio  = len(selected) / n
        reward = max(0.0, reward - abs(ratio - target_ratio) * 0.5)

        if reward > best_reward:
            best_reward   = reward
            best_selected = selected

        # Cập nhật Policy (REINFORCE)
        policy.update(states, actions, reward)

    return best_selected


# ──────────────────────────────────────────────────────────────
# Sentiment (rule-based, tiếng Việt)
# ──────────────────────────────────────────────────────────────

_POS = {"đồng ý", "tốt", "hoàn thành", "thành công", "xuất sắc",
        "hiệu quả", "tuyệt", "ổn", "chấp nhận", "phê duyệt", "đúng hạn"}
_NEG = {"vấn đề", "lỗi", "trễ", "khó khăn", "thất bại",
        "chậm", "hủy", "từ chối", "lo ngại", "chưa xong", "tồn đọng"}


def _sentiment(text: str) -> str:
    tokens = set(_tokenize(text))
    pos = len(tokens & _POS)
    neg = len(tokens & _NEG)
    return "POSITIVE" if pos > neg else ("NEGATIVE" if neg > pos else "NEUTRAL")


def _efficiency_score(n_selected: int, n_total: int, reward: float) -> str:
    """
    Điểm hiệu quả cuộc họp:
        score = 0.6 × Jaccard_reward + 0.4 × compression_ratio
    """
    compression = 1.0 - (n_selected / max(n_total, 1))
    score = round((reward * 0.6 + compression * 0.4) * 100, 1)
    return f"{score}/100"


# ──────────────────────────────────────────────────────────────
# Route
# ──────────────────────────────────────────────────────────────

_ACTION_KEYWORDS = {
    "cần", "phải", "sẽ", "deadline", "hạn chót",
    "giao", "hoàn thành", "submit", "báo cáo", "xử lý", "kiểm tra",
}


@meeting_bp.route("/analyze-meeting", methods=["POST"])
def analyze_meeting():
    """
    Nhận file audio cuộc họp → trả về bản tóm tắt AI.

    Form-data:
        audio (file): .wav / .mp3 / .m4a / .ogg

    Response JSON:
        {
            "status":           "success",
            "transcript":       ["câu 1", "câu 2", ...],
            "summary":          "Đoạn tóm tắt trích xuất...",
            "key_decisions":    "Các quyết định chính...",
            "action_items":     "Việc cần làm...",
            "sentiment":        "POSITIVE | NEUTRAL | NEGATIVE",
            "efficiency_score": "xx.x/100"
        }
    """
    if "audio" not in request.files:
        return jsonify({"status": "error", "message": "Thiếu file audio (field: 'audio')"}), 400

    audio_file = request.files["audio"]
    suffix = os.path.splitext(audio_file.filename)[1] or ".wav"

    # Lưu tạm file audio
    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        audio_file.save(tmp.name)
        tmp_path = tmp.name

    try:
        # ── Bước 2: Speech-to-Text (Faster-Whisper) ──────────
        model    = _get_whisper()
        segments, _ = model.transcribe(tmp_path, language="vi", beam_size=5)
        sentences = [seg.text.strip() for seg in segments if seg.text.strip()]

        if not sentences:
            return jsonify({"status": "error", "message": "Không nhận diện được giọng nói"}), 422

        # ── Bước 3+4: RL Summarization ───────────────────────
        selected = _rl_summarize(sentences, n_episodes=40)
        summary  = " ".join(selected)

        # Key decisions: câu đầu + câu cuối (mang thông tin bao quát)
        n     = len(sentences)
        top_k = max(1, n // 5)
        key_decisions = " ".join(
            dict.fromkeys(sentences[:top_k] + sentences[-top_k:])
        )

        # Action items: câu chứa từ khóa hành động
        action_sents = [
            s for s in sentences
            if any(kw in s.lower() for kw in _ACTION_KEYWORDS)
        ]
        action_items = " ".join(action_sents) if action_sents else "Không có action items cụ thể."

        # ── Bước 5: Phân tích phụ trợ ────────────────────────
        full_text    = " ".join(sentences)
        sent         = _sentiment(full_text)
        final_reward = _jaccard(selected, sentences) if selected else 0.0
        eff          = _efficiency_score(len(selected), n, final_reward)

        # ── Bước 5: Output ────────────────────────────────────
        return jsonify({
            "status":           "success",
            "transcript":       sentences,
            "summary":          summary,
            "key_decisions":    key_decisions,
            "action_items":     action_items,
            "sentiment":        sent,
            "efficiency_score": eff,
        })

    except Exception as e:
        return jsonify({"status": "error", "message": str(e)}), 500

    finally:
        os.unlink(tmp_path)
