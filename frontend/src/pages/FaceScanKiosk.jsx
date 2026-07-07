import { useState, useEffect, useRef } from 'react';
import { timekeepingAPI, employeesAPI } from '../api/client.js';
import { useToast } from '../hooks/useToast.js';

export function FaceScanKiosk() {
	const toast = useToast();
	const videoRef = useRef(null);
	const [employees, setEmployees] = useState([]);
	const [selectedEmployee, setSelectedEmployee] = useState('');
	const [scanning, setScanning] = useState(false);
	const [scanSuccess, setScanSuccess] = useState(null);
	const [cameraStream, setCameraStream] = useState(null);
	const [cameraError, setCameraError] = useState(false);
	const streamRef = useRef(null);
	const isMountedRef = useRef(true);

	const startCamera = async () => {
		setCameraError(false);
		try {
			const stream = await navigator.mediaDevices.getUserMedia({
				video: { width: 640, height: 480, facingMode: "user" }
			});
			if (!isMountedRef.current) {
				stream.getTracks().forEach(track => track.stop());
				return;
			}
			setCameraStream(stream);
			streamRef.current = stream;
		} catch (err) {
			console.error("Camera access error:", err);
			setCameraError(true);
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
	}, [cameraStream, videoRef]);

	const stopCamera = () => {
		if (streamRef.current) {
			streamRef.current.getTracks().forEach(track => track.stop());
			streamRef.current = null;
		}
		setCameraStream(null);
	};

	const handleStartScan = async () => {
		if (!selectedEmployee) {
			toast.warning('Chú ý', 'Vui lòng chọn nhân viên để chấm công');
			return;
		}

		if (!cameraStream && !cameraError) {
			await startCamera();
		}

		setScanning(true);
		setScanSuccess(null);

		setTimeout(async () => {
			try {
				const emp = employees.find(e => e.code === selectedEmployee);
				const payload = {
					employee_code: selectedEmployee,
					timestamp: new Date().toISOString(),
					location_gps: "21.0285,105.7823",
					device_id: "KIOSK-GATE-01"
				};

				const res = await timekeepingAPI.checkin(payload);
				setScanSuccess({
					name: emp ? emp.full_name : selectedEmployee,
					code: selectedEmployee,
					time: new Date().toLocaleTimeString('vi-VN')
				});
				toast.success('Thành công', res.data?.message || 'Điểm danh thành công!');

				setTimeout(() => {
					setScanSuccess(null);
				}, 3000);
			} catch (err) {
				console.error("Checkin error:", err);
				const errMsg = err.response?.data?.message || 'Lỗi chấm công khuôn mặt';
				toast.error('Lỗi', errMsg);
			} finally {
				setScanning(false);
			}
		}, 3000);
	};

	return (
		<div className="page-container" style={{ maxWidth: 800, margin: '0 auto' }}>
			<div className="page-header" style={{ marginBottom: 24, textAlign: 'center' }}>
				<h1 className="page-title">Trạm Chấm Công Khuôn Mặt (AI Face Kiosk)</h1>
				<p className="page-subtitle" style={{ color: 'var(--text-muted)' }}>
					Giả lập luồng nhận diện khuôn mặt Edge AI của nhân viên gửi trực tiếp về hệ thống
				</p>
			</div>

			<div className="card" style={{ padding: 24, display: 'flex', flexDirection: 'column', alignItems: 'center', background: 'var(--bg-card)', borderRadius: 'var(--radius-lg)' }}>
				<div style={{ width: '100%', maxWidth: 400, marginBottom: 20 }}>
					<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Chọn nhân viên chấm công (Thử nghiệm)</label>
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
				</div>

				<div style={{ position: 'relative', width: '100%', maxWidth: 480, height: 360, background: '#111', borderRadius: 'var(--radius-lg)', overflow: 'hidden', display: 'flex', justifyContent: 'center', alignItems: 'center', border: '2px solid var(--border-color)', boxShadow: '0 8px 30px rgba(0,0,0,0.5)' }}>
					{cameraStream ? (
						<video
							ref={videoRef}
							autoPlay
							playsInline
							muted
							style={{ width: '100%', height: '100%', objectFit: 'cover' }}
						/>
					) : (
						<div style={{ textAlign: 'center', padding: 20 }}>
							<div className="pulse-circle" style={{ width: 80, height: 80, borderRadius: '50%', border: '4px solid var(--accent)', margin: '0 auto 16px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
								👤
							</div>
							<p style={{ color: 'var(--text-muted)', fontSize: 14 }}>Nhấn Bắt đầu quét để mở camera</p>
						</div>
					)}

					{scanning && (
						<div style={{ position: 'absolute', inset: 0, border: '4px solid var(--accent)', pointerEvents: 'none', animation: 'pulseBorder 1.5s infinite' }}>
							<div className="scanner-line" style={{ position: 'absolute', top: 0, left: 0, right: 0, height: 4, background: 'var(--accent)', boxShadow: '0 0 15px var(--accent)', animation: 'scanMove 2s linear infinite' }} />
							<div style={{ position: 'absolute', bottom: 16, left: '50%', transform: 'translateX(-50%)', background: 'rgba(0,0,0,0.75)', color: 'var(--accent)', padding: '6px 16px', borderRadius: 20, fontSize: 12, fontWeight: 700, letterSpacing: 1 }}>
								ĐANG QUÉT KHUÔN MẶT...
							</div>
						</div>
					)}

					{scanSuccess && !scanning && (
						<div style={{ position: 'absolute', inset: 0, background: 'rgba(9, 30, 24, 0.9)', display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', animation: 'fadeIn 0.3s ease' }}>
							<div style={{ width: 64, height: 64, borderRadius: '50%', background: '#2e7d32', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 32, color: '#fff', marginBottom: 16, boxShadow: '0 0 20px rgba(46,125,50,0.5)' }}>
								✓
							</div>
							<h3 style={{ fontSize: 18, fontWeight: 700, color: '#4caf50', marginBottom: 4 }}>XÁC NHẬN THÀNH CÔNG</h3>
							<p style={{ fontSize: 16, fontWeight: 600, color: 'var(--text-primary)', marginBottom: 2 }}>{scanSuccess.name}</p>
							<p style={{ fontSize: 13, color: 'var(--text-muted)' }}>Mã: {scanSuccess.code}</p>
							<p style={{ fontSize: 13, color: 'var(--text-muted)', marginTop: 8 }}>Thời gian: {scanSuccess.time}</p>
						</div>
					)}
				</div>

				<div style={{ marginTop: 24, display: 'flex', gap: 12 }}>
					<button
						className="btn btn-primary"
						onClick={handleStartScan}
						disabled={scanning}
						style={{ padding: '12px 24px', fontSize: 15, fontWeight: 700, display: 'flex', alignItems: 'center', gap: 8 }}
					>
						{scanning ? '🔄 Đang xác thực...' : '📸 Bắt đầu chấm công'}
					</button>
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
