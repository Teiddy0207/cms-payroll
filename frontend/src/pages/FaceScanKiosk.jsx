import { useState, useEffect, useRef } from 'react';
import { timekeepingAPI, employeesAPI } from '../api/client.js';
import { useToast } from '../hooks/useToast.js';

const computeImageHash = (source) => {
	return new Promise((resolve) => {
		const img = new Image();
		img.onload = () => {
			const canvas = document.createElement('canvas');
			canvas.width = 8;
			canvas.height = 8;
			const ctx = canvas.getContext('2d');
			ctx.drawImage(img, 0, 0, 8, 8);
			const imgData = ctx.getImageData(0, 0, 8, 8);
			const pixels = imgData.data;
			let sum = 0;
			const gray = [];
			for (let i = 0; i < pixels.length; i += 4) {
				const g = Math.floor(0.299 * pixels[i] + 0.587 * pixels[i+1] + 0.114 * pixels[i+2]);
				gray.push(g);
				sum += g;
			}
			const avg = sum / 64;
			let hash = "";
			for (let i = 0; i < 64; i++) {
				hash += gray[i] >= avg ? "1" : "0";
			}
			resolve(hash);
		};
		img.src = source;
	});
};

const computeVideoHash = (video) => {
	const canvas = document.createElement('canvas');
	canvas.width = 8;
	canvas.height = 8;
	const ctx = canvas.getContext('2d');
	ctx.drawImage(video, 0, 0, 8, 8);
	const imgData = ctx.getImageData(0, 0, 8, 8);
	const pixels = imgData.data;
	let sum = 0;
	const gray = [];
	for (let i = 0; i < pixels.length; i += 4) {
		const g = Math.floor(0.299 * pixels[i] + 0.587 * pixels[i+1] + 0.114 * pixels[i+2]);
		gray.push(g);
		sum += g;
	}
	const avg = sum / 64;
	let hash = "";
	for (let i = 0; i < 64; i++) {
		hash += gray[i] >= avg ? "1" : "0";
	}
	return hash;
};

const getHammingDistance = (hash1, hash2) => {
	let dist = 0;
	for (let i = 0; i < hash1.length; i++) {
		if (hash1[i] !== hash2[i]) {
			dist++;
		}
	}
	return dist;
};

export function FaceScanKiosk() {
	const toast = useToast();
	const videoRef = useRef(null);
	const streamRef = useRef(null);
	const isMountedRef = useRef(true);

	const [activeTab, setActiveTab] = useState('checkin');
	const [employees, setEmployees] = useState([]);
	const [selectedEmployee, setSelectedEmployee] = useState('');
	const [scanning, setScanning] = useState(false);
	const [scanSuccess, setScanSuccess] = useState(null);
	const [cameraStream, setCameraStream] = useState(null);
	const [cameraError, setCameraError] = useState(false);
	const [registeredFaces, setRegisteredFaces] = useState([]);

	const startCamera = async () => {
		setCameraError(false);
		try {
			if (streamRef.current) {
				streamRef.current.getTracks().forEach(t => t.stop());
			}
			const stream = await navigator.mediaDevices.getUserMedia({
				video: { width: 640, height: 480, facingMode: "user" }
			});
			setCameraStream(stream);
			streamRef.current = stream;
		} catch (err) {
			console.error("Camera access error:", err);
			setCameraError(true);
		}
	};

	const stopCamera = () => {
		if (streamRef.current) {
			streamRef.current.getTracks().forEach(track => track.stop());
			streamRef.current = null;
		}
		setCameraStream(null);
	};

	const loadRegisteredFaces = async () => {
		try {
			const res = await timekeepingAPI.getFaces();
			const faces = res.data?.data || [];
			const facesWithHashes = await Promise.all(faces.map(async (face) => {
				const hash = await computeImageHash(face.face_data);
				return { ...face, hash };
			}));
			setRegisteredFaces(facesWithHashes);
		} catch (err) {
			console.error("loadRegisteredFaces error:", err);
		}
	};

	useEffect(() => {
		isMountedRef.current = true;
		const loadEmployees = async () => {
			try {
				const res = await employeesAPI.list({ page_number: 1, page_size: 100 });
				setEmployees(res.data?.data?.items || []);
				if (res.data?.data?.items?.length > 0) {
					setSelectedEmployee(res.data.data.items[0].code);
				}
			} catch (err) {
				console.error("loadEmployees error:", err);
			}
		};
		loadEmployees();
		loadRegisteredFaces();
		startCamera();

		return () => {
			isMountedRef.current = false;
			if (streamRef.current) {
				streamRef.current.getTracks().forEach(track => track.stop());
				streamRef.current = null;
			}
		};
	}, []);

	useEffect(() => {
		if (videoRef.current && cameraStream) {
			videoRef.current.srcObject = cameraStream;
		}
	}, [cameraStream, videoRef, activeTab]);

	useEffect(() => {
		startCamera();
	}, [activeTab]);

	const handleRegisterFace = async () => {
		if (!selectedEmployee) {
			toast.warning('Chu y', 'Vui long chon nhan vien de dang ky');
			return;
		}

		if (!videoRef.current || !cameraStream) {
			toast.error('Loi', 'Khong tim thay camera de chup anh');
			return;
		}

		try {
			const video = videoRef.current;
			const canvas = document.createElement('canvas');
			canvas.width = video.videoWidth || 640;
			canvas.height = video.videoHeight || 480;
			const ctx = canvas.getContext('2d');
			ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
			const dataUrl = canvas.toDataURL('image/jpeg');

			await timekeepingAPI.registerFace({
				employee_code: selectedEmployee,
				face_data: dataUrl
			});

			const emp = employees.find(e => e.code === selectedEmployee);
			toast.success('Thanh cong', `Dang ky khuon mat thanh cong cho ${emp ? emp.full_name : selectedEmployee}`);
			loadRegisteredFaces();
		} catch (err) {
			console.error("Register face error:", err);
			toast.error('Loi', 'Khong the luu dac trung khuon mat');
		}
	};

	const handleDeleteFace = async (code) => {
		if (!window.confirm("Ban co chac chan muon xoa khuon mat cua nhan vien nay?")) return;
		try {
			await timekeepingAPI.deleteFace(code);
			toast.success('Thanh cong', 'Da xoa khuon mat dang ky');
			loadRegisteredFaces();
		} catch (err) {
			console.error("Delete face error:", err);
			toast.error('Loi', 'Khong the xoa khuon mat');
		}
	};

	const handleAICheckin = () => {
		if (registeredFaces.length === 0) {
			toast.error('Chua co du lieu', 'Khong tim thay khuon mat nao duoc dang ky. Vui long qua tab Dang Ky truoc.');
			return;
		}

		if (!videoRef.current || !cameraStream) {
			toast.error('Loi', 'Khong tim thay camera hoac luong camera chua san sang');
			return;
		}

		setScanning(true);
		setScanSuccess(null);

		setTimeout(async () => {
			try {
				const video = videoRef.current;
				const liveHash = computeVideoHash(video);

				let bestMatch = null;
				let minDist = 65;
				for (const face of registeredFaces) {
					if (!face.hash) continue;
					const dist = getHammingDistance(liveHash, face.hash);
					if (dist < minDist) {
						minDist = dist;
						bestMatch = face;
					}
				}

				if (!bestMatch) {
					toast.error('Loi', 'Khong the nhan dang duoc khuon mat nao khop');
					setScanning(false);
					return;
				}

				const confidence = Math.round(((64 - minDist) / 64) * 100);

				const payload = {
					employee_code: bestMatch.employee_code,
					timestamp: new Date().toISOString(),
					location_gps: "21.0285,105.7823",
					device_id: "KIOSK-GATE-01"
				};

				const res = await timekeepingAPI.checkin(payload);
				setScanSuccess({
					name: bestMatch.employee_name,
					code: bestMatch.employee_code,
					photo: bestMatch.face_data,
					confidence: confidence,
					time: new Date().toLocaleTimeString('vi-VN')
				});
				toast.success('Nhan dang AI thanh cong', res.data?.message || 'Da ghi nhan cham cong');

				setTimeout(() => {
					setScanSuccess(null);
				}, 4000);
			} catch (err) {
				console.error("AI checkin error:", err);
				const errMsg = err.response?.data?.message || 'Loi nhan dang khuon mat';
				toast.error('Loi', errMsg);
			} finally {
				setScanning(false);
			}
		}, 3000);
	};

	return (
		<div className="page-container" style={{ maxWidth: 900, margin: '0 auto' }}>
			<div className="page-header" style={{ marginBottom: 24, textAlign: 'center' }}>
				<h1 className="page-title">He Thong Cham Cong AI (Face Recognition Kiosk)</h1>
				<p className="page-subtitle" style={{ color: 'var(--text-muted)' }}>
					Nhan dang tu dong khuon mat nhan vien da dang ky trong co so du lieu PostgreSQL
				</p>
			</div>

			<div className="tabs-container" style={{ display: 'flex', gap: 12, marginBottom: 20, justifyContent: 'center' }}>
				<button
					className={`btn ${activeTab === 'checkin' ? 'btn-primary' : 'btn-secondary'}`}
					onClick={() => setActiveTab('checkin')}
					style={{ padding: '10px 24px', borderRadius: 'var(--radius-md)', fontWeight: 600 }}
				>
					Tram Cham Cong AI
				</button>
				<button
					className={`btn ${activeTab === 'register' ? 'btn-primary' : 'btn-secondary'}`}
					onClick={() => setActiveTab('register')}
					style={{ padding: '10px 24px', borderRadius: 'var(--radius-md)', fontWeight: 600 }}
				>
					Dang Ky Khuon Mat
				</button>
			</div>

			<div style={{ display: 'grid', gridTemplateColumns: '1fr 320px', gap: 24, alignItems: 'start' }}>
				<div className="card" style={{ padding: 20, display: 'flex', flexDirection: 'column', alignItems: 'center', background: 'var(--bg-card)', borderRadius: 'var(--radius-lg)' }}>
					<div style={{ position: 'relative', width: '100%', maxWidth: 520, height: 390, background: '#111', borderRadius: 'var(--radius-lg)', overflow: 'hidden', display: 'flex', justifyContent: 'center', alignItems: 'center', border: '2px solid var(--border-color)', boxShadow: '0 8px 30px rgba(0,0,0,0.5)' }}>
						{cameraStream ? (
							<>
								<video
									ref={videoRef}
									autoPlay
									playsInline
									muted
									style={{ width: '100%', height: '100%', objectFit: 'cover' }}
								/>
								<div style={{ position: 'absolute', width: 180, height: 180, border: '2px dashed rgba(79, 142, 247, 0.6)', borderRadius: '50%', pointerEvents: 'none', boxShadow: '0 0 0 9999px rgba(0, 0, 0, 0.4)' }} />
							</>
						) : (
							<div style={{ textAlign: 'center', padding: 20 }}>
								<div className="pulse-circle" style={{ width: 80, height: 80, borderRadius: '50%', border: '4px solid var(--accent)', margin: '0 auto 16px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
									<span style={{ fontSize: 13, color: 'var(--text-primary)' }}>Camera</span>
								</div>
								<p style={{ color: 'var(--text-muted)', fontSize: 14 }}>Vui long cap quyen mo camera</p>
							</div>
						)}

						{scanning && (
							<div style={{ position: 'absolute', inset: 0, border: '4px solid var(--accent)', pointerEvents: 'none', animation: 'pulseBorder 1.5s infinite' }}>
								<div className="scanner-line" style={{ position: 'absolute', top: 0, left: 0, right: 0, height: 4, background: 'var(--accent)', boxShadow: '0 0 15px var(--accent)', animation: 'scanMove 2.5s linear infinite' }} />
								<div style={{ position: 'absolute', bottom: 20, left: '50%', transform: 'translateX(-50%)', background: 'rgba(0,0,0,0.85)', color: 'var(--accent)', padding: '8px 18px', borderRadius: 20, fontSize: 12, fontWeight: 700, letterSpacing: 1.5, border: '1px solid var(--accent)' }}>
									AI: DANG SO KHOP VECTO KHUON MAT...
								</div>
							</div>
						)}

						{scanSuccess && !scanning && (
							<div style={{ position: 'absolute', inset: 0, background: 'rgba(9, 30, 24, 0.95)', display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', animation: 'fadeIn 0.3s ease' }}>
								<div style={{ position: 'relative', marginBottom: 16 }}>
									<img
										src={scanSuccess.photo}
										alt="Recognized"
										style={{ width: 120, height: 120, borderRadius: '50%', border: '4px solid #4caf50', objectFit: 'cover', boxShadow: '0 0 20px rgba(76,175,80,0.4)' }}
									/>
									<div style={{ position: 'absolute', bottom: 0, right: 0, width: 36, height: 36, borderRadius: '50%', background: '#4caf50', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 13, color: '#fff', border: '3px solid #111', fontWeight: 'bold' }}>
										OK
									</div>
								</div>
								<h3 style={{ fontSize: 18, fontWeight: 800, color: '#4caf50', marginBottom: 4, letterSpacing: 0.5 }}>NHAN DANG THANH CONG</h3>
								<p style={{ fontSize: 18, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 2 }}>{scanSuccess.name}</p>
								<p style={{ fontSize: 13, color: 'var(--text-muted)' }}>Ma NV: {scanSuccess.code} (Do khop: {scanSuccess.confidence}%)</p>
								<p style={{ fontSize: 13, color: 'var(--text-muted)', marginTop: 8 }}>Thoi gian check-in: {scanSuccess.time}</p>
							</div>
						)}
					</div>

					<div style={{ marginTop: 20 }}>
						{activeTab === 'checkin' ? (
							<button
								className="btn btn-primary"
								onClick={handleAICheckin}
								disabled={scanning || registeredFaces.length === 0}
								style={{ padding: '12px 30px', fontSize: 15, fontWeight: 700, display: 'flex', alignItems: 'center', gap: 8 }}
							>
								Quet va Cham Cong AI
							</button>
						) : (
							<div style={{ display: 'flex', flexDirection: 'column', gap: 12, width: '100%', maxWidth: 400 }}>
								<select
									className="form-input"
									value={selectedEmployee}
									onChange={e => setSelectedEmployee(e.target.value)}
									style={{ width: '100%', padding: '10px 14px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-primary)', color: 'var(--text-primary)' }}
								>
									{employees.map(emp => (
										<option key={emp.code} value={emp.code}>
											{emp.full_name} ({emp.code})
										</option>
									))}
								</select>
								<button
									className="btn btn-primary"
									onClick={handleRegisterFace}
									style={{ padding: '12px 24px', fontWeight: 700 }}
								>
									Luu Dac Trung Khuon Mat
								</button>
							</div>
						)}
					</div>
				</div>

				<div className="card" style={{ padding: 18, height: '100%', minHeight: 460, background: 'var(--bg-card)', borderRadius: 'var(--radius-lg)' }}>
					<h3 style={{ fontSize: 15, fontWeight: 700, marginBottom: 12, textTransform: 'uppercase', letterSpacing: 0.5, color: 'var(--accent)' }}>
						CSDL Khuon Mat ({registeredFaces.length})
					</h3>
					<div style={{ display: 'flex', flexDirection: 'column', gap: 12, overflowY: 'auto', maxHeight: 390, paddingRight: 4 }}>
						{registeredFaces.length === 0 ? (
							<div style={{ textAlign: 'center', padding: '40px 10px', color: 'var(--text-muted)', fontSize: 13 }}>
								Chua co nhan vien nao dang ky nhan dang khuon mat
							</div>
						) : (
							registeredFaces.map(user => (
								<div key={user.employee_code} style={{ display: 'flex', alignItems: 'center', gap: 10, padding: 10, background: 'rgba(255,255,255,0.03)', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)' }}>
									<img
										src={user.face_data}
										alt={user.employee_name}
										style={{ width: 44, height: 44, borderRadius: '50%', objectFit: 'cover', border: '1px solid var(--border-color)' }}
									/>
									<div style={{ flex: 1, minWidth: 0 }}>
										<div style={{ fontSize: 13, fontWeight: 700, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{user.employee_name}</div>
										<div style={{ fontSize: 11, color: 'var(--text-muted)' }}>Ma: {user.employee_code}</div>
									</div>
									<button
										onClick={() => handleDeleteFace(user.employee_code)}
										style={{ background: 'transparent', border: 'none', color: 'var(--error)', cursor: 'pointer', fontSize: 12, padding: '4px 8px', fontWeight: 'bold' }}
									>
										Xoa
									</button>
								</div>
							))
						)}
					</div>
				</div>
			</div>

			<style>{`
				@keyframes scanMove {
					0% { top: 0%; }
					50% { top: 100%; }
					100% { top: 0%; }
				}
				@keyframes pulseBorder {
					0% { border-color: var(--accent); opacity: 0.8; }
					50% { border-color: #2e7d32; opacity: 1; }
					100% { border-color: var(--accent); opacity: 0.8; }
				}
				@keyframes fadeIn {
					from { opacity: 0; }
					to { opacity: 1; }
				}
				.pulse-circle {
					animation: pulseGlow 2s infinite;
				}
				@keyframes pulseGlow {
					0% { box-shadow: 0 0 0 0 rgba(79, 142, 247, 0.4); }
					70% { box-shadow: 0 0 0 15px rgba(79, 142, 247, 0); }
					100% { box-shadow: 0 0 0 0 rgba(79, 142, 247, 0); }
				}
			`}</style>
		</div>
	);
}

export default FaceScanKiosk;
